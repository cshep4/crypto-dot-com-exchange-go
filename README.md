# [Crypto.com Exchange](https://crypto.com/exchange)

Go client for the [Crypto.com Spot Exchange REST API](https://exchange-docs.crypto.com/spot/index.html).

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
    GetInstruments(ctx context.Context) ([]Instrument, error)
    // GetTickers fetches the public tickers for an instrument (e.g. BTC_USDT).
    // instrument can be left blank to retrieve tickers for ALL instruments.
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
    GetAccountSummary(ctx context.Context, currency string) ([]Account, error)
    // CreateOrder creates a new BUY or SELL order on the Exchange.
    // This call is asynchronous, so the response is simply a confirmation of the request.
    // The user.order subscription can be used to check when the order is successfully created.
    CreateOrder(ctx context.Context, req CreateOrderRequest) (*CreateOrderResult, error)
}
```

| Method                           | Support |
:--------------------------------: | :-----: |
| private/get-account-summary      | ✅       |
| private/create-order             | ✅       |
| private/cancel-order             | ⚠️       |
| private/cancel-all-orders        | ⚠️       |
| private/get-order-history        | ⚠️       |
| private/get-open-orders          | ⚠️       |
| private/get-order-detail         | ⚠️       |
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

