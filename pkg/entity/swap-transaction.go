package entity

type SwapTransaction struct {
	Value         float64 `db:"value"`
	TxHash        string `db:"tx_hash"`
	Version       string `db:"version"`
	Exchange      string `db:"exchange"`
	AmountIn      float64 `db:"amount_in"`
	AmountOutMin  float64 `db:"amount_out_min"`
	// Deadline      string `db:"deadline"`
	ToAddress     string `db:"to_address"`
	TokenPathFrom string `db:"token_path_from"`
	TokenPathTo   string `db:"token_path_to"`
}
