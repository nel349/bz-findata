package decoder

import (
	// "errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	v2 "github.com/nel349/bz-findata/internal/dex/eth/uniswap/v2"
	v3 "github.com/nel349/bz-findata/internal/dex/eth/uniswap/v3"
	"github.com/nel349/bz-findata/pkg/entity"
)

func DecodeSwap(tx *types.Transaction, version string) ([]*entity.SwapTransaction, error) {
	data := tx.Data()

    // Check if this is a multicall
    methodID := fmt.Sprintf("%x", data[:4])
    if methodID == v3.Multicall.Hex() {
		fmt.Println("Multicall detected")
        return DecodeMulticall(data)
    }
    
    // For non-multicall transactions, wrap single transaction in array
    swapTransactionResult := &entity.SwapTransaction{}    
    err := DecodeSwapGeneric(data, version, swapTransactionResult)
    if err != nil {
        return nil, err
    }
    
    return []*entity.SwapTransaction{swapTransactionResult}, nil
}

func DecodeSwapGeneric(data []byte, version string, swapTransactionResult *entity.SwapTransaction) error {
	methodID := fmt.Sprintf("%x", data[:4])

	// Debug prints
	// fmt.Println("Raw data length:", len(data))
	// fmt.Printf("Raw data hex: 0x%x\n", data)
	//

	var swapMethod interface{}
	var ok bool
	var swapMethodName string
	if version == "V2" {
		swapMethod, ok = v2.GetV2MethodFromID(methodID)
		swapMethodName = swapMethod.(v2.UniswapV2SwapMethod).String()
	} else {
		swapMethod, ok = v3.GetV3MethodFromID(methodID)
		swapMethodName = swapMethod.(v3.UniswapV3Method).String()
	}
	if !ok {
		return fmt.Errorf("unknown swap method: %s", methodID)
	}

    // Update existing struct instead of creating new one
    swapTransactionResult.AmountIn = "0"
    swapTransactionResult.Version = version
    swapTransactionResult.MethodID = methodID
	swapTransactionResult.MethodName = swapMethodName
	swapTransactionResult.Exchange = "Uniswap"

	switch swapMethod := swapMethod.(type) {
	case v2.UniswapV2SwapMethod:
		// fmt.Println("V2 swap method")
		// Lets do a switch for all the v2 swap methods
		switch swapMethod {
		case v2.SwapExactTokensForTokens:
			DecodeSwapExactTokensForTokens(data, version, swapTransactionResult)
		case v2.SwapExactTokensForTokensSupportingFeeOnTransferTokens:
			DecodeSwapExactTokensForTokensSupportingFeeOnTransferTokens(data, version, swapTransactionResult)
		case v2.SwapExactTokensForETHSupportingFeeOnTransferTokens:
			DecodeSwapExactTokensForETHSupportingFeeOnTransferTokens(data, version, swapTransactionResult)
		case v2.AddLiquidityETH:
			DecodeAddLiquidityETH(data, version, swapTransactionResult)
		case v2.AddLiquidity:
			DecodeAddLiquidity(data, version, swapTransactionResult)
		case v2.RemoveLiquidityETHWithPermit:
			DecodeRemoveLiquidityETHWithPermit(data, version, swapTransactionResult)
		case v2.RemoveLiquidityETH:
			DecodeRemoveLiquidityETH(data, version, swapTransactionResult)
		case v2.RemoveLiquidityETHWithPermitSupportingFeeOnTransferTokens:
			DecodeRemoveLiquidityETHWithPermitSupportingFeeOnTransferTokens(data, version, swapTransactionResult)
		case v2.SwapExactETHForTokensSupportingFeeOnTransferTokens:
			DecodeSwapExactETHForTokensSupportingFeeOnTransferTokens(data, version, swapTransactionResult)
		case v2.SwapExactETHForTokens:
			DecodeSwapExactETHForTokens(data, version, swapTransactionResult)
		case v2.SwapTokensForExactTokens:
			DecodeSwapTokensForExactTokens(data, version, swapTransactionResult)
		case v2.SwapExactTokensForETH:
			DecodeSwapExactTokensForETH(data, swapTransactionResult)
		case v2.SwapETHForExactTokens:
			DecodeSwapETHForExactTokens(data, swapTransactionResult)
		case v2.RemoveLiquidity:
			DecodeRemoveLiquidity(data, swapTransactionResult)
		case v2.SwapTokensForExactETH:
			DecodeSwapTokensForExactETH(data, swapTransactionResult)
		default:
			fmt.Println("not supported yet")
		}

	case v3.UniswapV3Method:
		fmt.Println("V3 swap method")
		// Lets do a switch for all the v3 swap methods
		switch swapMethod {
		case v3.ExactInputSingle:
			DecodeExactInputSingle(data, swapTransactionResult)
		case v3.ExactInput:
			DecodeExactInput(data, swapTransactionResult)
		case v3.ExactOutputSingle:
			DecodeExactOutputSingle(data, swapTransactionResult)
		default:
			fmt.Println("not supported yet")
		}
	}

	return nil
}

