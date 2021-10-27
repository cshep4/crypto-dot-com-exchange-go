package cdcexchange

import (
	"context"
	"fmt"

	"github.com/cshep4/crypto-dot-com-exchange-go/internal/api"
)

const (
	methodGetInstruments = "public/get-instruments"
)

type (
	// InstrumentsResponse is the base response returned from the public/get-instruments API.
	InstrumentsResponse struct {
		// api.BaseResponse is the common response fields.
		api.BaseResponse
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
// Method: public/get-instruments
func (c *client) GetInstruments(ctx context.Context) ([]Instrument, error) {
	body := api.Request{
		ID:     c.idGenerator.Generate(),
		Method: methodGetInstruments,
		Nonce:  c.clock.Now().UnixMilli(),
	}

	var instrumentsResponse InstrumentsResponse
	statusCode, err := c.requester.Get(ctx, body, methodGetInstruments, &instrumentsResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to execute post request: %w", err)
	}

	if err := c.requester.CheckErrorResponse(statusCode, instrumentsResponse.Code); err != nil {
		return nil, fmt.Errorf("error received in response: %w", err)
	}

	return instrumentsResponse.Result.Instruments, nil
}
