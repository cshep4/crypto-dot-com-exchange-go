package cdcexchange

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
)

type (
	signatureRequest struct {
		apiKey    string
		secretKey string
		id        int64
		method    string
		timestamp int64
		params    map[string]interface{}
	}

	param struct {
		key string
		val interface{}
	}
)

func (c *client) generateSignature(req signatureRequest) (string, error) {
	paramStr := c.buildParamString(req.params)

	signaturePayload := fmt.Sprintf("%s%d%s%s%d", req.method, req.id, req.apiKey, paramStr, req.timestamp)

	h := hmac.New(sha256.New, []byte(req.secretKey))

	_, err := h.Write([]byte(signaturePayload))
	if err != nil {
		return "", fmt.Errorf("failed to write signature: %w", err)
	}

	return hex.EncodeToString(h.Sum(nil)), nil

}

func (c *client) buildParamString(params map[string]interface{}) string {
	if len(params) == 0 {
		return ""
	}

	var paramsString string

	for _, p := range c.sortParams(params) {
		paramsString = fmt.Sprintf("%s%s%v", paramsString, p.key, p.val)
	}

	return paramsString
}

func (c *client) sortParams(params map[string]interface{}) []param {
	p := make([]param, 0, len(params))

	for k, v := range params {
		p = append(p, param{key: k, val: v})
	}

	sort.Slice(p, func(i, j int) bool {
		return p[i].key < p[j].key
	})

	return p
}
