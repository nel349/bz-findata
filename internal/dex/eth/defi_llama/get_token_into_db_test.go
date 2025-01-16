package defi_llama

import (
	"fmt"
	"log"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver // NEEDED FOR SQLX
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/nel349/bz-findata/pkg/entity"
)

// Test to get token info from defi llama and store in database similar to tests
// in exchange_test.go


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

// Test to get token info from server and store in database
func TestGetTokenInfoFromServerAndStoreInDB(t *testing.T) {
	// load the test env file
	if err := godotenv.Load("../../../../test.env"); err != nil {
		log.Fatal("Error loading test.env file")
	}
	db := setupTestDB(t)

	tokenInfo, err := GetTokenMetadataFromDbOrDefiLlama(db, "0x64766392ad32a6c94b965b5bf655e07371c23a1d", 10*time.Second)
	if err != nil {
		t.Errorf("failed to get token info: %v", err)
	}

	// check the token info fields are not empty
	if tokenInfo.Address == "" || tokenInfo.Symbol == "" || tokenInfo.Decimals == 0 || tokenInfo.Price == 0 {
		t.Errorf("token info fields are empty")
	}

	// check the token info is stored in the database
	var tokenInfoFromDB entity.TokenInfo
	err = db.Get(&tokenInfoFromDB, "SELECT * FROM token_metadata WHERE address = ?", "0x64766392ad32a6c94b965b5bf655e07371c23a1d")
	if err != nil {
		t.Errorf("failed to get token info from database: %v", err)
	}

	// expected token info
	expectedTokenInfo := entity.TokenInfo{
		Address:  "0x64766392ad32a6c94b965b5bf655e07371c23a1d",
		Symbol:   "Yilongma",
		Decimals: 9,
	}
	
	// check the address is the same as expected
	if tokenInfoFromDB.Address != expectedTokenInfo.Address {
		t.Errorf("address is not the same as expected")
	}

	// check the symbol is the same as expected
	if tokenInfoFromDB.Symbol != expectedTokenInfo.Symbol {
		t.Errorf("symbol is not the same as expected")
	}

	// check the decimals is the same as expected
	if tokenInfoFromDB.Decimals != expectedTokenInfo.Decimals {
		t.Errorf("decimals is not the same as expected")
	}
}
