package cdcexchange

import (
	"context"
	"fmt"
	"time"

	"github.com/cshep4/crypto-dot-com-exchange-go/errors"
	"github.com/cshep4/crypto-dot-com-exchange-go/internal/api"
	"github.com/cshep4/crypto-dot-com-exchange-go/internal/auth"
)

const (
	methodGetOrderHistory = "private/get-order-history"
)

type (
	// GetOrderHistoryRequest is the request params sent for the private/get-order-history API.
	//
	// The maximum duration between Start and End is 24 hours.
	//
	// You will receive an INVALID_DATE_RANGE error if the difference exceeds the maximum duration.
	//
	// For users looking to pull longer historical order data, users can create a loop to make a request
	// for each 24-period from the desired start to end time.
	GetOrderHistoryRequest struct {
		// InstrumentName represents the currency pair for the orders (e.g. ETH_CRO or BTC_USDT).
		// if InstrumentName is omitted, all instruments will be returned.
		InstrumentName string `json:"instrument_name"`
		// Start is the start timestamp (milliseconds since the Unix epoch)
		// (Default: 24 hours ago)
		Start time.Time `json:"start_ts"`
		// End is the end timestamp (milliseconds since the Unix epoch)
		// (Default: now)
		End time.Time `json:"end_ts"`
		// PageSize represents maximum number of orders returned (for pagination)
		// (Default: 20, Max: 200)
		// if PageSize is 0, it will be set as 20 by default.
		PageSize int `json:"page_size"`
		// Page represents the page number (for pagination)
		// (0-based)
		Page int `json:"page"`
	}

	// GetOrderHistoryResponse is the base response returned from the private/get-order-history API.
	GetOrderHistoryResponse struct {
		// api.BaseResponse is the common response fields.
		api.BaseResponse
		// Result is the response attributes of the endpoint.
		Result GetOrderHistoryResult `json:"result"`
	}

	// GetOrderHistoryResult is the result returned from the private/get-order-history API.
	GetOrderHistoryResult struct {
		// OrderList is the array of orders.
		OrderList []Order `json:"order_list"`
	}
)

// GetOrderHistory gets the order history for a particular instrument.
//
// Pagination is handled using page size (Default: 20, Max: 200) & number (0-based).
// If paging is used, enumerate each page (starting with 0) until an empty order_list array appears in the response.
//
// req.InstrumentName can be left blank to get open orders for all instruments.
//
// Method: private/get-order-history
func (c *client) GetOrderHistory(ctx context.Context, req GetOrderHistoryRequest) (*GetOrderHistoryResult, error) {
	if req.PageSize < 0 {
		return nil, errors.InvalidParameterError{Parameter: "req.PageSize", Reason: "cannot be less than 0"}
	}
	if req.PageSize > 200 {
		return nil, errors.InvalidParameterError{Parameter: "req.PageSize", Reason: "cannot be greater than 200"}
	}

	var (
		id        = c.idGenerator.Generate()
		timestamp = c.clock.Now().UnixMilli()
		params    = make(map[string]interface{})
	)

	if req.InstrumentName != "" {
		params["instrument_name"] = req.InstrumentName
	}
	if req.PageSize != 0 {
		params["page_size"] = req.PageSize
	}
	if !req.Start.IsZero() {
		params["start_ts"] = req.Start.UnixMilli()
	}
	if !req.End.IsZero() {
		params["end_ts"] = req.End.UnixMilli()
	}
	params["page"] = req.Page

	signature, err := c.signatureGenerator.GenerateSignature(auth.SignatureRequest{
		APIKey:    c.apiKey,
		SecretKey: c.secretKey,
		ID:        id,
		Method:    methodGetOrderHistory,
		Timestamp: timestamp,
		Params:    params,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create signature: %w", err)
	}

	body := api.Request{
		ID:        id,
		Method:    methodGetOrderHistory,
		Nonce:     timestamp,
		Params:    params,
		Signature: signature,
		APIKey:    c.apiKey,
	}

	var getOrderHistoryResponse GetOrderHistoryResponse
	statusCode, err := c.requester.Post(ctx, body, methodGetOrderHistory, &getOrderHistoryResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to execute post request: %w", err)
	}

	if err := c.requester.CheckErrorResponse(statusCode, getOrderHistoryResponse.Code); err != nil {
		return nil, fmt.Errorf("error received in response: %w", err)
	}

	return &getOrderHistoryResponse.Result, nil
}
