package decoder

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/nel349/bz-findata/pkg/entity"
)

// Test DecodeSwapExactTokensForTokens
func TestDecodeSwapV2(t *testing.T) {

	data := common.FromHex("0x38ed17390000000000000000000000000000000000000000000000000108d3a3aa9f11e0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000056eb903b0d2e858905feb7f1f4ad73458243d5a900000000000000000000000000000000000000000000000000000000673576e70000000000000000000000000000000000000000000000000000000000000002000000000000000000000000699ec925118567b6475fe495327ba0a778234aaa000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2")
	version := "v2"
	// tokenInfo := &defi_llama.TokenInfo{
	// 	Address:  "0x699Ec925118567b6475Fe495327ba0a778234AaA",
	// 	Decimals: 9,
	// 	Symbol:   "DUCKY",
	// }

	swapTransactionResult := &entity.SwapTransaction{}

	t.Run("Test DecodeSwapExactTokensForTokens", func(t *testing.T) {
		// https://etherscan.io/tx/0x7c7accc2ce3f94da7d8304c61668713829a66f41fca42b61e684a89b77ce3f25
		err := DecodeSwapExactTokensForTokens(data, version, swapTransactionResult)
		checkSwapNotNil(t, err, swapTransactionResult)

		expected := &entity.SwapTransaction{
			AmountIn:      ConvertToFloat64("74542093747294688"),
			TokenPathFrom: "0x699Ec925118567b6475Fe495327ba0a778234AaA",
			TokenPathTo:   "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		}

		fmt.Printf("Swap Transaction: %v\n", swapTransactionResult)

		// Compare with small epsilon to account for floating point precision
		if expected.AmountIn != swapTransactionResult.AmountIn {
			t.Errorf("Amount In does not match expected value %v, got %v", expected.AmountIn, swapTransactionResult.AmountIn)
		}

		if toLowerCaseHex(expected.TokenPathFrom) != toLowerCaseHex(swapTransactionResult.TokenPathFrom) {
			t.Errorf("Token Path From does not match expected value %v, got %v", expected.TokenPathFrom, swapTransactionResult.TokenPathFrom)
		}

		if toLowerCaseHex(expected.TokenPathTo) != toLowerCaseHex(swapTransactionResult.TokenPathTo) {
			t.Errorf("Token Path To does not match expected value %v, got %v", expected.TokenPathTo, swapTransactionResult.TokenPathTo)
		}
	})

	t.Run("Test DecodeSwapExactTokensForETHSupportingFeeOnTransferTokens", func(t *testing.T) {

		data = common.FromHex("0x791ac9470000000000000000000000000000000000000000000422ca8b0a00a424ffffff00000000000000000000000000000000000000000000000007c28167a0c8547400000000000000000000000000000000000000000000000000000000000000a00000000000000000000000001ad0eb3d4e0b79c20f8b3af24b706ae3c8e6a201000000000000000000000000000000000000000000000000000000006736f2f30000000000000000000000000000000000000000000000000000000000000002000000000000000000000000f3c7cecf8cbc3066f9a87b310cebe198d00479ac000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2")
		// tokenInfo := &defi_llama.TokenInfo{
		// 	Address:  "0xF3c7CECF8cBC3066F9a87b310cEBE198d00479aC",
		// 	Decimals: 18,
		// 	Symbol:   "FEG",
		// }

		err := DecodeSwapExactTokensForETHSupportingFeeOnTransferTokens(data, version, swapTransactionResult)

		checkSwapNotNil(t, err, swapTransactionResult)
		rawAmount, ok := new(big.Int).SetString("4999999999999999999999999", 10) // got expected amount from tenderly dev mode
		// https://dashboard.tenderly.co/tx/mainnet/0x708b5ce2f7a6e6bf95ed92206955afddfe226bbd1911ff31f1dc9604a25fd93d

		if !ok {
			t.Fatal("failed to parse big.Int")
		}

		// expectedAmountInFloat64, _ := new(big.Float).Quo(
		// 	new(big.Float).SetInt(rawAmount),
		// 	new(big.Float).SetFloat64(math.Pow(10, float64(tokenInfo.Decimals))),
		// ).Float64()

		expected := &entity.SwapTransaction{
			AmountIn:      ConvertToFloat64(rawAmount.String()),
			TokenPathFrom: "0xF3c7CECF8cBC3066F9a87b310cEBE198d00479aC",
			TokenPathTo:   "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		}
		if expected.AmountIn != swapTransactionResult.AmountIn {
			t.Errorf("Amount In does not match expected value %v, got %v", expected.AmountIn, swapTransactionResult.AmountIn)
		}

		if toLowerCaseHex(expected.TokenPathFrom) != toLowerCaseHex(swapTransactionResult.TokenPathFrom) {
			t.Errorf("Token Path From does not match expected value %v, got %v", expected.TokenPathFrom, swapTransactionResult.TokenPathFrom)
		}

		if toLowerCaseHex(expected.TokenPathTo) != toLowerCaseHex(swapTransactionResult.TokenPathTo) {
			t.Errorf("Token Path To does not match expected value %v, got %v", expected.TokenPathTo, swapTransactionResult.TokenPathTo)
		}
	})

	t.Run("Test DecodeSwapExactTokensForTokensSupportingFeeOnTransferTokens", func(t *testing.T) {

		// tx hash: 0xcf39e1501430f75f7ee041781b62592c6ba8a3749e5b5f3813f086023607dc1b
		// https://etherscan.io/tx/0xcf39e1501430f75f7ee041781b62592c6ba8a3749e5b5f3813f086023607dc1b
		// https://dashboard.tenderly.co/tx/mainnet/0xcf39e1501430f75f7ee041781b62592c6ba8a3749e5b5f3813f086023607dc1b
		data = common.FromHex("0x5c11d7950000000000000000000000000000000000000000000ec068614236ee611fe9470000000000000000000000000000000000000000000000000766b47bedbc6e9d00000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000c54a957d2e1da5067c7ad32d38d3a2bc2524531c000000000000000000000000000000000000000000000000000001932e9df0b80000000000000000000000000000000000000000000000000000000000000002000000000000000000000000debcad12e9c454a7338b3ec0c8058eec688c79d5000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2")
		// tokenInfo := &defi_llama.TokenInfo{
		// 	Address:  "0xdEbcaD12E9C454a7338B3EC0c8058EeC688c79d5",
		// 	Decimals: 18,
		// 	Symbol:   "$PnutKing",
		// }

		err := DecodeSwapExactTokensForTokensSupportingFeeOnTransferTokens(data, version, swapTransactionResult)
		checkSwapNotNil(t, err, swapTransactionResult)

		expected := &entity.SwapTransaction{
			AmountIn:      ConvertToFloat64("17833581308923813721794887"),
			TokenPathFrom: "0xdEbcaD12E9C454a7338B3EC0c8058EeC688c79d5",
			TokenPathTo:   "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		}

		if expected.AmountIn != swapTransactionResult.AmountIn {
			t.Errorf("Amount In does not match expected value %v, got %v", expected.AmountIn, swapTransactionResult.AmountIn)
		}

		if toLowerCaseHex(expected.TokenPathFrom) != toLowerCaseHex(swapTransactionResult.TokenPathFrom) {
			t.Errorf("Token Path From does not match expected value %v, got %v", expected.TokenPathFrom, swapTransactionResult.TokenPathFrom)
		}

		if toLowerCaseHex(expected.TokenPathTo) != toLowerCaseHex(swapTransactionResult.TokenPathTo) {
			t.Errorf("Token Path To does not match expected value %v, got %v", expected.TokenPathTo, swapTransactionResult.TokenPathTo)
		}
	})

	t.Run("Test DecodeRemoveLiquidityETHWithPermit", func(t *testing.T) {

		data = common.FromHex("0xded9382a000000000000000000000000fe34cbcaef94a06a8fc1adce86d486f49af242ba00000000000000000000000000000000000000000000119707a239721536ed2a000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001072d3be58b71b386724c624890547499fe39b8900000000000000000000000000000000000000000000000000000000673e95b70000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001bac42c5965bedf9500e171015c523074ac3162bf093ffd0627d51691a19abbae6386a09c71e9bbd14d91db91f9844c161f72753b1c84cf974b5f4fcbcae3d9418")

		// https://dashboard.tenderly.co/tx/mainnet/0xfee06fd682ac1a32b9ebdb31dd9441e3edf6ca05887d5ba70b7bbe522782a05c
		err := DecodeRemoveLiquidityETHWithPermit(data, version, swapTransactionResult)
		checkSwapNotNil(t, err, swapTransactionResult)	

		expected := &entity.SwapTransaction{
			TokenPathFrom: "0xfe34cbcaef94a06a8fc1adce86d486f49af242ba",
			Liquidity:     "83066238629180748524842",
		}

		if expected.TokenPathFrom != swapTransactionResult.TokenPathFrom {
			t.Errorf("Token Path From does not match expected value %v, got %v", expected.TokenPathFrom, swapTransactionResult.TokenPathFrom)
		}

		if expected.Liquidity != swapTransactionResult.Liquidity {
			t.Errorf("Liquidity does not match expected value %v, got %v", expected.Liquidity, swapTransactionResult.Liquidity)
		}
	})

}

// V3 tests
func TestDecodeSwapV3(t *testing.T) {
	version := "v3"
	// Lets do DecodeExactInputSingle test

	swapTransactionResult := &entity.SwapTransaction{}
	t.Run("Test DecodeExactInputSingle", func(t *testing.T) {

		data := common.FromHex("0x414bf389000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2000000000000000000000000ee2a03aa6dacf51c18679c516ad5283d8e7c26370000000000000000000000000000000000000000000000000000000000000bb8000000000000000000000000f5213a6a2f0890321712520b8048d9886c1a9900000000000000000000000000000000000000000000000000000000006736f0e40000000000000000000000000000000000000000000000000b9eafe9ee6f4000000000000000000000000000000000000000000000000000000019fe199f2e100000000000000000000000000000000000000000000000000000000000000000")

		err := DecodeExactInputSingle(data, version, swapTransactionResult)
		checkSwapNotNil(t, err, swapTransactionResult)

		expected := &entity.SwapTransaction{
			AmountIn:      ConvertToFloat64("837300000000000000"),
			TokenPathFrom: "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			TokenPathTo:   "0xee2a03aa6dacf51c18679c516ad5283d8e7c2637",
			ToAddress:     "0xf5213a6a2f0890321712520b8048d9886c1a9900",
		}

		if expected.AmountIn != swapTransactionResult.AmountIn {
			t.Errorf("Amount In does not match expected value %v, got %v", expected.AmountIn, swapTransactionResult.AmountIn)
		}

		if toLowerCaseHex(expected.TokenPathFrom) != toLowerCaseHex(swapTransactionResult.TokenPathFrom) {
			t.Errorf("Token Path From does not match expected value %v, got %v", expected.TokenPathFrom, swapTransactionResult.TokenPathFrom)
		}

		if toLowerCaseHex(expected.TokenPathTo) != toLowerCaseHex(swapTransactionResult.TokenPathTo) {
			t.Errorf("Token Path To does not match expected value %v, got %v", expected.TokenPathTo, swapTransactionResult.TokenPathTo)
		}

		if expected.ToAddress != swapTransactionResult.ToAddress {
			t.Errorf("To Address does not match expected value %v, got %v", expected.ToAddress, swapTransactionResult.ToAddress)
		}
	})

	t.Run("Test DecodeExactInput", func(t *testing.T) {

		data := common.FromHex("0xc04b8d59000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000b0ba33566bd35bcb80738810b2868dc1ddd1f0e900000000000000000000000000000000000000000000000000000000673c1fcf000000000000000000000000000000000000000000000000066eced5a631d580000000000000000000000000000000000000000000000b38ca1ce396b5800000000000000000000000000000000000000000000000000000000000000000002bc02aaa39b223fe8d0a0e5c4f27ead9083c756cc200271038e382f74dfb84608f3c1f10187f6bef5951de93000000000000000000000000000000000000000000")

		err := DecodeExactInput(data, version, swapTransactionResult)
		checkSwapNotNil(t, err, swapTransactionResult)

		expected := &entity.SwapTransaction{
			AmountIn:      ConvertToFloat64("463535228677379456"),
			TokenPathFrom: "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			TokenPathTo:   "0x38e382f74dfb84608f3c1f10187f6bef5951de93",
			Fee: "0x002710",
		}

		if expected.AmountIn != swapTransactionResult.AmountIn {
			t.Errorf("Amount In does not match expected value %v, got %v", expected.AmountIn, swapTransactionResult.AmountIn)
		}

		if toLowerCaseHex(expected.TokenPathFrom) != toLowerCaseHex(swapTransactionResult.TokenPathFrom) {
			t.Errorf("Token Path From does not match expected value %v, got %v", expected.TokenPathFrom, swapTransactionResult.TokenPathFrom)
		}

		if toLowerCaseHex(expected.TokenPathTo) != toLowerCaseHex(swapTransactionResult.TokenPathTo) {
			t.Errorf("Token Path To does not match expected value %v, got %v", expected.TokenPathTo, swapTransactionResult.TokenPathTo)
		}

		if expected.Fee != swapTransactionResult.Fee {
			t.Errorf("Fee does not match expected value %v, got %v", expected.Fee, swapTransactionResult.Fee)
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
