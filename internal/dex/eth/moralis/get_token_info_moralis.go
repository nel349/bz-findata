package moralis

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/nel349/bz-findata/pkg/entity"
)

/*
	{
	"tokenName": "Omira",
	"tokenSymbol": "OMIRA",
	"tokenLogo": null,
	"tokenDecimals": "18",
	"nativePrice": {
		"value": "2377330701051",
		"decimals": 18,
		"name": "Ether",
		"symbol": "ETH",
		"address": "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"
	},
	"usdPrice": 0.008604496332884841,
	"usdPriceFormatted": "0.008604496332884841",
	"exchangeName": "Uniswap v2",
	"exchangeAddress": "0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f",
	"tokenAddress": "0x06113abcef9d163c026441b112e70c82ee1c4a79",
	"priceLastChangedAtBlock": "21324770",
	"blockTimestamp": "1733262623000",
	"possibleSpam": false,
	"verifiedContract": false,
	"pairAddress": "0x279ba98e72f5bea8eda92f1bf0a449c32b7c420f",
	"pairTotalLiquidityUsd": "179968.20",
	"24hrPercentChange": "40.94741756145716",
	"securityScore": null
	}

*/

type MoralisResponse struct {
	TokenName string `json:"tokenName"`
	TokenSymbol string `json:"tokenSymbol"`
	TokenDecimals string `json:"tokenDecimals"`
	NativePrice struct {
		Value string `json:"value"`
		Decimals uint8 `json:"decimals"`
	} `json:"nativePrice"`
	UsdPrice float64 `json:"usdPrice"`
	UsdPriceFormatted string `json:"usdPriceFormatted"`
}

func GetTokenInfoFromMoralis(tokenAddress string) (entity.TokenInfo, error) {
	url := fmt.Sprintf("https://deep-index.moralis.io/api/v2.2/erc20/%s/price?chain=eth&include=percent_change", tokenAddress)

	req, _ := http.NewRequest("GET", url, nil)
  
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-API-Key", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJub25jZSI6ImI2NDViMmE1LTk1M2MtNDIxMi05MTZlLWNlODBmMmZhZWJmYSIsIm9yZ0lkIjoiNDE5MTM5IiwidXNlcklkIjoiNDMxMDM2IiwidHlwZUlkIjoiNmE5YTQ0ODYtYmFiYy00ZDA2LTlhNTItMTYxZWM5YWY4OWY4IiwidHlwZSI6IlBST0pFQ1QiLCJpYXQiOjE3MzMxNjgyMzgsImV4cCI6NDg4ODkyODIzOH0.hI-6bBK6JLZJjefKWVxM0mXSdWeDJiL2EK1w-XtwEN8")
  
	res, _ := http.DefaultClient.Do(req)

	if res.StatusCode != 200 {
		return entity.TokenInfo{}, fmt.Errorf("failed to get token info from Moralis")
	}
  
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	
	var response MoralisResponse
	err := json.Unmarshal(body, &response)
	if err != nil {
		return entity.TokenInfo{}, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Extract data from the response
	// key := fmt.Sprintf("ethereum:%s", tokenAddress)
	
	return entity.TokenInfo{
		Address:  tokenAddress,
		Decimals: response.NativePrice.Decimals,
		Symbol:   response.TokenSymbol,
		Price:    response.UsdPrice,
	}, nil
}
