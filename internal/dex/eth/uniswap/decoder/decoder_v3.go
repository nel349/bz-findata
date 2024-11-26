package decoder

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/nel349/bz-findata/pkg/entity"
)

/*
	Decoders for Uniswap V3 methods
*/

/*
	DecodeExactInputSingle
	Example:

	https://etherscan.io/tx/0x507937f90d3a14f357da57c582aad7e7f0a4fde516a6bae805b9278194714afb

Function: exactInputSingle(tuple params)

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

	MethodID: 0x414bf389
	[0]:  000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2
	[1]:  000000000000000000000000ee2a03aa6dacf51c18679c516ad5283d8e7c2637
	[2]:  0000000000000000000000000000000000000000000000000000000000000bb8
	[3]:  000000000000000000000000f5213a6a2f0890321712520b8048d9886c1a9900
	[4]:  000000000000000000000000000000000000000000000000000000006736f0e4
	[5]:  0000000000000000000000000000000000000000000000000b9eafe9ee6f4000
	[6]:  000000000000000000000000000000000000000000000000000019fe199f2e10
	[7]:  0000000000000000000000000000000000000000000000000000000000000000

	{
	  "params": {
	    "tokenIn": "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
	    "tokenOut": "0xee2a03aa6dacf51c18679c516ad5283d8e7c2637",
	    "fee": "3000",
	    "recipient": "0xf5213a6a2f0890321712520b8048d9886c1a9900",
	    "deadline": "1731653860",
	    "amountIn": "837300000000000000",
	    "amountOutMinimum": "28579142250000",
	    "sqrtPriceLimitX96": "0"
	  }
	}
*/
func DecodeExactInputSingle(data []byte, version string, swapTransactionResult *entity.SwapTransaction) (error) {

	data = data[4:]
    // [0]: tokenIn
    tokenIn := fmt.Sprintf("0x%s", common.Bytes2Hex(data[0:32])[24:])

    // [1]: tokenOut
    tokenOut := fmt.Sprintf("0x%s", common.Bytes2Hex(data[32:64])[24:])

    // [2]: fee is 3000 == 0xbb8

	
    // [3]: recipient e.g 0xf5213a6a2f0890321712520b8048d9886c1a9900
    recipient := fmt.Sprintf("0x%s", common.Bytes2Hex(data[96:128])[24:])

    // [4]: deadline
    // [5]: amountIn e.g 0xb9eafe9ee6f4000 == 837300000000000000 and padding
    amountIn := new(big.Int).SetBytes(data[160:192])

    // [6]: amountOutMinimum
    // [7]: sqrtPriceLimitX96

	swapTransactionResult.TokenPathFrom = tokenIn
	swapTransactionResult.TokenPathTo = tokenOut
	swapTransactionResult.ToAddress = recipient
	swapTransactionResult.AmountIn = amountIn.String()

    return nil
}


/*
	DecodeExactInput

	   struct ExactInputParams {
        bytes path;
        address recipient;
        uint256 deadline;
        uint256 amountIn;
        uint256 amountOutMinimum;
    }

	Function: exactInput(tuple params)

	MethodID: 0xc04b8d59
	[0]:  0000000000000000000000000000000000000000000000000000000000000020
	[1]:  00000000000000000000000000000000000000000000000000000000000000a0
	[2]:  000000000000000000000000b0ba33566bd35bcb80738810b2868dc1ddd1f0e9
	[3]:  00000000000000000000000000000000000000000000000000000000673c1fcf
	[4]:  000000000000000000000000000000000000000000000000066eced5a631d580
	[5]:  000000000000000000000000000000000000000000000b38ca1ce396b5800000
	[6]:  000000000000000000000000000000000000000000000000000000000000002b
	[7]:  c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200271038e382f74dfb84608f
	[8]:  3c1f10187f6bef5951de93000000000000000000000000000000000000000000


	{
		"params": {
			"path": "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc200271038e382f74dfb84608f3c1f10187f6bef5951de93",
			"recipient": "0xb0ba33566bd35bcb80738810b2868dc1ddd1f0e9",
			"deadline": "1731993551",
			"amountIn": "463535228677379456",
			"amountOutMinimum": "52993612745225271246848"
		}
	}

		// 0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2  // First token address (WETH)
		// 002710                                        // Fee tier (0.0027 = 0.27%)
		// 0x38e382f74dfb84608f3c1f10187f6bef5951de93  // Second token address
*/
func DecodeExactInput(data []byte, version string, swapTransactionResult *entity.SwapTransaction) (error) {
	data = data[4:]

	// [4] amountIn 000000000000000000000000000000000000000000000000066eced5a631d580
	amountIn := new(big.Int).SetBytes(data[144:160])
	
	// path [7] and [8]

	// first token address
	firstTokenAddress := fmt.Sprintf("0x%s", common.Bytes2Hex(data[224:244]))
	
	// fee 
	fee := fmt.Sprintf("0x%s", common.Bytes2Hex(data[244:247]))

	// second token address
	secondTokenAddress := fmt.Sprintf("0x%s", common.Bytes2Hex(data[247:267]))

	swapTransactionResult.AmountIn = amountIn.String()
	swapTransactionResult.TokenPathFrom = firstTokenAddress
	swapTransactionResult.TokenPathTo = secondTokenAddress
	swapTransactionResult.Fee = fee

	return nil
}
