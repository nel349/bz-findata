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
	OrderID string  `db:"order_id"`
	Size    float64 `db:"size"`
	Price   float64 `db:"price"`
}

func NewService(db *sqlx.DB) *Service {
	return &Service{db: db}
}

func (s *Service) GetLargestOrdersInLastNHours(ctx context.Context, hours int, limit int) ([]Order, error) {
	query := `
		SELECT order_id, size, price
		FROM orders
		WHERE timestamp > ?
		ORDER BY size * price DESC
		LIMIT ?
	`
	var orders []Order
	err := s.db.SelectContext(ctx, &orders, query, time.Now().Add(-time.Duration(hours)*time.Hour), limit)
	return orders, err
}