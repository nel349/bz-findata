package decoder

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
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
}
