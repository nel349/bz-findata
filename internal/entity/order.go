package entity

// Order model data of response from exchange
type Order struct {
	Type      string
	Timestamp int64
	ProductID string
	OrderID   string
	Funds     float64
	Side      string
	Size      float64
	Price     float64
	OrderType string
	ClientOID string
}