/*
* function for SwapExactTokensForTokens
// e.g
Function: swapExactTokensForTokens(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline)

MethodID: 0x38ed1739
[0]:  0000000000000000000000000000000000000000000000000108d3a3aa9f11e0
[1]:  0000000000000000000000000000000000000000000000000000000000000000
[2]:  00000000000000000000000000000000000000000000000000000000000000a0
[3]:  00000000000000000000000056eb903b0d2e858905feb7f1f4ad73458243d5a9
[4]:  00000000000000000000000000000000000000000000000000000000673576e7
[5]:  0000000000000000000000000000000000000000000000000000000000000002
[6]:  000000000000000000000000699ec925118567b6475fe495327ba0a778234aaa
[7]:  000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2
*/
func DecodeSwapExactTokensForTokens(data []byte, version string, swapTransactionResult *entity.SwapTransaction) error {
	// Skip first 4 bytes (method ID)
	data = data[4:]

	// print method signature
	// fmt.Printf("Method Signature: %s\n", v2.SwapExactTokensForTokens)

	// [0] amountIn (uint256)
	// fmt.Printf("Raw amountIn bytes: %x\n", data[:32])
	amountIn := new(big.Int).SetBytes(data[:32])
	// fmt.Printf("Amount In: %v\n", amountIn)

	// amountInDecimal := new(big.Float).Quo(
	// 	new(big.Float).SetInt(amountIn),
	// 	new(big.Float).SetFloat64(math.Pow(10, float64(tokenInfo.Decimals))),
	// )
	// amountInFloat64, _ := amountInDecimal.Float64()

	// [1] amountOutMin (uint256)
	amountOutMin := new(big.Int).SetBytes(data[32:64])
	// amountOutMinFloat64, _ := new(big.Float).Quo(
	// 	new(big.Float).SetInt(amountOutMin),
	// 	new(big.Float).SetFloat64(math.Pow(10, float64(tokenInfo.Decimals))),
	// ).Float64()

	// [2] path offset (points to [5])
	// offset := new(big.Int).SetBytes(data[64:96])

	// [3] to address (address)
	toAddress := fmt.Sprintf("%x", data[96:128])[24:] // Take last 20 bytes

	// [4] deadline (uint256)
	// deadline := new(big.Int).SetBytes(data[128:160])

	// [5] path length
	// pathOffset := 160  // 5 * 32 bytes
	// pathLength := new(big.Int).SetBytes(data[160:192]).Uint64()

	// [6] and [7] are token addresses in the path
	tokenPathFrom := fmt.Sprintf("0x%s", common.Bytes2Hex(data[192:224])[24:]) // First token in path
	tokenPathTo := fmt.Sprintf("0x%s", common.Bytes2Hex(data[224:256])[24:])   // Second token in path

	// fmt.Printf("Amount In: %v tokens\n", amountIn)
	// fmt.Printf("To Address: 0x%s\n", toAddress)
	// fmt.Printf("Token Path:\n")
	// fmt.Printf("  From: 0x%s\n", tokenPathFrom)
	// fmt.Printf("  To: 0x%s\n", tokenPathTo)

	swapTransactionResult.AmountIn = amountIn.String()
	swapTransactionResult.AmountOutMin = amountOutMin.String()
	swapTransactionResult.ToAddress = toAddress
	swapTransactionResult.TokenPathFrom = tokenPathFrom
	swapTransactionResult.TokenPathTo = tokenPathTo

	return nil
}

