package api

import "encoding/json"

type (
	Request struct {
		ID        int64                  `json:"id"`
		Method    string                 `json:"method"`
		Nonce     int64                  `json:"nonce"`
		Params    map[string]interface{} `json:"params"`
		Signature string                 `json:"sig,omitempty"`
		APIKey    string                 `json:"api_key,omitempty"`
	}

	BaseResponse struct {
		ID     json.Number `json:"id"`
		Method string      `json:"method"`
		Code   json.Number `json:"code"`
	}
)
