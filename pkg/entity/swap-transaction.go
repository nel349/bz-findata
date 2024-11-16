package entity

import "math/big"

type SwapTransaction struct {
	Value         float64 `db:"value"`
	TxHash        string `db:"tx_hash"`
	Version       string `db:"version"`
	Exchange      string `db:"exchange"`
	AmountIn      *big.Int `db:"amount_in"`
	AmountOutMin  *big.Int `db:"amount_out_min"`
	// Deadline      string `db:"deadline"`
	ToAddress     string `db:"to_address"`
	TokenPathFrom string `db:"token_path_from"`
	TokenPathTo   string `db:"token_path_to"`
}
