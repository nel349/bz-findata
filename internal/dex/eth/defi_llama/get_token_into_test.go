package defi_llama

import (
	"fmt"
	"testing"
)

func TestGetTokenInfo(t *testing.T) {

	tokenInfo, err := GetTokenInfo("0x699Ec925118567b6475Fe495327ba0a778234AaA")
	if err != nil {
		t.Errorf("failed to get token info: %v", err)
	}

	// check the token info fields are not empty
	if tokenInfo.Address == "" || tokenInfo.Symbol == "" || tokenInfo.Decimals == 0 || tokenInfo.Price == 0 {
		t.Errorf("token info fields are empty")
	}

	fmt.Println(tokenInfo)
	fmt.Printf("%.9f\n", tokenInfo.Price)
	fmt.Println(tokenInfo.Symbol)
	fmt.Println(tokenInfo.Decimals)
}
