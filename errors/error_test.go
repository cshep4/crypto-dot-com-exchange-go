package errors

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewResponseError_Error(t *testing.T) {
	type args struct {
		httpStatusCode int
		code           int64
	}
	tests := []struct {
		name   string
		client http.Client
		args
		expectedHTTPStatusCode int
		expectedCode           int64
		expectedErr            error
	}{
		{
			name: "returns unexpected error",
			args: args{
				httpStatusCode: http.StatusInternalServerError,
				code:           -1,
			},
			expectedHTTPStatusCode: http.StatusInternalServerError,
			expectedCode:           -1,
			expectedErr:            ErrUnexpectedError,
		},
		{
			name: "returns 10001 SYS_ERROR",
			args: args{
				httpStatusCode: http.StatusTeapot,
				code:           10001,
			},
			expectedHTTPStatusCode: http.StatusTeapot,
			expectedCode:           10001,
			expectedErr:            ErrSystemError,
		},
		{
			name: "returns 100001 SYS_ERROR",
			args: args{
				httpStatusCode: http.StatusTeapot,
				code:           100001,
			},
			expectedHTTPStatusCode: http.StatusTeapot,
			expectedCode:           100001,
			expectedErr:            ErrSystemError,
		},
		{
			name: "returns 10002 UNAUTHORIZED",
			args: args{
				httpStatusCode: http.StatusTeapot,
				code:           10002,
			},
			expectedHTTPStatusCode: http.StatusTeapot,
			expectedCode:           10002,
			expectedErr:            ErrUnauthorized,
		},
		{
			name: "returns 10003 IP_ILLEGAL",
			args: args{
				httpStatusCode: http.StatusTeapot,
				code:           10003,
			},
			expectedHTTPStatusCode: http.StatusTeapot,
			expectedCode:           10003,
			expectedErr:            ErrIllegalIP,
		},
		{
			name: "returns 10004 BAD_REQUEST",
			args: args{
				httpStatusCode: http.StatusTeapot,
				code:           10004,
			},
			expectedHTTPStatusCode: http.StatusTeapot,
			expectedCode:           10004,
			expectedErr:            ErrBadRequest,
		},
		{
			name: "returns 10005 USER_TIER_INVALID",
			args: args{
				httpStatusCode: http.StatusTeapot,
				code:           10005,
			},
			expectedHTTPStatusCode: http.StatusTeapot,
			expectedCode:           10005,
			expectedErr:            ErrUserTierInvalid,
		},
		{
			name: "returns 10006 TOO_MANY_REQUESTS",
			args: args{
				httpStatusCode: http.StatusTeapot,
				code:           10006,
			},
			expectedHTTPStatusCode: http.StatusTeapot,
			expectedCode:           10006,
			expectedErr:            ErrTooManyRequests,
		},
		{
			name: "returns 10007 INVALID_NONCE",
			args: args{
				httpStatusCode: http.StatusTeapot,
				code:           10007,
			},
			expectedHTTPStatusCode: http.StatusTeapot,
			expectedCode:           10007,
			expectedErr:            ErrInvalidNonce,
		},
		{
			name: "returns 10008 METHOD_NOT_FOUND",
			args: args{
				httpStatusCode: http.StatusTeapot,
				code:           10008,
			},
			expectedHTTPStatusCode: http.StatusTeapot,
			expectedCode:           10008,
			expectedErr:            ErrMethodNotFound,
		},
		{
			name: "returns 10009 INVALID_DATE_RANGE",
			args: args{
				httpStatusCode: http.StatusTeapot,
				code:           10009,
			},
			expectedHTTPStatusCode: http.StatusTeapot,
			expectedCode:           10009,
			expectedErr:            ErrInvalidDateRange,
		},
		{
			name: "returns 20001 DUPLICATE_RECORD",
			args: args{
				httpStatusCode: http.StatusTeapot,
				code:           20001,
			},
			expectedHTTPStatusCode: http.StatusTeapot,
			expectedCode:           20001,
			expectedErr:            ErrDuplicateRecord,
		},
		{
			name: "returns 20002 NEGATIVE_BALANCE",
			args: args{
				httpStatusCode: http.StatusTeapot,
				code:           20002,
			},
			expectedHTTPStatusCode: http.StatusTeapot,
			expectedCode:           20002,
			expectedErr:            ErrNegativeBalance,
		},
		{
			name: "returns 30003 SYMBOL_NOT_FOUND",
			args: args{
				httpStatusCode: http.StatusTeapot,
				code:           30003,
			},
			expectedHTTPStatusCode: http.StatusTeapot,
			expectedCode:           30003,
			expectedErr:            ErrSymbolNotFound,
		},
		{
			name: "returns 30004 SIDE_NOT_SUPPORTED",
			args: args{
				httpStatusCode: http.StatusTeapot,
				code:           30004,
			},
			expectedHTTPStatusCode: http.StatusTeapot,
			expectedCode:           30004,
			expectedErr:            ErrSideNotSupported,
		},
		{
			name: "returns 30005 ORDERTYPE_NOT_SUPPORTED",
			args: args{
				httpStatusCode: http.StatusTeapot,
				code:           30005,
			},
			expectedHTTPStatusCode: http.StatusTeapot,
			expectedCode:           30005,
			expectedErr:            ErrOrderTypeNotSupported,
		},
		{
			name: "returns 30006 MIN_PRICE_VIOLATED",
			args: args{
				httpStatusCode: http.StatusTeapot,
				code:           30006,
			},
			expectedHTTPStatusCode: http.StatusTeapot,
			expectedCode:           30006,
			expectedErr:            ErrMinPriceViolated,
		},
		{
			name: "returns 30007 MAX_PRICE_VIOLATED",
			args: args{
				httpStatusCode: http.StatusTeapot,
				code:           30007,
			},
			expectedHTTPStatusCode: http.StatusTeapot,
			expectedCode:           30007,
			expectedErr:            ErrMaxPriceViolated,
		},
		{
			name: "returns 30008 MIN_QUANTITY_VIOLATED",
			args: args{
				httpStatusCode: http.StatusTeapot,
				code:           30008,
			},
			expectedHTTPStatusCode: http.StatusTeapot,
			expectedCode:           30008,
			expectedErr:            ErrMinQuantityViolated,
		},
		{
			name: "returns 30009 MAX_QUANTITY_VIOLATED",
			args: args{
				httpStatusCode: http.StatusTeapot,
				code:           30009,
			},
			expectedHTTPStatusCode: http.StatusTeapot,
			expectedCode:           30009,
			expectedErr:            ErrMaxQuantityViolated,
		},
		{
			name: "returns 30010 MISSING_ARGUMENT",
			args: args{
				httpStatusCode: http.StatusTeapot,
				code:           30010,
			},
			expectedHTTPStatusCode: http.StatusTeapot,
			expectedCode:           30010,
			expectedErr:            ErrMissingArgument,
		},
		{
			name: "returns 30013 INVALID_PRICE_PRECISION",
			args: args{
				httpStatusCode: http.StatusTeapot,
				code:           30013,
			},
			expectedHTTPStatusCode: http.StatusTeapot,
			expectedCode:           30013,
			expectedErr:            ErrInvalidPricePrecision,
		},
		{
			name: "returns 30014 INVALID_QUANTITY_PRECISION",
			args: args{
				httpStatusCode: http.StatusTeapot,
				code:           30014,
			},
			expectedHTTPStatusCode: http.StatusTeapot,
			expectedCode:           30014,
			expectedErr:            ErrInvalidQuantityPrecision,
		},
		{
			name: "returns 30016 MIN_NOTIONAL_VIOLATED",
			args: args{
				httpStatusCode: http.StatusTeapot,
				code:           30016,
			},
			expectedHTTPStatusCode: http.StatusTeapot,
			expectedCode:           30016,
			expectedErr:            ErrMinNotionalViolated,
		},
		{
			name: "returns 30017 MAX_NOTIONAL_VIOLATED",
			args: args{
				httpStatusCode: http.StatusTeapot,
				code:           30017,
			},
			expectedHTTPStatusCode: http.StatusTeapot,
			expectedCode:           30017,
			expectedErr:            ErrMaxNotionalViolated,
		},
		{
			name: "returns 30023 MIN_AMOUNT_VIOLATED",
			args: args{
				httpStatusCode: http.StatusTeapot,
				code:           30023,
			},
			expectedHTTPStatusCode: http.StatusTeapot,
			expectedCode:           30023,
			expectedErr:            ErrMinAmountViolated,
		},
		{
			name: "returns 30024 MAX_AMOUNT_VIOLATED",
			args: args{
				httpStatusCode: http.StatusTeapot,
				code:           30024,
			},
			expectedHTTPStatusCode: http.StatusTeapot,
			expectedCode:           30024,
			expectedErr:            ErrMaxAmountViolated,
		},
		{
			name: "returns 30025 AMOUNT_PRECISION_OVERFLOW",
			args: args{
				httpStatusCode: http.StatusTeapot,
				code:           30025,
			},
			expectedHTTPStatusCode: http.StatusTeapot,
			expectedCode:           30025,
			expectedErr:            ErrAmountPrecisionOverflow,
		},
		{
			name: "returns 40001 MG_INVALID_ACCOUNT_STATUS",
			args: args{
				httpStatusCode: http.StatusTeapot,
				code:           40001,
			},
			expectedHTTPStatusCode: http.StatusTeapot,
			expectedCode:           40001,
			expectedErr:            ErrMGInvalidAccountStatus,
		},
		{
			name: "returns 40002 MG_TRANSFER_ACTIVE_LOAN",
			args: args{
				httpStatusCode: http.StatusTeapot,
				code:           40002,
			},
			expectedHTTPStatusCode: http.StatusTeapot,
			expectedCode:           40002,
			expectedErr:            ErrMGTransferActiveLoan,
		},
		{
			name: "returns 40003 MG_INVALID_LOAN_CURRENCY",
			args: args{
				httpStatusCode: http.StatusTeapot,
				code:           40003,
			},
			expectedHTTPStatusCode: http.StatusTeapot,
			expectedCode:           40003,
			expectedErr:            ErrMGInvalidLoanCurrency,
		},
		{
			name: "returns 40004 MG_INVALID_REPAY_AMOUNT",
			args: args{
				httpStatusCode: http.StatusTeapot,
				code:           40004,
			},
			expectedHTTPStatusCode: http.StatusTeapot,
			expectedCode:           40004,
			expectedErr:            ErrMGInvalidRepayAmount,
		},
		{
			name: "returns 40005 MG_NO_ACTIVE_LOAN",
			args: args{
				httpStatusCode: http.StatusTeapot,
				code:           40005,
			},
			expectedHTTPStatusCode: http.StatusTeapot,
			expectedCode:           40005,
			expectedErr:            ErrMGNoActiveLoan,
		},
		{
			name: "returns 40006 MG_BLOCKED_BORROW",
			args: args{
				httpStatusCode: http.StatusTeapot,
				code:           40006,
			},
			expectedHTTPStatusCode: http.StatusTeapot,
			expectedCode:           40006,
			expectedErr:            ErrMGBlockedBorrow,
		},
		{
			name: "returns 40007 MG_BLOCKED_NEW_ORDER",
			args: args{
				httpStatusCode: http.StatusTeapot,
				code:           40007,
			},
			expectedHTTPStatusCode: http.StatusTeapot,
			expectedCode:           40007,
			expectedErr:            ErrMGBlockedNewOrder,
		},
		{
			name: "returns 50001 DW_CREDIT_LINE_NOT_MAINTAINED",
			args: args{
				httpStatusCode: http.StatusTeapot,
				code:           50001,
			},
			expectedHTTPStatusCode: http.StatusTeapot,
			expectedCode:           50001,
			expectedErr:            ErrMGCreditLineNotMaintained,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewResponseError(tt.httpStatusCode, tt.code)
			require.Error(t, err)

			var responseError ResponseError
			require.True(t, errors.As(err, &responseError))

			assert.True(t, errors.Is(err, tt.expectedErr))

			assert.Equal(t, tt.expectedHTTPStatusCode, responseError.HTTPStatusCode)
			assert.Equal(t, tt.expectedCode, responseError.Code)
			assert.Equal(t, tt.expectedErr, responseError.Err)
		})
	}
}

func TestNewResponseError_Success(t *testing.T) {
	type args struct {
		httpStatusCode int
		code           int64
	}
	tests := []struct {
		name   string
		client http.Client
		args
	}{
		{
			name: "returns nil for success response codes",
			args: args{
				httpStatusCode: http.StatusOK,
				code:           0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewResponseError(tt.httpStatusCode, tt.code)
			require.NoError(t, err)

			assert.Empty(t, err)
		})
	}
}
