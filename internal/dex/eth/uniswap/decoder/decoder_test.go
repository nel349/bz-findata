package decoder

import (
	"fmt"
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
			AmountIn:      "74542093747294688",
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
		// https://dashboard.tenderly.co/tx/mainnet/0x708b5ce2f7a6e6bf95ed92206955afddfe226bbd1911ff31f1dc9604a25fd93d

		// expectedAmountInFloat64, _ := new(big.Float).Quo(
		// 	new(big.Float).SetInt(rawAmount),
		// 	new(big.Float).SetFloat64(math.Pow(10, float64(tokenInfo.Decimals))),
		// ).Float64()

		expected := &entity.SwapTransaction{
			AmountIn:      "4999999999999999999999999",
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

		err := DecodeSwapExactTokensForTokensSupportingFeeOnTransferTokens(data, version, swapTransactionResult)
		checkSwapNotNil(t, err, swapTransactionResult)

		expected := &entity.SwapTransaction{
			AmountIn:      "17833581308923813721794887",
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

	t.Run("Test DecodeSwapExactETHForTokensSupportingFeeOnTransferTokens", func(t *testing.T) {

		data = common.FromHex("0xb6f9de950000000000000000000000000000000000000000000000000004f2cf373add62000000000000000000000000000000000000000000000000000000000000008000000000000000000000000076db926b75e225af64b954c95fef653926ea7965000000000000000000000000000000000000000000000000000001938f5069630000000000000000000000000000000000000000000000000000000000000002000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2000000000000000000000000790336af90933aa7bd10d4534db6909507098440")

		err := DecodeSwapExactETHForTokensSupportingFeeOnTransferTokens(data, version, swapTransactionResult)
		checkSwapNotNil(t, err, swapTransactionResult)

		expected := &entity.SwapTransaction{
			AmountOutMin:  "1392871705599330",
			TokenPathFrom: "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			TokenPathTo:   "0x790336af90933aa7bd10d4534db6909507098440",
		}

		if expected.AmountOutMin != swapTransactionResult.AmountOutMin {
			t.Errorf("Amount Out Min does not match expected value %v, got %v", expected.AmountOutMin, swapTransactionResult.AmountOutMin)
		}

		if toLowerCaseHex(expected.TokenPathFrom) != toLowerCaseHex(swapTransactionResult.TokenPathFrom) {
			t.Errorf("Token Path From does not match expected value %v, got %v", expected.TokenPathFrom, swapTransactionResult.TokenPathFrom)
		}

		if toLowerCaseHex(expected.TokenPathTo) != toLowerCaseHex(swapTransactionResult.TokenPathTo) {
			t.Errorf("Token Path To does not match expected value %v, got %v", expected.TokenPathTo, swapTransactionResult.TokenPathTo)
		}
	})

	t.Run("Test DecodeSwapTokensForExactTokens", func(t *testing.T) {

		data = common.FromHex("0x8803dbee0000000000000000000000000000000000000000000009b6c2b9818b505e200000000000000000000000000000000000000000000000000006e9f733c51ed70300000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000f1da51404a3c42dc46fcc6924944fab21fd7e9b900000000000000000000000000000000000000000000000000000000674fb9610000000000000000000000000000000000000000000000000000000000000002000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2000000000000000000000000425087bf4969f45818c225ae30f8560ce518582e")
		err := DecodeSwapTokensForExactTokens(data, version, swapTransactionResult)
		checkSwapNotNil(t, err, swapTransactionResult)

		expected := &entity.SwapTransaction{
			AmountOut:     "45872637155791343591424",
			AmountInMax:   "498201035523675907",
			TokenPathFrom: "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			TokenPathTo:   "0x425087bf4969f45818c225ae30f8560ce518582e",
		}

		if expected.AmountOut != swapTransactionResult.AmountOut {
			t.Errorf("Amount Out does not match expected value %v, got %v", expected.AmountOut, swapTransactionResult.AmountOut)
		}

		if expected.AmountInMax != swapTransactionResult.AmountInMax {
			t.Errorf("Amount In Max does not match expected value %v, got %v", expected.AmountInMax, swapTransactionResult.AmountInMax)
		}

		if toLowerCaseHex(expected.TokenPathFrom) != toLowerCaseHex(swapTransactionResult.TokenPathFrom) {
			t.Errorf("Token Path From does not match expected value %v, got %v", expected.TokenPathFrom, swapTransactionResult.TokenPathFrom)
		}

		if toLowerCaseHex(expected.TokenPathTo) != toLowerCaseHex(swapTransactionResult.TokenPathTo) {
			t.Errorf("Token Path To does not match expected value %v, got %v", expected.TokenPathTo, swapTransactionResult.TokenPathTo)
		}
	})

	t.Run("Test DecodeSwapExactTokensForETH", func(t *testing.T) {

		data = common.FromHex("0x18cbafe5000000000000000000000000000000000000000002386e78249bcbd50c1c0000000000000000000000000000000000000000000000000000011b475f5f08b23a00000000000000000000000000000000000000000000000000000000000000a00000000000000000000000008c5bd953f3e2756d5270abe14ed8d9c7f574b42a00000000000000000000000000000000000000000000000000000000ce9f5c2a0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000b0ac2b5a73da0e67a8e5489ba922b3f8d582e058000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2")

		err := DecodeSwapExactTokensForETH(data, swapTransactionResult)
		checkSwapNotNil(t, err, swapTransactionResult)

		expected := &entity.SwapTransaction{
			AmountIn:      "687191542101440000000000000",
			AmountOutMin:  "79735893350986298",
			TokenPathFrom: "0xb0ac2b5a73da0e67a8e5489ba922b3f8d582e058",
			TokenPathTo:   "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		}

		if expected.AmountIn != swapTransactionResult.AmountIn {
			t.Errorf("Amount In does not match expected value %v, got %v", expected.AmountIn, swapTransactionResult.AmountIn)
		}

		if expected.AmountOutMin != swapTransactionResult.AmountOutMin {
			t.Errorf("Amount Out Min does not match expected value %v, got %v", expected.AmountOutMin, swapTransactionResult.AmountOutMin)
		}

		if toLowerCaseHex(expected.TokenPathFrom) != toLowerCaseHex(swapTransactionResult.TokenPathFrom) {
			t.Errorf("Token Path From does not match expected value %v, got %v", expected.TokenPathFrom, swapTransactionResult.TokenPathFrom)
		}

		if toLowerCaseHex(expected.TokenPathTo) != toLowerCaseHex(swapTransactionResult.TokenPathTo) {
			t.Errorf("Token Path To does not match expected value %v, got %v", expected.TokenPathTo, swapTransactionResult.TokenPathTo)
		}
	})

	t.Run("Test DecodeSwapETHForExactTokens", func(t *testing.T) {

		data = common.FromHex("0xfb3bdb41000000000000000000000000000000000000000000000000000003aa29f07d58000000000000000000000000000000000000000000000000000000000000008000000000000000000000000011a8286db1eb78e4d0d3409be5ddafac2966cc0b000000000000000000000000000000000000000000000000000000006750f4800000000000000000000000000000000000000000000000000000000000000002000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2000000000000000000000000e0805c80588913c1c2c89ea4a8dcf485d4038a3e")

		err := DecodeSwapETHForExactTokens(data, swapTransactionResult)
		checkSwapNotNil(t, err, swapTransactionResult)

		expected := &entity.SwapTransaction{
			AmountOut:     "4029382950232",
			TokenPathFrom: "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			TokenPathTo:   "0xe0805c80588913c1c2c89ea4a8dcf485d4038a3e",
		}

		if expected.AmountOut != swapTransactionResult.AmountOut {
			t.Errorf("Amount Out does not match expected value %v, got %v", expected.AmountOut, swapTransactionResult.AmountOut)
		}

		if toLowerCaseHex(expected.TokenPathFrom) != toLowerCaseHex(swapTransactionResult.TokenPathFrom) {
			t.Errorf("Token Path From does not match expected value %v, got %v", expected.TokenPathFrom, swapTransactionResult.TokenPathFrom)
		}

		if toLowerCaseHex(expected.TokenPathTo) != toLowerCaseHex(swapTransactionResult.TokenPathTo) {
			t.Errorf("Token Path To does not match expected value %v, got %v", expected.TokenPathTo, swapTransactionResult.TokenPathTo)
		}
	})

	t.Run("Test DecodeRemoveLiquidity", func(t *testing.T) {
		data = common.FromHex("0xbaa2abde00000000000000000000000099ec69f6624abd625782e2127f7ca23432aab7a1000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000000000000000000000000000000048db1aa0d89d0000000000000000000000000000000000000000000000000000084dca2559ac80000000000000000000000000000000000000000000000000027ce754fedea24000000000000000000000000ba8ed95797f9bf37c99f564d9aa26eeb1851bf1f00000000000000000000000000000000000000000000000000000000675164f2")

		err := DecodeRemoveLiquidity(data, swapTransactionResult)
		checkSwapNotNil(t, err, swapTransactionResult)

		expected := &entity.SwapTransaction{
			TokenA:         "0x99ec69f6624abd625782e2127f7ca23432aab7a1",
			TokenB:         "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			Liquidity:     "1281694108584400",
			AmountAMin:    "146083151190728",
			AmountBMin:    "11204527339203108",
		}

		if expected.TokenA != swapTransactionResult.TokenA {
			t.Errorf("Token Path From does not match expected value %v, got %v", expected.TokenA, swapTransactionResult.TokenA)
		}

		if expected.TokenB != swapTransactionResult.TokenB {
			t.Errorf("Token Path To does not match expected value %v, got %v", expected.TokenB, swapTransactionResult.TokenB)
		}

		if expected.Liquidity != swapTransactionResult.Liquidity {
			t.Errorf("Liquidity does not match expected value %v, got %v", expected.Liquidity, swapTransactionResult.Liquidity)
		}

		if expected.AmountAMin != swapTransactionResult.AmountAMin {
			t.Errorf("Amount A Min does not match expected value %v, got %v", expected.AmountAMin, swapTransactionResult.AmountAMin)
		}

		if expected.AmountBMin != swapTransactionResult.AmountBMin {
			t.Errorf("Amount B Min does not match expected value %v, got %v", expected.AmountBMin, swapTransactionResult.AmountBMin)
		}
	})

	t.Run("Test DecodeSwapExactETHForTokens", func(t *testing.T) {
		data = common.FromHex("0x7ff36ab500000000000000000000000000000000000000000000000589a74db4fba713f10000000000000000000000000000000000000000000000000000000000000080000000000000000000000000659402c9c3eb08503ab3ecbaff5f384f6f7d62af00000000000000000000000000000000000000000000000000000000675163430000000000000000000000000000000000000000000000000000000000000002000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc20000000000000000000000006019dcb2d0b3e0d1d8b0ce8d16191b3a4f93703d")

		err := DecodeSwapExactETHForTokens(data, version, swapTransactionResult)
		checkSwapNotNil(t, err, swapTransactionResult)

		expected := &entity.SwapTransaction{
			AmountOutMin:  "102152702512566047729",
			TokenPathFrom: "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			TokenPathTo:   "0x6019dcb2d0b3e0d1d8b0ce8d16191b3a4f93703d",
		}

		if expected.AmountOutMin != swapTransactionResult.AmountOutMin {
			t.Errorf("Amount Out Min does not match expected value %v, got %v", expected.AmountOutMin, swapTransactionResult.AmountOutMin)
		}

		if toLowerCaseHex(expected.TokenPathFrom) != toLowerCaseHex(swapTransactionResult.TokenPathFrom) {
			t.Errorf("Token Path From does not match expected value %v, got %v", expected.TokenPathFrom, swapTransactionResult.TokenPathFrom)
		}

		if toLowerCaseHex(expected.TokenPathTo) != toLowerCaseHex(swapTransactionResult.TokenPathTo) {
			t.Errorf("Token Path To does not match expected value %v, got %v", expected.TokenPathTo, swapTransactionResult.TokenPathTo)
		}
	})

	t.Run("Test DecodeSwapTokensForExactETH", func(t *testing.T) {
		data = common.FromHex("0x4a25d94a0000000000000000000000000000000000000000000000000473c8321cbe3b23ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000f604ec052a32eb53aaa3b104993bc4d5b6132a52000000000000000000000000000000000000000000000000000000006753615900000000000000000000000000000000000000000000000000000000000000020000000000000000000000006a4402a535d74bd0c9cdb5ce2d51822fc9f6620e000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2")

		err := DecodeSwapTokensForExactETH(data, swapTransactionResult)
		checkSwapNotNil(t, err, swapTransactionResult)

		expected := &entity.SwapTransaction{
			AmountOut:     "320820116029586211",
			AmountInMax:   "115792089237316195423570985008687907853269984665640564039457584007913129639935",
			TokenPathFrom: "0x6a4402a535d74bd0c9cdb5ce2d51822fc9f6620e",
			TokenPathTo:   "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		}

		if expected.AmountOut != swapTransactionResult.AmountOut {
			t.Errorf("Amount Out does not match expected value %v, got %v", expected.AmountOut, swapTransactionResult.AmountOut)
		}

		if expected.AmountInMax != swapTransactionResult.AmountInMax {
			t.Errorf("Amount In Max does not match expected value %v, got %v", expected.AmountInMax, swapTransactionResult.AmountInMax)
		}

		if toLowerCaseHex(expected.TokenPathFrom) != toLowerCaseHex(swapTransactionResult.TokenPathFrom) {
			t.Errorf("Token Path From does not match expected value %v, got %v", expected.TokenPathFrom, swapTransactionResult.TokenPathFrom)
		}

		if toLowerCaseHex(expected.TokenPathTo) != toLowerCaseHex(swapTransactionResult.TokenPathTo) {
			t.Errorf("Token Path To does not match expected value %v, got %v", expected.TokenPathTo, swapTransactionResult.TokenPathTo)
		}
	})

	// test DecodeAddLiquidityETH
	t.Run("Test DecodeAddLiquidityETH", func(t *testing.T) {
		//https://dashboard.tenderly.co/tx/mainnet/0xfff632109695d0e5ba72799b447b185ac18a994371259106ba546ca187d33479
		data = common.FromHex("0xf305d719000000000000000000000000728f30fa2f100742c7949d1961804fa8e0b1387d000000000000000000000000000000000000000000000389cbf4dbec81f400000000000000000000000000000000000000000000000003854489651071f1800000000000000000000000000000000000000000000000000006e25df1763291ff00000000000000000000000093f9a1ee83522bd206cccc875ccb1026b69dbaa0000000000000000000000000000000000000000000000000000000006753825c")

		err := DecodeAddLiquidityETH(data, version, swapTransactionResult)
		checkSwapNotNil(t, err, swapTransactionResult)

		expected := &entity.SwapTransaction{
			TokenPathTo:        "0x728f30fa2f100742c7949d1961804fa8e0b1387d",
			AmountTokenDesired: "16709000000000000000000",
			AmountTokenMin:     "16625455000000000000000",
			AmountETHMin:       "496062200615703039",
		}

		if expected.TokenPathTo != swapTransactionResult.TokenPathTo {
			t.Errorf("Token Path To does not match expected value %v, got %v", expected.TokenPathTo, swapTransactionResult.TokenPathTo)
		}

		if expected.AmountTokenDesired != swapTransactionResult.AmountTokenDesired {
			t.Errorf("Amount Token Desired does not match expected value %v, got %v", expected.AmountTokenDesired, swapTransactionResult.AmountTokenDesired)
		}

		if expected.AmountTokenMin != swapTransactionResult.AmountTokenMin {
			t.Errorf("Amount Token Min does not match expected value %v, got %v", expected.AmountTokenMin, swapTransactionResult.AmountTokenMin)
		}

		if expected.AmountETHMin != swapTransactionResult.AmountETHMin {
			t.Errorf("Amount ETH Min does not match expected value %v, got %v", expected.AmountETHMin, swapTransactionResult.AmountETHMin)
		}

	})

	t.Run("Test DecodeAddLiquidity", func(t *testing.T) {

		data := common.FromHex("0xe8e33700000000000000000000000000910812c44ed2a3b611e4b051d9d83a88d652e2dd000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2000000000000000000000000000000000000000000003f5b6a1d988446cc08550000000000000000000000000000000000000000000000000a11bc5997c59c50000000000000000000000000000000000000000000003f0a5143d908bc33f8ee0000000000000000000000000000000000000000000000000a04d8d92517d2920000000000000000000000007349bd4f327e697814ff96ac48861b44562ed7220000000000000000000000000000000000000000000000000000000067538242")

		err := DecodeAddLiquidity(data, version, swapTransactionResult)
		checkSwapNotNil(t, err, swapTransactionResult)

		expected := &entity.SwapTransaction{
			TokenA:         "0x910812c44ed2a3b611e4b051d9d83a88d652e2dd",
			TokenB:         "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			AmountADesired: "299195388566931453511765",
			AmountBDesired: "725568107967781968",
			AmountAMin:     "297699411624096796244206",
			AmountBMin:     "721940267427943058",
		}

		if expected.TokenA != swapTransactionResult.TokenA {
			t.Errorf("Token A does not match expected value %v, got %v", expected.TokenA, swapTransactionResult.TokenA)
		}

		if expected.TokenB != swapTransactionResult.TokenB {
			t.Errorf("Token B does not match expected value %v, got %v", expected.TokenB, swapTransactionResult.TokenB)
		}

		if expected.AmountADesired != swapTransactionResult.AmountADesired {
			t.Errorf("Amount A Desired does not match expected value %v, got %v", expected.AmountADesired, swapTransactionResult.AmountADesired)
		}

		if expected.AmountBDesired != swapTransactionResult.AmountBDesired {
			t.Errorf("Amount B Desired does not match expected value %v, got %v", expected.AmountBDesired, swapTransactionResult.AmountBDesired)
		}

		if expected.AmountAMin != swapTransactionResult.AmountAMin {
			t.Errorf("Amount A Min does not match expected value %v, got %v", expected.AmountAMin, swapTransactionResult.AmountAMin)
		}

		if expected.AmountBMin != swapTransactionResult.AmountBMin {
			t.Errorf("Amount B Min does not match expected value %v, got %v", expected.AmountBMin, swapTransactionResult.AmountBMin)
		}
	})


	/*
		Function: addLiquidityETH(address token, uint256 amountTokenDesired, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline)

		MethodID: 0xf305d719
		[0]:  00000000000000000000000009e4a242cdc18a54c34015f96a3c07058f8e1967
		[1]:  0000000000000000000000000000000000000000033b2e3c9fd0803ce8000000
		[2]:  000000000000000000000000000000000000000003370b7214c67f98c3000000
		[3]:  0000000000000000000000000000000000000000000000008a1580485b230000
		[4]:  00000000000000000000000013b2ef18dc68f4fbbcd515b0bfe59e2f594c91df
		[5]:  0000000000000000000000000000000000000000000000000000000067591285

		{
			"token": "0x09e4a242cdc18a54c34015f96a3c07058f8e1967",
			"amountTokenDesired": "1000000000000000000000000000",
			"amountTokenMin": "995000000000000000000000000",
			"amountETHMin": "9950000000000000000",
			"to": "0x13b2ef18dc68f4fbbcd515b0bfe59e2f594c91df",
			"deadline": "1733890693"
		}
	*/
	t.Run("Test DecodeaAddLiquidityETH", func(t *testing.T) {

		data := common.FromHex("0xf305d71900000000000000000000000009e4a242cdc18a54c34015f96a3c07058f8e19670000000000000000000000000000000000000000033b2e3c9fd0803ce8000000000000000000000000000000000000000000000003370b7214c67f98c30000000000000000000000000000000000000000000000000000008a1580485b23000000000000000000000000000013b2ef18dc68f4fbbcd515b0bfe59e2f594c91df0000000000000000000000000000000000000000000000000000000067591285")

		err := DecodeAddLiquidityETH(data, version, swapTransactionResult)
		checkSwapNotNil(t, err, swapTransactionResult)

		expected := &entity.SwapTransaction{
			TokenPathTo:        "0x09e4a242cdc18a54c34015f96a3c07058f8e1967",
			AmountTokenDesired: "1000000000000000000000000000",
			AmountTokenMin: "995000000000000000000000000",
			AmountETHMin: "9950000000000000000",
		}

		if expected.TokenPathTo != swapTransactionResult.TokenPathTo {
			t.Errorf("Token Path To does not match expected value %v, got %v", expected.TokenPathTo, swapTransactionResult.TokenPathTo)
		}

		if expected.AmountTokenDesired != swapTransactionResult.AmountTokenDesired {
			t.Errorf("Amount Token Desired does not match expected value %v, got %v", expected.AmountTokenDesired, swapTransactionResult.AmountTokenDesired)
		}

		if expected.AmountTokenMin != swapTransactionResult.AmountTokenMin {
			t.Errorf("Amount Token Min does not match expected value %v, got %v", expected.AmountTokenMin, swapTransactionResult.AmountTokenMin)
		}

		if expected.AmountETHMin != swapTransactionResult.AmountETHMin {
			t.Errorf("Amount ETH Min does not match expected value %v, got %v", expected.AmountETHMin, swapTransactionResult.AmountETHMin)
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
