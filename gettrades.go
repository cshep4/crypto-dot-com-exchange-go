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
	methodGetTrades = "private/get-trades"
)

type (
	// GetTradesRequest is the request params sent for the private/get-trades API.
	//
	// The maximum duration between Start and End is 24 hours.
	//
	// You will receive an INVALID_DATE_RANGE error if the difference exceeds the maximum duration.
	//
	// For users looking to pull longer historical trade data, users can create a loop to make a request
	// for each 24-period from the desired start to end time.
	GetTradesRequest struct {
		// InstrumentName represents the currency pair for the trades (e.g. ETH_CRO or BTC_USDT).
		// if InstrumentName is omitted, all instruments will be returned.
		InstrumentName string `json:"instrument_name"`
		// Start is the start timestamp (milliseconds since the Unix epoch)
		// (Default: 24 hours ago)
		Start time.Time `json:"start_ts"`
		// End is the end timestamp (milliseconds since the Unix epoch)
		// (Default: now)
		End time.Time `json:"end_ts"`
		// PageSize represents maximum number of trades returned (for pagination)
		// (Default: 20, Max: 200)
		// if PageSize is 0, it will be set as 20 by default.
		PageSize int `json:"page_size"`
		// Page represents the page number (for pagination)
		// (0-based)
		Page int `json:"page"`
	}

	// GetTradesResponse is the base response returned from the private/get-trades API.
	GetTradesResponse struct {
		// api.BaseResponse is the common response fields.
		api.BaseResponse
		// Result is the response attributes of the endpoint.
		Result GetTradesResult `json:"result"`
	}

	// GetTradesResult is the result returned from the private/get-trades API.
	GetTradesResult struct {
		// TradeList is the array of trades.
		TradeList []Trade `json:"trade_list"`
	}
)

// GetTrades gets all executed trades for a particular instrument.
//
// Pagination is handled using page size (Default: 20, Max: 200) & number (0-based).
// If paging is used, enumerate each page (starting with 0) until an empty trade_list array appears in the response.
//
// req.InstrumentName can be left blank to get executed trades for all instruments.
//
// Method: private/get-trades
func (c *client) GetTrades(ctx context.Context, req GetTradesRequest) ([]Trade, error) {
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
		Method:    methodGetTrades,
		Timestamp: timestamp,
		Params:    params,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create signature: %w", err)
	}

	body := api.Request{
		ID:        id,
		Method:    methodGetTrades,
		Nonce:     timestamp,
		Params:    params,
		Signature: signature,
		APIKey:    c.apiKey,
	}

	var getTradesResponse GetTradesResponse
	statusCode, err := c.requester.Post(ctx, body, methodGetTrades, &getTradesResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to execute post request: %w", err)
	}

	if err := c.requester.CheckErrorResponse(statusCode, getTradesResponse.Code); err != nil {
		return nil, fmt.Errorf("error received in response: %w", err)
	}

	return getTradesResponse.Result.TradeList, nil
}
