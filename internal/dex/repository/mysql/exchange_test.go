package mysql

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver
	"github.com/jmoiron/sqlx"
	"github.com/nel349/bz-findata/pkg/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewDexExchangeRepository(t *testing.T) {
	type args struct {
		db *sqlx.DB
	}
	tests := []struct {
		name string
		args args
		want *dexExchangeRepo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDexExchangeRepository(tt.args.db)
			if got.db != tt.want.db {
				t.Errorf("NewDexExchangeRepository() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func setupTestDB(t *testing.T) *sqlx.DB {
	// Connect to MySQL server without specifying a database
	dsn := "root:root@tcp(localhost:3306)/?parseTime=true"
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		t.Fatalf("Failed to connect to MySQL server: %v Did you start the mysql server?", err)
	}

	// Create a new test database
	dbName := "testdb"
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Connect to the test database
	dsnWithDB := fmt.Sprintf("root:root@tcp(localhost:3306)/%s?parseTime=true", dbName)
	db, err = sqlx.Connect("mysql", dsnWithDB)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Create the necessary tables
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS swap_transactions (
			tx_hash varchar(66) NOT NULL,
			version varchar(8) NOT NULL,
			exchange varchar(100) NOT NULL,
			amount_in varchar(100) NOT NULL,
			to_address varchar(42) NOT NULL,
			token_path_from varchar(42) NOT NULL,
			token_path_to varchar(42) NOT NULL,
			value float NOT NULL DEFAULT 0,
			amount_token_desired varchar(100) NULL,
			amount_token_min varchar(100) NULL,
			amount_eth_min varchar(100) NULL,
			method_id varchar(10) NULL,
			method_name varchar(100) NULL,
			liquidity varchar(100) NULL,
			token_a varchar(42) NULL,
			token_b varchar(42) NULL,
			amount_a_desired varchar(100) NULL,
			amount_b_desired varchar(100) NULL,
			amount_a_min varchar(100) NULL,
			amount_b_min varchar(100) NULL,
			PRIMARY KEY (tx_hash)
		) ENGINE=InnoDB;
	`)
	if err != nil {
		t.Fatalf("Failed to create tables: %v", err)
	}

	// Lets create a table for token info
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS token_metadata (
			address varchar(42) NOT NULL,
			symbol varchar(100) NOT NULL,
			decimals int NOT NULL,
			price float NOT NULL,
			last_updated datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (address)
		) ENGINE=InnoDB;
	`)
	if err != nil {
		t.Fatalf("Failed to create tables: %v", err)
	}

	// Cleanup: Drop the database after the test
	t.Cleanup(func() {
		db.Exec("DROP DATABASE " + dbName)
		db.Close()
	})

	return db
}

