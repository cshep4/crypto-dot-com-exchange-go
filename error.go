package cdcexchange

import (
	"errors"
	"fmt"
)

var (
	ErrUnexpectedError           = errors.New("unexpected error")
	ErrSystemError               = errors.New("system error")
	ErrUnauthorized              = errors.New("request not authenticated or key/signature is incorrect")
	ErrIllegalIP                 = errors.New("ip address not whitelisted")
	ErrBadRequest                = errors.New("missing required fields")
	ErrUserTierInvalid           = errors.New("disallowed based on user tier")
	ErrTooManyRequests           = errors.New("requests have exceeded rate limits")
	ErrInvalidNonce              = errors.New("nonce value differs by more than 30 seconds from server")
	ErrMethodNotFound            = errors.New("invalid method specified")
	ErrInvalidDateRange          = errors.New("invalid date range")
	ErrDuplicateRecord           = errors.New("duplicated record")
	ErrNegativeBalance           = errors.New("insufficient balance")
	ErrSymbolNotFound            = errors.New("invalid instrument_name specified")
	ErrSideNotSupported          = errors.New("invalid side specified")
	ErrOrderTypeNotSupported     = errors.New("invalid type specified")
	ErrMinPriceViolated          = errors.New("price is lower than the minimum")
	ErrMaxPriceViolated          = errors.New("price is higher than the maximum")
	ErrMinQuantityViolated       = errors.New("quantity is lower than the minimum")
	ErrMaxQuantityViolated       = errors.New("quantity is higher than the maximum")
	ErrMissingArgument           = errors.New("required argument is blank or missing")
	ErrInvalidPricePrecision     = errors.New("too many decimal places for price")
	ErrInvalidQuantityPrecision  = errors.New("too many decimal places for quantity")
	ErrMinNotionalViolated       = errors.New("the notional amount is less than the minimum")
	ErrMaxNotionalViolated       = errors.New("the notional amount exceeds the maximum")
	ErrMinAmountViolated         = errors.New("amount is less than the minimum")
	ErrMaxAmountViolated         = errors.New("amount exceeds the maximum")
	ErrAmountPrecisionOverflow   = errors.New("amount precision exceeds the maximum")
	ErrMGInvalidAccountStatus    = errors.New("operation has failed due to your account's status. please try again later")
	ErrMGTransferActiveLoan      = errors.New("transfer has failed due to holding an active loan. please repay your loan and try again later")
	ErrMGInvalidLoanCurrency     = errors.New("currency is not same as loan currency of active loan")
	ErrMGInvalidRepayAmount      = errors.New("only supporting full repayment of all margin loans")
	ErrMGNoActiveLoan            = errors.New("no active loan")
	ErrMGBlockedBorrow           = errors.New("borrow has been suspended. please try again later")
	ErrMGBlockedNewOrder         = errors.New("placing new order has been suspended. please try again later")
	ErrMGCreditLineNotMaintained = errors.New("please ensure your credit line is maintained and try again later")
)

type InvalidParameterError struct {
	Parameter string
	Reason    string
}

func (ipe InvalidParameterError) Error() string {
	return fmt.Sprintf("invalid parameter: %s %s", ipe.Parameter, ipe.Reason)
}

type ResponseError struct {
	Code int64
	Err  error
}

func (re ResponseError) Error() string {
	return fmt.Sprintf("%d: %v", re.Code, re.Err)
}

func (re ResponseError) Unwrap() error {
	return re.Err
}

func newResponseError(code int64) error {
	err := ResponseError{Code: code}

	switch code {
	case 0:
		return nil
	case 10001, 100001:
		err.Err = ErrSystemError
	case 10002:
		err.Err = ErrUnauthorized
	case 10003:
		err.Err = ErrIllegalIP
	case 10004:
		err.Err = ErrBadRequest
	case 10005:
		err.Err = ErrUserTierInvalid
	case 10006:
		err.Err = ErrTooManyRequests
	case 10007:
		err.Err = ErrInvalidNonce
	case 10008:
		err.Err = ErrMethodNotFound
	case 10009:
		err.Err = ErrInvalidDateRange
	case 20001:
		err.Err = ErrDuplicateRecord
	case 20002:
		err.Err = ErrNegativeBalance
	case 30003:
		err.Err = ErrSymbolNotFound
	case 30004:
		err.Err = ErrSideNotSupported
	case 30005:
		err.Err = ErrOrderTypeNotSupported
	case 30006:
		err.Err = ErrMinPriceViolated
	case 30007:
		err.Err = ErrMaxPriceViolated
	case 30008:
		err.Err = ErrMinQuantityViolated
	case 30009:
		err.Err = ErrMaxQuantityViolated
	case 30010:
		err.Err = ErrMissingArgument
	case 30013:
		err.Err = ErrInvalidPricePrecision
	case 30014:
		err.Err = ErrInvalidQuantityPrecision
	case 30016:
		err.Err = ErrMinNotionalViolated
	case 30017:
		err.Err = ErrMaxNotionalViolated
	case 30023:
		err.Err = ErrMinAmountViolated
	case 30024:
		err.Err = ErrMaxAmountViolated
	case 30025:
		err.Err = ErrAmountPrecisionOverflow
	case 40001:
		err.Err = ErrMGInvalidAccountStatus
	case 40002:
		err.Err = ErrMGTransferActiveLoan
	case 40003:
		err.Err = ErrMGInvalidLoanCurrency
	case 40004:
		err.Err = ErrMGInvalidRepayAmount
	case 40005:
		err.Err = ErrMGNoActiveLoan
	case 40006:
		err.Err = ErrMGBlockedBorrow
	case 40007:
		err.Err = ErrMGBlockedNewOrder
	case 50001:
		err.Err = ErrMGCreditLineNotMaintained
	default:
		err.Err = ErrUnexpectedError
	}

	return err
}
