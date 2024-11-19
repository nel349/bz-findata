package decoder

import "math/big"

func ConvertToFloat64(str string) float64 {
    f := new(big.Float)
    f.SetString(str)
    result, _ := f.Float64()
    return result
}

func GetEthValue(value *big.Int) float64 {
	result, _ := new(big.Float).Quo(
		new(big.Float).SetInt(value),
		new(big.Float).SetFloat64(1e18),
	).Float64()

	return result
}