/*
Test Save Swap - RemoveLiquidity

https://dashboard.tenderly.co/tx/mainnet/0x1141dee16423b087413a9dff4635f00ce5067b54cdd736d89e9d7f2ce5916fa8

	{
		"tokenA": "0x6b175474e89094c44da98b954eedeac495271d0f",
		"tokenB": "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		"liquidity": "5688426355532103909",
		"amountAMin": "702110776270740190514",
		"amountBMin": "179779514399370049",
		"to": "0xc47e5d32f7be0cc171740ebbb3f26f78488cd22f",
		"deadline": "1734127406"
	}
*/
func Test_dexExchangeRepo_SaveSwap(t *testing.T) {
	// Setup a test database connection
	db := setupTestDB(t)

	data := common.FromHex("0xbaa2abde0000000000000000000000006b175474e89094c44da98b954eedeac495271d0f000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc20000000000000000000000000000000000000000000000004ef159a1bc7600e50000000000000000000000000000000000000000000000260fbe8f136f3af532000000000000000000000000000000000000000000000000027eb4840daaab41000000000000000000000000c47e5d32f7be0cc171740ebbb3f26f78488cd22f00000000000000000000000000000000000000000000000000000000675caf2e")

	// Create a mock transaction
	tx := types.NewTransaction(
		0, // nonce
		common.HexToAddress("0xdc0488a855ca075c5f7c6b9567fa0b84a9d97ffd5bb4bea913e9840a402b5b79"), // to address
		big.NewInt(0), // value
		0,             // gas limit
		big.NewInt(0), // gas price
		data,          // data
	)

	fmt.Println("tx: ", tx.Hash().Hex())

	// Create a mock SwapTransaction entity
	swapTransaction := &entity.SwapTransaction{
		TxHash:     tx.Hash().Hex(),
		TokenA:     "0x6b175474e89094c44da98b954eedeac495271d0f",
		TokenB:     "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		Liquidity:  "5688426355532103909",
		AmountAMin: "702110776270740190514",
		AmountBMin: "179779514399370049",
	}

	// Create a dexExchangeRepo instance
	repo := &dexExchangeRepo{db: db}

	// Insert the swap transaction
	err := repo.SaveSwap(context.Background(), tx, "V2")
	if err != nil {
		t.Errorf("Failed to insert swap transaction: %v", err)
	}

	// Verify the insertion
	var insertedSwap entity.SwapTransaction
	err = db.Get(&insertedSwap, "SELECT * FROM swap_transactions WHERE tx_hash = ?", swapTransaction.TxHash)
	if err != nil {
		t.Errorf("Failed to retrieve inserted swap transaction: %v", err)
	}

	fmt.Println("Token A: ", insertedSwap.TokenA)
	fmt.Println("Token B: ", insertedSwap.TokenB)
	fmt.Println("Liquidity: ", insertedSwap.Liquidity)
	fmt.Println("Amount A Min: ", insertedSwap.AmountAMin)
	fmt.Println("Amount B Min: ", insertedSwap.AmountBMin)
	fmt.Println("To Address: ", insertedSwap.ToAddress)
	fmt.Println("Tx Hash: ", insertedSwap.TxHash)

	// method name
	fmt.Println("Method Name: ", insertedSwap.MethodName)

	if insertedSwap.TokenA != swapTransaction.TokenA {
		t.Errorf("Token A does not match expected value %v, got %v", swapTransaction.TokenA, insertedSwap.TokenA)
	}

	if insertedSwap.TokenB != swapTransaction.TokenB {
		t.Errorf("Token B does not match expected value %v, got %v", swapTransaction.TokenB, insertedSwap.TokenB)
	}

	if insertedSwap.Liquidity != swapTransaction.Liquidity {
		t.Errorf("Liquidity does not match expected value %v, got %v", swapTransaction.Liquidity, insertedSwap.Liquidity)
	}

	if insertedSwap.AmountAMin != swapTransaction.AmountAMin {
		t.Errorf("Amount A Min does not match expected value %v, got %v", swapTransaction.AmountAMin, insertedSwap.AmountAMin)
	}

	if insertedSwap.AmountBMin != swapTransaction.AmountBMin {
		t.Errorf("Amount B Min does not match expected value %v, got %v", swapTransaction.AmountBMin, insertedSwap.AmountBMin)
	}

}

