package v2

import "strings"

// SwapMethod represents different Uniswap swap method types
type UniswapV2SwapMethod string

const (
    SwapExactTokensForTokens                              UniswapV2SwapMethod = "38ed1739"
    SwapTokensForExactTokens                             UniswapV2SwapMethod = "8803dbee"
    SwapExactETHForTokens                                UniswapV2SwapMethod = "7ff36ab5"
    SwapTokensForExactETH                                UniswapV2SwapMethod = "4a25d94a"
    SwapExactTokensForETH                                UniswapV2SwapMethod = "18cbafe5"
    SwapETHForExactTokens                                UniswapV2SwapMethod = "fb3bdb41"
    SwapExactTokensForTokensSupportingFeeOnTransferTokens UniswapV2SwapMethod = "5c11d795"
    SwapExactETHForTokensSupportingFeeOnTransferTokens    UniswapV2SwapMethod = "b6f9de95"
    SwapExactTokensForETHSupportingFeeOnTransferTokens    UniswapV2SwapMethod = "791ac947"
)

// GetV2MethodFromID returns the UniswapV2SwapMethod for a given method ID
func GetV2MethodFromID(methodID string) (UniswapV2SwapMethod, bool) {
    // Remove 0x prefix if present
    methodID = strings.TrimPrefix(methodID, "0x")
    
    // Check each method
    switch methodID {
    case string(SwapExactTokensForTokens):
        return SwapExactTokensForTokens, true
    case string(SwapTokensForExactTokens):
        return SwapTokensForExactTokens, true
    case string(SwapExactETHForTokens):
        return SwapExactETHForTokens, true
    case string(SwapTokensForExactETH):
        return SwapTokensForExactETH, true
    case string(SwapExactTokensForETH):
        return SwapExactTokensForETH, true
    case string(SwapETHForExactTokens):
        return SwapETHForExactTokens, true
    case string(SwapExactTokensForTokensSupportingFeeOnTransferTokens):
        return SwapExactTokensForTokensSupportingFeeOnTransferTokens, true
    case string(SwapExactETHForTokensSupportingFeeOnTransferTokens):
        return SwapExactETHForTokensSupportingFeeOnTransferTokens, true
    case string(SwapExactTokensForETHSupportingFeeOnTransferTokens):
        return SwapExactTokensForETHSupportingFeeOnTransferTokens, true
    default:
        return "", false
    }
}

// String returns the human-readable name of the swap method
func (s UniswapV2SwapMethod) String() string {
    switch s {
    case SwapExactTokensForTokens:
        return "SwapExactTokensForTokens"
    case SwapTokensForExactTokens:
        return "SwapTokensForExactTokens"
    case SwapExactETHForTokens:
        return "SwapExactETHForTokens"
    case SwapTokensForExactETH:
        return "SwapTokensForExactETH"
    case SwapExactTokensForETH:
        return "SwapExactTokensForETH"
    case SwapETHForExactTokens:
        return "SwapETHForExactTokens"
    case SwapExactTokensForTokensSupportingFeeOnTransferTokens:
        return "SwapExactTokensForTokensSupportingFeeOnTransferTokens"
    case SwapExactETHForTokensSupportingFeeOnTransferTokens:
        return "SwapExactETHForTokensSupportingFeeOnTransferTokens"
    case SwapExactTokensForETHSupportingFeeOnTransferTokens:
        return "SwapExactTokensForETHSupportingFeeOnTransferTokens"
    default:
        return "Unknown"
    }
}

// IsETHInput returns true if the method takes ETH as input
func (s UniswapV2SwapMethod) IsETHInput() bool {
    switch s {
    case SwapExactETHForTokens,
         SwapETHForExactTokens,
         SwapExactETHForTokensSupportingFeeOnTransferTokens:
        return true
    default:
        return false
    }
}

// IsETHOutput returns true if the method outputs ETH
func (s UniswapV2SwapMethod) IsETHOutput() bool {
    switch s {
    case SwapTokensForExactETH,
         SwapExactTokensForETH,
         SwapExactTokensForETHSupportingFeeOnTransferTokens:
        return true
    default:
        return false
    }
}