/*
Function: swapExactTokensForETHSupportingFeeOnTransferTokens(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline)

MethodID: 0x791ac947
[0]:  0000000000000000000000000000000000000000000422ca8b0a00a424ffffff
[1]:  00000000000000000000000000000000000000000000000007c28167a0c85474
[2]:  00000000000000000000000000000000000000000000000000000000000000a0
[3]:  0000000000000000000000001ad0eb3d4e0b79c20f8b3af24b706ae3c8e6a201
[4]:  000000000000000000000000000000000000000000000000000000006736f2f3
[5]:  0000000000000000000000000000000000000000000000000000000000000002
[6]:  000000000000000000000000f3c7cecf8cbc3066f9a87b310cebe198d00479ac
[7]:  000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2
[0]: 32 bytes - amountIn (uint256)
[1]: 32 bytes - amountOutMin (uint256)
[2]: 32 bytes - path array offset (uint256)
[3]: 32 bytes - to address (20 bytes + 12 bytes padding)
[4]: 32 bytes - deadline (uint256)
[5]: 32 bytes - path array length (uint256)
[6]: 32 bytes - first token address (20 bytes + 12 bytes padding)
[7]: 32 bytes - second token address (20 bytes + 12 bytes padding)
*/
func DecodeSwapExactTokensForETHSupportingFeeOnTransferTokens(
	data []byte,
	version string,
	swapTransactionResult *entity.SwapTransaction,
) error {

	data = data[4:]

	amountIn := new(big.Int).SetBytes(data[:32])
	// amountInFloat64, _ := new(big.Float).Quo(
	// 	new(big.Float).SetInt(amountIn),
	// 	new(big.Float).SetFloat64(math.Pow(10, float64(tokenInfo.Decimals))),
	// ).Float64()

	// [5] is where the path array starts (160 bytes from start)
	// [6] and [7] are token addresses in the path
	tokenPathFrom := fmt.Sprintf("0x%s", common.Bytes2Hex(data[192:224])[24:]) // First token in path
	tokenPathTo := fmt.Sprintf("0x%s", common.Bytes2Hex(data[224:256])[24:])   // Second token in path

	// fmt.Printf("Amount In: %v tokens\n", amountIn)
	// fmt.Printf("Token Path:\n")
	// fmt.Printf("  From: 0x%s\n", tokenPathFrom)
	// fmt.Printf("  To: 0x%s\n", tokenPathTo)

	swapTransactionResult.AmountIn = amountIn.String()
	swapTransactionResult.TokenPathFrom = tokenPathFrom
	swapTransactionResult.TokenPathTo = tokenPathTo

	return nil
}

/*
https://etherscan.io/tx/0xcf39e1501430f75f7ee041781b62592c6ba8a3749e5b5f3813f086023607dc1b
Function: swapExactTokensForTokensSupportingFeeOnTransferTokens(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline)

MethodID: 0x5c11d795
[0]:  0000000000000000000000000000000000000000000ec068614236ee611fe947
[1]:  0000000000000000000000000000000000000000000000000766b47bedbc6e9d
[2]:  00000000000000000000000000000000000000000000000000000000000000a0
[3]:  000000000000000000000000c54a957d2e1da5067c7ad32d38d3a2bc2524531c
[4]:  000000000000000000000000000000000000000000000000000001932e9df0b8
[5]:  0000000000000000000000000000000000000000000000000000000000000002
[6]:  000000000000000000000000debcad12e9c454a7338b3ec0c8058eec688c79d5
[7]:  000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2

	{
	  "amountIn": "17833581308923813721794887",
	  "amountOutMin": "533312050252508829",
	  "path": [
	    "0xdebcad12e9c454a7338b3ec0c8058eec688c79d5",
	    "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"
	  ],
	  "to": "0xc54a957d2e1da5067c7ad32d38d3a2bc2524531c",
	  "deadline": "1731653923000"
	}
*/
func DecodeSwapExactTokensForTokensSupportingFeeOnTransferTokens(
	data []byte,
	version string,
	swapTransactionResult *entity.SwapTransaction,
) error {
	data = data[4:]

	amountIn := new(big.Int).SetBytes(data[:32])
	// amountInFloat64, _ := new(big.Float).Quo(
	// 	new(big.Float).SetInt(amountIn),
	// 	new(big.Float).SetFloat64(math.Pow(10, float64(fromTokenInfo.Decimals))),
	// ).Float64()

	// [6] and [7] are token addresses in the path
	tokenPathFrom := fmt.Sprintf("0x%s", common.Bytes2Hex(data[192:224])[24:]) // First token in path
	tokenPathTo := fmt.Sprintf("0x%s", common.Bytes2Hex(data[224:256])[24:])   // Second token in path

	swapTransactionResult.AmountIn = amountIn.String()
	swapTransactionResult.TokenPathFrom = tokenPathFrom
	swapTransactionResult.TokenPathTo = tokenPathTo
	return nil
}

