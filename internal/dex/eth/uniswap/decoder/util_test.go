package decoder

import (
	"math"
	"math/big"
	"testing"
)

func TestGetUsdValueFromToken(t *testing.T) {
	amountIn, ok := new(big.Int).SetString("91745882892710", 10) // 	
	if !ok {
		t.Fatal("failed to parse big.Int")
	}

	tokenPrice := 0.026534231
	decimals := 9

	result := GetUsdValueFromToken(amountIn, tokenPrice, decimals)

	expected := 2434.40644997 // 91745.882892710 * 0.026534231

	// lets do some rounding
	if math.Abs(result-expected) > 0.00001 {
		t.Errorf("Result does not match expected value %v, got %v", expected, result)
	}
}

func TestGetUsdValueFromEth(t *testing.T) {
	amountIn, ok := new(big.Int).SetString("1000000000000000000", 10)
	if !ok {
		t.Fatal("failed to parse big.Int")
	}

	ethPrice := 1949.99

	result := GetUsdValueFromEth(amountIn, ethPrice)

	expected := 1949.99

	if result != expected {
		t.Errorf("Result does not match expected value %v, got %v", expected, result)
	}
}

// func TestGetTokenFromRawToUsd(t *testing.T) {
// 	amountIn := "1000000000000000000"
// 	tokenPrice := 0.026534231
// 	decimals := 9

// 	result := GetTokenFromRawToUsd(amountIn, tokenPrice, decimals)

// 	expected := 26.534231

// 	if result != expected {
// 		t.Errorf("Result does not match expected value %v, got %v", expected, result)
// 	}
// }
