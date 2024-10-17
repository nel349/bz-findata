package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/dmitryburov/go-coinbase-socket/internal/entity"
	"github.com/jmoiron/sqlx"
)

type exchangeRepo struct {
	db *sqlx.DB
}

// NewExchangeRepository created exchange repository
func NewExchangeRepository(db *sqlx.DB) *exchangeRepo {
	return &exchangeRepo{db}
}

func (e *exchangeRepo) CreateTick(ctx context.Context, message entity.Message) error {
	ctxReq, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if message.Ticker != nil {
		_, err := e.db.NamedExecContext(
			ctxReq,
			"INSERT INTO ticks (symbol, timestamp, bid, ask) VALUES (:symbol, :timestamp, :bid, :ask)",
			message.Ticker,
		)
		return err
	} else if message.Order != nil {
		// Handle order data if needed
		// For now, we'll just log that we received an order
		fmt.Printf("Received order: %+v", message.Order)
		return nil
	}

	return fmt.Errorf("message contains neither ticker nor order data")
}

func (e *exchangeRepo) CreateOrder(ctx context.Context, message entity.Message) error {
	ctxReq, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if message.Order != nil {
		_, err := e.db.NamedExecContext(
			ctxReq,
			"INSERT INTO orders (symbol, timestamp, bid, ask) VALUES (:symbol, :timestamp, :bid, :ask)",
			message.Order)
		return err
	}

	return fmt.Errorf("message contains neither ticker nor order data")
}
