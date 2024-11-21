package entity

type SwapTransaction struct {
	Value        float64 `db:"value"`
	TxHash       string  `db:"tx_hash"`
	Version      string  `db:"version"`
	Exchange     string  `db:"exchange"`
	AmountIn     float64 `db:"amount_in"`
	AmountOutMin float64 `db:"amount_out_min"`
	// Deadline      string `db:"deadline"`
	ToAddress     string `db:"to_address"`
	TokenPathFrom string `db:"token_path_from"`
	TokenPathTo   string `db:"token_path_to"`
	MethodID      string `db:"method_id"`
	MethodName    string `db:"method_name"`

	// Uniswap V2 add liquidity
	AmountTokenDesired float64 `db:"amount_token_desired"`
	AmountTokenMin     float64 `db:"amount_token_min"`
	AmountETHMin       float64 `db:"amount_eth_min"`

	Liquidity          string  `db:"liquidity"`

	// Uniswap V3
	Fee string `db:"fee"`
}
