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
)

func TestClient_GetTickers_Error(t *testing.T) {
	const (
		apiKey    = "some api key"
		secretKey = "some secret key"
	)

	tests := []struct {
		name        string
		client      http.Client
		expectedErr error
	}{
		{
			name: "returns error if error from server",
			client: http.Client{
				Transport: roundTripper{
					err: errors.New("some error"),
				},
			},
			expectedErr: errors.New("some error"),
		},
		{
			name: "returns 10001 SYS_ERROR",
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: cdcexchange.TickerResponse{
						APIBaseResponse: cdcexchange.APIBaseResponse{
							Code: "10001",
						},
					},
				},
			},
			expectedErr: cdcexchange.ResponseError{
				Code: 10001,
				Err:  cdcexchange.ErrSystemError,
			},
		},
		{
			name: "returns 100001 SYS_ERROR",
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: cdcexchange.TickerResponse{
						APIBaseResponse: cdcexchange.APIBaseResponse{
							Code: "100001",
						},
					},
				},
			},
			expectedErr: cdcexchange.ResponseError{
				Code: 100001,
				Err:  cdcexchange.ErrSystemError,
			},
		},
		{
			name: "returns 10002 UNAUTHORIZED",
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: cdcexchange.TickerResponse{
						APIBaseResponse: cdcexchange.APIBaseResponse{
							Code: "10002",
						},
					},
				},
			},
			expectedErr: cdcexchange.ResponseError{
				Code: 10002,
				Err:  cdcexchange.ErrUnauthorized,
			},
		},
		{
			name: "returns 10003 IP_ILLEGAL",
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: cdcexchange.TickerResponse{
						APIBaseResponse: cdcexchange.APIBaseResponse{
							Code: "10003",
						},
					},
				},
			},
			expectedErr: cdcexchange.ResponseError{
				Code: 10003,
				Err:  cdcexchange.ErrIllegalIP,
			},
		},
		{
			name: "returns 10004 BAD_REQUEST",
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: cdcexchange.TickerResponse{
						APIBaseResponse: cdcexchange.APIBaseResponse{
							Code: "10004",
						},
					},
				},
			},
			expectedErr: cdcexchange.ResponseError{
				Code: 10004,
				Err:  cdcexchange.ErrBadRequest,
			},
		},
		{
			name: "returns 10005 USER_TIER_INVALID",
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: cdcexchange.TickerResponse{
						APIBaseResponse: cdcexchange.APIBaseResponse{
							Code: "10005",
						},
					},
				},
			},
			expectedErr: cdcexchange.ResponseError{
				Code: 10005,
				Err:  cdcexchange.ErrUserTierInvalid,
			},
		},
		{
			name: "returns 10006 TOO_MANY_REQUESTS",
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: cdcexchange.TickerResponse{
						APIBaseResponse: cdcexchange.APIBaseResponse{
							Code: "10006",
						},
					},
				},
			},
			expectedErr: cdcexchange.ResponseError{
				Code: 10006,
				Err:  cdcexchange.ErrTooManyRequests,
			},
		},
		{
			name: "returns 10007 INVALID_NONCE",
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: cdcexchange.TickerResponse{
						APIBaseResponse: cdcexchange.APIBaseResponse{
							Code: "10007",
						},
					},
				},
			},
			expectedErr: cdcexchange.ResponseError{
				Code: 10007,
				Err:  cdcexchange.ErrInvalidNonce,
			},
		},
		{
			name: "returns 10008 METHOD_NOT_FOUND",
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: cdcexchange.TickerResponse{
						APIBaseResponse: cdcexchange.APIBaseResponse{
							Code: "10008",
						},
					},
				},
			},
			expectedErr: cdcexchange.ResponseError{
				Code: 10008,
				Err:  cdcexchange.ErrMethodNotFound,
			},
		},
		{
			name: "returns 10009 INVALID_DATE_RANGE",
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: cdcexchange.TickerResponse{
						APIBaseResponse: cdcexchange.APIBaseResponse{
							Code: "10009",
						},
					},
				},
			},
			expectedErr: cdcexchange.ResponseError{
				Code: 10009,
				Err:  cdcexchange.ErrInvalidDateRange,
			},
		},
		{
			name: "returns 20001 DUPLICATE_RECORD",
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: cdcexchange.TickerResponse{
						APIBaseResponse: cdcexchange.APIBaseResponse{
							Code: "20001",
						},
					},
				},
			},
			expectedErr: cdcexchange.ResponseError{
				Code: 20001,
				Err:  cdcexchange.ErrDuplicateRecord,
			},
		},
		{
			name: "returns 20002 NEGATIVE_BALANCE",
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: cdcexchange.TickerResponse{
						APIBaseResponse: cdcexchange.APIBaseResponse{
							Code: "20002",
						},
					},
				},
			},
			expectedErr: cdcexchange.ResponseError{
				Code: 20002,
				Err:  cdcexchange.ErrNegativeBalance,
			},
		},
		{
			name: "returns 30003 SYMBOL_NOT_FOUND",
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: cdcexchange.TickerResponse{
						APIBaseResponse: cdcexchange.APIBaseResponse{
							Code: "30003",
						},
					},
				},
			},
			expectedErr: cdcexchange.ResponseError{
				Code: 30003,
				Err:  cdcexchange.ErrSymbolNotFound,
			},
		},
		{
			name: "returns 30004 SIDE_NOT_SUPPORTED",
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: cdcexchange.TickerResponse{
						APIBaseResponse: cdcexchange.APIBaseResponse{
							Code: "30004",
						},
					},
				},
			},
			expectedErr: cdcexchange.ResponseError{
				Code: 30004,
				Err:  cdcexchange.ErrSideNotSupported,
			},
		},
		{
			name: "returns 30005 ORDERTYPE_NOT_SUPPORTED",
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: cdcexchange.TickerResponse{
						APIBaseResponse: cdcexchange.APIBaseResponse{
							Code: "30005",
						},
					},
				},
			},
			expectedErr: cdcexchange.ResponseError{
				Code: 30005,
				Err:  cdcexchange.ErrOrderTypeNotSupported,
			},
		},
		{
			name: "returns 30006 MIN_PRICE_VIOLATED",
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: cdcexchange.TickerResponse{
						APIBaseResponse: cdcexchange.APIBaseResponse{
							Code: "30006",
						},
					},
				},
			},
			expectedErr: cdcexchange.ResponseError{
				Code: 30006,
				Err:  cdcexchange.ErrMinPriceViolated,
			},
		},
		{
			name: "returns 30007 MAX_PRICE_VIOLATED",
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: cdcexchange.TickerResponse{
						APIBaseResponse: cdcexchange.APIBaseResponse{
							Code: "30007",
						},
					},
				},
			},
			expectedErr: cdcexchange.ResponseError{
				Code: 30007,
				Err:  cdcexchange.ErrMaxPriceViolated,
			},
		},
		{
			name: "returns 30008 MIN_QUANTITY_VIOLATED",
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: cdcexchange.TickerResponse{
						APIBaseResponse: cdcexchange.APIBaseResponse{
							Code: "30008",
						},
					},
				},
			},
			expectedErr: cdcexchange.ResponseError{
				Code: 30008,
				Err:  cdcexchange.ErrMinQuantityViolated,
			},
		},
		{
			name: "returns 30009 MAX_QUANTITY_VIOLATED",
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: cdcexchange.TickerResponse{
						APIBaseResponse: cdcexchange.APIBaseResponse{
							Code: "30009",
						},
					},
				},
			},
			expectedErr: cdcexchange.ResponseError{
				Code: 30009,
				Err:  cdcexchange.ErrMaxQuantityViolated,
			},
		},
		{
			name: "returns 30010 MISSING_ARGUMENT",
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: cdcexchange.TickerResponse{
						APIBaseResponse: cdcexchange.APIBaseResponse{
							Code: "30010",
						},
					},
				},
			},
			expectedErr: cdcexchange.ResponseError{
				Code: 30010,
				Err:  cdcexchange.ErrMissingArgument,
			},
		},
		{
			name: "returns 30013 INVALID_PRICE_PRECISION",
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: cdcexchange.TickerResponse{
						APIBaseResponse: cdcexchange.APIBaseResponse{
							Code: "30013",
						},
					},
				},
			},
			expectedErr: cdcexchange.ResponseError{
				Code: 30013,
				Err:  cdcexchange.ErrInvalidPricePrecision,
			},
		},
		{
			name: "returns 30014 INVALID_QUANTITY_PRECISION",
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: cdcexchange.TickerResponse{
						APIBaseResponse: cdcexchange.APIBaseResponse{
							Code: "30014",
						},
					},
				},
			},
			expectedErr: cdcexchange.ResponseError{
				Code: 30014,
				Err:  cdcexchange.ErrInvalidQuantityPrecision,
			},
		},
		{
			name: "returns 30016 MIN_NOTIONAL_VIOLATED",
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: cdcexchange.TickerResponse{
						APIBaseResponse: cdcexchange.APIBaseResponse{
							Code: "30016",
						},
					},
				},
			},
			expectedErr: cdcexchange.ResponseError{
				Code: 30016,
				Err:  cdcexchange.ErrMinNotionalViolated,
			},
		},
		{
			name: "returns 30017 MAX_NOTIONAL_VIOLATED",
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: cdcexchange.TickerResponse{
						APIBaseResponse: cdcexchange.APIBaseResponse{
							Code: "30017",
						},
					},
				},
			},
			expectedErr: cdcexchange.ResponseError{
				Code: 30017,
				Err:  cdcexchange.ErrMaxNotionalViolated,
			},
		},
		{
			name: "returns 30023 MIN_AMOUNT_VIOLATED",
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: cdcexchange.TickerResponse{
						APIBaseResponse: cdcexchange.APIBaseResponse{
							Code: "30023",
						},
					},
				},
			},
			expectedErr: cdcexchange.ResponseError{
				Code: 30023,
				Err:  cdcexchange.ErrMinAmountViolated,
			},
		},
		{
			name: "returns 30024 MAX_AMOUNT_VIOLATED",
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: cdcexchange.TickerResponse{
						APIBaseResponse: cdcexchange.APIBaseResponse{
							Code: "30024",
						},
					},
				},
			},
			expectedErr: cdcexchange.ResponseError{
				Code: 30024,
				Err:  cdcexchange.ErrMaxAmountViolated,
			},
		},
		{
			name: "returns 30025 AMOUNT_PRECISION_OVERFLOW",
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: cdcexchange.TickerResponse{
						APIBaseResponse: cdcexchange.APIBaseResponse{
							Code: "30025",
						},
					},
				},
			},
			expectedErr: cdcexchange.ResponseError{
				Code: 30025,
				Err:  cdcexchange.ErrAmountPrecisionOverflow,
			},
		},
		{
			name: "returns 40001 MG_INVALID_ACCOUNT_STATUS",
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: cdcexchange.TickerResponse{
						APIBaseResponse: cdcexchange.APIBaseResponse{
							Code: "40001",
						},
					},
				},
			},
			expectedErr: cdcexchange.ResponseError{
				Code: 40001,
				Err:  cdcexchange.ErrMGInvalidAccountStatus,
			},
		},
		{
			name: "returns 40002 MG_TRANSFER_ACTIVE_LOAN",
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: cdcexchange.TickerResponse{
						APIBaseResponse: cdcexchange.APIBaseResponse{
							Code: "40002",
						},
					},
				},
			},
			expectedErr: cdcexchange.ResponseError{
				Code: 40002,
				Err:  cdcexchange.ErrMGTransferActiveLoan,
			},
		},
		{
			name: "returns 40003 MG_INVALID_LOAN_CURRENCY",
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: cdcexchange.TickerResponse{
						APIBaseResponse: cdcexchange.APIBaseResponse{
							Code: "40003",
						},
					},
				},
			},
			expectedErr: cdcexchange.ResponseError{
				Code: 40003,
				Err:  cdcexchange.ErrMGInvalidLoanCurrency,
			},
		},
		{
			name: "returns 40004 MG_INVALID_REPAY_AMOUNT",
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: cdcexchange.TickerResponse{
						APIBaseResponse: cdcexchange.APIBaseResponse{
							Code: "40004",
						},
					},
				},
			},
			expectedErr: cdcexchange.ResponseError{
				Code: 40004,
				Err:  cdcexchange.ErrMGInvalidRepayAmount,
			},
		},
		{
			name: "returns 40005 MG_NO_ACTIVE_LOAN",
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: cdcexchange.TickerResponse{
						APIBaseResponse: cdcexchange.APIBaseResponse{
							Code: "40005",
						},
					},
				},
			},
			expectedErr: cdcexchange.ResponseError{
				Code: 40005,
				Err:  cdcexchange.ErrMGNoActiveLoan,
			},
		},
		{
			name: "returns 40006 MG_BLOCKED_BORROW",
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: cdcexchange.TickerResponse{
						APIBaseResponse: cdcexchange.APIBaseResponse{
							Code: "40006",
						},
					},
				},
			},
			expectedErr: cdcexchange.ResponseError{
				Code: 40006,
				Err:  cdcexchange.ErrMGBlockedBorrow,
			},
		},
		{
			name: "returns 40007 MG_BLOCKED_NEW_ORDER",
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: cdcexchange.TickerResponse{
						APIBaseResponse: cdcexchange.APIBaseResponse{
							Code: "40007",
						},
					},
				},
			},
			expectedErr: cdcexchange.ResponseError{
				Code: 40007,
				Err:  cdcexchange.ErrMGBlockedNewOrder,
			},
		},
		{
			name: "returns 50001 DW_CREDIT_LINE_NOT_MAINTAINED",
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: cdcexchange.TickerResponse{
						APIBaseResponse: cdcexchange.APIBaseResponse{
							Code: "50001",
						},
					},
				},
			},
			expectedErr: cdcexchange.ResponseError{
				Code: 50001,
				Err:  cdcexchange.ErrMGCreditLineNotMaintained,
			},
		},
		{
			name: "returns unexpected error",
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: cdcexchange.TickerResponse{
						APIBaseResponse: cdcexchange.APIBaseResponse{
							Code: "1",
						},
					},
				},
			},
			expectedErr: cdcexchange.ResponseError{
				Code: 1,
				Err:  cdcexchange.ErrUnexpectedError,
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

			if errors.Is(err, cdcexchange.ResponseError{}) {
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
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
	now := time.Now()

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

				res := cdcexchange.SingleTickerResponse{
					Result: cdcexchange.SingleTickerResult{
						Data: cdcexchange.Ticker{Instrument: instrument},
					},
				}

				require.NoError(t, json.NewEncoder(w).Encode(res))
			},
			expectedResult: []cdcexchange.Ticker{{Instrument: instrument}},
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

				res := cdcexchange.TickerResponse{
					Result: cdcexchange.TickerResult{
						Data: []cdcexchange.Ticker{{Instrument: instrument}},
					},
				}

				require.NoError(t, json.NewEncoder(w).Encode(res))
			},
			expectedResult: []cdcexchange.Ticker{{Instrument: instrument}},
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
