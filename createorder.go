package cdcexchange

import (
	"context"
	"fmt"

	"github.com/cshep4/crypto-dot-com-exchange-go/internal/api"
	"github.com/cshep4/crypto-dot-com-exchange-go/internal/auth"
)

const (
	methodCreateOrder = "private/create-order"

	OrderSideBuy  OrderSide = "BUY"
	OrderSideSell OrderSide = "SELL"

	OrderTypeLimit           OrderType = "LIMIT"
	OrderTypeMarket          OrderType = "MARKET"
	OrderTypeStopLoss        OrderType = "STOP_LOSS"
	OrderTypeStopLimit       OrderType = "STOP_LIMIT"
	OrderTypeTakeProfit      OrderType = "TAKE_PROFIT"
	OrderTypeTakeProfitLimit OrderType = "TAKE_PROFIT_LIMIT"

	TimeInForceGoodTilCancelled  TimeInForce = "GOOD_TILL_CANCEL"
	TimeInForceFillOrKill        TimeInForce = "FILL_OR_KILL"
	TimeInForceImmediateOrCancel TimeInForce = "IMMEDIATE_OR_CANCEL"

	ExecInstPostOnly ExecInst = "POST_ONLY"
)

type (
	// OrderSide is the side of the order (BUY/SELL).
	OrderSide string
	// OrderType is the type of order (e.g. LIMIT, MARKET, etc).
	OrderType string
	// TimeInForce represents how long the order should be active before being cancelled.
	TimeInForce string
	// ExecInst for Limit Orders Only (POST_ONLY or left blank).
	ExecInst string

	// CreateOrderRequest is the request params sent for the private/create-order API.
	// Mandatory parameters based on order type:
	// ------------------+------+-----------------------------------------
	// Type 			 | Side | Additional Mandatory Parameters
	// ------------------+------+-----------------------------------------
	// LIMIT 			 | Both | quantity, price
	// MARKET 			 | BUY  | notional or quantity, mutually exclusive
	// MARKET 			 | SELL | quantity
	// STOP_LIMIT 		 | Both | price, quantity, trigger_price
	// TAKE_PROFIT_LIMIT | Both | price, quantity, trigger_price
	// STOP_LOSS 		 | BUY  | notional, trigger_price
	// STOP_LOSS 		 | SELL | quantity, trigger_price
	// TAKE_PROFIT 	  	 | BUY  | notional, trigger_price
	// TAKE_PROFIT 	  	 | SELL | quantity, trigger_price
	// ------------------+------+-----------------------------------------
	CreateOrderRequest struct {
		// InstrumentName represents the currency pair to trade (e.g. ETH_CRO or BTC_USDT).
		InstrumentName string `json:"instrument_name"`
		// Side represents whether the order is buy or sell.
		Side OrderSide `json:"side"`
		// Type represents the type of order.
		Type OrderType `json:"type"`
		// Price determines the price of which the trade should be executed.
		// For LIMIT and STOP_LIMIT orders only.
		Price float64 `json:"price"`
		// Quantity is the quantity to be sold
		// For LIMIT, MARKET, STOP_LOSS, TAKE_PROFIT orders only.
		Quantity float64 `json:"quantity"`
		// Notional is the amount to spend.
		// For MARKET (BUY), STOP_LOSS (BUY), TAKE_PROFIT (BUY) orders only.
		Notional float64 `json:"notional"`
		// ClientOID is the optional Client order ID.
		ClientOID string `json:"client_oid"`
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

	// CreateOrderResponse is the base response returned from the private/create-order API.
	CreateOrderResponse struct {
		// api.BaseResponse is the common response fields.
		api.BaseResponse
		// Result is the response attributes of the endpoint.
		Result CreateOrderResult `json:"result"`
	}

	// CreateOrderResult is the result returned from the private/create-order API.
	CreateOrderResult struct {
		// OrderID is the newly created order ID.
		OrderID int64 `json:"order_id"`
		// ClientOID is the optional Client order ID (if provided in request).
		ClientOID string `json:"client_oid"`
	}
)

// CreateOrder creates a new BUY or SELL order on the Exchange.
//
// This call is asynchronous, so the response is simply a confirmation of the request.
//
// The user.order subscription can be used to check when the order is successfully created.
//
// Method: private/create-order
func (c *client) CreateOrder(ctx context.Context, req CreateOrderRequest) (*CreateOrderResult, error) {
	var (
		id        = c.idGenerator.Generate()
		timestamp = c.clock.Now().UnixMilli()
		params    = make(map[string]interface{})
	)

	if req.InstrumentName != "" {
		params["instrument_name"] = req.InstrumentName
	}
	if req.Side != "" {
		params["side"] = req.Side
	}
	if req.Type != "" {
		params["type"] = req.Type
	}
	if req.Price != 0 {
		params["price"] = req.Price
	}
	if req.Quantity != 0 {
		params["quantity"] = req.Quantity
	}
	if req.Notional != 0 {
		params["notional"] = req.Notional
	}
	if req.ClientOID != "" {
		params["client_oid"] = req.ClientOID
	}
	if req.TimeInForce != "" {
		params["time_in_force"] = req.TimeInForce
	}
	if req.ExecInst != "" {
		params["exec_inst"] = req.ExecInst
	}
	if req.TriggerPrice != 0 {
		params["trigger_price"] = req.TriggerPrice
	}

	signature, err := c.signatureGenerator.GenerateSignature(auth.SignatureRequest{
		APIKey:    c.apiKey,
		SecretKey: c.secretKey,
		ID:        id,
		Method:    methodCreateOrder,
		Timestamp: timestamp,
		Params:    params,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create signature: %w", err)
	}

	body := api.Request{
		ID:        id,
		Method:    methodCreateOrder,
		Nonce:     timestamp,
		Params:    params,
		Signature: signature,
		APIKey:    c.apiKey,
	}

	var createOrderResponse CreateOrderResponse
	statusCode, err := c.requester.Post(ctx, body, methodCreateOrder, &createOrderResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to execute post request: %w", err)
	}

	if err := c.requester.CheckErrorResponse(statusCode, createOrderResponse.Code); err != nil {
		return nil, fmt.Errorf("error received in response: %w", err)
	}

	return &createOrderResponse.Result, nil
}
