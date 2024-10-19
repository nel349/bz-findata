package entity

// Order model data of response from exchange
type Order struct {
	Type      string
	Timestamp int64
	ProductID string  `db:"product_id"`
	OrderID   string  `db:"order_id"`
	Funds     float64 `db:"funds"`
	Side      string  `db:"side"`
	Size      float64 `db:"size"`
	Price     float64 `db:"price"`
	OrderType string  `db:"order_type"`
	ClientOID string  `db:"client_oid"`
	Sequence  int     `db:"sequence"`
}
