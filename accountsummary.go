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
	methodGetAccountSummary = "private/get-account-summary"
)

type (
	// AccountSummaryResponse is the base response returned from the private/get-account-summary API.
	AccountSummaryResponse struct {
		// APIBaseResponse is the common response fields.
		APIBaseResponse
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

	signature, err := c.generateSignature(signatureRequest{
		apiKey:    c.apiKey,
		secretKey: c.secretKey,
		id:        id,
		method:    methodGetAccountSummary,
		timestamp: timestamp,
		params:    params,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create signature: %w", err)
	}

	body := APIRequest{
		ID:        id,
		Method:    methodGetAccountSummary,
		Nonce:     timestamp,
		Params:    params,
		Signature: signature,
		APIKey:    c.apiKey,
	}

	var accountSummaryResponse AccountSummaryResponse
	statusCode, err := c.postRequest(ctx, body, methodGetAccountSummary, &accountSummaryResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to execute post request: %w", err)
	}

	if statusCode > 299 {
		code, err := accountSummaryResponse.Code.Int64()
		if err != nil {
			return nil, fmt.Errorf("invalid http status code: %d - response code: %v", statusCode, accountSummaryResponse.Code)
		}
		return nil, newResponseError(code)
	}

	return accountSummaryResponse.Result.Accounts, nil
}

func (c *client) postRequest(ctx context.Context, body APIRequest, method string, response interface{}) (int, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s%s", c.baseURL, method), bytes.NewBuffer(b))
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
