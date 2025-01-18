package analysis

import (
	"context"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/supabase-community/supabase-go"
)

type Service struct {
	db *sqlx.DB
	supabaseClient *supabase.Client
}

type Order struct {
	OrderID   string    `json:"order_id" db:"order_id"` // Use both 'json' and 'db' tags
	Price     float64   `json:"price" db:"price"`
	ProductID string    `json:"product_id,omitempty" db:"product_id,omitempty"`
	Type      string    `json:"type,omitempty" db:"type,omitempty"`
	Timestamp int64     `json:"timestamp,omitempty" db:"timestamp,omitempty"`
}

type ReceivedOrder struct {
	Order
	OrderType string `json:"order_type,omitempty"`
	Size      float64   `json:"size" db:"size"`
	Side      string    `json:"side" db:"side"`
}

type OpenOrder struct {
	Order
	RemainingSize float64 `json:"remaining_size" db:"remaining_size"`
	Side          string  `json:"side" db:"side"`
}

type DoneOrder struct {
	Order
	RemainingSize float64 `json:"remaining_size"`
	Side          string  `json:"side" db:"side"`
	Reason        string  `json:"reason" db:"reason"`
}

type MatchOrder struct {
	Order
	Size          float64 `json:"size" db:"size"`
	RemainingSize float64 `json:"remaining_size" db:"remaining_size"`
	Side          string  `json:"side" db:"side"`
}

func NewService(db *sqlx.DB, supabaseClient *supabase.Client) *Service {
	return &Service{db: db, supabaseClient: supabaseClient}
}

func (s *Service) GetLargestReceivedOrdersInLastNHours(ctx context.Context, hours int, limit int) ([]ReceivedOrder, error) {
	query := `
		SELECT type, product_id, order_id, size, price, side, timestamp
		FROM orders
		WHERE timestamp > ?
		AND type = 'received'
		ORDER BY size DESC
		LIMIT ?
	`
	var orders []ReceivedOrder
	err := s.db.SelectContext(ctx, &orders, query, time.Now().Add(-time.Duration(hours)*time.Hour).UnixNano(), limit)
	if err != nil {
		log.Println("error selecting orders from db", err)
	}
	return orders, err
}

// Get the largest open orders in last N hours
func (s *Service) GetLargestOpenOrdersInLastNHours(ctx context.Context, hours int, limit int) ([]OpenOrder, error) {
	query := `
		SELECT order_id, type, product_id, price, remaining_size, side, timestamp
		FROM orders
		WHERE timestamp > ?
		AND type = 'open'
		ORDER BY remaining_size DESC
		LIMIT ?
	`
	var orders []OpenOrder
	err := s.db.SelectContext(ctx, &orders, query, time.Now().Add(-time.Duration(hours)*time.Hour).UnixNano(), limit)
	return orders, err
}

// Get the largest match orders in last N hours
func (s *Service) GetLargestMatchOrdersInLastNHours(ctx context.Context, hours, limit int) ([]MatchOrder, error) {

	query := `
		SELECT order_id, type, product_id, price, remaining_size, side, timestamp, size
		FROM orders
		WHERE timestamp > ?
		AND type = 'match'
		ORDER BY size DESC
		LIMIT ?
	`
	var orders []MatchOrder
	err := s.db.SelectContext(ctx, &orders, query, time.Now().Add(-time.Duration(hours)*time.Hour).UnixNano(), limit)

	return orders, err
}

// Generic store method that handles different order types
func (s *Service) storeOrdersInSupabase(orders interface{}, tableName string) error {
    // Store in Supabase with the specified table
    _, err := s.supabaseClient.From(tableName).Insert(orders, false, "", "", "").ExecuteTo(&orders)
    if err != nil {
        log.Printf("error inserting orders to supabase table %s: %v", tableName, err)
        return err
    }
    return nil
}

// Helper methods for different order types
func (s *Service) StoreReceivedOrdersInSupabase(ctx context.Context, hours int, limit int) error {
    orders, err := s.GetLargestReceivedOrdersInLastNHours(ctx, hours, limit)
    if err != nil {
        return err
    }
    
    // Convert timestamps for Supabase compatibility
    for i := range orders {
        orders[i].Timestamp = orders[i].Timestamp / 1e9
    }
    return s.storeOrdersInSupabase(orders, "orders")
}

func (s *Service) StoreMatchOrdersInSupabase(ctx context.Context, hours int, limit int) error {
    orders, err := s.GetLargestMatchOrdersInLastNHours(ctx, hours, limit)
    if err != nil {
        return err
    }
    
    // Convert timestamps for Supabase compatibility
    for i := range orders {
        orders[i].Timestamp = orders[i].Timestamp / 1e9
    }
    return s.storeOrdersInSupabase(orders, "orders")
}

func (s *Service) StoreOpenOrdersInSupabase(ctx context.Context, hours int, limit int) error {
    orders, err := s.GetLargestOpenOrdersInLastNHours(ctx, hours, limit)
    if err != nil {
        return err
    }
    
    // Convert timestamps for Supabase compatibility
    for i := range orders {
        orders[i].Timestamp = orders[i].Timestamp / 1e9
    }
    return s.storeOrdersInSupabase(orders, "orders")
}
