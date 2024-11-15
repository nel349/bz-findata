package decoder

import (
	// "errors"
	"fmt"
	"math"
	"math/big"

	"github.com/jmoiron/sqlx"
	"github.com/nel349/bz-findata/internal/dex/eth/defi_llama"
	"github.com/nel349/bz-findata/internal/dex/eth/uniswap/v2"
	"github.com/nel349/bz-findata/internal/dex/eth/uniswap/v3"
	"github.com/nel349/bz-findata/pkg/entity"
)

func DecodeSwap(data []byte, version string, db *sqlx.DB, tokenInfo *defi_llama.TokenInfo) (*entity.SwapTransaction, error) {
	methodID := fmt.Sprintf("%x", data[:4])

	// Debug prints
	// fmt.Println("Raw data length:", len(data))
	// fmt.Printf("Raw data hex: 0x%x\n", data)
	//

	var swapMethod interface{}
	var ok bool
	if version == "V2" {
		swapMethod, ok = v2.GetV2MethodFromID(methodID)
	} else {
		swapMethod, ok = v3.GetV3MethodFromID(methodID)
	}
	if !ok {
		return nil, fmt.Errorf("unknown swap method: %s", methodID)
	}

	switch swapMethod := swapMethod.(type) {
	case v2.UniswapV2SwapMethod:
		fmt.Println("V2 swap method")
        // Lets do a switch for all the v2 swap methods
        switch swapMethod {
            case v2.SwapExactTokensForTokens:
                return DecodeSwapExactTokensForTokens(data, version, tokenInfo)
            case v2.SwapTokensForExactTokens:
                return DecodeSwapExactTokensForTokens(data, version, tokenInfo)
            // Add more cases for other v2 swap methods as needed
        }

	case v3.UniswapV3Method:
		fmt.Println("V3 swap method")
	}

	// Debug prints
	fmt.Printf("Swap Method: %s\n", swapMethod)
	fmt.Println("First 4 bytes (method signature):", methodID)
	return &entity.SwapTransaction{}, nil
}

/** function for SwapExactTokensForTokens
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
func DecodeSwapExactTokensForTokens(data []byte, version string, tokenInfo *defi_llama.TokenInfo) (*entity.SwapTransaction, error) {
    // Skip first 4 bytes (method ID)
    data = data[4:]
    
    // [0] amountIn (uint256)
	fmt.Printf("Raw amountIn bytes: %x\n", data[:32])
    amountIn := new(big.Int).SetBytes(data[:32])
	fmt.Printf("Amount In: %v\n", amountIn)

    amountInDecimal := new(big.Float).Quo(
        new(big.Float).SetInt(amountIn),
        new(big.Float).SetFloat64(math.Pow(10, float64(tokenInfo.Decimals))),
    )
	amountInFloat64, _ := amountInDecimal.Float64()
    
    // [1] amountOutMin (uint256)
    amountOutMin := new(big.Int).SetBytes(data[32:64])
	amountOutMinFloat64, _ := new(big.Float).Quo(
		new(big.Float).SetInt(amountOutMin),
		new(big.Float).SetFloat64(math.Pow(10, float64(tokenInfo.Decimals))),
	).Float64()
    
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
    tokenPathFrom := fmt.Sprintf("%x", data[192:224])[24:] // First token in path
    tokenPathTo := fmt.Sprintf("%x", data[224:256])[24:]   // Second token in path

    fmt.Printf("Amount In: %v tokens\n", amountInDecimal)
    fmt.Printf("To Address: 0x%s\n", toAddress)
    fmt.Printf("Token Path:\n")
    fmt.Printf("  From: 0x%s\n", tokenPathFrom)
    fmt.Printf("  To: 0x%s\n", tokenPathTo)

    return &entity.SwapTransaction{
        AmountIn:      amountInFloat64,
        AmountOutMin:  amountOutMinFloat64,
        ToAddress:     toAddress,
        // Deadline:      deadline.String(),
        TokenPathFrom: tokenPathFrom,
        TokenPathTo:   tokenPathTo,
    }, nil
}