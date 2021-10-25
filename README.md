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

The table of listed APIs that are supported by this package.

**✅ - API is supported**

**⚠️ - API is not yet supported (hopefully should be available soon!)**

### Common API

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

| Method                             | Support |
:----------------------------------: | :-----: |
| private/deriv/transfer             | ⚠️       |
| private/deriv/get-transfer-history | ⚠️       |

### Sub-account API

| Method                                  | Support |
:---------------------------------------: | :-----: |
| private/subaccount/get-sub-accounts     | ⚠️       |
| private/subaccount/get-transfer-history | ⚠️       |
| private/subaccount/transfer             | ⚠️       |

### Websocket Heartbeats

| Method                   | Support |
:------------------------: | :-----: |
| public/respond-heartbeat | ⚠️       |

### Websocket Subscriptions

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

