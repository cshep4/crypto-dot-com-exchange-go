package cdcexchange

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cshep4/crypto-dot-com-exchange-go/internal/api"
)

const (
	methodGetTicker = "public/get-ticker"
)

type (
	// TickerResponse is the base response returned from the public/get-ticker API
	// when no instrument is specified.
	TickerResponse struct {
		// api.BaseResponse is the common response fields.
		api.BaseResponse
		// Result is the response attributes of the endpoint.
		Result TickerResult `json:"result"`
	}

	TickerResult struct {
		// Data is the returned ticker data for all instruments.
		Data []Ticker `json:"data"`
	}

	// SingleTickerResponse is the base response returned from the public/get-ticker API
	// when an instrument is specified.
	SingleTickerResponse struct {
		// api.BaseResponse is the common response fields.
		api.BaseResponse
		// Result is the response attributes of the endpoint.
		Result SingleTickerResult `json:"result"`
	}

	// SingleTickerResult is the result returned from the public/get-ticker API
	// when an instrument is specified.
	SingleTickerResult struct {
		// Data is the returned ticker data for specified instrument.
		Data Ticker `json:"data"`
	}

	// Ticker represents ticker details of a specific currency pair.
	Ticker struct {
		// Instrument is the instrument name (e.g. BTC_USDT, ETH_CRO, etc).
		Instrument string `json:"i"`
		// BidPrice is the current best bid price, 0 if there aren't any bids.
		BidPrice float64 `json:"b"`
		// AskPrice is the current best ask price, 0 if there aren't any asks.
		AskPrice float64 `json:"k"`
		// LatestTradePrice is the price of the latest trade, 0 if there weren't any trades.
		LatestTradePrice float64 `json:"a"`
		// Timestamp is the unix timestamp of the data.
		Timestamp float64 `json:"t"`
		// Volume24H is the total 24h traded volume.
		Volume24H float64 `json:"v"`
		// PriceHigh24h is the price of the 24h highest trade, 0 if there weren't any trades.
		PriceHigh24h float64 `json:"h"`
		// PriceLow24h is the price of the 24h lowest trade, 0 if there weren't any trades.
		PriceLow24h float64 `json:"l"`
		// PriceChange24h is the 24-hour price change, 0 if there weren't any trades.
		PriceChange24h float64 `json:"c"`
	}
)

// GetTickers fetches the public tickers for an instrument (e.g. BTC_USDT).
//
// instrument can be left blank to retrieve tickers for ALL instruments.
//
// Method: public/get-ticker
func (c *client) GetTickers(ctx context.Context, instrument string) ([]Ticker, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s%s", c.requester.BaseURL, methodGetTicker), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// if instrument is omitted, ALL tickers are returned.
	if instrument != "" {
		q := req.URL.Query()
		q.Add("instrument_name", instrument)
		req.URL.RawQuery = q.Encode()
	}

	res, err := c.requester.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer res.Body.Close()

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var (
		tickers []Ticker
		code    json.Number
	)

	if instrument == "" {
		var tickerResponse TickerResponse
		if err := json.Unmarshal(resBytes, &tickerResponse); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
		}

		tickers = tickerResponse.Result.Data
		code = tickerResponse.Code
	} else {
		var tickerResponse SingleTickerResponse
		if err := json.Unmarshal(resBytes, &tickerResponse); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
		}

		tickers = []Ticker{tickerResponse.Result.Data}
		code = tickerResponse.Code
	}

	if err := c.requester.CheckErrorResponse(res.StatusCode, code); err != nil {
		return nil, fmt.Errorf("error received in response: %w", err)
	}

	return tickers, nil
}
