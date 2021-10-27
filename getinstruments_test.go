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
	id_mocks "github.com/cshep4/crypto-dot-com-exchange-go/internal/mocks/id"
)

func TestClient_GetInstruments_Error(t *testing.T) {
	const (
		apiKey    = "some api key"
		secretKey = "some secret key"
		id        = int64(1234)
	)
	testErr := errors.New("some error")

	tests := []struct {
		name        string
		client      http.Client
		responseErr error
		expectedErr error
	}{
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
				idGenerator = id_mocks.NewMockIDGenerator(ctrl)
				now         = time.Now()
				clock       = clockwork.NewFakeClockAt(now)
			)

			client, err := cdcexchange.New(apiKey, secretKey,
				cdcexchange.WithIDGenerator(idGenerator),
				cdcexchange.WithClock(clock),
				cdcexchange.WithHTTPClient(&tt.client),
			)
			require.NoError(t, err)

			idGenerator.EXPECT().Generate().Return(id)

			instruments, err := client.GetInstruments(ctx)
			require.Error(t, err)

			assert.Empty(t, instruments)

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

func TestClient_GetIstruments_Success(t *testing.T) {
	const (
		apiKey     = "some api key"
		secretKey  = "some secret key"
		id         = int64(1234)
		instrument = "some instrument"
	)
	now := time.Now()

	tests := []struct {
		name           string
		handlerFunc    func(w http.ResponseWriter, r *http.Request)
		expectedResult []cdcexchange.Instrument
	}{
		{
			name: "returns instruments",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				assert.Contains(t, r.URL.Path, cdcexchange.MethodGetInstruments)
				assert.Equal(t, http.MethodGet, r.Method)
				t.Cleanup(func() { require.NoError(t, r.Body.Close()) })

				var body api.Request
				require.NoError(t, json.NewDecoder(r.Body).Decode(&body))

				assert.Equal(t, cdcexchange.MethodGetInstruments, body.Method)
				assert.Equal(t, id, body.ID)
				assert.Equal(t, now.UnixMilli(), body.Nonce)
				assert.Empty(t, body.APIKey)
				assert.Empty(t, body.Signature)
				assert.Empty(t, map[string]interface{}{}, body.Params)

				res := cdcexchange.InstrumentsResponse{
					Result: cdcexchange.InstrumentResult{
						Instruments: []cdcexchange.Instrument{{InstrumentName: instrument}},
					},
				}

				require.NoError(t, json.NewEncoder(w).Encode(res))
			},
			expectedResult: []cdcexchange.Instrument{{InstrumentName: instrument}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			t.Cleanup(ctrl.Finish)

			var (
				idGenerator = id_mocks.NewMockIDGenerator(ctrl)
				clock       = clockwork.NewFakeClockAt(now)
			)

			s := httptest.NewServer(http.HandlerFunc(tt.handlerFunc))
			t.Cleanup(s.Close)

			client, err := cdcexchange.New(apiKey, secretKey,
				cdcexchange.WithIDGenerator(idGenerator),
				cdcexchange.WithClock(clock),
				cdcexchange.WithHTTPClient(s.Client()),
				cdcexchange.WithBaseURL(fmt.Sprintf("%s/", s.URL)),
			)
			require.NoError(t, err)

			idGenerator.EXPECT().Generate().Return(id)

			instruments, err := client.GetInstruments(ctx)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedResult, instruments)
		})
	}
}