/*
Function: addLiquidityETH(address token, uint256 amountTokenDesired, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline)

MethodID: 0xf305d719
[0]:  0000000000000000000000006de1b3605a5e587e969e08166e3e5c5bfc4b1a16
[1]:  0000000000000000000000000000000000000000033b2e3c9fd0803ce8000000
[2]:  0000000000000000000000000000000000000000033b2e3c9fd0803ce8000000
[3]:  000000000000000000000000000000000000000000000001158e460913d00000
[4]:  0000000000000000000000002a853910205bbc1a879809441ce12e813c9eb018
[5]:  00000000000000000000000000000000000000000000000000000000673c2107

	{
	  "token": "0x6de1b3605a5e587e969e08166e3e5c5bfc4b1a16",
	  "amountTokenDesired": "1000000000000000000000000000",
	  "amountTokenMin": "1000000000000000000000000000",
	  "amountETHMin": "20000000000000000000",
	  "to": "0x2a853910205bbc1a879809441ce12e813c9eb018",
	  "deadline": "1731993863"
	}
*/
func DecodeAddLiquidityETH(data []byte, version string, swapTransactionResult *entity.SwapTransaction) error {
	data = data[4:]

	// [0] token address (20 bytes + 12 bytes padding)
	tokenAddress := fmt.Sprintf("0x%s", common.Bytes2Hex(data[0:32])[24:])

	amountTokenDesired := new(big.Int).SetBytes(data[32:64]).String()
	amountTokenMin := new(big.Int).SetBytes(data[64:96]).String()
	amountETHMin := new(big.Int).SetBytes(data[96:128]).String()

	swapTransactionResult.AmountTokenDesired = amountTokenDesired
	swapTransactionResult.AmountTokenMin = amountTokenMin
	swapTransactionResult.AmountETHMin = amountETHMin
	swapTransactionResult.TokenPathTo = tokenAddress // Address of the token to pair with ETH
	return nil
}

/*
	Function: removeLiquidityETHWithPermit(address token, uint256 liquidity, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline, bool approveMax, uint8 v, bytes32 r, bytes32 s)

	MethodID: 0xded9382a
	[0]:  000000000000000000000000fe34cbcaef94a06a8fc1adce86d486f49af242ba
	[1]:  00000000000000000000000000000000000000000000119707a239721536ed2a
	[2]:  0000000000000000000000000000000000000000000000000000000000000000
	[3]:  0000000000000000000000000000000000000000000000000000000000000000
	[4]:  0000000000000000000000001072d3be58b71b386724c624890547499fe39b89
	[5]:  00000000000000000000000000000000000000000000000000000000673e95b7
	[6]:  0000000000000000000000000000000000000000000000000000000000000000
	[7]:  000000000000000000000000000000000000000000000000000000000000001b
	[8]:  ac42c5965bedf9500e171015c523074ac3162bf093ffd0627d51691a19abbae6
	[9]:  386a09c71e9bbd14d91db91f9844c161f72753b1c84cf974b5f4fcbcae3d9418
*/
func DecodeRemoveLiquidityETHWithPermit(data []byte, version string, swapTransactionResult *entity.SwapTransaction) error {
	data = data[4:]


	// [0] token address (20 bytes + 12 bytes padding)
	tokenAddress := fmt.Sprintf("0x%s", common.Bytes2Hex(data[12:32]))

	// [1] liquidity (uint256)
	liquidity := fmt.Sprintf("%d", new(big.Int).SetBytes(data[54:64]))

	swapTransactionResult.TokenPathFrom = tokenAddress
	swapTransactionResult.Liquidity = liquidity
	return nil
}

/*
Function: removeLiquidityETH(address token, uint256 liquidity, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline)

MethodID: 0x02751cec
[0]:  0000000000000000000000000f5c78f152152dda52a2ea45b0a8c10733010748
[1]:  000000000000000000000000000000000000000000000001a2ef2e3a848cf74c
[2]:  0000000000000000000000000000000000000000000008dacc72a69e84b90445
[3]:  0000000000000000000000000000000000000000000000000059949fb27347c8
[4]:  000000000000000000000000e1428dc0202f1f7d61d3116269d45d4ebbe7502d
[5]:  00000000000000000000000000000000000000000000000000000000673ea577
*/
func DecodeRemoveLiquidityETH(data []byte, version string, swapTransactionResult *entity.SwapTransaction) error {
	// use same logic as removeLiquidityETHWithPermit as we only need to decode the token address and liquidity
	return DecodeRemoveLiquidityETHWithPermit(data, version, swapTransactionResult)
}

