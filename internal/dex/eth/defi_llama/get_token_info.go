package defi_llama

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

type TokenInfo struct {
	Address  string
	Decimals uint8
	Symbol   string
	Price    float64
}

func GetTokenInfo(tokenAddress string) (TokenInfo, error) {

	url := fmt.Sprintf("https://coins.llama.fi/prices/current/ethereum:%s?searchWidth=4h", tokenAddress)

	// Make the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return TokenInfo{}, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close() // Ensure the response body is closed

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return TokenInfo{}, fmt.Errorf("failed to fetch data: %s", resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TokenInfo{}, fmt.Errorf("failed to read response body: %w", err)
	}

    var response DefiLlamaResponse
    err = json.Unmarshal(body, &response)
    if err != nil {
        return TokenInfo{}, fmt.Errorf("failed to parse JSON: %w", err)
    }

    // Extract data from the response
    key := fmt.Sprintf("ethereum:%s", tokenAddress)
    if tokenData, exists := response.Coins[key]; exists {
        return TokenInfo{
            Address:  tokenAddress,
            Decimals: tokenData.Decimals,
            Symbol:   tokenData.Symbol,
            Price:    tokenData.Price,
        }, nil
    }

	return TokenInfo{}, fmt.Errorf("token data not found in response")
}
