# [Crypto.com Exchange](https://crypto.com/exchange)

Go client for the [Crypto.com Spot Exchange REST API](https://exchange-docs.crypto.com/spot/index.html).

- [Installation](#installation)
- [Setup](#setup)
- [Optional Configurations](#optional-configurations)
  - [UAT Sandbox Environment](#uat-sandbox-environment)
  - [Production Environment](#production-environment)
  - [Custom HTTP Client](#custom-http-client)
- [Supported API](#supported-api-official-docs)
    - [Common API](#common-api)
    - [Spot Trading API](#spot-trading-api)
    - [Margin Trading API](#margin-trading-api)
    - [Derivatives Transfer API](#derivatives-transfer-api)
    - [Sub-account API](#sub-account-api)
    - [Websocket](#websocket)
        - [Websocket Heartbeats](#websocket-heartbeats)
        - [Websocket Subscriptions](#websocket-subscriptions)
- [Errors](#errors)
  - [Response Codes](#response-codes)


## Installation

To import the package, run:

    go get github.com/cshep4/crypto-dot-com-exchange-go

## Setup

An instance of the client can be created like this:

```go
import (
    cdcexchange "github.com/cshep4/crypto-dot-com-exchange-go"
)

client, err := cdcexchange.New("<api_key>", "<secret_key>")
if err != nil {
    return err
}

// optional, configuration can be updated with UpdateConfig
err = client.UpdateConfig("<api_key>", "<secret_key>")
if err != nil {
    return err
}
```

## Optional Configurations

### UAT Sandbox Environment

The client can be configured to make requests against the UAT Sandbox environment using the `WithUATEnvironment` functional option like so:

```go
import (
    cdcexchange "github.com/cshep4/crypto-dot-com-exchange-go"
)

client, err := cdcexchange.New("<api_key>", "<secret_key>", 
    cdcexchange.WithUATEnvironment(),
)
if err != nil {
    return err
}
```

### Production Environment

The client can be configured to make requests against the Production environment using the `WithProductionEnvironment` functional option. Clients will make requests against this environment by default even if this is not specified.

```go
import (
    cdcexchange "github.com/cshep4/crypto-dot-com-exchange-go"
)

client, err := cdcexchange.New("<api_key>", "<secret_key>", 
    cdcexchange.WithProductionEnvironment(),
)
if err != nil {
    return err
}
```

### Custom HTTP Client

The client can be configured to use a custom HTTP client using the `WithHTTPClient` functional option. This can be used to create custom timeouts, enable tracing, etc. This is initialised like so:

```go
import (
    "net/http"

    cdcexchange "github.com/cshep4/crypto-dot-com-exchange-go"
)

client, err := cdcexchange.New("<api_key>", "<secret_key>",
    cdcexchange.WithHTTPClient(&http.Client{
        Timeout: 15 * time.Second,
    }),
)
if err != nil {
    return err
}
```


## Supported API ([Official Docs](https://exchange-docs.crypto.com/spot/index.html)):

The supported APIs for each module are listed below.

**✅ - API is supported**

**⚠️ - API is not yet supported (hopefully should be available soon!)**

Each API module is separated into separate interfaces, however the `CryptoDotComExchange` interface can be used to access all methods:
```go
// CryptoDotComExchange is a Crypto.com Exchange client for all available APIs.
type CryptoDotComExchange interface {
    // UpdateConfig can be used to update the configuration of the client object.
    // (e.g. change api key, secret key, environment, etc).
    UpdateConfig(apiKey string, secretKey string, opts ...ClientOption) error
    CommonAPI
    SpotTradingAPI
    MarginTradingAPI
    DerivativesTransferAPI
    SubAccountAPI
    Websocket
}
```

Client interfaces can be found in [client.go](client.go).

### Common API

```go
// CommonAPI is a Crypto.com Exchange client for Common API.
type CommonAPI interface {
    // GetInstruments provides information on all supported instruments (e.g. BTC_USDT).
    // Method: public/get-instruments
    GetInstruments(ctx context.Context) ([]Instrument, error)
    // GetTickers fetches the public tickers for an instrument (e.g. BTC_USDT).
    // instrument can be left blank to retrieve tickers for ALL instruments.
    // Method: public/get-ticker
    GetTickers(ctx context.Context, instrument string) ([]Ticker, error)
}
```

| Method                           | Support |
:--------------------------------: | :-----: |
| public/auth                      | ⚠️ |
| public/get-instruments           | ✅ |
| public/get-book                  | ⚠️ |
| public/get-candlestick           | ⚠️ |
| public/get-ticker                | ✅ |
| public/get-trades                | ⚠️ |
| private/set-cancel-on-disconnect | ⚠️ |
| private/get-cancel-on-disconnect | ⚠️ |
| private/create-withdrawal        | ⚠️ |
| private/get-withdrawal-history   | ⚠️ |
| private/get-deposit-history      | ⚠️ |
| private/get-deposit-address      | ⚠️ |

### Spot Trading API

```go
// SpotTradingAPI is a Crypto.com Exchange client for Spot Trading API.
type SpotTradingAPI interface {
    // GetAccountSummary returns the account balance of a user for a particular token.
    // currency can be left blank to retrieve balances for ALL tokens.
    // Method: private/get-account-summary
    GetAccountSummary(ctx context.Context, currency string) ([]Account, error)
    // CreateOrder creates a new BUY or SELL order on the Exchange.
    // This call is asynchronous, so the response is simply a confirmation of the request.
    // The user.order subscription can be used to check when the order is successfully created.
    // Method: private/create-order
    CreateOrder(ctx context.Context, req CreateOrderRequest) (*CreateOrderResult, error)
    // GetOpenOrders gets all open orders for a particular instrument
    // Pagination is handled using page size (Default: 20, Max: 200) & number (0-based).
    // req.InstrumentName can be left blank to get open orders for all instruments.
    // Method: private/get-open-orders
    GetOpenOrders(ctx context.Context, req GetOpenOrdersRequest) (*GetOpenOrdersResult, error)
    // GetOrderDetail gets details of an order for a particular order ID
    // Method: get-order-detail
    GetOrderDetail(ctx context.Context, orderID string) (*GetOrderDetailResult, error)
}
```

| Method                           | Support |
:--------------------------------: | :-----: |
| private/get-account-summary      | ✅       |
| private/create-order             | ✅       |
| private/cancel-order             | ⚠️       |
| private/cancel-all-orders        | ⚠️       |
| private/get-order-history        | ⚠️       |
| private/get-open-orders          | ✅       |
| private/get-order-detail         | ✅       |
| private/get-trades               | ⚠️       |

### Margin Trading API

```go
// MarginTradingAPI is a Crypto.com Exchange client for Margin Trading API.
type MarginTradingAPI interface {
}
```

| Method                                 | Support |
:--------------------------------------: | :-----: |
| public/margin/get-transfer-currencies  | ⚠️       |
| public/margin/get-loan-currencies      | ⚠️       |
| private/margin/get-user-config         | ⚠️       |
| private/margin/get-account-summary     | ⚠️       |
| private/margin/transfer                | ⚠️       |
| private/margin/borrow                  | ⚠️       |
| private/margin/repay                   | ⚠️       |
| private/margin/get-transfer-history    | ⚠️       |
| private/margin/get-borrow-history      | ⚠️       |
| private/margin/get-interest-history    | ⚠️       |
| private/margin/get-repay-history       | ⚠️       |
| private/margin/get-liquidation-history | ⚠️       |
| private/margin/get-liquidation-orders  | ⚠️       |
| private/margin/create-order            | ⚠️       |
| private/margin/cancel-order            | ⚠️       |
| private/margin/cancel-all-orders       | ⚠️       |
| private/margin/get-order-history       | ⚠️       |
| private/margin/get-open-orders         | ⚠️       |
| private/margin/get-order-detail        | ⚠️       |
| private/margin/get-trades              | ⚠️       |

### Derivatives Transfer API

```go
// DerivativesTransferAPI is a Crypto.com Exchange client for Derivatives Transfer API.
type DerivativesTransferAPI interface {
}
```

| Method                             | Support |
:----------------------------------: | :-----: |
| private/deriv/transfer             | ⚠️       |
| private/deriv/get-transfer-history | ⚠️       |

### Sub-account API

```go
// SubAccountAPI is a Crypto.com Exchange client for Sub-account API.
type SubAccountAPI interface {
}
```

| Method                                  | Support |
:---------------------------------------: | :-----: |
| private/subaccount/get-sub-accounts     | ⚠️       |
| private/subaccount/get-transfer-history | ⚠️       |
| private/subaccount/transfer             | ⚠️       |

### Websocket

```go
// Websocket is a Crypto.com Exchange client websocket methods & channels.
type Websocket interface {
}
```

#### Websocket Heartbeats

| Method                   | Support |
:------------------------: | :-----: |
| public/respond-heartbeat | ⚠️       |

#### Websocket Subscriptions

| Channel                                  | Support |
:----------------------------------------: | :-----: |
| user.order.{instrument_name}             | ⚠️       |
| user.trade.{instrument_name}             | ⚠️       |
| user.balance                             | ⚠️       |
| user.margin.order.{instrument_name}      | ⚠️       |
| user.margin.trade.{instrument_name}      | ⚠️       |
| user.margin.balance                      | ⚠️       |
| book.{instrument_name}.{depth}           | ⚠️       |
| ticker.{instrument_name}                 | ⚠️       |
| trade.{instrument_name}                  | ⚠️       |
| candlestick.{interval}.{instrument_name} | ⚠️       |


## Errors

Custom errors are returned based on the HTTP status code and reason codes returned in the API response.
Official documentation on reason codes can be found [here](https://exchange-docs.crypto.com/spot/index.html#response-and-reason-codes).
All errors returned from the client can be found in the [errors](/errors) package.

Custom error handling on client errors can be implemented like so:

```go
import (
    "errors"

    cdcerrors "github.com/cshep4/crypto-dot-com-exchange-go/errors"
    cdcexchange "github.com/cshep4/crypto-dot-com-exchange-go"
)

...

res, err := client.GetAccountSummary(ctx, "CRO")
if err != nil {
    switch {
    case errors.Is(err, cdcerrors.ErrSystemError):
        // handle system error
    case errors.Is(err, cdcerrors.ErrUnauthorized):
        // handle unauthorized error
    case errors.Is(err, cdcerrors.ErrIllegalIP):
        // handle illegal IP error
    ...

    }

    return err
}
```

Resoonse errors can also be cast to retrieve the response code and HTTP status code:

```go
import (
    "errors"
    "log"
    
    cdcerrors "github.com/cshep4/crypto-dot-com-exchange-go/errors"
    cdcexchange "github.com/cshep4/crypto-dot-com-exchange-go"
)

...

res, err := client.GetAccountSummary(ctx, "CRO")
if err != nil {
    var responseError cdcerrors.ResponseError
    if errors.Is(err, &responseError) {
        // response code
        log.Println(responseError.Code)
		
        // HTTP status code
        log.Println(responseError.HTTPStatusCode)

        // underlying error
        log.Println(responseError.Err)
    }
	
    return err
}
```

### Response Codes

|Code   | HTTP Status | Client Error                 | Message Code                  | Description                                                                                    |
:-----: | :---------: | :--------------------------: | :---------------------------: | :--------------------------------------------------------------------------------------------: |
| 0     | 200         | nil                          | --                            | Success                                                                                        |
| 10001 | 500         | ErrSystemError               | SYS_ERROR                     | Malformed request, (E.g. not using application/json for REST)                                  |
| 10002 | 401         | ErrUnauthorized              | UNAUTHORIZED                  | Not authenticated, or key/signature incorrect                                                  |
| 10003 | 401         | ErrIllegalIP                 | IP_ILLEGAL                    | IP address not whitelisted                                                                     |
| 10004 | 400         | ErrBadRequest                | BAD_REQUEST                   | Missing required fields                                                                        |
| 10005 | 401         | ErrUserTierInvalid           | USER_TIER_INVALID             | Disallowed based on user tier                                                                  |
| 10006 | 429         | ErrTooManyRequests           | TOO_MANY_REQUESTS             | Requests have exceeded rate limits                                                             |
| 10007 | 400         | ErrInvalidNonce              | INVALID_NONCE                 | Nonce value differs by more than 30 seconds from server                                        |
| 10008 | 400         | ErrMethodNotFound            | METHOD_NOT_FOUND              | Invalid method specified                                                                       |
| 10009 | 400         | ErrInvalidDateRange          | INVALID_DATE_RANGE            | Invalid date range                                                                             |
| 20001 | 400         | ErrDuplicateRecord           | DUPLICATE_RECORD              | Duplicated record                                                                              |
| 20002 | 400         | ErrNegativeBalance           | NEGATIVE_BALANCE              | Insufficient balance                                                                           |
| 30003 | 400         | ErrSymbolNotFound            | SYMBOL_NOT_FOUND              | Invalid instrument_name specified                                                              |
| 30004 | 400         | ErrSideNotSupported          | SIDE_NOT_SUPPORTED            | Invalid side specified                                                                         |
| 30005 | 400         | ErrOrderTypeNotSupported     | ORDERTYPE_NOT_SUPPORTED       | Invalid type specified                                                                         |
| 30006 | 400         | ErrMinPriceViolated          | MIN_PRICE_VIOLATED            | Price is lower than the minimum                                                                |
| 30007 | 400         | ErrMaxPriceViolated          | MAX_PRICE_VIOLATED            | Price is higher than the maximum                                                               |
| 30008 | 400         | ErrMinQuantityViolated       | MIN_QUANTITY_VIOLATED         | Quantity is lower than the minimum                                                             |
| 30009 | 400         | ErrMaxQuantityViolated       | MAX_QUANTITY_VIOLATED         | Quantity is higher than the maximum                                                            |
| 30010 | 400         | ErrMissingArgument           | MISSING_ARGUMENT              | Required argument is blank or missing                                                          |
| 30013 | 400         | ErrInvalidPricePrecision     | INVALID_PRICE_PRECISION       | Too many decimal places for Price                                                              |
| 30014 | 400         | ErrInvalidQuantityPrecision  | INVALID_QUANTITY_PRECISION    | Too many decimal places for Quantity                                                           |
| 30016 | 400         | ErrMinNotionalViolated       | MIN_NOTIONAL_VIOLATED         | The notional amount is less than the minimum                                                   |
| 30017 | 400         | ErrMaxNotionalViolated       | MAX_NOTIONAL_VIOLATED         | The notional amount exceeds the maximum                                                        |
| 30023 | 400         | ErrMinAmountViolated         | MIN_AMOUNT_VIOLATED           | Amount is lower than the minimum                                                               |
| 30024 | 400         | ErrMaxAmountViolated         | MAX_AMOUNT_VIOLATED           | Amount is higher than the maximum                                                              |
| 30025 | 400         | ErrAmountPrecisionOverflow   | AMOUNT_PRECISION_OVERFLOW     | Amount precision exceeds the maximum                                                           |
| 40001 | 400         | ErrMGInvalidAccountStatus    | MG_INVALID_ACCOUNT_STATUS     | Operation has failed due to your account's status. Please try again later.                     |
| 40002 | 400         | ErrMGTransferActiveLoan      | MG_TRANSFER_ACTIVE_LOAN       | Transfer has failed due to holding an active loan. Please repay your loan and try again later. |
| 40003 | 400         | ErrMGInvalidLoanCurrency     | MG_INVALID_LOAN_CURRENCY      | Currency is not same as loan currency of active loan                                           |
| 40004 | 400         | ErrMGInvalidRepayAmount      | MG_INVALID_REPAY_AMOUNT       | Only supporting full repayment of all margin loans                                             |
| 40005 | 400         | ErrMGNoActiveLoan            | MG_NO_ACTIVE_LOAN             | No active loan                                                                                 |
| 40006 | 400         | ErrMGBlockedBorrow           | MG_BLOCKED_BORROW             | Borrow has been suspended. Please try again later.                                             |
| 40007 | 400         | ErrMGBlockedNewOrder         | MG_BLOCKED_NEW_ORDER          | Placing new order has been suspended. Please try again later.                                  |
| 50001 | 400         | ErrMGCreditLineNotMaintained | DW_CREDIT_LINE_NOT_MAINTAINED | Please ensure your credit line is maintained and try again later.                              |