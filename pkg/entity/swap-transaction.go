package entity

type SwapTransaction struct {
	TxHash        string `db:"tx_hash"`
	Version       string `db:"version"`
	Exchange      string `db:"exchange"`
	AmountIn      string `db:"amount_in"`
	ToAddress     string `db:"to_address"`
	// Deadline      string `db:"deadline"`
	TokenPathFrom string `db:"token_path_from"`
	TokenPathTo   string `db:"token_path_to"`
}
