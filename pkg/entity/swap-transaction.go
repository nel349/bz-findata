package entity

type SwapTransaction struct {
	Value        float64 `json:"value" db:"value"`
	TxHash       string  `json:"tx_hash" db:"tx_hash"`
	Version      string  `json:"version" db:"version"`
	Exchange     string  `json:"exchange" db:"exchange"`
	AmountIn     string  `json:"amount_in" db:"amount_in"`
	AmountOutMin string  `json:"amount_out_min" db:"amount_out_min"`
	// Deadline      string `json:"deadline" db:"deadline"`
	ToAddress     string `json:"to_address" db:"to_address"`
	TokenPathFrom string `json:"token_path_from" db:"token_path_from"`
	TokenPathTo   string `json:"token_path_to" db:"token_path_to"`
	MethodID      string `json:"method_id" db:"method_id"`
	MethodName    string `json:"method_name" db:"method_name"`
	LastUpdated   string `json:"last_updated" db:"last_updated"`

	// Uniswap V2 add/remove liquidity
	AmountTokenDesired string `json:"amount_token_desired" db:"amount_token_desired"`
	AmountTokenMin     string `json:"amount_token_min" db:"amount_token_min"`
	AmountETHMin       string `json:"amount_eth_min" db:"amount_eth_min"`
	TokenA             string `json:"token_a" db:"token_a"`
	TokenB             string `json:"token_b" db:"token_b"`
	AmountADesired     string `json:"amount_a_desired" db:"amount_a_desired"`
	AmountBDesired     string `json:"amount_b_desired" db:"amount_b_desired"`
	Liquidity          string `json:"liquidity" db:"liquidity"`
	AmountAMin         string `json:"amount_a_min" db:"amount_a_min"`
	AmountBMin         string `json:"amount_b_min" db:"amount_b_min"`

	// v2 swapTokensForExactTokens
	AmountOut   string `json:"amount_out" db:"amount_out"`
	AmountInMax string `json:"amount_in_max" db:"amount_in_max"`

	// Uniswap V3
	Fee string `json:"fee" db:"fee"`

	// Uniswap V3 Multicall
	NumberOfCalls int      `json:"number_of_calls,omitempty" db:"-"`
	CallsData     []string `json:"calls_data,omitempty" db:"-"`
}
