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

func TestClient_GetOpenOrders_Error(t *testing.T) {
	const (
		apiKey    = "some api key"
		secretKey = "some secret key"
		id        = int64(1234)
	)
	testErr := errors.New("some error")

	type args struct {
		req cdcexchange.GetOpenOrdersRequest
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
			name: "returns error when page size is less than 0",
			args: args{
				req: cdcexchange.GetOpenOrdersRequest{
					PageSize: -1,
				},
			},
			expectedErr: cdcerrors.InvalidParameterError{
				Parameter: "req.PageSize",
				Reason:    "cannot be less than 0",
			},
		},
		{
			name: "returns error when page size is greater than 200",
			args: args{
				req: cdcexchange.GetOpenOrdersRequest{
					PageSize: 201,
				},
			},
			expectedErr: cdcerrors.InvalidParameterError{
				Parameter: "req.PageSize",
				Reason:    "cannot be greater than 200",
			},
		},
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

			if tt.req.PageSize >= 0 && tt.req.PageSize < 200 {
				idGenerator.EXPECT().Generate().Return(id)
				signatureGenerator.EXPECT().GenerateSignature(auth.SignatureRequest{
					APIKey:    apiKey,
					SecretKey: secretKey,
					ID:        id,
					Method:    cdcexchange.MethodGetOpenOrders,
					Timestamp: now.UnixMilli(),
					Params:    map[string]interface{}{"page": 0},
				}).Return("signature", tt.signatureErr)
			}

			res, err := client.GetOpenOrders(ctx, tt.req)
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

func TestClient_GetOpenOrders_Success(t *testing.T) {
	const (
		apiKey    = "some api key"
		secretKey = "some secret key"
		id        = int64(1234)
		signature = "some signature"

		instrument   = "some instrument"
		clientOID    = "some client oid"
	)
	now := time.Now()

	type args struct {
		req cdcexchange.GetOpenOrdersRequest
	}
	tests := []struct {
		name        string
		handlerFunc func(w http.ResponseWriter, r *http.Request)
		args
		expectedParams map[string]interface{}
		expectedResult cdcexchange.GetOpenOrdersResult
	}{
		{
			name: "successfully gets all orders for an instrument",
			args: args{
				req: cdcexchange.GetOpenOrdersRequest{
					InstrumentName: instrument,
					PageSize:       100,
					Page:           1,
				},
			},
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				assert.Contains(t, r.URL.Path, cdcexchange.MethodGetOpenOrders)
				t.Cleanup(func() { require.NoError(t, r.Body.Close()) })

				var body api.Request
				require.NoError(t, json.NewDecoder(r.Body).Decode(&body))

				assert.Equal(t, cdcexchange.MethodGetOpenOrders, body.Method)
				assert.Equal(t, id, body.ID)
				assert.Equal(t, apiKey, body.APIKey)
				assert.Equal(t, now.UnixMilli(), body.Nonce)
				assert.Equal(t, signature, body.Signature)
				assert.Equal(t, instrument, body.Params["instrument_name"])
				assert.Equal(t, float64(100), body.Params["page_size"])
				assert.Equal(t, float64(1), body.Params["page"])

				res := cdcexchange.GetOpenOrdersResponse{
					BaseResponse: api.BaseResponse{},
					Result: cdcexchange.GetOpenOrdersResult{
						Count: 1234,
						OrderList: []cdcexchange.Order{
							{
								ClientOID: clientOID,
							},
						},
					},
				}

				require.NoError(t, json.NewEncoder(w).Encode(res))
			},
			expectedParams: map[string]interface{}{
				"instrument_name": instrument,
				"page_size":       100,
				"page":            1,
			},
			expectedResult: cdcexchange.GetOpenOrdersResult{
				Count: 1234,
				OrderList: []cdcexchange.Order{
					{ClientOID: clientOID},
				},
			},
		},
		{
			name: "successfully gets all orders for all instrument",
			args: args{
				req: cdcexchange.GetOpenOrdersRequest{},
			},
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				assert.Contains(t, r.URL.Path, cdcexchange.MethodGetOpenOrders)
				t.Cleanup(func() { require.NoError(t, r.Body.Close()) })

				var body api.Request
				require.NoError(t, json.NewDecoder(r.Body).Decode(&body))

				assert.Equal(t, cdcexchange.MethodGetOpenOrders, body.Method)
				assert.Equal(t, id, body.ID)
				assert.Equal(t, apiKey, body.APIKey)
				assert.Equal(t, now.UnixMilli(), body.Nonce)
				assert.Equal(t, signature, body.Signature)
				assert.Equal(t, float64(0), body.Params["page"])

				res := cdcexchange.GetOpenOrdersResponse{
					BaseResponse: api.BaseResponse{},
					Result: cdcexchange.GetOpenOrdersResult{
						Count: 1234,
						OrderList: []cdcexchange.Order{
							{
								ClientOID: clientOID,
							},
						},
					},
				}

				require.NoError(t, json.NewEncoder(w).Encode(res))
			},
			expectedParams: map[string]interface{}{
				"page": 0,
			},
			expectedResult: cdcexchange.GetOpenOrdersResult{
				Count: 1234,
				OrderList: []cdcexchange.Order{
					{ClientOID: clientOID},
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
				Method:    cdcexchange.MethodGetOpenOrders,
				Timestamp: now.UnixMilli(),
				Params:    tt.expectedParams,
			}).Return(signature, nil)

			res, err := client.GetOpenOrders(ctx, tt.req)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedResult, *res)
		})
	}
}
