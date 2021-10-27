package cdcexchange

import (
	"context"
	"fmt"

	"github.com/cshep4/crypto-dot-com-exchange-go/errors"
	"github.com/cshep4/crypto-dot-com-exchange-go/internal/api"
	"github.com/cshep4/crypto-dot-com-exchange-go/internal/auth"
	"github.com/cshep4/crypto-dot-com-exchange-go/internal/time"
)

const (
	methodGetOrderDetail = "private/get-order-detail"
)

type (
	// GetOrderDetailResponse is the base response returned from the private/get-order-detail API.
	GetOrderDetailResponse struct {
		// api.BaseResponse is the common response fields.
		api.BaseResponse
		// Result is the response attributes of the endpoint.
		Result GetOrderDetailResult `json:"result"`
	}

	// GetOrderDetailResult is the result returned from the private/get-order-detail API.
	GetOrderDetailResult struct {
		// TradeList is a list of trades for the order (if any).
		TradeList []Trade `json:"trade_list"`
		// OrderInfo is the detailed information about the order.
		OrderInfo Order `json:"order_info"`
	}

	// Trade represents the details of a specific trade.
	Trade struct {
		// Side represents whether the trade is buy or sell.
		Side OrderSide `json:"side"`
		// InstrumentName represents the currency pair to trade (e.g. ETH_CRO or BTC_USDT).
		InstrumentName string `json:"instrument_name"`
		// Fee is the trade fee.
		Fee float64 `json:"fee"`
		// TradeID is the unique identifier for the trade.
		TradeID string `json:"trade_id"`
		// CreateTime is the trade creation time.
		CreateTime time.Time `json:"create_time"`
		// TradedPrice is the executed trade price
		TradedPrice float64 `json:"traded_price"`
		// TradedQuantity is the executed trade quantity
		TradedQuantity float64 `json:"traded_quantity"`
		// FeeCurrency is the currency used for the fees (e.g. CRO).
		FeeCurrency string `json:"fee_currency"`
		// OrderID is the unique identifier for the order.
		OrderID string `json:"order_id"`
	}
)

// GetOrderDetail gets details of an order for a particular order ID
func (c *client) GetOrderDetail(ctx context.Context, orderID string) (*GetOrderDetailResult, error) {
	if orderID == "" {
		return nil, errors.InvalidParameterError{Parameter: "orderID", Reason: "cannot be empty"}
	}

	var (
		id        = c.idGenerator.Generate()
		timestamp = c.clock.Now().UnixMilli()
		params    = make(map[string]interface{})
	)

	params["order_id"] = orderID

	signature, err := c.signatureGenerator.GenerateSignature(auth.SignatureRequest{
		APIKey:    c.apiKey,
		SecretKey: c.secretKey,
		ID:        id,
		Method:    methodGetOrderDetail,
		Timestamp: timestamp,
		Params:    params,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create signature: %w", err)
	}

	body := api.Request{
		ID:        id,
		Method:    methodGetOrderDetail,
		Nonce:     timestamp,
		Params:    params,
		Signature: signature,
		APIKey:    c.apiKey,
	}

	var getOrderDetailResponse GetOrderDetailResponse
	statusCode, err := c.requester.Post(ctx, body, methodGetOrderDetail, &getOrderDetailResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to execute post request: %w", err)
	}

	if err := c.requester.CheckErrorResponse(statusCode, getOrderDetailResponse.Code); err != nil {
		return nil, fmt.Errorf("error received in response: %w", err)
	}

	return &getOrderDetailResponse.Result, nil
}
