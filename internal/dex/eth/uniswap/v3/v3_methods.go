package v3

import "strings"

// UniswapV3Method represents different Uniswap V3 swap method types
type UniswapV3Method string

const (
    // Core swap methods
    ExactInput         UniswapV3Method = "c04b8d59"
    ExactInputSingle   UniswapV3Method = "414bf389"
    ExactOutput        UniswapV3Method = "f28c0498"
    ExactOutputSingle  UniswapV3Method = "db3e2198"
    
    // Multicall methods often used for swaps
    Multicall          UniswapV3Method = "ac9650d8"
    MulticallWithValue UniswapV3Method = "5ae401dc"
)

// GetV3MethodFromID returns the UniswapV3Method for a given method ID
func GetV3MethodFromID(methodID string) (UniswapV3Method, bool) {
    methodID = strings.TrimPrefix(methodID, "0x")
    
    switch methodID {
    case string(ExactInput):
        return ExactInput, true
    case string(ExactInputSingle):
        return ExactInputSingle, true
    case string(ExactOutput):
        return ExactOutput, true
    case string(ExactOutputSingle):
        return ExactOutputSingle, true
    case string(Multicall):
        return Multicall, true
    case string(MulticallWithValue):
        return MulticallWithValue, true
    default:
        return "", false
    }
}

// String returns the human-readable name of the swap method
func (m UniswapV3Method) String() string {
    switch m {
    case ExactInput:
        return "ExactInput"
    case ExactInputSingle:
        return "ExactInputSingle"
    case ExactOutput:
        return "ExactOutput"
    case ExactOutputSingle:
        return "ExactOutputSingle"
    case Multicall:
        return "Multicall"
    case MulticallWithValue:
        return "MulticallWithValue"
    default:
        return "Unknown"
    }
}

// IsMulticall returns true if the method is a multicall variant
func (m UniswapV3Method) IsMulticall() bool {
    return m == Multicall || m == MulticallWithValue
}

// RequiresValue returns true if the method accepts ETH value
func (m UniswapV3Method) RequiresValue() bool {
    switch m {
    case MulticallWithValue:
        return true
    case ExactInput, ExactInputSingle, ExactOutput, ExactOutputSingle:
        // These can accept value depending on the path
        return true
    default:
        return false
    }
}