func TestSaveSwapMethods(t *testing.T) {
	// Setup test database
	db := setupTestDB(t)

	// Insert test token metadata
	setupTestTokenMetadata(t, db)

	testCases := []struct {
		name       string
		txData     []byte
		version    string
		methodName string
		tokenA     string
		tokenB     string
		value      float64
	}{
		{
			name: "Test RemoveLiquidity",
			txData: common.FromHex("0xbaa2abde" + // RemoveLiquidity method ID
				"0000000000000000000000006b175474e89094c44da98b954eedeac495271d0f" + // tokenA (DAI)
				"000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2" + // tokenB (WETH)
				"0000000000000000000000000000000000000000000000004ef159a1bc7600e5" + // liquidity
				"0000000000000000000000000000000000000000000000260fbe8f136f3af532" + // amountAMin
				"000000000000000000000000000000000000000000000000027eb4840daaab41" + // amountBMin
				"000000000000000000000000c47e5d32f7be0cc171740ebbb3f26f78488cd22f" + // to
				"00000000000000000000000000000000000000000000000000000000675caf2e"), // deadline
			version:    "V2",
			methodName: "RemoveLiquidity",
			tokenA:     "0x6b175474e89094c44da98b954eedeac495271d0f", // DAI
			tokenB:     "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2", // WETH
		},
		{
			name: "Test AddLiquidity",
			txData: common.FromHex("0xe8e33700" + // AddLiquidity method ID
				"0000000000000000000000006b175474e89094c44da98b954eedeac495271d0f" + // tokenA (DAI)
				"000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2" + // tokenB (WETH)
				"0000000000000000000000000000000000000000000000008ac7230489e80000" + // amountADesired (10 DAI)
				"0000000000000000000000000000000000000000000000000de0b6b3a7640000" + // amountBDesired (1 WETH)
				"0000000000000000000000000000000000000000000000008ac7230489e80000" + // amountAMin
				"0000000000000000000000000000000000000000000000000de0b6b3a7640000" + // amountBMin
				"000000000000000000000000c47e5d32f7be0cc171740ebbb3f26f78488cd22f" + // to
				"00000000000000000000000000000000000000000000000000000000675caf2e"), // deadline
			version:    "V2",
			methodName: "AddLiquidity",
			tokenA:     "0x6b175474e89094c44da98b954eedeac495271d0f", // DAI
			tokenB:     "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2", // WETH
		},
		{
			name: "Test SwapExactETHForTokens",
			txData: common.FromHex("0x7ff36ab5" + // SwapExactETHForTokens method ID
				"0000000000000000000000000000000000000000000000000de0b6b3a7640000" + // amountOutMin
				"0000000000000000000000000000000000000000000000000000000000000080" + // path offset
				"000000000000000000000000c47e5d32f7be0cc171740ebbb3f26f78488cd22f" + // to
				"00000000000000000000000000000000000000000000000000000000675caf2e" + // deadline
				"0000000000000000000000000000000000000000000000000000000000000002" + // path length
				"000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2" + // WETH
				"0000000000000000000000006b175474e89094c44da98b954eedeac495271d0f"), // DAI
			version:    "V2",
			methodName: "SwapExactETHForTokens",
			tokenA:     "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2", // WETH
			tokenB:     "0x6b175474e89094c44da98b954eedeac495271d0f", // DAI
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create transaction with test data
			tx := types.NewTransaction(
				0, // nonce
				common.HexToAddress("0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"), // Uniswap V2 Router
				big.NewInt(0),           // value
				500000,                  // gas limit
				big.NewInt(50000000000), // gas price
				tc.txData,               // data
			)

			// Create repository instance
			repo := NewDexExchangeRepository(db)

			// Execute SaveSwap
			err := repo.SaveSwap(context.Background(), tx, tc.version)
			if err != nil {
				t.Fatalf("Failed to save swap: %v", err)
			}

			// Verify the saved transaction
			var saved entity.SwapTransaction
			err = db.Get(&saved, "SELECT * FROM swap_transactions WHERE tx_hash = ?", tx.Hash().Hex())
			if err != nil {
				t.Fatalf("Failed to retrieve saved transaction: %v", err)
			}

			// Assertions
			assert.Equal(t, tc.methodName, saved.MethodName)
			assert.Equal(t, tc.version, saved.Version)
			assert.NotZero(t, saved.Value) // Value should be calculated

			// Method-specific assertions
			switch tc.methodName {
			case "addLiquidity":
				assert.Equal(t, tc.tokenA, saved.TokenA)
				assert.Equal(t, tc.tokenB, saved.TokenB)
				assert.NotEmpty(t, saved.AmountADesired)
				assert.NotEmpty(t, saved.AmountBDesired)
			case "removeLiquidity":
				assert.Equal(t, tc.tokenA, saved.TokenA)
				assert.Equal(t, tc.tokenB, saved.TokenB)
				assert.NotEmpty(t, saved.AmountAMin)
				assert.NotEmpty(t, saved.AmountBMin)
			case "swapExactETHForTokens":
				assert.Equal(t, tc.tokenA, saved.TokenPathFrom)
				assert.Equal(t, tc.tokenB, saved.TokenPathTo)
				assert.NotEmpty(t, saved.AmountIn)
			}

			// Print detailed information for debugging
			// t.Logf("Transaction details for %s:", tc.name)
			// t.Logf("Method Name: %s", saved.MethodName)
			// t.Logf("Version: %s", saved.Version)
			// t.Logf("Value: %f", saved.Value)
			// t.Logf("Token A: %s", saved.TokenA)
			// t.Logf("Token B: %s", saved.TokenB)
		})
	}
}

func setupTestTokenMetadata(t *testing.T, db *sqlx.DB) {
	// Insert test token metadata
	tokens := []struct {
		address  string
		symbol   string
		decimals int
		price    float64
	}{
		{
			address:  "0x6b175474e89094c44da98b954eedeac495271d0f", // DAI
			symbol:   "DAI",
			decimals: 18,
			price:    1.0, // 1 USD
		},
		{
			address:  "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2", // WETH
			symbol:   "WETH",
			decimals: 18,
			price:    2000.0, // 2000 USD
		},
	}

	for _, token := range tokens {
		_, err := db.Exec(`
			INSERT INTO token_metadata (address, symbol, decimals, price)
			VALUES (?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
				symbol = VALUES(symbol),
				decimals = VALUES(decimals),
				price = VALUES(price)
		`,
			token.address,
			token.symbol,
			token.decimals,
			token.price,
		)
		if err != nil {
			t.Fatalf("Failed to insert token metadata: %v", err)
		}
	}
}
