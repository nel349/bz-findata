package decoder

import (
	"fmt"
	"math"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/nel349/bz-findata/internal/dex/eth/defi_llama"
	"github.com/nel349/bz-findata/pkg/entity"
)

// Test DecodeSwapExactTokensForTokens
func TestDecodeSwapExactTokensForTokens(t *testing.T) {

	data := common.FromHex("0x38ed17390000000000000000000000000000000000000000000000000108d3a3aa9f11e0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000056eb903b0d2e858905feb7f1f4ad73458243d5a900000000000000000000000000000000000000000000000000000000673576e70000000000000000000000000000000000000000000000000000000000000002000000000000000000000000699ec925118567b6475fe495327ba0a778234aaa000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2")
	version := "v2"
	tokenInfo := &defi_llama.TokenInfo{
		Address:  "0x699Ec925118567b6475Fe495327ba0a778234AaA",
		Decimals: 9,
		Symbol:   "DUCKY",
	}

	t.Run("Test DecodeSwapExactTokensForTokens", func(t *testing.T) {
		swapTransaction, err := DecodeSwapExactTokensForTokens(data, version, tokenInfo)

		if err != nil {
			t.Errorf("Error decoding swap: %v", err)
			return
		}

		expected := &entity.SwapTransaction{
			AmountIn: 74542093.747294688,
			// ToAddress:     "0x699ec925118567b6475fe495327ba0a778234aaa",
			// TokenPathFrom: "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			// TokenPathTo:   "0x699ec925118567b6475fe495327ba0a778234aaa",
		}

		if swapTransaction == nil {
			t.Errorf("Swap Transaction is nil")
			return
		}
		fmt.Printf("Swap Transaction: %v\n", swapTransaction)

		// Compare with small epsilon to account for floating point precision
		if math.Abs(expected.AmountIn-swapTransaction.AmountIn) > 0.000001 {
			t.Errorf("Amount In does not match expected value %v, got %v", expected.AmountIn, swapTransaction.AmountIn)
		}

	})

}
