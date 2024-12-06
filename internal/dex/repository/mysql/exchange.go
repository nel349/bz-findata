package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jmoiron/sqlx"
	"github.com/nel349/bz-findata/internal/dex/eth/defi_llama"
	"github.com/nel349/bz-findata/internal/dex/eth/uniswap/decoder"
	v2 "github.com/nel349/bz-findata/internal/dex/eth/uniswap/v2"
	v3 "github.com/nel349/bz-findata/internal/dex/eth/uniswap/v3"
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
	ctxReq, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	swapTransactions, err := decoder.DecodeSwap(tx, version)
	if err != nil {
		fmt.Println("Error decoding swap", "error", err)
		return err
	}

	// Process each transaction
	for _, swapTransaction := range swapTransactions {

		// should be in use case interface. but for now, here
		tokenInfoFrom, err := defi_llama.GetTokenMetadataFromDbOrDefiLlama(e.db, swapTransaction.TokenPathFrom, 15*time.Minute)
		if err != nil {
			fmt.Println("Error getting token metadata", "error", err)
			// fmt.Println("Token path from", swapTransaction.TokenPathFrom)
			return err
		}

		if version == "V2" {
			// if it is a token input
			if _, ok := v2.GetV2MethodFromID(swapTransaction.MethodID); ok {
				tokenAmount := decoder.ConvertToBigInt(swapTransaction.AmountIn)
				swapTransaction.Value = decoder.GetUsdValueFromToken(tokenAmount, tokenInfoFrom.Price, int(tokenInfoFrom.Decimals))
			}
		}

		if version == "V3" {
			// if it is a token input
			if _, ok := v3.GetV3MethodFromID(swapTransaction.MethodID); ok {
				tokenAmount := decoder.ConvertToBigInt(swapTransaction.AmountIn)
				swapTransaction.Value = decoder.GetUsdValueFromToken(tokenAmount, tokenInfoFrom.Price, int(tokenInfoFrom.Decimals))

				// Debug
				// fmt.Printf("DEBUG: Token amount: %s\n, price: %.9f\n, value: %.9f\n", swapTransaction.AmountIn, tokenInfoFrom.Price, swapTransaction.Value)
			}
		}

		// fmt.Printf("TokenInfoFrom decimals: %d\n, symbol: %s\n, price: %.9f\n", tokenInfoFrom.Decimals, tokenInfoFrom.Symbol, tokenInfoFrom.Price)

		if tx != nil {
			query := `
		INSERT INTO swap_transactions (
			value,
			tx_hash, 
			version, 
			exchange, 
			amount_in, 
			to_address, 
			token_path_from, 
			token_path_to,
			amount_token_desired,
			amount_token_min,
			amount_eth_min,
			method_id,
			method_name,
			liquidity
		) VALUES (
			:value,
			:tx_hash, 
			:version, 
			:exchange, 
			:amount_in, 
			:to_address, 
			:token_path_from, 
			:token_path_to,
			:amount_token_desired,
			:amount_token_min,
			:amount_eth_min,
			:method_id,
			:method_name,
			:liquidity
		)`
			_, err := e.db.NamedExecContext(
				ctxReq,
				query,
				entity.SwapTransaction{
					Value:              swapTransaction.Value,
					TxHash:             tx.Hash().Hex(),
					Version:            version,
					Exchange:           swapTransaction.Exchange,
					AmountIn:           swapTransaction.AmountIn,
					ToAddress:          swapTransaction.ToAddress,
					TokenPathFrom:      swapTransaction.TokenPathFrom,
					TokenPathTo:        swapTransaction.TokenPathTo,
					AmountTokenDesired: swapTransaction.AmountTokenDesired,
					AmountTokenMin:     swapTransaction.AmountTokenMin,
					AmountETHMin:       swapTransaction.AmountETHMin,
					MethodID:           swapTransaction.MethodID,
					MethodName:         swapTransaction.MethodName,
					Liquidity:          swapTransaction.Liquidity,
				},
			)
			if err != nil {
				fmt.Println("Error inserting swap", "error", err)
			}
			fmt.Printf("Swap inserted successfully: %s\n", swapTransaction.TxHash)

		}
	}

	return fmt.Errorf("transaction is nil")
}
