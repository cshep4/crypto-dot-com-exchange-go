package cdcexchange

import (
	"context"
	"fmt"
	
	"github.com/cshep4/crypto-dot-com-exchange-go/internal/api"
	"github.com/cshep4/crypto-dot-com-exchange-go/internal/auth"
)

const (
	methodGetAccountSummary = "private/get-account-summary"
)

type (
	// AccountSummaryResponse is the base response returned from the private/get-account-summary API.
	AccountSummaryResponse struct {
		// api.BaseResponse is the common response fields.
		api.BaseResponse
		// Result is the response attributes of the endpoint.
		Result AccountSummaryResult `json:"result"`
	}

	// AccountSummaryResult is the result returned from the private/get-account-summary API.
	AccountSummaryResult struct {
		// Accounts is the returned account data.
		Accounts []Account `json:"accounts"`
	}

	// Account represents balance details of a specific token.
	Account struct {
		// Balance is the total balance (Available + Order + Stake).
		Balance float64 `json:"balance"`
		// Available is the available balance (e.g. not in orders, or locked, etc.).
		Available float64 `json:"available"`
		// Order is the balance locked in orders.
		Order float64 `json:"order"`
		// Stake is the balance locked for staking (typically only used for CRO).
		Stake float64 `json:"stake"`
		// Currency is the symbol for the currency (e.g. CRO).
		Currency string `json:"currency"`
	}
)

// GetAccountSummary returns the account balance of a user for a particular token.
// currency can be left blank to retrieve balances for ALL tokens.
func (c *client) GetAccountSummary(ctx context.Context, currency string) ([]Account, error) {
	var (
		id        = c.idGenerator.Generate()
		timestamp = c.clock.Now().UnixMilli()
		params    = make(map[string]interface{})
	)

	// if currency is omitted, ALL currencies are returned.
	if currency != "" {
		params["currency"] = currency
	}

	signature, err := c.signatureGenerator.GenerateSignature(auth.SignatureRequest{
		APIKey:    c.apiKey,
		SecretKey: c.secretKey,
		ID:        id,
		Method:    methodGetAccountSummary,
		Timestamp: timestamp,
		Params:    params,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create signature: %w", err)
	}

	body := api.Request{
		ID:        id,
		Method:    methodGetAccountSummary,
		Nonce:     timestamp,
		Params:    params,
		Signature: signature,
		APIKey:    c.apiKey,
	}

	var accountSummaryResponse AccountSummaryResponse
	statusCode, err := c.requester.Post(ctx, body, methodGetAccountSummary, &accountSummaryResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to execute post request: %w", err)
	}

	if err := c.requester.CheckErrorResponse(statusCode, accountSummaryResponse.Code); err != nil {
		return nil, fmt.Errorf("error received in response: %w", err)
	}

	return accountSummaryResponse.Result.Accounts, nil
}
