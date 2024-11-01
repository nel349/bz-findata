package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/nel349/bz-findata/pkg/entity"
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
	}

	return fmt.Errorf("message should be ticker")
}

func (e *exchangeRepo) CreateOrder(ctx context.Context, message entity.Message) error {
	ctxReq, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if message.Order != nil {
		_, err := e.db.NamedExecContext(
			ctxReq,
			`INSERT INTO orders (type, product_id, timestamp, order_id, funds, side, size, price, order_type, client_oid, sequence, remaining_size, reason) 
			 VALUES (:type, :product_id, :timestamp, :order_id, :funds, :side, :size, :price, :order_type, :client_oid, IFNULL(:sequence, 0), :remaining_size, :reason)`,
			message.Order)
		if err != nil {
			fmt.Println("Error inserting order", "error", err)
			return err
		}
		fmt.Println("Order inserted successfully")
		return nil
	}

	return fmt.Errorf("message should be order")
}
