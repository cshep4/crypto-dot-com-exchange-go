package cdcexchange_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cdcexchange "github.com/cshep4/crypto-dot-com-exchange-go"
	cdcerrors "github.com/cshep4/crypto-dot-com-exchange-go/errors"
	"github.com/cshep4/crypto-dot-com-exchange-go/internal/api"
	"github.com/cshep4/crypto-dot-com-exchange-go/internal/auth"
	id_mocks "github.com/cshep4/crypto-dot-com-exchange-go/internal/mocks/id"
	signature_mocks "github.com/cshep4/crypto-dot-com-exchange-go/internal/mocks/signature"
	cdctime "github.com/cshep4/crypto-dot-com-exchange-go/internal/time"
)

func TestClient_GetOrderDetail_Error(t *testing.T) {
	const (
		apiKey    = "some api key"
		secretKey = "some secret key"
		id        = int64(1234)
		orderID   = "some order id"
	)
	testErr := errors.New("some error")

	type args struct {
		orderID string
	}
	tests := []struct {
		name string
		args
		client       http.Client
		signatureErr error
		responseErr  error
		expectedErr  error
	}{
		{
			name: "returns error when order id is empty",
			args: args{
				orderID: "",
			},
			expectedErr: cdcerrors.InvalidParameterError{
				Parameter: "orderID",
				Reason:    "cannot be empty",
			},
		},
		{
			name: "returns error given error generating signature",
			args: args{
				orderID: orderID,
			},
			signatureErr: testErr,
			expectedErr:  testErr,
		},
		{
			name: "returns error given error making request",
			args: args{
				orderID: orderID,
			},
			client: http.Client{
				Transport: roundTripper{
					err: testErr,
				},
			},
			expectedErr: testErr,
		},
		{
			name: "returns error given error response",
			args: args{
				orderID: orderID,
			},
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: api.BaseResponse{
						Code: "10003",
					},
				},
			},
			responseErr: nil,
			expectedErr: cdcerrors.ResponseError{
				Code:           10003,
				HTTPStatusCode: http.StatusTeapot,
				Err:            cdcerrors.ErrIllegalIP,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			t.Cleanup(ctrl.Finish)

			var (
				idGenerator        = id_mocks.NewMockIDGenerator(ctrl)
				signatureGenerator = signature_mocks.NewMockSignatureGenerator(ctrl)
				now                = time.Now()
				clock              = clockwork.NewFakeClockAt(now)
			)

			client, err := cdcexchange.New(apiKey, secretKey,
				cdcexchange.WithIDGenerator(idGenerator),
				cdcexchange.WithClock(clock),
				cdcexchange.WithHTTPClient(&tt.client),
				cdcexchange.WithSignatureGenerator(signatureGenerator),
			)
			require.NoError(t, err)

			if tt.orderID != "" {
				idGenerator.EXPECT().Generate().Return(id)
				signatureGenerator.EXPECT().GenerateSignature(auth.SignatureRequest{
					APIKey:    apiKey,
					SecretKey: secretKey,
					ID:        id,
					Method:    cdcexchange.MethodGetOrderDetail,
					Timestamp: now.UnixMilli(),
					Params:    map[string]interface{}{"order_id": orderID},
				}).Return("signature", tt.signatureErr)
			}

			res, err := client.GetOrderDetail(ctx, tt.orderID)
			require.Error(t, err)

			assert.Empty(t, res)

			assert.True(t, errors.Is(err, tt.expectedErr))

			var expectedResponseError cdcerrors.ResponseError
			if errors.As(tt.expectedErr, &expectedResponseError) {
				var responseError cdcerrors.ResponseError
				require.True(t, errors.As(err, &responseError))

				assert.Equal(t, expectedResponseError.Code, responseError.Code)
				assert.Equal(t, expectedResponseError.HTTPStatusCode, responseError.HTTPStatusCode)
				assert.Equal(t, expectedResponseError.Err, responseError.Err)

				assert.True(t, errors.Is(err, expectedResponseError.Err))
			}
		})
	}
}

