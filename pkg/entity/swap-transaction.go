package entity

type SwapTransaction struct {
	Value        float64 `db:"value"`
	TxHash       string  `db:"tx_hash"`
	Version      string  `db:"version"`
	Exchange     string  `db:"exchange"`
	AmountIn     string `db:"amount_in"`
	AmountOutMin string `db:"amount_out_min"`
	// Deadline      string `db:"deadline"`
	ToAddress     string `db:"to_address"`
	TokenPathFrom string `db:"token_path_from"`
	TokenPathTo   string `db:"token_path_to"`
	MethodID      string `db:"method_id"`
	MethodName    string `db:"method_name"`

	// Uniswap V2 add liquidity
	AmountTokenDesired string `db:"amount_token_desired"`
	AmountTokenMin     string `db:"amount_token_min"`
	AmountETHMin       string `db:"amount_eth_min"`

	// Uniswap V2 remove liquidity
	Liquidity          string  `db:"liquidity"`
	AmountAMin         string  `db:"amount_a_min"`
	AmountBMin         string  `db:"amount_b_min"`

	// v2 swapTokensForExactTokens
	AmountOut        string `db:"amount_out"`
	AmountInMax      string `db:"amount_in_max"`

	// Uniswap V3
	Fee string `db:"fee"`

	// Uniswap V3 Multicall
	NumberOfCalls int `db:"number_of_calls"`
	CallsData []string 
}
