package cdcexchange

import (
	"net/http"

	"github.com/jonboulle/clockwork"

	"github.com/cshep4/crypto-dot-com-exchange-go/errors"
	"github.com/cshep4/crypto-dot-com-exchange-go/internal/auth"
	"github.com/cshep4/crypto-dot-com-exchange-go/internal/id"
)

const (
	UATSandboxBaseURL = uatSandboxBaseURL
	ProductionBaseURL = productionBaseURL
)

func (c client) BaseURL() string {
	return c.requester.BaseURL
}

func (c client) APIKey() string {
	return c.apiKey
}

func (c client) SecretKey() string {
	return c.secretKey
}

func (c client) HTTPClient() *http.Client {
	return c.requester.Client
}

func WithIDGenerator(idGenerator id.IDGenerator) ClientOption {
	return func(c *client) error {
		if idGenerator == nil {
			return errors.InvalidParameterError{Parameter: "idGenerator", Reason: "cannot be empty"}
		}

		c.idGenerator = idGenerator
		return nil
	}
}

func WithSignatureGenerator(signatureGenerator auth.SignatureGenerator) ClientOption {
	return func(c *client) error {
		if signatureGenerator == nil {
			return errors.InvalidParameterError{Parameter: "signatureGenerator", Reason: "cannot be empty"}
		}

		c.signatureGenerator = signatureGenerator
		return nil
	}
}

func WithClock(clock clockwork.Clock) ClientOption {
	return func(c *client) error {
		if clock == nil {
			return errors.InvalidParameterError{Parameter: "clock", Reason: "cannot be empty"}
		}

		c.clock = clock
		return nil
	}
}

func WithBaseURL(url string) ClientOption {
	return func(c *client) error {
		if url == "" {
			return errors.InvalidParameterError{Parameter: "url", Reason: "cannot be empty"}
		}

		c.requester.BaseURL = url
		return nil
	}
}
