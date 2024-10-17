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
	Type      string       `json:"type"`
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
	Time      string  `json:"time"`
	Sequence  int64   `json:"sequence"`
	OrderID   string  `json:"order_id"`
	Size      string  `json:"size,omitempty"`  // Only for limit orders
	Price     string  `json:"price,omitempty"` // Only for limit orders
	Funds     string  `json:"funds,omitempty"` // Only for market orders
	Side      string  `json:"side"`
	OrderType string  `json:"order_type"`
	ClientOID string  `json:"client-oid"` // Note the hyphen in the JSON tag
}

type HeartbeatResponse struct {
	Type        string    `json:"type"`
	Sequence    int64     `json:"sequence"`
	LastTradeID int64     `json:"last_trade_id"`
	ProductID   string    `json:"product_id"`
	Time        time.Time `json:"time"`
}
type ResponseType int

const (
	Error ResponseType = iota
	Subscriptions
	Unsubscribe
	Heartbeat
	Ticker
	Level2
	Received
)

var responseTypeNames = [...]string{"error", "subscriptions", "unsubscribe", "heartbeat", "ticker", "level2", "received"}

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

// Update the ToReceivedOrder method to handle string conversions
func (r *ReceivedOrderResponse) ToReceivedOrder() (*entity.Order, error) {
	size, err := strconv.ParseFloat(r.Size, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid size: %w", err)
	}

	price, err := strconv.ParseFloat(r.Price, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid price: %w", err)
	}

	funds, err := strconv.ParseFloat(r.Funds, 64)
	if err != nil && r.Funds != "" {
		return nil, fmt.Errorf("invalid funds: %w", err)
	}

	return &entity.Order{
		OrderID:   r.OrderID,
		OrderType: r.OrderType,
		Size:      size,
		Price:     price,
		Funds:     funds,
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
		return nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	switch baseResponse.Type {
	case Ticker.String():
		var tickerResponse TickerResponse
		if err := json.Unmarshal(message, &tickerResponse); err != nil {
			return nil, err
		}
		return &tickerResponse, nil

	case Heartbeat.String():
		var heartbeatResponse HeartbeatResponse
		if err := json.Unmarshal(message, &heartbeatResponse); err != nil {
			return nil, err
		}
		return &heartbeatResponse, nil

	case Received.String():
		fmt.Printf("Message: %s\n", message)
		var orderResponse ReceivedOrderResponse
		if err := json.Unmarshal(message, &orderResponse); err != nil {
			return nil, err
		}
		return &orderResponse, nil
	// Add other cases as needed
	default:
		fmt.Printf("Unknown response type: %s and pr", baseResponse.Type)
		return &baseResponse, nil
	}
}
