package cdcexchange

import (
	"context"
	"net/http"

	"github.com/jonboulle/clockwork"

	"github.com/cshep4/crypto-dot-com-exchange-go/errors"
	"github.com/cshep4/crypto-dot-com-exchange-go/internal/api"
	"github.com/cshep4/crypto-dot-com-exchange-go/internal/auth"
	"github.com/cshep4/crypto-dot-com-exchange-go/internal/id"
)

const (
	EnvironmentUATSandbox Environment = "uat_sandbox"
	EnvironmentProduction Environment = "production"

	uatSandboxBaseURL = "https://uat-api.3ona.co/v2/"
	productionBaseURL = "https://api.crypto.com/v2/"
)

type (
	// CryptoDotComExchange is a Crypto.com Exchange client for all available APIs.
	CryptoDotComExchange interface {
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

	// CommonAPI is a Crypto.com Exchange client for Common API.
	CommonAPI interface {
		// GetInstruments provides information on all supported instruments (e.g. BTC_USDT).
		//
		// Method: public/get-instruments
		GetInstruments(ctx context.Context) ([]Instrument, error)
		// GetTickers fetches the public tickers for an instrument (e.g. BTC_USDT).
		//
		// instrument can be left blank to retrieve tickers for ALL instruments.
		//
		// Method: public/get-ticker
		GetTickers(ctx context.Context, instrument string) ([]Ticker, error)
	}

	// SpotTradingAPI is a Crypto.com Exchange client for Spot Trading API.
	SpotTradingAPI interface {
		// GetAccountSummary returns the account balance of a user for a particular token.
		//
		// currency can be left blank to retrieve balances for ALL tokens.
		//
		// Method: private/get-account-summary
		GetAccountSummary(ctx context.Context, currency string) ([]Account, error)
		// CreateOrder creates a new BUY or SELL order on the Exchange.
		//
		// This call is asynchronous, so the response is simply a confirmation of the request.
		//
		// The user.order subscription can be used to check when the order is successfully created.
		//
		// Method: private/create-order
		CreateOrder(ctx context.Context, req CreateOrderRequest) (*CreateOrderResult, error)
		// CancelOrder cancels an existing order on the Exchange.
		//
		// This call is asynchronous, so the response is simply a confirmation of the request.
		//
		// The user.order subscription can be used to check when the order is successfully cancelled.
		//
		// Method: private/cancel-order
		CancelOrder(ctx context.Context, instrumentName string, orderID string) error
		// CancelAllOrders cancels  all orders for a particular instrument/pair.
		//
		// This call is asynchronous, so the response is simply a confirmation of the request.
		//
		// The user.order subscription can be used to check when the order is successfully cancelled.
		//
		// Method: private/cancel-all-orders
		CancelAllOrders(ctx context.Context, instrumentName string) error
		// GetOrderHistory gets the order history for a particular instrument.
		//
		// Pagination is handled using page size (Default: 20, Max: 200) & number (0-based).
		// If paging is used, enumerate each page (starting with 0) until an empty order_list array appears in the response.
		//
		// req.InstrumentName can be left blank to get open orders for all instruments.
		//
		// Method: private/get-order-history
		GetOrderHistory(ctx context.Context, req GetOrderHistoryRequest) (*GetOrderHistoryResult, error)
		// GetOpenOrders gets all open orders for a particular instrument.
		//
		// Pagination is handled using page size (Default: 20, Max: 200) & number (0-based).
		//
		// req.InstrumentName can be left blank to get open orders for all instruments.
		//
		// Method: private/get-open-orders
		GetOpenOrders(ctx context.Context, req GetOpenOrdersRequest) (*GetOpenOrdersResult, error)
		// GetOrderDetail gets details of an order for a particular order ID.
		//
		// Method: private/get-order-detail
		GetOrderDetail(ctx context.Context, orderID string) (*GetOrderDetailResult, error)
	}

	// MarginTradingAPI is a Crypto.com Exchange client for Margin Trading API.
	MarginTradingAPI interface {
	}

	// DerivativesTransferAPI is a Crypto.com Exchange client for Derivatives Transfer API.
	DerivativesTransferAPI interface {
	}

	// SubAccountAPI is a Crypto.com Exchange client for Sub-account API.
	SubAccountAPI interface {
	}

	// Websocket is a Crypto.com Exchange client websocket methods & channels.
	Websocket interface {
	}

	// Environment represents the environment against which calls are made.
	Environment string

	// ClientOption represents optional configurations for the client.
	ClientOption func(*client) error

	// client is a concrete implementation of CryptoDotComExchange.
	client struct {
		apiKey             string
		secretKey          string
		clock              clockwork.Clock
		idGenerator        id.IDGenerator
		signatureGenerator auth.SignatureGenerator
		requester          api.Requester
	}
)

// New will construct a new instance of client.
func New(apiKey string, secretKey string, opts ...ClientOption) (*client, error) {
	c := &client{
		idGenerator:        &id.Generator{},
		signatureGenerator: &auth.Generator{},
		clock:              clockwork.NewRealClock(),
		requester: api.Requester{
			Client:  http.DefaultClient,
			BaseURL: productionBaseURL,
		},
	}

	if err := c.UpdateConfig(apiKey, secretKey, opts...); err != nil {
		return nil, err
	}

	return c, nil
}

// UpdateConfig can be used to update the configuration of the client object.
// (e.g. change api key, secret key, environment, etc).
func (c *client) UpdateConfig(apiKey string, secretKey string, opts ...ClientOption) error {
	switch {
	case apiKey == "":
		return errors.InvalidParameterError{Parameter: "apiKey", Reason: "cannot be empty"}
	case secretKey == "":
		return errors.InvalidParameterError{Parameter: "secretKey", Reason: "cannot be empty"}
	}

	c.apiKey = apiKey
	c.secretKey = secretKey

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return err
		}
	}

	return nil
}

// WithProductionEnvironment will initialise the client to make requests against the production environment.
// This is the default setting.
func WithProductionEnvironment() ClientOption {
	return func(c *client) error {
		c.requester.BaseURL = productionBaseURL
		return nil
	}
}

// WithUATEnvironment will initialise the client to make requests against the UAT sandbox environment.
func WithUATEnvironment() ClientOption {
	return func(c *client) error {
		c.requester.BaseURL = uatSandboxBaseURL
		return nil
	}
}

// WithHTTPClient will allow the client to be initialised with a custom http client.
// Can be used to create custom timeouts, enable tracing, etc.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *client) error {
		if httpClient == nil {
			return errors.InvalidParameterError{Parameter: "httpClient", Reason: "cannot be empty"}
		}

		c.requester.Client = httpClient
		return nil
	}
}
