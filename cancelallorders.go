package cdcexchange

import (
	"context"
	"fmt"

	"github.com/cshep4/crypto-dot-com-exchange-go/errors"
	"github.com/cshep4/crypto-dot-com-exchange-go/internal/api"
	"github.com/cshep4/crypto-dot-com-exchange-go/internal/auth"
)

const methodCancelAllOrders = "private/cancel-all-orders"

// CancelAllOrdersResponse is the base response returned from the private/cancel-all-orders API.
type CancelAllOrdersResponse struct {
	// api.BaseResponse is the common response fields.
	api.BaseResponse
}

// CancelAllOrders cancels  all orders for a particular instrument/pair.
// This call is asynchronous, so the response is simply a confirmation of the request.
// The user.order subscription can be used to check when the order is successfully cancelled.
// Method: private/cancel-all-orders
func (c *client) CancelAllOrders(ctx context.Context, instrumentName string) error {
	if instrumentName == "" {
		return errors.InvalidParameterError{Parameter: "instrumentName", Reason: "cannot be empty"}
	}

	var (
		id        = c.idGenerator.Generate()
		timestamp = c.clock.Now().UnixMilli()
		params    = make(map[string]interface{})
	)

	params["instrument_name"] = instrumentName

	signature, err := c.signatureGenerator.GenerateSignature(auth.SignatureRequest{
		APIKey:    c.apiKey,
		SecretKey: c.secretKey,
		ID:        id,
		Method:    methodCancelAllOrders,
		Timestamp: timestamp,
		Params:    params,
	})
	if err != nil {
		return fmt.Errorf("failed to cancel signature: %w", err)
	}

	body := api.Request{
		ID:        id,
		Method:    methodCancelAllOrders,
		Nonce:     timestamp,
		Params:    params,
		Signature: signature,
		APIKey:    c.apiKey,
	}

	var cancelAllOrdersResponse CancelAllOrdersResponse
	statusCode, err := c.requester.Post(ctx, body, methodCancelAllOrders, &cancelAllOrdersResponse)
	if err != nil {
		return fmt.Errorf("failed to execute post request: %w", err)
	}

	if err := c.requester.CheckErrorResponse(statusCode, cancelAllOrdersResponse.Code); err != nil {
		return fmt.Errorf("error received in response: %w", err)
	}

	return nil
}
