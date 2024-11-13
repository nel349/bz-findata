package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jmoiron/sqlx"
)

type SwapTransaction struct {
	TxHash string `db:"tx_hash"`
	Version string `db:"version"`
}

type dexExchangeRepo struct {
	db *sqlx.DB
}

// NewDexExchangeRepository created exchange repository
func NewDexExchangeRepository(db *sqlx.DB) *dexExchangeRepo {
	return &dexExchangeRepo{db}
}

func (e *dexExchangeRepo) SaveSwap(ctx context.Context, tx *types.Transaction, version string) error {
	ctxReq, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if tx != nil {
		_, err := e.db.NamedExecContext(
			ctxReq,
			`INSERT INTO swap_transactions (tx_hash, version) VALUES (:tx_hash, :version)`,
			SwapTransaction{
				TxHash: tx.Hash().Hex(),
				Version: version,
			},
		)	
		if err != nil {
			fmt.Println("Error inserting swap", "error", err)
			return err
		}
		fmt.Println("Swap inserted successfully")
		return nil
	}

	return fmt.Errorf("transaction is nil")
}
