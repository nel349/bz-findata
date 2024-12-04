package defi_llama

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/nel349/bz-findata/internal/dex/eth/moralis"
	"github.com/nel349/bz-findata/pkg/entity"
)

/**
	Get token info from defi llama api
	https://coins.llama.fi/prices/current/ethereum:0x699Ec925118567b6475Fe495327ba0a778234AaA?searchWidth=4h
	Response:
{
  "coins": {
    "ethereum:0x699Ec925118567b6475Fe495327ba0a778234AaA": {
      "decimals": 9,
      "price": 0.00003981222012071399,
      "symbol": "DUCKY",
      "timestamp": 1731647306,
      "confidence": 0.98
    }
  }
}
*/

type DefiLlamaResponse struct {
	Coins map[string]struct {
		Decimals   uint8   `json:"decimals"`
		Price      float64 `json:"price"`
		Symbol     string  `json:"symbol"`
		Timestamp  int64   `json:"timestamp"`
		Confidence float64 `json:"confidence"`
	} `json:"coins"`
}

func GetTokenInfo(tokenAddress string) (entity.TokenInfo, error) {

	url := fmt.Sprintf("https://coins.llama.fi/prices/current/ethereum:%s?searchWidth=4h", tokenAddress)

	// Make the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return entity.TokenInfo{}, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close() // Ensure the response body is closed

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return entity.TokenInfo{}, fmt.Errorf("failed to fetch data: %s", resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return entity.TokenInfo{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var response DefiLlamaResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return entity.TokenInfo{}, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Extract data from the response
	key := fmt.Sprintf("ethereum:%s", tokenAddress)
	if tokenData, exists := response.Coins[key]; exists {
		return entity.TokenInfo{
			Address:  tokenAddress,
			Decimals: tokenData.Decimals,
			Symbol:   tokenData.Symbol,
			Price:    tokenData.Price,
		}, nil
	}

	return entity.TokenInfo{}, fmt.Errorf("token data not found in response")
}

func GetTokenMetadataFromDbOrDefiLlama(db *sqlx.DB, tokenAddress string, updateInterval time.Duration) (entity.TokenInfo, error) {
	// Try to get from database first
	var tokenInfo entity.TokenInfo
	err := db.Get(&tokenInfo, "SELECT * FROM token_metadata WHERE BINARY address = ?", strings.ToLower(tokenAddress))
	if err == nil {
        timeSinceUpdate := time.Since(tokenInfo.LastUpdated)
        // fmt.Printf("Debug: Token %s last updated: %v\n", tokenAddress, tokenInfo.LastUpdated)
        // fmt.Printf("Debug: Time since last update: %v\n", timeSinceUpdate)
        // fmt.Printf("Debug: Update interval: %v\n", updateInterval)
        // fmt.Printf("Debug: Is stale? %v\n", timeSinceUpdate >= updateInterval)
        
        if timeSinceUpdate < updateInterval {
            fmt.Println("Found token metadata in database: ", tokenInfo)
            return tokenInfo, nil
        }
		fmt.Println("Token metadata is stale, updating from API with address: ", strings.ToLower(tokenAddress))
	} else {
		fmt.Println("Token metadata not found in database, fetching from API with address: ", strings.ToLower(tokenAddress))
		fmt.Println("Error: ", err)
	}

	// Fetch from defi llama api if data is stale or not found
	tokenInfo, err = GetTokenInfo(tokenAddress)
	if err != nil {
		// fetch from moralis as fallback
		tokenInfo, err = moralis.GetTokenInfoFromMoralis(tokenAddress)
		if err != nil {
			return entity.TokenInfo{},
				fmt.Errorf("failed to get token %s info from moralis: %w", tokenAddress, err)
		}

		fmt.Println("Fetched token metadata from moralis: ", tokenInfo)
	}

	// Store in database
	_, err = db.Exec(`
    INSERT INTO token_metadata (address, decimals, symbol, price, last_updated) 
    VALUES (?, ?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE
		address = VALUES(address),
		decimals = VALUES(decimals),
		symbol = VALUES(symbol),
		price = VALUES(price),
		last_updated = VALUES(last_updated)
	`,
		strings.ToLower(tokenInfo.Address), tokenInfo.Decimals, tokenInfo.Symbol, tokenInfo.Price, time.Now(),
	)
	if err != nil {
		return entity.TokenInfo{}, fmt.Errorf("failed to store token metadata: %w", err)
	}
	// fmt.Println("Stored or updated token metadata in database: ", tokenInfo)
	return tokenInfo, nil
}

func GetWETHPrice(db *sqlx.DB) (float64, error) {
	tokenInfo, err := GetTokenMetadataFromDbOrDefiLlama(db, "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2", 15*time.Minute)
	if err != nil {
		return 0, fmt.Errorf("failed to get WETH price: %w", err)
	}

	return tokenInfo.Price, nil
}
