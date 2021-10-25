package cdcexchange

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	methodGetInstruments = "public/get-instruments"
)

type (
	// InstrumentsResponse is the base response returned from the public/get-instruments API.
	InstrumentsResponse struct {
		// APIBaseResponse is the common response fields.
		APIBaseResponse
		// Result is the response attributes of the endpoint.
		Result InstrumentResult `json:"result"`
	}

	// InstrumentResult is the result returned from the public/get-instruments API.
	InstrumentResult struct {
		// Instruments is a list of the returned instruments.
		Instruments []Instrument `json:"instruments"`
	}

	// Instrument represents details of a specific currency pair
	Instrument struct {
		// InstrumentName represents the name of the instrument (e.g. BTC_USDT).
		InstrumentName string `json:"instrument_name"`
		// QuoteCurrency represents the quote currency (e.g. USDT).
		QuoteCurrency string `json:"quote_currency"`
		// BaseCurrency represents the base currency (e.g. BTC).
		BaseCurrency string `json:"base_currency"`
		// PriceDecimals is the maximum decimal places for specifying price.
		PriceDecimals int `json:"price_decimals"`
		// QuantityDecimals is the maximum decimal places for specifying quantity.
		QuantityDecimals int `json:"quantity_decimals"`
		// MarginTradingEnabled represents whether margin trading is enabled for the instrument.
		MarginTradingEnabled bool `json:"margin_trading_enabled"`
	}
)

// GetInstruments provides information on all supported instruments (e.g. BTC_USDT).
func (c *client) GetInstruments(ctx context.Context) ([]Instrument, error) {
	body := APIRequest{
		ID:     c.idGenerator.Generate(),
		Method: methodGetInstruments,
		Nonce:  c.clock.Now().UnixMilli(),
	}

	var instrumentsResponse InstrumentsResponse
	statusCode, err := c.getRequest(ctx, body, methodGetInstruments, &instrumentsResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to execute post request: %w", err)
	}

	if statusCode > 299 {
		code, err := instrumentsResponse.Code.Int64()
		if err != nil {
			return nil, fmt.Errorf("invalid http status code: %d - response code: %v", statusCode, instrumentsResponse.Code)
		}
		return nil, newResponseError(code)
	}

	return instrumentsResponse.Result.Instruments, nil
}

func (c *client) getRequest(ctx context.Context, body APIRequest, method string, response interface{}) (int, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s%s", c.baseURL, method), bytes.NewBuffer(b))
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to do request: %w", err)
	}
	defer res.Body.Close()

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(resBytes, &response); err != nil {
		return 0, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return res.StatusCode, nil
}