func DecodeRemoveLiquidityETHWithPermitSupportingFeeOnTransferTokens(data []byte, version string, swapTransactionResult *entity.SwapTransaction) error {
	// use same logic as removeLiquidityETHWithPermit as we only need to decode the token address and liquidity
	return DecodeRemoveLiquidityETHWithPermit(data, version, swapTransactionResult)
}

/*
	#####Receiving Eth for tokens####	
*/


/*
	https://dashboard.tenderly.co/tx/mainnet/0xf798d58c2018c4c6b7f93c9373bd6f67e20a88793b7c7f9f679f768df0efb88c
	Function: swapExactETHForTokensSupportingFeeOnTransferTokens(uint256 amountOutMin, address[] path, address to, uint256 deadline)

	MethodID: 0xb6f9de95
	[0]:  0000000000000000000000000000000000000000000000000004f2cf373add62
	[1]:  0000000000000000000000000000000000000000000000000000000000000080
	[2]:  00000000000000000000000076db926b75e225af64b954c95fef653926ea7965
	[3]:  000000000000000000000000000000000000000000000000000001938f506963
	[4]:  0000000000000000000000000000000000000000000000000000000000000002
	[5]:  000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2
	[6]:  000000000000000000000000790336af90933aa7bd10d4534db6909507098440
	{
		"amountOutMin": "1392871705599330",
		"path": [
			"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			"0x790336af90933aa7bd10d4534db6909507098440"
		],
		"to": "0x76db926b75e225af64b954c95fef653926ea7965",
		"deadline": "1733276232035"
	}
*/
func DecodeSwapExactETHForTokensSupportingFeeOnTransferTokens(data []byte, version string, swapTransactionResult *entity.SwapTransaction) error {

	// skip first 4 bytes (method ID)
	data = data[4:]

	// [0] amountOutMin
	amountOutMin := new(big.Int).SetBytes(data[:32]).String()

	// [5] and [6] are token addresses in the path
	tokenPathFrom := fmt.Sprintf("0x%s", common.Bytes2Hex(data[160:192])[24:]) // First token in path
	tokenPathTo := fmt.Sprintf("0x%s", common.Bytes2Hex(data[192:224])[24:])   // Second token in path

	swapTransactionResult.AmountOutMin = amountOutMin
	swapTransactionResult.AmountIn = amountOutMin // same as amountOutMin for this function
	swapTransactionResult.TokenPathFrom = tokenPathFrom
	swapTransactionResult.TokenPathTo = tokenPathTo
	return nil
}

/*
	https://dashboard.tenderly.co/tx/mainnet/0x7d310e4465220787b594afd033014cc97c0d6ae4fccf290e3aa4110e5a6f26c8
	Function: swapTokensForExactTokens(uint256 amountOut, uint256 amountInMax, address[] path, address to, uint256 deadline)

	MethodID: 0x8803dbee
	[0]:  0000000000000000000000000000000000000000000009b6c2b9818b505e2000
	[1]:  00000000000000000000000000000000000000000000000006e9f733c51ed703
	[2]:  00000000000000000000000000000000000000000000000000000000000000a0
	[3]:  000000000000000000000000f1da51404a3c42dc46fcc6924944fab21fd7e9b9
	[4]:  00000000000000000000000000000000000000000000000000000000674fb961
	[5]:  0000000000000000000000000000000000000000000000000000000000000002
	[6]:  000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2
	[7]:  000000000000000000000000425087bf4969f45818c225ae30f8560ce518582e


	{
		"amountOut": "45872637155791343591424",
		"amountInMax": "498201035523675907",
		"path": [
			"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			"0x425087bf4969f45818c225ae30f8560ce518582e"
		],
		"to": "0xf1da51404a3c42dc46fcc6924944fab21fd7e9b9",
		"deadline": "1733278049"
	}

*/

func DecodeSwapTokensForExactTokens(data []byte, version string, swapTransactionResult *entity.SwapTransaction) error {

	data = data[4:]

	// [0] amountOut
	amountOut := new(big.Int).SetBytes(data[:32]).String()

	// [1] amountInMax
	amountInMax := new(big.Int).SetBytes(data[32:64]).String()

	// [6] and [7] are token addresses in the path
	tokenPathFrom := fmt.Sprintf("0x%s", common.Bytes2Hex(data[192:224])[24:]) // First token in path
	tokenPathTo := fmt.Sprintf("0x%s", common.Bytes2Hex(data[224:256])[24:])   // Second token in path

	swapTransactionResult.AmountOut = amountOut
	swapTransactionResult.AmountInMax = amountInMax
	swapTransactionResult.AmountIn = amountInMax // same as amountInMax for this function
	swapTransactionResult.TokenPathFrom = tokenPathFrom
	swapTransactionResult.TokenPathTo = tokenPathTo
	return nil
}