func TestClient_GetOrderDetail_Success(t *testing.T) {
	const (
		apiKey    = "some api key"
		secretKey = "some secret key"
		id        = int64(1234)
		signature = "some signature"
		orderID   = "some order id"

		clientOID = "some client oid"
	)
	now := time.Now().Round(time.Second)

	type args struct {
		orderID string
	}
	tests := []struct {
		name        string
		handlerFunc func(w http.ResponseWriter, r *http.Request)
		args
		expectedParams map[string]interface{}
		expectedResult cdcexchange.GetOrderDetailResult
	}{
		{
			name: "successfully gets order details",
			args: args{
				orderID: orderID,
			},
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				assert.Contains(t, r.URL.Path, cdcexchange.MethodGetOrderDetail)
				t.Cleanup(func() { require.NoError(t, r.Body.Close()) })

				var body api.Request
				require.NoError(t, json.NewDecoder(r.Body).Decode(&body))

				assert.Equal(t, cdcexchange.MethodGetOrderDetail, body.Method)
				assert.Equal(t, id, body.ID)
				assert.Equal(t, apiKey, body.APIKey)
				assert.Equal(t, now.UnixMilli(), body.Nonce)
				assert.Equal(t, signature, body.Signature)
				assert.Equal(t, orderID, body.Params["order_id"])

				res := fmt.Sprintf(`{
				  "id": 11,
				  "method": "private/get-order-detail",
				  "code": 0,
				  "result": {
					"trade_list": [
					  {
						"side": "BUY",
						"instrument_name": "ETH_CRO",
						"fee": 0.007,
						"trade_id": "371303044218155296",
						"create_time": %d,
						"traded_price": 7,
						"traded_quantity": 7,
						"fee_currency": "CRO",
						"order_id": "%s"
					  }
					],
					"order_info": {
					  "status": "FILLED",
					  "side": "BUY",
					  "order_id": "%s",
					  "client_oid": "%s",
					  "create_time": %d,
					  "update_time": %d,
					  "type": "LIMIT",
					  "instrument_name": "ETH_CRO",
					  "cumulative_quantity": 7,
					  "cumulative_value": 7,
					  "avg_price": 7,
					  "fee_currency": "CRO",
					  "time_in_force": "GOOD_TILL_CANCEL",
					  "exec_inst": "POST_ONLY"
					}
				  }
				}`, now.UnixMilli(), orderID, orderID, clientOID, now.UnixMilli(), now.UnixMilli())

				_, err := w.Write([]byte(res))
				require.NoError(t, err)
			},
			expectedParams: map[string]interface{}{
				"order_id": orderID,
			},
			expectedResult: cdcexchange.GetOrderDetailResult{
				TradeList: []cdcexchange.Trade{
					{
						Side:           cdcexchange.OrderSideBuy,
						InstrumentName: "ETH_CRO",
						Fee:            0.007,
						TradeID:        "371303044218155296",
						CreateTime:     cdctime.Time(now),
						TradedPrice:    7,
						TradedQuantity: 7,
						FeeCurrency:    "CRO",
						OrderID:        orderID,
					},
				},
				OrderInfo: cdcexchange.Order{
					Status:             cdcexchange.OrderStatusFilled,
					Side:               cdcexchange.OrderSideBuy,
					OrderID:            orderID,
					ClientOID:          clientOID,
					CreateTime:         cdctime.Time(now),
					UpdateTime:         cdctime.Time(now),
					OrderType:          cdcexchange.OrderTypeLimit,
					InstrumentName:     "ETH_CRO",
					CumulativeQuantity: 7,
					CumulativeValue:    7,
					AvgPrice:           7,
					FeeCurrency:        "CRO",
					TimeInForce:        cdcexchange.TimeInForceGoodTilCancelled,
					ExecInst:           cdcexchange.ExecInstPostOnly,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			t.Cleanup(ctrl.Finish)

			var (
				signatureGenerator = signature_mocks.NewMockSignatureGenerator(ctrl)
				idGenerator        = id_mocks.NewMockIDGenerator(ctrl)
				clock              = clockwork.NewFakeClockAt(now)
			)

			s := httptest.NewServer(http.HandlerFunc(tt.handlerFunc))
			t.Cleanup(s.Close)

			client, err := cdcexchange.New(apiKey, secretKey,
				cdcexchange.WithIDGenerator(idGenerator),
				cdcexchange.WithClock(clock),
				cdcexchange.WithHTTPClient(s.Client()),
				cdcexchange.WithBaseURL(fmt.Sprintf("%s/", s.URL)),
				cdcexchange.WithSignatureGenerator(signatureGenerator),
			)
			require.NoError(t, err)

			idGenerator.EXPECT().Generate().Return(id)
			signatureGenerator.EXPECT().GenerateSignature(auth.SignatureRequest{
				APIKey:    apiKey,
				SecretKey: secretKey,
				ID:        id,
				Method:    cdcexchange.MethodGetOrderDetail,
				Timestamp: now.UnixMilli(),
				Params:    tt.expectedParams,
			}).Return(signature, nil)

			res, err := client.GetOrderDetail(ctx, tt.orderID)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedResult, *res)
		})
	}
}
