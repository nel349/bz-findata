package moralis

import (
	"fmt"
	"log"
	"testing"

	"github.com/joho/godotenv"
)

func TestGetTokenInfo(t *testing.T) {

	// load the test env file
	if err := godotenv.Load("../../../../test.env"); err != nil {
		log.Fatal("Error loading test.env file")
	}

	tokenInfo, err := GetTokenInfoFromMoralis("0x06113abcef9d163c026441b112e70c82ee1c4a79")
	if err != nil {
		t.Errorf("failed to get token info: %v", err)
	}

	// check the token info fields are not empty
	if tokenInfo.Address == "" || tokenInfo.Symbol == "" || tokenInfo.Decimals == 0 || tokenInfo.Price == 0 {
		t.Errorf("token info fields are empty")
	}

	// check the following contract address is correct
	if tokenInfo.Address != "0x06113abcef9d163c026441b112e70c82ee1c4a79" {
		t.Errorf("token address is incorrect")
	}

	// check the following symbol is correct
	if tokenInfo.Symbol != "OMIRA" {
		t.Errorf("token symbol is incorrect")
	}

	// check the following decimals is correct
	if tokenInfo.Decimals != 18 {
		t.Errorf("token decimals is incorrect")
	}

	fmt.Println(tokenInfo)
	fmt.Printf("%.9f\n", tokenInfo.Price)
	fmt.Println(tokenInfo.Symbol)
	fmt.Println(tokenInfo.Decimals)
}