/*
	https://dashboard.tenderly.co/tx/mainnet/0x13153a9aa376a5414e3e03f3fd811dc3f4043dab005b3c121905965e309b397d
	Function: swapExactTokensForETH(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline)

	MethodID: 0x18cbafe5
	[0]:  000000000000000000000000000000000000000002386e78249bcbd50c1c0000
	[1]:  000000000000000000000000000000000000000000000000011b475f5f08b23a
	[2]:  00000000000000000000000000000000000000000000000000000000000000a0
	[3]:  0000000000000000000000008c5bd953f3e2756d5270abe14ed8d9c7f574b42a
	[4]:  00000000000000000000000000000000000000000000000000000000ce9f5c2a
	[5]:  0000000000000000000000000000000000000000000000000000000000000002
	[6]:  000000000000000000000000b0ac2b5a73da0e67a8e5489ba922b3f8d582e058
	[7]:  000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2

	{
		"amountIn": "687191542101440000000000000",
		"amountOutMin": "79735893350986298",
		"path": [
			"0xb0ac2b5a73da0e67a8e5489ba922b3f8d582e058",
			"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"
		],
		"to": "0x8c5bd953f3e2756d5270abe14ed8d9c7f574b42a",
		"deadline": "3466550314"
	}
*/

func DecodeSwapExactTokensForETH(data []byte, swapTransactionResult *entity.SwapTransaction) error {

	data = data[4:]

	// [0] amountIn	
	amountIn := new(big.Int).SetBytes(data[:32]).String()

	// [1] amountOutMin
	amountOutMin := new(big.Int).SetBytes(data[32:64]).String()


	// [6] and [7] are token addresses in the path
	tokenPathFrom := fmt.Sprintf("0x%s", common.Bytes2Hex(data[192:224])[24:]) // First token in path
	tokenPathTo := fmt.Sprintf("0x%s", common.Bytes2Hex(data[224:256])[24:])   // Second token in path

	swapTransactionResult.AmountIn = amountIn
	swapTransactionResult.AmountOutMin = amountOutMin
	swapTransactionResult.TokenPathFrom = tokenPathFrom
	swapTransactionResult.TokenPathTo = tokenPathTo
	return nil
}

/*
	https://dashboard.tenderly.co/tx/mainnet/0x8736fd7b40cc01eaf145cb274a33ef4657af81dd9f48fad5bad99fb1bdd64088
	Function: swapETHForExactTokens(uint256 amountOut, address[] path, address to, uint256 deadline)

	MethodID: 0xfb3bdb41
	[0]:  000000000000000000000000000000000000000000000000000003aa29f07d58
	[1]:  0000000000000000000000000000000000000000000000000000000000000080
	[2]:  00000000000000000000000011a8286db1eb78e4d0d3409be5ddafac2966cc0b
	[3]:  000000000000000000000000000000000000000000000000000000006750f480
	[4]:  0000000000000000000000000000000000000000000000000000000000000002
	[5]:  000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2
	[6]:  000000000000000000000000e0805c80588913c1c2c89ea4a8dcf485d4038a3e

	{
		"amountOut": "4029382950232",
		"path": [
			"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			"0xe0805c80588913c1c2c89ea4a8dcf485d4038a3e"
		],
		"to": "0x11a8286db1eb78e4d0d3409be5ddafac2966cc0b",
		"deadline": "1733358720"
	}

*/

func DecodeSwapETHForExactTokens(data []byte, swapTransactionResult *entity.SwapTransaction) error {

	data = data[4:]

	// [0] amountOut
	amountOut := new(big.Int).SetBytes(data[:32]).String()

	// [5] and [6] are token addresses in the path
	tokenPathFrom := fmt.Sprintf("0x%s", common.Bytes2Hex(data[160:192])[24:]) // First token in path
	tokenPathTo := fmt.Sprintf("0x%s", common.Bytes2Hex(data[192:224])[24:])   // Second token in path

	swapTransactionResult.AmountOut = amountOut
	swapTransactionResult.TokenPathFrom = tokenPathFrom
	swapTransactionResult.TokenPathTo = tokenPathTo
	return nil
}

