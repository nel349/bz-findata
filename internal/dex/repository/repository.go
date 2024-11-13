package repository

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jmoiron/sqlx"
	"github.com/nel349/bz-findata/internal/dex/repository/mysql"
)

// Exchange method implementations
type DexExchange interface {
	SaveSwap(ctx context.Context, tx *types.Transaction, version string) error
}

// This could contain multiple exchange repositories
type DexRepositories struct {
	DexExchange
}

func NewDexRepositories(db *sqlx.DB) *DexRepositories {
	return &DexRepositories{
		DexExchange: mysql.NewDexExchangeRepository(db),
	}
}
