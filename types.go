package cdcexchange

import "encoding/json"

type (
	APIRequest struct {
		ID        int64                  `json:"id"`
		Method    string                 `json:"method"`
		Nonce     int64                  `json:"nonce"`
		Params    map[string]interface{} `json:"params"`
		Signature string                 `json:"sig,omitempty"`
		APIKey    string                 `json:"api_key,omitempty"`
	}

	APIBaseResponse struct {
		ID     json.Number  `json:"id"`
		Method string `json:"method"`
		Code   json.Number  `json:"code"`
	}
)