/*

	https://dashboard.tenderly.co/tx/mainnet/0x93127ec9b1c13815c6aca9363c737fcb5f874c6f8ddfc6a795b01308d5f17225
	Function: removeLiquidity(address tokenA, address tokenB, uint256 liquidity, uint256 amountAMin, uint256 amountBMin, address to, uint256 deadline)

	MethodID: 0xbaa2abde
	[0]:  00000000000000000000000099ec69f6624abd625782e2127f7ca23432aab7a1
	[1]:  000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2
	[2]:  00000000000000000000000000000000000000000000000000048db1aa0d89d0
	[3]:  000000000000000000000000000000000000000000000000000084dca2559ac8
	[4]:  0000000000000000000000000000000000000000000000000027ce754fedea24
	[5]:  000000000000000000000000ba8ed95797f9bf37c99f564d9aa26eeb1851bf1f
	[6]:  00000000000000000000000000000000000000000000000000000000675164f2

	{
		"tokenA": "0x99ec69f6624abd625782e2127f7ca23432aab7a1",
		"tokenB": "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		"liquidity": "1281694108584400",
		"amountAMin": "146083151190728",
		"amountBMin": "11204527339203108",
		"to": "0xba8ed95797f9bf37c99f564d9aa26eeb1851bf1f",
		"deadline": "1733387506"
	}
*/

// TODO: consult return values for this function to get the actual amount of tokens removed
func DecodeRemoveLiquidity(data []byte, swapTransactionResult *entity.SwapTransaction) error {

	data = data[4:]

	// [0] tokenA
	tokenA := fmt.Sprintf("0x%s", common.Bytes2Hex(data[:32])[24:])

	// [1] tokenB
	tokenB := fmt.Sprintf("0x%s", common.Bytes2Hex(data[32:64])[24:])


	// [2] liquidity
	liquidity := new(big.Int).SetBytes(data[64:96]).String()

	// [3] amountAMin
	amountAMin := new(big.Int).SetBytes(data[96:128]).String()

	// [4] amountBMin
	amountBMin := new(big.Int).SetBytes(data[128:160]).String()

	swapTransactionResult.TokenA = tokenA
	swapTransactionResult.TokenB = tokenB
	swapTransactionResult.Liquidity = liquidity
	swapTransactionResult.AmountAMin = amountAMin
	swapTransactionResult.AmountBMin = amountBMin
	return nil
}

/*

	https://github.com/nel349/bz-findata/issues/27

	https://dashboard.tenderly.co/tx/mainnet/0x84ccd996bc71ab8bae0e20ed4c0f3a87a8e822f430425b17800bfa37b5b5cb1e
	Function: swapExactETHForTokens(uint256 amountOutMin, address[] path, address to, uint256 deadline)

	MethodID: 0x7ff36ab5
	[0]:  00000000000000000000000000000000000000000000000589a74db4fba713f1
	[1]:  0000000000000000000000000000000000000000000000000000000000000080
	[2]:  000000000000000000000000659402c9c3eb08503ab3ecbaff5f384f6f7d62af
	[3]:  0000000000000000000000000000000000000000000000000000000067516343
	[4]:  0000000000000000000000000000000000000000000000000000000000000002
	[5]:  000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2
	[6]:  0000000000000000000000006019dcb2d0b3e0d1d8b0ce8d16191b3a4f93703d

	{
		"amountOutMin": "102152702512566047729",
		"path": [
			"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			"0x6019dcb2d0b3e0d1d8b0ce8d16191b3a4f93703d"
		],
		"to": "0x659402c9c3eb08503ab3ecbaff5f384f6f7d62af",
		"deadline": "1733387075"
	}
*/

func DecodeSwapExactETHForTokens(data []byte, version string, swapTransactionResult *entity.SwapTransaction) error {

	data = data[4:]

	// [0] amountOutMin
	amountOutMin := new(big.Int).SetBytes(data[:32]).String()

	// [5] and [6] are token addresses in the path
	tokenPathFrom := fmt.Sprintf("0x%s", common.Bytes2Hex(data[160:192])[24:]) // First token in path
	tokenPathTo := fmt.Sprintf("0x%s", common.Bytes2Hex(data[192:224])[24:])   // Second token in path

	swapTransactionResult.AmountOutMin = amountOutMin
	swapTransactionResult.TokenPathFrom = tokenPathFrom
	swapTransactionResult.TokenPathTo = tokenPathTo
	return nil
}

