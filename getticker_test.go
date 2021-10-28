package cdcexchange_test

import (
	"context"
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
	cdctime "github.com/cshep4/crypto-dot-com-exchange-go/internal/time"
)

func TestClient_GetTickers_Error(t *testing.T) {
	const (
		apiKey    = "some api key"
		secretKey = "some secret key"
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
				now   = time.Now()
				clock = clockwork.NewFakeClockAt(now)
			)

			client, err := cdcexchange.New(apiKey, secretKey,
				cdcexchange.WithClock(clock),
				cdcexchange.WithHTTPClient(&tt.client),
			)
			require.NoError(t, err)

			tickers, err := client.GetTickers(ctx, "some instrument")
			require.Error(t, err)

			assert.Empty(t, tickers)

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

func TestClient_GetTickers_Success(t *testing.T) {
	const (
		apiKey     = "some api key"
		secretKey  = "some secret key"
		instrument = "some instrument"
	)
	now := time.Now().Round(time.Second)

	type args struct {
		instrument string
	}
	tests := []struct {
		name        string
		handlerFunc func(w http.ResponseWriter, r *http.Request)
		args
		expectedResult []cdcexchange.Ticker
	}{
		{
			name: "returns tickers for specific instrument",
			args: args{
				instrument: instrument,
			},
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				assert.Contains(t, r.URL.Path, cdcexchange.MethodGetTicker)
				assert.Equal(t, http.MethodGet, r.Method)
				t.Cleanup(func() { require.NoError(t, r.Body.Close()) })

				require.Empty(t, r.Body)

				instrumentName := r.URL.Query().Get("instrument_name")
				assert.Equal(t, instrument, instrumentName)

				res := fmt.Sprintf(`{
							"id": 0,
							"method":"",
							"code":0,
							"result":{
								"data": {
									"i": "%s",
									"t": %d
								}
							}
						}`, instrument, now.UnixMilli())

				_, err := w.Write([]byte(res))
				require.NoError(t, err)
			},
			expectedResult: []cdcexchange.Ticker{{
				Instrument: instrument,
				Timestamp:  cdctime.Time(now),
			}},
		},
		{
			name: "returns all tickers",
			args: args{
				instrument: "",
			},
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				assert.Contains(t, r.URL.Path, cdcexchange.MethodGetTicker)
				assert.Equal(t, http.MethodGet, r.Method)
				t.Cleanup(func() { require.NoError(t, r.Body.Close()) })

				require.Empty(t, r.Body)

				assert.False(t, r.URL.Query().Has("instrument_name"))

				res := fmt.Sprintf(`{
							"id": 0,
							"method":"",
							"code":0,
							"result":{
								"data": [{
									"i": "%s",
									"t": %d
								}]
							}
						}`, instrument, now.UnixMilli())

				_, err := w.Write([]byte(res))
				require.NoError(t, err)
			},
			expectedResult: []cdcexchange.Ticker{{
				Instrument: instrument,
				Timestamp:  cdctime.Time(now),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			t.Cleanup(ctrl.Finish)

			var (
				clock = clockwork.NewFakeClockAt(now)
			)

			s := httptest.NewServer(http.HandlerFunc(tt.handlerFunc))
			t.Cleanup(s.Close)

			client, err := cdcexchange.New(apiKey, secretKey,
				cdcexchange.WithClock(clock),
				cdcexchange.WithHTTPClient(s.Client()),
				cdcexchange.WithBaseURL(fmt.Sprintf("%s/", s.URL)),
			)
			require.NoError(t, err)

			tickers, err := client.GetTickers(ctx, tt.instrument)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedResult, tickers)
		})
	}
}
