package entity

// Order model data of response from exchange
type Order struct {
	Type      string  `db:"type"`
	Timestamp int64   `db:"timestamp"`
	ProductID string  `db:"product_id"`
	OrderID   string  `db:"order_id"`
	Funds     float64 `db:"funds"`
	Side      string  `db:"side"`
	Size      float64 `db:"size"`
	Price     float64 `db:"price"`
	OrderType string  `db:"order_type"`
	ClientOID   string  `db:"client_oid"`
	Sequence    int     `db:"sequence"`
	RemainingSize float64 `db:"remaining_size"`
	Reason        string  `db:"reason"`
	TradeID       int64   `db:"trade_id"`
	MakerOrderID  string  `db:"maker_order_id"`
	TakerOrderID  string  `db:"taker_order_id"`
}	
