package decoder

import (
	"math"
	"math/big"
)

func ConvertToFloat64(str string) float64 {
	f := new(big.Float)
	f.SetString(str)
	result, _ := f.Float64()
	return result
}

func ConvertToBigInt(str string) *big.Int {
	f := new(big.Int)
	f.SetString(str, 10)
	return f
}

func GetEthValue(value *big.Int) float64 {
	result, _ := new(big.Float).Quo(
		new(big.Float).SetInt(value),
		new(big.Float).SetFloat64(1e18),
	).Float64()

	return result
}

func GetUsdValueFromEth(value *big.Int, ethPrice float64) float64 {

	tokenValue := new(big.Float).SetInt(value)
	scale := big.NewFloat(math.Pow(10, 18)) // 18 decimals for ETH
	tokenValue = tokenValue.Quo(tokenValue, scale)

	result, _ := new(big.Float).Mul(
		tokenValue,
		big.NewFloat(ethPrice),
	).Float64()

	return result
}

func GetUsdValueFromToken(value *big.Int, tokenPrice float64, decimals int) float64 {
	tokenValue := new(big.Float).SetInt(value)
	scale := big.NewFloat(math.Pow(10, float64(decimals)))
	tokenValue = tokenValue.Quo(tokenValue, scale)

	result, _ := new(big.Float).Mul(
        tokenValue,
        big.NewFloat(tokenPrice),
    ).Float64()

    return result
}

// GetTokenAmountFromRaw converts a raw token amount to a string with the token amount in the correct format
func GetTokenAmountFromRawToBigFloat(value string, decimals int) *big.Float {
	tokenValue, _ := new(big.Float).SetString(value)
	scale := big.NewFloat(math.Pow(10, float64(decimals)))
	tokenValue = tokenValue.Quo(tokenValue, scale)
	return tokenValue
}

func GetTokenFromRawToUsd(value string, tokenPrice float64, decimals int) float64 {
	amountInFloat := GetTokenAmountFromRawToBigFloat(value, decimals)
	tokenPriceFloat := new(big.Float).SetFloat64(tokenPrice)
	result, _ := amountInFloat.Mul(amountInFloat, tokenPriceFloat).Float64()
	return result
}