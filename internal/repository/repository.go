package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/nel349/bz-findata/pkg/entity"
	"github.com/nel349/bz-findata/internal/repository/mysql"
)

// Exchange methods implementation
type Exchange interface {
	// CreateTick write in storage ticker data
	CreateTick(ctx context.Context, message entity.Message) error
	// CreateOrder write in storage order data
	CreateOrder(ctx context.Context, message entity.Message) error
}

// Repositories of based interface for repository layout
type Repositories struct {
	Exchange
}

// NewRepositories init repository layout
func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		Exchange: mysql.NewExchangeRepository(db),
	}
}
