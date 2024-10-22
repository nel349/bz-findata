package analysis

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type Service struct {
	db *sqlx.DB
}

type Order struct {
	OrderID   string    `db:"order_id"`
	Price     float64   `db:"price"`
	ProductID string    `db:"product_id,omitempty"`
	Type      string    `db:"type,omitempty"`
	Timestamp int64     `db:"timestamp,omitempty"`
}

type ReceivedOrder struct {
	Order
	OrderType string `db:"order_type,omitempty"`
	Size      float64   `db:"size"`
	Side      string    `db:"side"`
}

type OpenOrder struct {
	Order
	RemainingSize float64 `db:"remaining_size"`
	Side          string  `db:"side"`
}

type DoneOrder struct {
	Order
	RemainingSize float64 `db:"remaining_size"`
	Side          string  `db:"side"`
	Reason        string  `db:"reason"`
}

func NewService(db *sqlx.DB) *Service {
	return &Service{db: db}
}

func (s *Service) GetLargestOrdersInLastNHours(ctx context.Context, hours int, limit int) ([]Order, error) {
	query := `
		SELECT order_id, type, product_id, price
		FROM orders
		WHERE timestamp > ?
		ORDER BY size * price DESC
		LIMIT ?
	`
	var orders []Order
	err := s.db.SelectContext(ctx, &orders, query, time.Now().Add(-time.Duration(hours)*time.Hour), limit)
	return orders, err
}

func (s *Service) GetLargestReceivedOrdersInLastNHours(ctx context.Context, hours int, limit int) ([]ReceivedOrder, error) {
	query := `
		SELECT type, product_id, order_id, size, price, side
		FROM orders
		WHERE timestamp > ?
		AND type = 'received'
		ORDER BY size DESC
		LIMIT ?
	`
	var orders []ReceivedOrder
	err := s.db.SelectContext(ctx, &orders, query, time.Now().Add(-time.Duration(hours)*time.Hour), limit)
	return orders, err
}

// Get the largest open orders in last N hours
func (s *Service) GetLargestOpenOrdersInLastNHours(ctx context.Context, hours int, limit int) ([]OpenOrder, error) {
	query := `
		SELECT order_id, type, product_id, price, remaining_size, side
		FROM orders
		WHERE timestamp > ?
		AND type = 'open'
		ORDER BY remaining_size DESC
		LIMIT ?
	`
	var orders []OpenOrder
	err := s.db.SelectContext(ctx, &orders, query, time.Now().Add(-time.Duration(hours)*time.Hour), limit)
	return orders, err
}
