package cdcexchange

import (
	"net/http"

	"github.com/jonboulle/clockwork"

	"github.com/cshep4/crypto-dot-com-exchange-go/internal/id"
)

const (
	UATSandboxBaseURL = uatSandboxBaseURL
	ProductionBaseURL = productionBaseURL
)

type Client = client

func (c client) BaseURL() string {
	return c.baseURL
}

func (c client) APIKey() string {
	return c.apiKey
}

func (c client) SecretKey() string {
	return c.secretKey
}

func (c client) HTTPClient() *http.Client {
	return c.client
}

func WithIDGenerator(idGenerator id.IDGenerator) ClientOption {
	return func(c *client) error {
		if idGenerator == nil {
			return InvalidParameterError{Parameter: "idGenerator", Reason: "cannot be empty"}
		}

		c.idGenerator = idGenerator
		return nil
	}
}

func WithClock(clock clockwork.Clock) ClientOption {
	return func(c *client) error {
		if clock == nil {
			return InvalidParameterError{Parameter: "clock", Reason: "cannot be empty"}
		}

		c.clock = clock
		return nil
	}
}

func WithBaseURL(url string) ClientOption {
	return func(c *client) error {
		if url == "" {
			return InvalidParameterError{Parameter: "url", Reason: "cannot be empty"}
		}

		c.baseURL = url
		return nil
	}
}
