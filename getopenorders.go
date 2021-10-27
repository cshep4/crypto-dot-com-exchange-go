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
	methodGetOpenOrders = "private/get-open-orders"

	OrderStatusActive    OrderStatus = "ACTIVE"
	OrderStatusCancelled OrderStatus = "CANCELED"
	OrderStatusFilled    OrderStatus = "FILLED"
	OrderStatusRejected  OrderStatus = "REJECTED"
	OrderStatusExpired   OrderStatus = "EXPIRED"
)

type (
	// OrderStatus is the current status of the order.
	OrderStatus string

	// GetOpenOrdersRequest is the request params sent for the private/get-open-orders API.
	GetOpenOrdersRequest struct {
		// InstrumentName represents the currency pair for the orders (e.g. ETH_CRO or BTC_USDT).
		// if InstrumentName is omitted, all instruments will be returned.
		InstrumentName string `json:"instrument_name"`
		// PageSize represents maximum number of orders returned (for pagination)
		// (Default: 20, Max: 200)
		// if PageSize is 0, it will be set as 20 by default.
		PageSize int `json:"page_size"`
		// Page represents the page number (for pagination)
		// (0-based)
		Page int `json:"page"`
	}

	// GetOpenOrdersResponse is the base response returned from the private/get-open-orders API.
	GetOpenOrdersResponse struct {
		// api.BaseResponse is the common response fields.
		api.BaseResponse
		// Result is the response attributes of the endpoint.
		Result GetOpenOrdersResult `json:"result"`
	}

	// GetOpenOrdersResult is the result returned from the private/get-open-orders API.
	GetOpenOrdersResult struct {
		// Count is the total count of orders.
		Count int `json:"order_id"`
		// OrderList is the array of open orders.
		OrderList []Order `json:"order_list"`
	}

	// Order represents the details of a specific order.
	// Note: To detect a 'partial filled' status, look for status as ACTIVE and cumulative_quantity > 0.
	Order struct {
		// Status is the status of the order, can be ACTIVE, CANCELED, FILLED, REJECTED or EXPIRED.
		Status OrderStatus `json:"status"`
		// Reason is the reason code for rejected orders (see "Response and Reason Codes").
		Reason string `json:"reason"`
		// Side represents whether the order is buy or sell.
		Side OrderSide `json:"side"`
		// Price is the price specified in the order.
		Price float64 `json:"price"`
		// Quantity	is the quantity specified in the order.
		Quantity float64 `json:"quantity"`
		// OrderID is the unique identifier for the order.
		OrderID string `json:"order_id"`
		// ClientOID is the optional Client order ID (if provided in request).
		ClientOID string `json:"client_oid"`
		// CreateTime is the order creation time.
		CreateTime time.Time `json:"create_time"`
		// UpdateTime is the order update time.
		UpdateTime time.Time `json:"update_time"`
		// Type represents the type of order.
		OrderType OrderType `json:"type"`
		// InstrumentName represents the currency pair to trade (e.g. ETH_CRO or BTC_USDT).
		InstrumentName string `json:"instrument_name"`
		// CumulativeQuantity is the cumulative-executed quantity (for partially filled orders).
		CumulativeQuantity float64 `json:"cumulative_quantity"`
		// CumulativeValue is the cumulative-executed value (for partially filled orders).
		CumulativeValue float64 `json:"cumulative_value"`
		// AvgPrice is the average filled price. If none is filled, 0 is returned.
		AvgPrice float64 `json:"avg_price"`
		// FeeCurrency is the currency used for the fees (e.g. CRO).
		FeeCurrency string `json:"fee_currency"`
		// TimeInForce represents how long the order should be active before being cancelled.
		// (Limit Orders Only) Options are:
		//  - GOOD_TILL_CANCEL (Default if unspecified)
		//  - FILL_OR_KILL
		//  - IMMEDIATE_OR_CANCEL
		TimeInForce TimeInForce `json:"time_in_force"`
		// (Limit Orders Only) Options are:
		// - POST_ONLY
		// - Or leave empty
		ExecInst ExecInst `json:"exec_inst"`
		// TriggerPrice is the price at which the order is triggered.
		// Used with STOP_LOSS, STOP_LIMIT, TAKE_PROFIT, and TAKE_PROFIT_LIMIT orders.
		TriggerPrice float64 `json:"trigger_price"`
	}
)

// GetOpenOrders gets all open orders for a particular instrument
// Pagination is handled using page size (Default: 20, Max: 200) & number (0-based).
// req.InstrumentName can be left blank to get open orders for all instruments.
// Method: private/get-open-orders
func (c *client) GetOpenOrders(ctx context.Context, req GetOpenOrdersRequest) (*GetOpenOrdersResult, error) {
	if req.PageSize < 0 {
		return nil, errors.InvalidParameterError{Parameter: "req.PageSize", Reason: "cannot be less than 0"}
	}
	if req.PageSize > 200 {
		return nil, errors.InvalidParameterError{Parameter: "req.PageSize", Reason: "cannot be greater than 200"}
	}

	var (
		id        = c.idGenerator.Generate()
		timestamp = c.clock.Now().UnixMilli()
		params    = make(map[string]interface{})
	)

	if req.InstrumentName != "" {
		params["instrument_name"] = req.InstrumentName
	}
	if req.PageSize != 0 {
		params["page_size"] = req.PageSize
	}
	params["page"] = req.Page

	signature, err := c.signatureGenerator.GenerateSignature(auth.SignatureRequest{
		APIKey:    c.apiKey,
		SecretKey: c.secretKey,
		ID:        id,
		Method:    methodGetOpenOrders,
		Timestamp: timestamp,
		Params:    params,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create signature: %w", err)
	}

	body := api.Request{
		ID:        id,
		Method:    methodGetOpenOrders,
		Nonce:     timestamp,
		Params:    params,
		Signature: signature,
		APIKey:    c.apiKey,
	}

	var getOpenOrdersResponse GetOpenOrdersResponse
	statusCode, err := c.requester.Post(ctx, body, methodGetOpenOrders, &getOpenOrdersResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to execute post request: %w", err)
	}

	if err := c.requester.CheckErrorResponse(statusCode, getOpenOrdersResponse.Code); err != nil {
		return nil, fmt.Errorf("error received in response: %w", err)
	}

	return &getOpenOrdersResponse.Result, nil
}
