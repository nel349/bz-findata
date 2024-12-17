package defi_llama

import (
	"fmt"
	"testing"
)

func TestGetTokenInfo(t *testing.T) {

	tokenInfo, err := GetTokenInfo("0xdF574c24545E5FfEcb9a659c229253D4111d87e1")
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
