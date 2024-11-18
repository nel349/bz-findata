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
func DecodeExactInputSingle(data []byte, version string) (*entity.SwapTransaction, error) {

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
    return &entity.SwapTransaction{
        AmountIn:       amountIn,
        TokenPathFrom:  tokenIn,
        TokenPathTo:    tokenOut,
        ToAddress:      recipient,
    }, nil
}
