package decoder

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/nel349/bz-findata/pkg/entity"
)

// V3 tests
func TestDecodeSwapV3(t *testing.T) {
	// Lets do DecodeExactInputSingle test

	swapTransactionResult := &entity.SwapTransaction{}
	t.Run("Test DecodeExactInputSingle", func(t *testing.T) {

		data := common.FromHex("0x414bf389000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2000000000000000000000000ee2a03aa6dacf51c18679c516ad5283d8e7c26370000000000000000000000000000000000000000000000000000000000000bb8000000000000000000000000f5213a6a2f0890321712520b8048d9886c1a9900000000000000000000000000000000000000000000000000000000006736f0e40000000000000000000000000000000000000000000000000b9eafe9ee6f4000000000000000000000000000000000000000000000000000000019fe199f2e100000000000000000000000000000000000000000000000000000000000000000")

		err := DecodeExactInputSingle(data, swapTransactionResult)
		checkSwapNotNil(t, err, swapTransactionResult)

		expected := &entity.SwapTransaction{
			AmountIn:      "837300000000000000",
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

		err := DecodeExactInput(data, swapTransactionResult)
		checkSwapNotNil(t, err, swapTransactionResult)

		expected := &entity.SwapTransaction{
			AmountIn:      "463535228677379456",
			TokenPathFrom: "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			TokenPathTo:   "0x38e382f74dfb84608f3c1f10187f6bef5951de93",
			Fee:           "0x002710",
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

	t.Run("Test DecodeExactOutputSingle", func(t *testing.T) {

		data := common.FromHex("0xdb3e2198000000000000000000000000d31a59c85ae9d8edefec411d448f90841571b89c000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc20000000000000000000000000000000000000000000000000000000000000bb8000000000000000000000000f5213a6a2f0890321712520b8048d9886c1a990000000000000000000000000000000000000000000000000000000000674fad840000000000000000000000000000000000000000000000012b824fb44600c0000000000000000000000000000000000000000000000000000000004e3f884eef0000000000000000000000000000000000000000000000000000000000000000")

		err := DecodeExactOutputSingle(data, swapTransactionResult)
		checkSwapNotNil(t, err, swapTransactionResult)

		expected := &entity.SwapTransaction{
			AmountInMax:   "336073346799",
			AmountIn:      "336073346799",
			TokenPathFrom: "0xd31a59c85ae9d8edefec411d448f90841571b89c",
			TokenPathTo:   "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		}

		if expected.AmountInMax != swapTransactionResult.AmountInMax {
			t.Errorf("Amount In Max does not match expected value %v, got %v", expected.AmountInMax, swapTransactionResult.AmountInMax)
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

	t.Run("Test DecodeMulticallArray", func(t *testing.T) {

		data := common.FromHex("0xac9650d800000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000001800000000000000000000000000000000000000000000000000000000000000104db3e2198000000000000000000000000960692640ac4986ffce41620b7e3aa03cf1a0e8f000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc20000000000000000000000000000000000000000000000000000000000000bb8000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000674fa9ef0000000000000000000000000000000000000000000000000214e8348c4f00000000000000000000000000000000000000000000000002138669df859f1c439f000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004449404b7c0000000000000000000000000000000000000000000000000214e8348c4f00000000000000000000000000000f648e13fbe57251d367b08f42e8557d7461637a00000000000000000000000000000000000000000000000000000000")

		err := DecodeDataArray(data, swapTransactionResult)
		if err != nil {
			t.Errorf("Error decoding multicall array: %v", err)
		}

		expected := &entity.SwapTransaction{
			NumberOfCalls: 2,
			CallsData: []string{
				"0xdb3e2198000000000000000000000000960692640ac4986ffce41620b7e3aa03cf1a0e8f000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc20000000000000000000000000000000000000000000000000000000000000bb8000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000674fa9ef0000000000000000000000000000000000000000000000000214e8348c4f00000000000000000000000000000000000000000000000002138669df859f1c439f0000000000000000000000000000000000000000000000000000000000000000",
				"0x49404b7c0000000000000000000000000000000000000000000000000214e8348c4f00000000000000000000000000000f648e13fbe57251d367b08f42e8557d7461637a",
			},
		}

		if expected.NumberOfCalls != swapTransactionResult.NumberOfCalls {
			t.Errorf("Number Of Calls does not match expected value %v, got %v", expected.NumberOfCalls, swapTransactionResult.NumberOfCalls)
		}

		if expected.CallsData[0] != swapTransactionResult.CallsData[0] {
			t.Errorf("Calls Data does not match expected value %v, got %v", expected.CallsData[0], swapTransactionResult.CallsData[0])
		}

		if expected.CallsData[1] != swapTransactionResult.CallsData[1] {
			t.Errorf("Calls Data does not match expected value %v, got %v", expected.CallsData[1], swapTransactionResult.CallsData[1])
		}

	})

	t.Run("Test DecodeMulticall", func(t *testing.T) {

		data := common.FromHex("0xac9650d800000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000001800000000000000000000000000000000000000000000000000000000000000104db3e2198000000000000000000000000960692640ac4986ffce41620b7e3aa03cf1a0e8f000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc20000000000000000000000000000000000000000000000000000000000000bb8000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000674fa9ef0000000000000000000000000000000000000000000000000214e8348c4f00000000000000000000000000000000000000000000000002138669df859f1c439f000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004449404b7c0000000000000000000000000000000000000000000000000214e8348c4f00000000000000000000000000000f648e13fbe57251d367b08f42e8557d7461637a00000000000000000000000000000000000000000000000000000000")

		swapTransactions, err := DecodeMulticall(data)
		if err != nil {
			t.Fatalf("DecodeMulticall failed: %v", err)
		}
		if len(swapTransactions) == 0 {
			t.Fatal("No swap transactions decoded")
		}

		// two calls in the multicall

		// first call is ExactOutputSingle
		expectedTransaction1 := &entity.SwapTransaction{
			AmountIn:      "9804906621378401944479",
			TokenPathFrom: "0x960692640ac4986ffce41620b7e3aa03cf1a0e8f",
			TokenPathTo:   "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		}

		if expectedTransaction1.AmountIn != swapTransactions[0].AmountIn {
			t.Errorf("Amount In does not match expected value %v, got %v", expectedTransaction1.AmountIn, swapTransactions[0].AmountIn)
		}

		if toLowerCaseHex(expectedTransaction1.TokenPathFrom) != toLowerCaseHex(swapTransactions[0].TokenPathFrom) {
			t.Errorf("Token Path From does not match expected value %v, got %v", expectedTransaction1.TokenPathFrom, swapTransactions[0].TokenPathFrom)
		}

		if toLowerCaseHex(expectedTransaction1.TokenPathTo) != toLowerCaseHex(swapTransactions[0].TokenPathTo) {
			t.Errorf("Token Path To does not match expected value %v, got %v", expectedTransaction1.TokenPathTo, swapTransactions[0].TokenPathTo)
		}

		// TODO: add test for second call SweepToken
	})

	/*
		https://dashboard.tenderly.co/tx/mainnet/0x3dd149b18dd892a585577dd5abce476b554eac984c380cd8055b35dd05a42fb0
		Function: multicall(bytes[] data)

		MethodID: 0xac9650d8
		[0]:  0000000000000000000000000000000000000000000000000000000000000020
		[1]:  0000000000000000000000000000000000000000000000000000000000000001
		[2]:  0000000000000000000000000000000000000000000000000000000000000020
		[3]:  0000000000000000000000000000000000000000000000000000000000000104
		[4]:  414bf389000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead908
		[5]:  3c756cc2000000000000000000000000dac17f958d2ee523a2206206994597c1
		[6]:  3d831ec700000000000000000000000000000000000000000000000000000000
		[7]:  000001f4000000000000000000000000ead659a621741e7f4773637c6b738d0c
		[8]:  16a9ddb000000000000000000000000000000000000000000000000000000000
		[9]:  675365d200000000000000000000000000000000000000000000000045639182
		[10]: 44f4000000000000000000000000000000000000000000000000000000000004
		[11]: 91c21cd700000000000000000000000000000000000000000000000000000000
		[12]: 0000000000000000000000000000000000000000000000000000000000000000

		Expected:
		NumberOfCalls: 1
		"data": [
			"0x414bf389000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2000000000000000000000000dac17f958d2ee523a2206206994597c13d831ec700000000000000000000000000000000000000000000000000000000000001f4000000000000000000000000ead659a621741e7f4773637c6b738d0c16a9ddb000000000000000000000000000000000000000000000000000000000675365d20000000000000000000000000000000000000000000000004563918244f400000000000000000000000000000000000000000000000000000000000491c21cd70000000000000000000000000000000000000000000000000000000000000000"
		]

		0x414bf389
		000000000000000000000000
		c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2
		000000000000000000000000
		dac17f958d2ee523a2206206994597c13d831ec7
		0000000000000000000000000000000000000000000000000000000000000
		1f4
		000000000000000000000000
		ead659a621741e7f4773637c6b738d0c16a9ddb
		000000000000000000000000000000000000000000000000000000000
		675365d2
		000000000000000000000000000000000000000000000000
		4563918244f4
		00000000000000000000000000000000000000000000000000000000000
		491c21cd7
		0000000000000000000000000000000000000000000000000000000000000000


		struct ExactInputSingleParams {
			address tokenIn;
			address tokenOut;
			uint24 fee;
			address recipient;
			uint256 deadline;
			uint256 amountIn;
			uint256 amountOutMinimum;
			uint160 sqrtPriceLimitX96;
		}
	*/
	t.Run("Test DecodeMulticall", func(t *testing.T) {
		version := "v3"
		data := common.FromHex("0xac9650d80000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000104414bf389000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2000000000000000000000000dac17f958d2ee523a2206206994597c13d831ec700000000000000000000000000000000000000000000000000000000000001f4000000000000000000000000ead659a621741e7f4773637c6b738d0c16a9ddb000000000000000000000000000000000000000000000000000000000675365d20000000000000000000000000000000000000000000000004563918244f400000000000000000000000000000000000000000000000000000000000491c21cd7000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")

		tx := types.NewTransaction(
			0, // nonce
			common.Address{}, // to address
			big.NewInt(0), // value
			0, // gas limit
			big.NewInt(0), // gas price
			data, // data
		)

		swapTransactions, err := DecodeSwap(
			tx,
			version,
		)
		checkSwapNotNil(t, err, swapTransactions[0])

		expectedTransaction1 := &entity.SwapTransaction{
			MethodID: "414bf389",
			TokenPathFrom: "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			TokenPathTo:   "0xdac17f958d2ee523a2206206994597c13d831ec7",
			AmountIn: "4563918244f4",
			AmountOut: "491c21cd7",

		}

		if expectedTransaction1.MethodID != swapTransactions[0].MethodID {
			t.Errorf("Method ID does not match expected value %v, got %v", expectedTransaction1.MethodID, swapTransactions[0].MethodID)
		}

		if toLowerCaseHex(expectedTransaction1.TokenPathFrom) != toLowerCaseHex(swapTransactions[0].TokenPathFrom) {
			t.Errorf("Token Path From does not match expected value %v, got %v", expectedTransaction1.TokenPathFrom, swapTransactions[0].TokenPathFrom)
		}

		if toLowerCaseHex(expectedTransaction1.TokenPathTo) != toLowerCaseHex(swapTransactions[0].TokenPathTo) {
			t.Errorf("Token Path To does not match expected value %v, got %v", expectedTransaction1.TokenPathTo, swapTransactions[0].TokenPathTo)
		}

		if expectedTransaction1.AmountIn != swapTransactions[0].AmountIn {
			t.Errorf("Amount In does not match expected value %v, got %v", expectedTransaction1.AmountIn, swapTransactions[0].AmountIn)
		}

		if expectedTransaction1.AmountOut != swapTransactions[0].AmountOut {
			t.Errorf("Amount Out does not match expected value %v, got %v", expectedTransaction1.AmountOut, swapTransactions[0].AmountOut)
		}

	})

	t.Run("Test DecodeExactOutput", func(t *testing.T) {

		data := common.FromHex("0xf28c0498000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000007d14b142cad1379e85682f4b2006cdfed38988d30000000000000000000000000000000000000000000000000000000067576d6400000000000000000000000000000000000000000000000000000002b8c73e8000000000000000000000000000000000000000000000000000000001a9fa43ae000000000000000000000000000000000000000000000000000000000000004277e06c9eccf2e797fd462a92b6d7642ef85b0a44000bb8c02aaa39b223fe8d0a0e5c4f27ead9083c756cc20001f4a0b86991c6218b36c1d19d4a2e9eb0ce3606eb48000000000000000000000000000000000000000000000000000000000000")

		err := DecodeExactOutput(data, swapTransactionResult)
		checkSwapNotNil(t, err, swapTransactionResult)

		expected := &entity.SwapTransaction{
			TokenPathFrom: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
			TokenPathTo:   "0x77e06c9eccf2e797fd462a92b6d7642ef85b0a44",
			ToAddress:    "0x7d14b142cad1379e85682f4b2006cdfed38988d3",
			AmountOut:    "11690000000",
			AmountInMax: "7146718126",
		}

		if expected.TokenPathFrom != swapTransactionResult.TokenPathFrom {
			t.Errorf("Token Path From does not match expected value %v, got %v", expected.TokenPathFrom, swapTransactionResult.TokenPathFrom)
		}

		if expected.TokenPathTo != swapTransactionResult.TokenPathTo {
			t.Errorf("Token Path To does not match expected value %v, got %v", expected.TokenPathTo, swapTransactionResult.TokenPathTo)
		}

		if expected.ToAddress != swapTransactionResult.ToAddress {
			t.Errorf("To Address does not match expected value %v, got %v", expected.ToAddress, swapTransactionResult.ToAddress)
		}

		if expected.AmountOut != swapTransactionResult.AmountOut {
			t.Errorf("Amount Out does not match expected value %v, got %v", expected.AmountOut, swapTransactionResult.AmountOut)
		}

		if expected.AmountInMax != swapTransactionResult.AmountInMax {
			t.Errorf("Amount In Max does not match expected value %v, got %v", expected.AmountInMax, swapTransactionResult.AmountInMax)
		}

	})

}