/*
	https://dashboard.tenderly.co/tx/mainnet/0xaae2b315f040eb9f342e7f561add662b881792ccef3c09120645c652f81bad8e
	Function: swapTokensForExactETH(uint256 amountOut, uint256 amountInMax, address[] path, address to, uint256 deadline)

	MethodID: 0x4a25d94a
	[0]:  0000000000000000000000000000000000000000000000000473c8321cbe3b23
	[1]:  ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff
	[2]:  00000000000000000000000000000000000000000000000000000000000000a0
	[3]:  000000000000000000000000f604ec052a32eb53aaa3b104993bc4d5b6132a52
	[4]:  0000000000000000000000000000000000000000000000000000000067536159
	[5]:  0000000000000000000000000000000000000000000000000000000000000002
	[6]:  0000000000000000000000006a4402a535d74bd0c9cdb5ce2d51822fc9f6620e
	[7]:  000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2


	{
		"amountOut": "320820116029586211",
		"amountInMax": "115792089237316195423570985008687907853269984665640564039457584007913129639935",
		"path": [
			"0x6a4402a535d74bd0c9cdb5ce2d51822fc9f6620e",
			"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"
		],
		"to": "0xf604ec052a32eb53aaa3b104993bc4d5b6132a52",
		"deadline": "1733517657"
	}
*/

func DecodeSwapTokensForExactETH(data []byte, swapTransactionResult *entity.SwapTransaction) error {

	data = data[4:]

	// [0] amountOut
	amountOut := new(big.Int).SetBytes(data[:32]).String()

	// [1] amountInMax
	amountInMax := new(big.Int).SetBytes(data[32:64]).String()

	// [6] and [7] are token addresses in the path
	tokenPathFrom := fmt.Sprintf("0x%s", common.Bytes2Hex(data[192:224])[24:]) // First token in path
	tokenPathTo := fmt.Sprintf("0x%s", common.Bytes2Hex(data[224:256])[24:])   // Second token in path

	swapTransactionResult.AmountOut = amountOut
	swapTransactionResult.AmountInMax = amountInMax
	swapTransactionResult.TokenPathFrom = tokenPathFrom
	swapTransactionResult.TokenPathTo = tokenPathTo
	return nil
}

/*
	https://dashboard.tenderly.co/tx/mainnet/0x27d8fea3897b0d0acd15900b77bbb85a894256a65b23c9b9cf900526bb760caf
	Function: addLiquidity(address tokenA, address tokenB, uint256 amountADesired, uint256 amountBDesired, uint256 amountAMin, uint256 amountBMin, address to, uint256 deadline)

	MethodID: 0xe8e33700
	[0]:  000000000000000000000000910812c44ed2a3b611e4b051d9d83a88d652e2dd
	[1]:  000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2
	[2]:  000000000000000000000000000000000000000000003f5b6a1d988446cc0855
	[3]:  0000000000000000000000000000000000000000000000000a11bc5997c59c50
	[4]:  000000000000000000000000000000000000000000003f0a5143d908bc33f8ee
	[5]:  0000000000000000000000000000000000000000000000000a04d8d92517d292
	[6]:  0000000000000000000000007349bd4f327e697814ff96ac48861b44562ed722
	[7]:  0000000000000000000000000000000000000000000000000000000067538242

	{
		"tokenA": "0x910812c44ed2a3b611e4b051d9d83a88d652e2dd",
		"tokenB": "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		"amountADesired": "299195388566931453511765",
		"amountBDesired": "725568107967781968",
		"amountAMin": "297699411624096796244206",
		"amountBMin": "721940267427943058",
		"to": "0x7349bd4f327e697814ff96ac48861b44562ed722",
		"deadline": "1733526082"
	}
*/

func DecodeAddLiquidity(data []byte, version string, swapTransactionResult *entity.SwapTransaction) error {

	data = data[4:]

	// [0] tokenA
	tokenA := fmt.Sprintf("0x%s", common.Bytes2Hex(data[:32])[24:])

	// [1] tokenB
	tokenB := fmt.Sprintf("0x%s", common.Bytes2Hex(data[32:64])[24:])

	// [2] amountADesired
	amountADesired := new(big.Int).SetBytes(data[64:96]).String()

	// [3] amountBDesired
	amountBDesired := new(big.Int).SetBytes(data[96:128]).String()

	// [4] amountAMin
	amountAMin := new(big.Int).SetBytes(data[128:160]).String()

	// [5] amountBMin
	amountBMin := new(big.Int).SetBytes(data[160:192]).String()

	swapTransactionResult.TokenA = tokenA
	swapTransactionResult.TokenB = tokenB
	swapTransactionResult.AmountADesired = amountADesired
	swapTransactionResult.AmountBDesired = amountBDesired
	swapTransactionResult.AmountAMin = amountAMin
	swapTransactionResult.AmountBMin = amountBMin
	return nil
}
