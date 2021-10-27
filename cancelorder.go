package cdcexchange

import (
	"context"
	"fmt"

	"github.com/cshep4/crypto-dot-com-exchange-go/errors"
	"github.com/cshep4/crypto-dot-com-exchange-go/internal/api"
	"github.com/cshep4/crypto-dot-com-exchange-go/internal/auth"
)

const methodCancelOrder = "private/cancel-order"

// CancelOrderResponse is the base response returned from the private/cancel-order API.
type CancelOrderResponse struct {
	// api.BaseResponse is the common response fields.
	api.BaseResponse
}

// CancelOrder cancels an existing order on the Exchange.
//
// This call is asynchronous, so the response is simply a confirmation of the request.
//
// The user.order subscription can be used to check when the order is successfully cancelled.
//
// Method: private/cancel-order
func (c *client) CancelOrder(ctx context.Context, instrumentName string, orderID string) error {
	if instrumentName == "" {
		return errors.InvalidParameterError{Parameter: "instrumentName", Reason: "cannot be empty"}
	}
	if orderID == "" {
		return errors.InvalidParameterError{Parameter: "orderID", Reason: "cannot be empty"}
	}

	var (
		id        = c.idGenerator.Generate()
		timestamp = c.clock.Now().UnixMilli()
		params    = make(map[string]interface{})
	)

	params["instrument_name"] = instrumentName
	params["order_id"] = orderID

	signature, err := c.signatureGenerator.GenerateSignature(auth.SignatureRequest{
		APIKey:    c.apiKey,
		SecretKey: c.secretKey,
		ID:        id,
		Method:    methodCancelOrder,
		Timestamp: timestamp,
		Params:    params,
	})
	if err != nil {
		return fmt.Errorf("failed to cancel signature: %w", err)
	}

	body := api.Request{
		ID:        id,
		Method:    methodCancelOrder,
		Nonce:     timestamp,
		Params:    params,
		Signature: signature,
		APIKey:    c.apiKey,
	}

	var cancelOrderResponse CancelOrderResponse
	statusCode, err := c.requester.Post(ctx, body, methodCancelOrder, &cancelOrderResponse)
	if err != nil {
		return fmt.Errorf("failed to execute post request: %w", err)
	}

	if err := c.requester.CheckErrorResponse(statusCode, cancelOrderResponse.Code); err != nil {
		return fmt.Errorf("error received in response: %w", err)
	}

	return nil
}
