package coinbase

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/dmitryburov/go-coinbase-socket/internal/entity"
)

// Base Response struct
type Response struct {
	Type      ResponseType `json:"type,string"`
	Message   string       `json:"message,omitempty"`
	Reason    string       `json:"reason,omitempty"`
	ProductID string       `json:"product_id"`
}

// TickerResponse for ticker data
type TickerResponse struct {
	Response
	BestBid string `json:"best_bid"`
	BestAsk string `json:"best_ask"`
}

// ReceivedOrderResponse for received order data
type ReceivedOrderResponse struct {
	Response
	Time       string  `json:"time"`
	Sequence   int64   `json:"sequence"`
	OrderID    string  `json:"order_id"`
	Size       float64  `json:"size,omitempty"`  // Only for limit orders
	Price      float64  `json:"price,omitempty"` // Only for limit orders
	Funds      float64  `json:"funds,omitempty"` // Only for market orders
	Side       string  `json:"side"`
	OrderType  string  `json:"order_type"`
	ClientOID  string  `json:"client-oid"` // Note the hyphen in the JSON tag
}

type ResponseType int

const (
	Error ResponseType = iota
	Subscriptions
	Unsubscribe
	Heartbeat
	Ticker
	Level2
)

var responseTypeNames = [...]string{"error", "subscriptions", "unsubscribe", "heartbeat", "ticker", "level2"}

func (r ResponseType) String() string {
	return responseTypeNames[r]
}

// Conversion methods
func (r *TickerResponse) ToTicker() (*entity.Ticker, error) {
	bid, err := strconv.ParseFloat(r.BestBid, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid bid: %w", err)
	}

	ask, err := strconv.ParseFloat(r.BestAsk, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid ask: %w", err)
	}

	return &entity.Ticker{
		Timestamp: time.Now().UnixNano(),
		Bid:       bid,
		Ask:       ask,
		Symbol:    r.ProductID,
	}, nil
}

func (r *ReceivedOrderResponse) ToReceivedOrder() (*entity.Order, error) {
	return &entity.Order{
		OrderID:   r.OrderID,
		OrderType: r.OrderType,
		Funds:     r.Funds,
		Side:      r.Side,
		ClientOID: r.ClientOID,
		ProductID: r.ProductID,
		// Set other fields as needed
	}, nil
}

func (r *ResponseType) UnmarshalJSON(v []byte) error {
	str := string(v)
	for i, name := range responseTypeNames {
		if name == str {
			*r = ResponseType(i)
			return nil
		}
	}

	return fmt.Errorf("invalid locality type %q", str)
}

func ParseResponse(message []byte) (interface{}, error) {
	var baseResponse Response
	if err := json.Unmarshal(message, &baseResponse); err != nil {
		return nil, err
	}

	switch baseResponse.Type {
	case Ticker:
		var tickerResponse TickerResponse
		if err := json.Unmarshal(message, &tickerResponse); err != nil {
			return nil, err
		}
		return &tickerResponse, nil
	// case ReceivedOrder:
	// 	var orderResponse ReceivedOrderResponse
	// 	if err := json.Unmarshal(message, &orderResponse); err != nil {
	// 		return nil, err
	// 	}
	// 	return &orderResponse, nil
	// Add other cases as needed
	default:
		return &baseResponse, nil
	}
}
