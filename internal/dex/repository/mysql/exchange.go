package mysql

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jmoiron/sqlx"
	"github.com/nel349/bz-findata/internal/dex/eth/defi_llama"
	"github.com/nel349/bz-findata/internal/dex/eth/uniswap/decoder"
	v2 "github.com/nel349/bz-findata/internal/dex/eth/uniswap/v2"
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
	ctxReq, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	swapTransactions, err := decoder.DecodeSwap(tx, version)
	if err != nil {
		fmt.Println("Error decoding swap", "error", err)
		return err
	}

	// Process each transaction
	for _, swapTransaction := range swapTransactions {

		// Get token metadata based on operation type
		var tokenInfoFrom, tokenInfoA, tokenInfoB entity.TokenInfo
		var err error

		if swapTransaction.MethodName == v2.AddLiquidity.String() || swapTransaction.MethodName == v2.RemoveLiquidity.String() {
			// For liquidity operations, get both token A and B metadata
			tokenInfoA, err = defi_llama.GetTokenMetadataFromDbOrDefiLlama(e.db, swapTransaction.TokenA, 15*time.Minute)
			if err != nil {
				fmt.Println("Error getting token A metadata", "error", err)
				return err
			}

			tokenInfoB, err = defi_llama.GetTokenMetadataFromDbOrDefiLlama(e.db, swapTransaction.TokenB, 15*time.Minute)
			if err != nil {
				fmt.Println("Error getting token B metadata", "error", err)
				return err
			}
		} else {
			// For swap operations, get token metadata as before
			tokenInfoFrom, err = defi_llama.GetTokenMetadataFromDbOrDefiLlama(e.db, swapTransaction.TokenPathFrom, 15*time.Minute)
			if err != nil {
				fmt.Println("Error getting token metadata", "error", err)
				return err
			}
		}

		// Calculate value based on operation type
		switch swapTransaction.MethodName {
		case v2.AddLiquidity.String():
			// Calculate combined value from both tokens
			amountADesired := decoder.ConvertToBigInt(swapTransaction.AmountADesired)
			amountBDesired := decoder.ConvertToBigInt(swapTransaction.AmountBDesired)

			valueA := decoder.GetUsdValueFromToken(amountADesired, tokenInfoA.Price, int(tokenInfoA.Decimals))
			valueB := decoder.GetUsdValueFromToken(amountBDesired, tokenInfoB.Price, int(tokenInfoB.Decimals))

			swapTransaction.Value = valueA + valueB

		case v2.RemoveLiquidity.String():
			amountAToken := decoder.ConvertToBigInt(swapTransaction.AmountAMin)
			amountBToken := decoder.ConvertToBigInt(swapTransaction.AmountBMin)

			valueA := decoder.GetUsdValueFromToken(amountAToken, tokenInfoA.Price, int(tokenInfoA.Decimals))
			valueB := decoder.GetUsdValueFromToken(amountBToken, tokenInfoB.Price, int(tokenInfoB.Decimals))

			swapTransaction.Value = valueA + valueB

		case v2.RemoveLiquidityETHWithPermitSupportingFeeOnTransferTokens.String():
			amountAToken := decoder.ConvertToBigInt(swapTransaction.AmountTokenMin)
			amountBToken := decoder.ConvertToBigInt(swapTransaction.AmountETHMin)

			valueA := decoder.GetUsdValueFromToken(amountAToken, tokenInfoA.Price, int(tokenInfoA.Decimals))
			valueB := decoder.GetUsdValueFromToken(amountBToken, tokenInfoB.Price, int(tokenInfoB.Decimals))

			swapTransaction.Value = valueA + valueB


		case v2.RemoveLiquidityETH.String(): // uses the liquidity to get the value of the token
			liquidity := decoder.ConvertToBigInt(swapTransaction.Liquidity)
			swapTransaction.Value = decoder.GetUsdValueFromToken(liquidity, tokenInfoFrom.Price, int(tokenInfoFrom.Decimals))

		default:
			if version == "V2" || version == "V3" {
				var tokenAmount *big.Int
				var isETHInputMethod bool

				// Check if method uses ETH as input
				if method, ok := v2.GetV2MethodFromID(swapTransaction.MethodID); ok {
					isETHInputMethod = method.IsETHInput()
				}

				if isETHInputMethod {
					// For ETH input methods, use transaction value
					tokenAmount = tx.Value()
					// Get WETH price for value calculation
					tokenInfoFrom, err = defi_llama.GetTokenMetadataFromDbOrDefiLlama(
						e.db,
						"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2", // WETH address
						15*time.Minute,
					)
				} else {
					// For token input methods, use decoded amount
					tokenAmount = decoder.ConvertToBigInt(swapTransaction.AmountIn)
					tokenInfoFrom, err = defi_llama.GetTokenMetadataFromDbOrDefiLlama(
						e.db,
						swapTransaction.TokenPathFrom,
						15*time.Minute,
					)
				}

				if err != nil {
					return err
				}

				swapTransaction.Value = decoder.GetUsdValueFromToken(
					tokenAmount,
					tokenInfoFrom.Price,
					int(tokenInfoFrom.Decimals),
				)
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
			amount_out_min,
			method_id,
			method_name,
			liquidity,
			token_a,
			token_b,
			amount_a_desired,
			amount_b_desired,
			amount_a_min,
			amount_b_min
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
			:amount_out_min,
			:method_id,
			:method_name,
			:liquidity,
			:token_a,
			:token_b,
			:amount_a_desired,
			:amount_b_desired,
			:amount_a_min,
			:amount_b_min
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
					TokenA:             swapTransaction.TokenA,
					TokenB:             swapTransaction.TokenB,
					AmountADesired:     swapTransaction.AmountADesired,
					AmountBDesired:     swapTransaction.AmountBDesired,
					AmountAMin:         swapTransaction.AmountAMin,
					AmountBMin:         swapTransaction.AmountBMin,
				},
			)
			if err != nil {
				fmt.Printf("error inserting swap: %v", err)
			}
			fmt.Printf("Swap inserted successfully: %s\n", swapTransaction.TxHash)

		}
	}

	return nil
}
