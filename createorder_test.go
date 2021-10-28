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
)

func TestClient_CreateOrder_Error(t *testing.T) {
	const (
		apiKey    = "some api key"
		secretKey = "some secret key"
		id        = int64(1234)
	)
	testErr := errors.New("some error")

	tests := []struct {
		name         string
		client       http.Client
		signatureErr error
		responseErr  error
		expectedErr  error
	}{
		{
			name:         "returns error given error generating signature",
			signatureErr: testErr,
			expectedErr:  testErr,
		},
		{
			name: "returns error given error making request",
			client: http.Client{
				Transport: roundTripper{
					err: testErr,
				},
			},
			expectedErr: testErr,
		},
		{
			name: "returns error given error response",
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

			idGenerator.EXPECT().Generate().Return(id)
			signatureGenerator.EXPECT().GenerateSignature(auth.SignatureRequest{
				APIKey:    apiKey,
				SecretKey: secretKey,
				ID:        id,
				Method:    cdcexchange.MethodCreateOrder,
				Timestamp: now.UnixMilli(),
				Params:    map[string]interface{}{},
			}).Return("signature", tt.signatureErr)

			res, err := client.CreateOrder(ctx, cdcexchange.CreateOrderRequest{})
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

func TestClient_CreateOrder_Success(t *testing.T) {
	const (
		apiKey    = "some api key"
		secretKey = "some secret key"
		id        = int64(1234)
		signature = "some signature"

		instrument   = "some instrument"
		orderSide    = cdcexchange.OrderSideBuy
		orderType    = cdcexchange.OrderTypeMarket
		price        = 1.234
		quantity     = 5.678
		notional     = 9.012
		clientOID    = "some client oid"
		timeInForce  = cdcexchange.TimeInForceGoodTilCancelled
		execInst     = cdcexchange.ExecInstPostOnly
		triggerPrice = 3.456

		orderID = "5678"
	)
	now := time.Now()

	type args struct {
		req cdcexchange.CreateOrderRequest
	}
	tests := []struct {
		name        string
		handlerFunc func(w http.ResponseWriter, r *http.Request)
		args
		expectedParams map[string]interface{}
		expectedResult cdcexchange.CreateOrderResult
	}{
		{
			name: "successfully creates an order",
			args: args{
				req: cdcexchange.CreateOrderRequest{
					InstrumentName: instrument,
					Side:           orderSide,
					Type:           orderType,
					Price:          price,
					Quantity:       quantity,
					Notional:       notional,
					ClientOID:      clientOID,
					TimeInForce:    timeInForce,
					ExecInst:       execInst,
					TriggerPrice:   triggerPrice,
				},
			},
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				assert.Contains(t, r.URL.Path, cdcexchange.MethodCreateOrder)
				t.Cleanup(func() { require.NoError(t, r.Body.Close()) })

				var body api.Request
				require.NoError(t, json.NewDecoder(r.Body).Decode(&body))

				assert.Equal(t, cdcexchange.MethodCreateOrder, body.Method)
				assert.Equal(t, id, body.ID)
				assert.Equal(t, apiKey, body.APIKey)
				assert.Equal(t, now.UnixMilli(), body.Nonce)
				assert.Equal(t, signature, body.Signature)
				assert.Equal(t, instrument, body.Params["instrument_name"])
				assert.Equal(t, string(orderSide), body.Params["side"])
				assert.Equal(t, string(orderType), body.Params["type"])
				assert.Equal(t, price, body.Params["price"])
				assert.Equal(t, quantity, body.Params["quantity"])
				assert.Equal(t, notional, body.Params["notional"])
				assert.Equal(t, clientOID, body.Params["client_oid"])
				assert.Equal(t, string(timeInForce), body.Params["time_in_force"])
				assert.Equal(t, string(execInst), body.Params["exec_inst"])
				assert.Equal(t, triggerPrice, body.Params["trigger_price"])

				res := cdcexchange.CreateOrderResponse{
					BaseResponse: api.BaseResponse{},
					Result: cdcexchange.CreateOrderResult{
						ClientOID: clientOID,
						OrderID:   orderID,
					},
				}

				require.NoError(t, json.NewEncoder(w).Encode(res))
			},
			expectedParams: map[string]interface{}{
				"instrument_name": instrument,
				"side":            orderSide,
				"type":            orderType,
				"price":           price,
				"quantity":        quantity,
				"notional":        notional,
				"client_oid":      clientOID,
				"time_in_force":   timeInForce,
				"exec_inst":       execInst,
				"trigger_price":   triggerPrice,
			},
			expectedResult: cdcexchange.CreateOrderResult{
				ClientOID: clientOID,
				OrderID:   orderID,
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
				Method:    cdcexchange.MethodCreateOrder,
				Timestamp: now.UnixMilli(),
				Params:    tt.expectedParams,
			}).Return(signature, nil)

			res, err := client.CreateOrder(ctx, tt.req)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedResult, *res)
		})
	}
}
