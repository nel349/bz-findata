package entity

// Order model data of response from exchange
type Order struct {
	Timestamp int64
	ProductID string
	OrderID   string
	Funds     float64
	Side      string
	OrderType string
	ClientOID string
}
