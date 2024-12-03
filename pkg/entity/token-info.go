package entity

import "time"

type TokenInfo struct {
	Address  string  `db:"address"`
	Decimals   uint8   `db:"decimals"`
	Symbol     string  `db:"symbol"`
	Price      float64 `db:"price"`
	LastUpdated time.Time `db:"last_updated"`
}