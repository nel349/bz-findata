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
		t.Fatalf("Failed to connect to MySQL server: %v", err)
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
		0,                          // nonce
		common.HexToAddress("0xdc0488a855ca075c5f7c6b9567fa0b84a9d97ffd5bb4bea913e9840a402b5b79"), // to address
		big.NewInt(0),              // value
		0,                          // gas limit
		big.NewInt(0),              // gas price
		data,                        // data
	)

	fmt.Println("tx: ", tx.Hash().Hex())

	// Create a mock SwapTransaction entity
	swapTransaction := &entity.SwapTransaction{
		TxHash:             tx.Hash().Hex(),
		TokenA:             "0x6b175474e89094c44da98b954eedeac495271d0f",
		TokenB:             "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		Liquidity:          "5688426355532103909",
		AmountAMin:         "702110776270740190514",
		AmountBMin:         "179779514399370049",
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
