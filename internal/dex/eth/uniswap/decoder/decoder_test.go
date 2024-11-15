package decoder

import (
	"fmt"
	"math"
	"math/big"
	"strings"
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
		checkSwapNotNil(t, err, swapTransaction)
		expected := &entity.SwapTransaction{
			AmountIn: 74542093.747294688,
			TokenPathFrom: "0x699Ec925118567b6475Fe495327ba0a778234AaA",
			TokenPathTo:   "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		}

		fmt.Printf("Swap Transaction: %v\n", swapTransaction)

		// Compare with small epsilon to account for floating point precision
		if math.Abs(expected.AmountIn-swapTransaction.AmountIn) > 0.000001 {
			t.Errorf("Amount In does not match expected value %v, got %v", expected.AmountIn, swapTransaction.AmountIn)
		}

		if toLowerCaseHex(expected.TokenPathFrom) != toLowerCaseHex(swapTransaction.TokenPathFrom) {
			t.Errorf("Token Path From does not match expected value %v, got %v", expected.TokenPathFrom, swapTransaction.TokenPathFrom)
		}

		if toLowerCaseHex(expected.TokenPathTo) != toLowerCaseHex(swapTransaction.TokenPathTo) {
			t.Errorf("Token Path To does not match expected value %v, got %v", expected.TokenPathTo, swapTransaction.TokenPathTo)
		}
	})

	t.Run("Test DecodeSwapExactTokensForETHSupportingFeeOnTransferTokens", func(t *testing.T) {

		data = common.FromHex("0x791ac9470000000000000000000000000000000000000000000422ca8b0a00a424ffffff00000000000000000000000000000000000000000000000007c28167a0c8547400000000000000000000000000000000000000000000000000000000000000a00000000000000000000000001ad0eb3d4e0b79c20f8b3af24b706ae3c8e6a201000000000000000000000000000000000000000000000000000000006736f2f30000000000000000000000000000000000000000000000000000000000000002000000000000000000000000f3c7cecf8cbc3066f9a87b310cebe198d00479ac000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2")
		tokenInfo := &defi_llama.TokenInfo{
			Address:  "0xF3c7CECF8cBC3066F9a87b310cEBE198d00479aC",
			Decimals: 18,
			Symbol:   "FEG",
		}

		swapTransaction, err := DecodeSwapExactTokensForETHSupportingFeeOnTransferTokens(data, version, tokenInfo)

		checkSwapNotNil(t, err, swapTransaction)
		rawAmount, ok := new(big.Int).SetString("4999999999999999999999999", 10) // got expected amount from tenderly dev mode
		// https://dashboard.tenderly.co/tx/mainnet/0x708b5ce2f7a6e6bf95ed92206955afddfe226bbd1911ff31f1dc9604a25fd93d

		if !ok {
			t.Fatal("failed to parse big.Int")
		}

		expectedAmountInFloat64, _ := new(big.Float).Quo(
			new(big.Float).SetInt(rawAmount),
			new(big.Float).SetFloat64(math.Pow(10, float64(tokenInfo.Decimals))),
		).Float64()

		expected := &entity.SwapTransaction{
			AmountIn: expectedAmountInFloat64,
			TokenPathFrom: "0xF3c7CECF8cBC3066F9a87b310cEBE198d00479aC",
			TokenPathTo:   "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		}
		if math.Abs(expected.AmountIn-swapTransaction.AmountIn) > 0.000001 {
			t.Errorf("Amount In does not match expected value %v, got %v", expected.AmountIn, swapTransaction.AmountIn)
		}

		if toLowerCaseHex(expected.TokenPathFrom) != toLowerCaseHex(swapTransaction.TokenPathFrom) {
			t.Errorf("Token Path From does not match expected value %v, got %v", expected.TokenPathFrom, swapTransaction.TokenPathFrom)
		}

		if toLowerCaseHex(expected.TokenPathTo) != toLowerCaseHex(swapTransaction.TokenPathTo) {
			t.Errorf("Token Path To does not match expected value %v, got %v", expected.TokenPathTo, swapTransaction.TokenPathTo)
		}
	})

}

func checkSwapNotNil(t *testing.T, err error, swapTransaction *entity.SwapTransaction) {
	if err != nil {
		t.Errorf("Error decoding swap: %v", err)
		return
	}
	if swapTransaction == nil {
		t.Errorf("Swap Transaction is nil")
	}
}

// toLowerCaseHex converts a hex string to lowercase
func toLowerCaseHex(hex string) string {
	return strings.ToLower(hex)
}
