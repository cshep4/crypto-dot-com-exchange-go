package cdcexchange_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cshep4/crypto-dot-com-exchange-go/internal/api"
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
	"github.com/cshep4/crypto-dot-com-exchange-go/internal/auth"
	id_mocks "github.com/cshep4/crypto-dot-com-exchange-go/internal/mocks/id"
	signature_mocks "github.com/cshep4/crypto-dot-com-exchange-go/internal/mocks/signature"
)

func TestClient_CancelOrder_Error(t *testing.T) {
	const (
		apiKey         = "some api key"
		secretKey      = "some secret key"
		id             = int64(1234)
		orderID        = "some order id"
		instrumentName = "some instrument name"
	)
	testErr := errors.New("some error")

	type args struct {
		instrumentName string
		orderID        string
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
			name: "returns error when instrument name is empty",
			args: args{
				instrumentName: "",
			},
			expectedErr: cdcerrors.InvalidParameterError{
				Parameter: "instrumentName",
				Reason:    "cannot be empty",
			},
		},
		{
			name: "returns error when order id is empty",
			args: args{
				orderID:        "",
				instrumentName: instrumentName,
			},
			expectedErr: cdcerrors.InvalidParameterError{
				Parameter: "orderID",
				Reason:    "cannot be empty",
			},
		},
		{
			name: "returns error given error generating signature",
			args: args{
				orderID:        orderID,
				instrumentName: instrumentName,
			},
			signatureErr: testErr,
			expectedErr:  testErr,
		},
		{
			name: "returns error given error making request",
			args: args{
				orderID:        orderID,
				instrumentName: instrumentName,
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
				orderID:        orderID,
				instrumentName: instrumentName,
			},
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: api.BaseResponse{
						Code: "10003",
					},
				},
			},
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

			if tt.orderID != "" && tt.instrumentName != "" {
				idGenerator.EXPECT().Generate().Return(id)
				signatureGenerator.EXPECT().GenerateSignature(auth.SignatureRequest{
					APIKey:    apiKey,
					SecretKey: secretKey,
					ID:        id,
					Method:    cdcexchange.MethodCancelOrder,
					Timestamp: now.UnixMilli(),
					Params: map[string]interface{}{
						"instrument_name": instrumentName,
						"order_id":        orderID,
					},
				}).Return("signature", tt.signatureErr)
			}

			err = client.CancelOrder(ctx, tt.instrumentName, tt.orderID)
			require.Error(t, err)

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

func TestClient_CancelOrder_Success(t *testing.T) {
	const (
		apiKey    = "some api key"
		secretKey = "some secret key"
		id        = int64(1234)
		signature = "some signature"

		orderID        = "some order id"
		instrumentName = "some instrument name"
	)
	now := time.Now()

	type args struct {
		instrumentName string
		orderID        string
	}
	tests := []struct {
		name        string
		handlerFunc func(w http.ResponseWriter, r *http.Request)
		args
		expectedParams map[string]interface{}
	}{
		{
			name: "successfully cancels an order",
			args: args{
				instrumentName: instrumentName,
				orderID:        orderID,
			},
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				assert.Contains(t, r.URL.Path, cdcexchange.MethodCancelOrder)
				t.Cleanup(func() { require.NoError(t, r.Body.Close()) })

				var body api.Request
				require.NoError(t, json.NewDecoder(r.Body).Decode(&body))

				assert.Equal(t, cdcexchange.MethodCancelOrder, body.Method)
				assert.Equal(t, id, body.ID)
				assert.Equal(t, apiKey, body.APIKey)
				assert.Equal(t, now.UnixMilli(), body.Nonce)
				assert.Equal(t, signature, body.Signature)
				assert.Equal(t, instrumentName, body.Params["instrument_name"])
				assert.Equal(t, orderID, body.Params["order_id"])

				res := cdcexchange.CancelOrderResponse{
					BaseResponse: api.BaseResponse{},
				}

				require.NoError(t, json.NewEncoder(w).Encode(res))
			},
			expectedParams: map[string]interface{}{
				"instrument_name": instrumentName,
				"order_id":            orderID,
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
				Method:    cdcexchange.MethodCancelOrder,
				Timestamp: now.UnixMilli(),
				Params:    tt.expectedParams,
			}).Return(signature, nil)

			err = client.CancelOrder(ctx, tt.instrumentName, tt.orderID)
			require.NoError(t, err)
		})
	}
}
