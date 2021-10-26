package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
)

type (
	SignatureRequest struct {
		APIKey    string
		SecretKey string
		ID        int64
		Method    string
		Timestamp int64
		Params    map[string]interface{}
	}

	SignatureGenerator interface {
		GenerateSignature(req SignatureRequest) (string, error)
	}

	param struct {
		key string
		val interface{}
	}

	Generator struct{}
)

func (g Generator) GenerateSignature(req SignatureRequest) (string, error) {
	paramStr := g.buildParamString(req.Params)

	signaturePayload := fmt.Sprintf("%s%d%s%s%d", req.Method, req.ID, req.APIKey, paramStr, req.Timestamp)

	h := hmac.New(sha256.New, []byte(req.SecretKey))

	_, err := h.Write([]byte(signaturePayload))
	if err != nil {
		return "", fmt.Errorf("failed to write signature: %w", err)
	}

	return hex.EncodeToString(h.Sum(nil)), nil

}

func (g Generator) buildParamString(params map[string]interface{}) string {
	if len(params) == 0 {
		return ""
	}

	var paramsString string

	for _, p := range g.sortParams(params) {
		paramsString = fmt.Sprintf("%s%s%v", paramsString, p.key, p.val)
	}

	return paramsString
}

func (Generator) sortParams(params map[string]interface{}) []param {
	p := make([]param, 0, len(params))

	for k, v := range params {
		p = append(p, param{key: k, val: v})
	}

	sort.Slice(p, func(i, j int) bool {
		return p[i].key < p[j].key
	})

	return p
}
