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

func TestClient_GetAccountSummary_Error(t *testing.T) {
	const (
		apiKey    = "some api key"
		secretKey = "some secret key"
		id        = int64(1234)
	)
	testErr := errors.New("some error")

	type args struct {
		currency string
	}
	tests := []struct {
		name   string
		client http.Client
		args
		signatureErr error
		responseErr  error
		expectedErr  error
	}{
		{
			name: "returns error given error generating signature",
			args: args{
				currency: "currency",
			},
			signatureErr: testErr,
			expectedErr:  testErr,
		},
		{
			name: "returns error given error making request",
			args: args{
				currency: "currency",
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
				currency: "currency",
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

			idGenerator.EXPECT().Generate().Return(id)
			signatureGenerator.EXPECT().GenerateSignature(auth.SignatureRequest{
				APIKey:    apiKey,
				SecretKey: secretKey,
				ID:        id,
				Method:    cdcexchange.MethodGetAccountSummary,
				Timestamp: now.UnixMilli(),
				Params:    map[string]interface{}{"currency": tt.currency},
			}).Return("signature", tt.signatureErr)

			accounts, err := client.GetAccountSummary(ctx, tt.currency)
			require.Error(t, err)

			assert.Empty(t, accounts)

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

func TestClient_GetAccountSummary_Success(t *testing.T) {
	const (
		apiKey    = "some api key"
		secretKey = "some secret key"
		id        = int64(1234)
		currency  = "some currency"
		signature = "some signature"
	)
	now := time.Now()

	type args struct {
		currency string
	}
	tests := []struct {
		name        string
		handlerFunc func(w http.ResponseWriter, r *http.Request)
		args
		expectedParams map[string]interface{}
		expectedResult []cdcexchange.Account
	}{
		{
			name: "returns account summary for all tickers",
			args: args{
				currency: "",
			},
			expectedParams: map[string]interface{}{},
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				assert.Contains(t, r.URL.Path, cdcexchange.MethodGetAccountSummary)
				t.Cleanup(func() { require.NoError(t, r.Body.Close()) })

				var body api.Request
				require.NoError(t, json.NewDecoder(r.Body).Decode(&body))

				assert.Equal(t, cdcexchange.MethodGetAccountSummary, body.Method)
				assert.Equal(t, id, body.ID)
				assert.Equal(t, apiKey, body.APIKey)
				assert.Equal(t, now.UnixMilli(), body.Nonce)
				assert.Equal(t, signature, body.Signature)
				assert.Equal(t, map[string]interface{}{}, body.Params)

				res := cdcexchange.AccountSummaryResponse{
					BaseResponse: api.BaseResponse{},
					Result: cdcexchange.AccountSummaryResult{
						Accounts: []cdcexchange.Account{{Currency: currency}},
					},
				}

				require.NoError(t, json.NewEncoder(w).Encode(res))
			},
			expectedResult: []cdcexchange.Account{{Currency: currency}},
		},
		{
			name: "returns account summary for currency",
			args: args{
				currency: currency,
			},
			expectedParams: map[string]interface{}{"currency": currency},
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				assert.Contains(t, r.URL.Path, cdcexchange.MethodGetAccountSummary)
				t.Cleanup(func() { require.NoError(t, r.Body.Close()) })

				var body api.Request
				require.NoError(t, json.NewDecoder(r.Body).Decode(&body))

				assert.Equal(t, cdcexchange.MethodGetAccountSummary, body.Method)
				assert.Equal(t, id, body.ID)
				assert.Equal(t, apiKey, body.APIKey)
				assert.Equal(t, now.UnixMilli(), body.Nonce)
				assert.Equal(t, signature, body.Signature)
				assert.Equal(t, map[string]interface{}{"currency": currency}, body.Params)

				res := cdcexchange.AccountSummaryResponse{
					BaseResponse: api.BaseResponse{},
					Result: cdcexchange.AccountSummaryResult{
						Accounts: []cdcexchange.Account{{Currency: currency}},
					},
				}

				require.NoError(t, json.NewEncoder(w).Encode(res))
			},
			expectedResult: []cdcexchange.Account{{Currency: currency}},
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
				Method:    cdcexchange.MethodGetAccountSummary,
				Timestamp: now.UnixMilli(),
				Params:    tt.expectedParams,
			}).Return(signature, nil)

			accounts, err := client.GetAccountSummary(ctx, tt.currency)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedResult, accounts)
		})
	}
}
