package decoder

import (
	"encoding/binary"
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
func DecodeExactInputSingle(data []byte, swapTransactionResult *entity.SwapTransaction) error {

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
func DecodeExactInput(data []byte, swapTransactionResult *entity.SwapTransaction) error {
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

/*
https://dashboard.tenderly.co/tx/mainnet/0x38a89aca46ef19f71dc1c39e2fe519e528d22cdb7fdbe740e3f298e61cdec322

Function: exactOutputSingle(tuple params)

MethodID: 0xdb3e2198
[0]:  000000000000000000000000d31a59c85ae9d8edefec411d448f90841571b89c
[1]:  000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2
[2]:  0000000000000000000000000000000000000000000000000000000000000bb8
[3]:  000000000000000000000000f5213a6a2f0890321712520b8048d9886c1a9900
[4]:  00000000000000000000000000000000000000000000000000000000674fad84
[5]:  0000000000000000000000000000000000000000000000012b824fb44600c000
[6]:  0000000000000000000000000000000000000000000000000000004e3f884eef
[7]:  0000000000000000000000000000000000000000000000000000000000000000
*/
func DecodeExactOutputSingle(data []byte, swapTransactionResult *entity.SwapTransaction) error {

	data = data[4:]

	// [0] tokenIn
	tokenIn := fmt.Sprintf("0x%s", common.Bytes2Hex(data[0:32])[24:])

	// [1] tokenOut
	tokenOut := fmt.Sprintf("0x%s", common.Bytes2Hex(data[32:64])[24:])

	// [6] amountInMaximum
	amountInMaximum := new(big.Int).SetBytes(data[192:224])

	swapTransactionResult.AmountInMax = amountInMaximum.String()
	swapTransactionResult.TokenPathFrom = tokenIn
	swapTransactionResult.TokenPathTo = tokenOut
	swapTransactionResult.AmountIn = amountInMaximum.String() // same as amountInMaximum for this function

	return nil
}

/*
Function: multicall(bytes[] data)

MethodID: 0xac9650d8
[0]:  0000000000000000000000000000000000000000000000000000000000000020
[1]:  0000000000000000000000000000000000000000000000000000000000000002
[2]:  0000000000000000000000000000000000000000000000000000000000000040
[3]:  0000000000000000000000000000000000000000000000000000000000000180
[4]:  0000000000000000000000000000000000000000000000000000000000000104
[5]:  db3e2198000000000000000000000000960692640ac4986ffce41620b7e3aa03
[6]:  cf1a0e8f000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead908
[7]:  3c756cc200000000000000000000000000000000000000000000000000000000
[8]:  00000bb800000000000000000000000000000000000000000000000000000000
[9]:  0000000000000000000000000000000000000000000000000000000000000000
[10]: 674fa9ef0000000000000000000000000000000000000000000000000214e834
[11]: 8c4f00000000000000000000000000000000000000000000000002138669df85
[12]: 9f1c439f00000000000000000000000000000000000000000000000000000000
[13]: 0000000000000000000000000000000000000000000000000000000000000000
[14]: 0000000000000000000000000000000000000000000000000000000000000044
[15]: 49404b7c0000000000000000000000000000000000000000000000000214e834
[16]: 8c4f00000000000000000000000000000f648e13fbe57251d367b08f42e8557d
[17]: 7461637a00000000000000000000000000000000000000000000000000000000

// first call is ExactOutputSingle
// second call is sweepToken

	{
		"data": [
			"0xdb3e2198000000000000000000000000960692640ac4986ffce41620b7e3aa03cf1a0e8f000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc20000000000000000000000000000000000000000000000000000000000000bb8000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000674fa9ef0000000000000000000000000000000000000000000000000214e8348c4f00000000000000000000000000000000000000000000000002138669df859f1c439f0000000000000000000000000000000000000000000000000000000000000000",
			"0x49404b7c0000000000000000000000000000000000000000000000000214e8348c4f00000000000000000000000000000f648e13fbe57251d367b08f42e8557d7461637a"
		]
	}
*/
func DecodeMulticall(data []byte) ([]*entity.SwapTransaction, error) {

	swapTransactionResult := &entity.SwapTransaction{}
	err := DecodeDataArray(data, swapTransactionResult)
	if err != nil {
		return nil, err
	}

	if swapTransactionResult.NumberOfCalls == 0 {
		return nil, fmt.Errorf("no calls found")
	}

	// create an array of swap transactions to be returned
	swapTransactions := make([]*entity.SwapTransaction, swapTransactionResult.NumberOfCalls)

	// lets iterate over the calls data and decode each one and add it to the swapTransactions array
	for i := 0; i < swapTransactionResult.NumberOfCalls; i++ {
		swapTransactions[i] = &entity.SwapTransaction{}  // Initialize each element
		callData := common.FromHex(swapTransactionResult.CallsData[i])
		err := DecodeSwapGeneric(callData, "V3", swapTransactions[i])
		if err != nil {
			return nil, err
		}
	}

	return swapTransactions, nil
}

func DecodeDataArray(data []byte, swapTransactionResult *entity.SwapTransaction) error {
    // ignore the first 4 bytes (methodID) for the multicall
    data = data[4:]

    // [0]: offset to array start (should be 32/0x20)
    arrayStartOffset := new(big.Int).SetBytes(data[0:32])
    if arrayStartOffset.Uint64() != 32 {
        return fmt.Errorf("invalid array start offset: expected 32, got %d", arrayStartOffset)
    }

    // [1]: number of calls in the array
    numberOfCalls := int(binary.BigEndian.Uint32(data[60:64]))
    swapTransactionResult.NumberOfCalls = numberOfCalls

    // Create a slice to store the offsets for each call
    callOffsets := make([]uint64, numberOfCalls)
    
    // Skip the first 64 bytes (array offset and length)
    data = data[64:]
    
    // Read the offsets for each call
    for i := 0; i < numberOfCalls; i++ {
        if len(data) < (i+1)*32 {
            return fmt.Errorf("data too short for offset %d", i)
        }
        offset := new(big.Int).SetBytes(data[i*32:(i+1)*32])
        callOffsets[i] = offset.Uint64()
    }

    // Now process each call's data
    if swapTransactionResult.CallsData == nil {
        swapTransactionResult.CallsData = make([]string, numberOfCalls)
    }

    for i := 0; i < numberOfCalls; i++ {
        startOffset := callOffsets[i]
        var endOffset uint64
        if i == numberOfCalls-1 {
            endOffset = uint64(len(data))
        } else {
            endOffset = callOffsets[i+1]
        }

        if startOffset < uint64(len(data)) && endOffset <= uint64(len(data)) {
            callData := data[startOffset:endOffset]
            
            // Read the length of this call's data
            callLength := new(big.Int).SetBytes(callData[0:32]).Uint64()
            
            // Skip the length prefix (32 bytes) and only take callLength bytes
            actualCallData := callData[32:32+callLength]
            
            swapTransactionResult.CallsData[i] = "0x" + common.Bytes2Hex(actualCallData)
        } else {
            return fmt.Errorf("invalid offset bounds: start=%d, end=%d, len=%d", startOffset, endOffset, len(data))
		}
	}

	return nil
}

/*
	https://dashboard.tenderly.co/tx/mainnet/0x4c431f853758f9f89e3b8eaa538ff7069e2fb0c53222db25e2fd93de344f781e?trace=0
	Function: exactOutput(tuple params)

	MethodID: 0xf28c0498
	[0]:  0000000000000000000000000000000000000000000000000000000000000020
	[1]:  00000000000000000000000000000000000000000000000000000000000000a0
	[2]:  0000000000000000000000007d14b142cad1379e85682f4b2006cdfed38988d3
	[3]:  0000000000000000000000000000000000000000000000000000000067576d64
	[4]:  00000000000000000000000000000000000000000000000000000002b8c73e80
	[5]:  00000000000000000000000000000000000000000000000000000001a9fa43ae
	[6]:  0000000000000000000000000000000000000000000000000000000000000042
	[7]:  77e06c9eccf2e797fd462a92b6d7642ef85b0a44000bb8c02aaa39b223fe8d0a
	[8]:  0e5c4f27ead9083c756cc20001f4a0b86991c6218b36c1d19d4a2e9eb0ce3606
	[9]:  eb48000000000000000000000000000000000000000000000000000000000000

	{
		"params": {
			"path": "0x77e06c9eccf2e797fd462a92b6d7642ef85b0a44000bb8c02aaa39b223fe8d0a0e5c4f27ead9083c756cc20001f4a0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
			"recipient": "0x7d14b142cad1379e85682f4b2006cdfed38988d3",
			"deadline": "1733782884",
			"amountOut": "11690000000",
			"amountInMaximum": "7146718126"
		}
	}

*/
func DecodeExactOutput(data []byte, swapTransactionResult *entity.SwapTransaction) error {

	data = data[4:]

	// [0] first offset to the the path offset
	firstOffset := new(big.Int).SetBytes(data[:32])

	// [1] recipient,deadline,amountOut,amountInMaximum offset
	secondOffset := new(big.Int).SetBytes(data[firstOffset.Uint64():firstOffset.Uint64() + 32])

	newstart := firstOffset.Uint64() + 32

	recipient := fmt.Sprintf("0x%s", common.Bytes2Hex(data[newstart:newstart + 32])[24:])
	// deadline := fmt.Sprintf("0x%s", common.Bytes2Hex(data[secondOffset.Uint64() + 32:secondOffset.Uint64() + 64])[24:])
	amountOut := new(big.Int).SetBytes(data[newstart + 64:newstart + 96])
	amountInMaximum := new(big.Int).SetBytes(data[newstart + 96:newstart + 128])

	swapTransactionResult.ToAddress = recipient
	// swapTransactionResult.Deadline = deadline
	swapTransactionResult.AmountOut = amountOut.String()
	swapTransactionResult.AmountInMax = amountInMaximum.String()

	// fmt.Println("recipient", recipient)
	// fmt.Println("amountOut", amountOut)
	// fmt.Println("amountInMaximum", amountInMaximum)

	newstart = secondOffset.Uint64() + 32

	// given the offsets, we can now get the next offset
	thirdOffset := new(big.Int).SetBytes(data[newstart:newstart + 32])

	// fmt.Println("firstOffset", firstOffset)
	// fmt.Println("secondOffset", secondOffset)
	// fmt.Println("thirdOffset", thirdOffset)


	newstart = newstart + 32

	// get the path
	path := common.Bytes2Hex(data[newstart: newstart + thirdOffset.Uint64()])
	// fmt.Printf("%s\n", path)
	// lets split the path into tokenIn and tokenOut
	// 0x77e06c9eccf2e797fd462a92b6d7642ef85b0a44000bb8c02aaa39b223fe8d0a0e5c4f27ead9083c756cc20001f4a0b86991c6218b36c1d19d4a2e9eb0ce3606eb48
	// tokenOut is the first 0x77e06c9eccf2e797fd462a92b6d7642ef85b0a44
	// tokenIn is the last 0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48

	tokenOut := fmt.Sprintf("0x%s", path[0:40])
	tokenIn := fmt.Sprintf("0x%s", path[len(path)-40:])

	swapTransactionResult.TokenPathTo = tokenOut
	swapTransactionResult.TokenPathFrom = tokenIn

	// fmt.Println("tokenIn:", tokenIn)
	// fmt.Println("tokenOut:", tokenOut)


	return nil
}


