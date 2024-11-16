package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jmoiron/sqlx"
	"github.com/nel349/bz-findata/internal/dex/eth/uniswap/decoder"
	"github.com/nel349/bz-findata/pkg/entity"
)

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

	swapTransaction, err := decoder.DecodeSwap(tx.Data(), version, e.db)
	if err != nil {
		fmt.Println("Error decoding swap", "error", err)
		return err
	}

	if tx != nil {
		query := `
		INSERT INTO swap_transactions (
			tx_hash, 
			version, 
			exchange, 
			amount_in, 
			to_address, 
			token_path_from, 
			token_path_to
		) VALUES (
			:tx_hash, 
			:version, 
			:exchange, 
			:amount_in, 
			:to_address, 
			:token_path_from, 
			:token_path_to
		)`
		_, err := e.db.NamedExecContext(
			ctxReq,
			query,
			entity.SwapTransaction{
				TxHash:    tx.Hash().Hex(),
				Version:   version,
				Exchange:  swapTransaction.Exchange,
				AmountIn:  swapTransaction.AmountIn,
				ToAddress: swapTransaction.ToAddress,
				// Deadline:      swapTransaction.Deadline,
				TokenPathFrom: swapTransaction.TokenPathFrom,
				TokenPathTo:   swapTransaction.TokenPathTo,
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
