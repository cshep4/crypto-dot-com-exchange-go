package cdcexchange

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cshep4/crypto-dot-com-exchange-go/internal/api"
	"github.com/cshep4/crypto-dot-com-exchange-go/internal/time"
)

const (
	methodGetBook = "public/get-book"
)

type (
	// BookResponse is the base response returned from the public/get-book API
	// when no instrument is specified.
	BookResponse struct {
		// api.BaseResponse is the common response fields.
		api.BaseResponse
		// Result is the response attributes of the endpoint.
		Result BookResult `json:"result"`
	}

	// BookResult is the result returned from the public/get-book API.
	BookResult struct {
		// Bids is an array of bids.
		// [0] = Price, [1] = Quantity, [2] = Number of Orders.
		Bids [][]float64 `json:"bids"`
		// Asks is an array of asks.
		// [0] = Price, [1] = Quantity, [2] = Number of Orders.
		Asks [][]float64 `json:"asks"`
		// Timestamp is the timestamp of the data.
		Timestamp time.Time `json:"t"`
	}
)

// GetBook fetches the public order book for a particular instrument and depth.
//
// Method: public/get-book
func (c *client) GetBook(ctx context.Context, instrument string, depth int) (*BookResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s%s", c.requester.BaseURL, methodGetBook), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()

	q.Add("instrument_name", instrument)

	if depth > 0 {
		q.Add("depth", fmt.Sprintf("%d", depth))
	}

	req.URL.RawQuery = q.Encode()

	res, err := c.requester.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer res.Body.Close()

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var bookResponse BookResponse
	if err := json.Unmarshal(resBytes, &bookResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	if err := c.requester.CheckErrorResponse(res.StatusCode, bookResponse.Code); err != nil {
		return nil, fmt.Errorf("error received in response: %w", err)
	}

	return &bookResponse.Result, nil
}
