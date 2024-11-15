package token

import (
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/nel349/bz-findata/internal/dex/eth/defi_llama"
)

type TokenMetadata struct {
	Address     string    `db:"address"`
	Decimals    uint8     `db:"decimals"`
	Symbol      string    `db:"symbol"`
	LastUpdated time.Time `db:"last_updated"`
}

func GetTokenDecimals(db *sqlx.DB, tokenInfo *defi_llama.TokenInfo) (uint8, error) {
	// Try to get from database first
	var metadata TokenMetadata
	err := db.Get(&metadata, "SELECT * FROM token_metadata WHERE address = ?", strings.ToLower(tokenInfo.Address))
	if err == nil {
		return metadata.Decimals, nil
	}

	// If not in database, fetch from contract and store
	decimals := uint8(tokenInfo.Decimals) // Default value
	// TODO: Get from defi llama api

	// Store in database
	_, err = db.Exec(`
        INSERT INTO token_metadata (address, decimals, symbol) 
        VALUES (?, ?, ?) 
        ON DUPLICATE KEY UPDATE decimals = ?, symbol = ?`,
		strings.ToLower(tokenInfo.Address), decimals, tokenInfo.Symbol, decimals, tokenInfo.Symbol,
	)
	if err != nil {
		return decimals, fmt.Errorf("failed to store token decimals: %w", err)
	}

	return decimals, nil
}
