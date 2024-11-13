package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jmoiron/sqlx"
	"github.com/nel349/bz-findata/config"
	"github.com/nel349/bz-findata/internal/dex/repository"
	"github.com/nel349/bz-findata/pkg/database/mysql"
)

const (
	// Uniswap V2 Router address
	UniswapRouterAddress = "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"

	// Uniswap V3 Router address
	UniswapV3RouterAddress = "0xE592427A0AEce92De3Edee1F18E0157C05861564"

	// Method signatures for swaps
	SwapExactETHForTokens              = "0x7ff36ab5"
	SwapExactETHForTokensSupportingFee = "0xb6f9de95"
)

type DexRepository struct {
	db *sqlx.DB
}

func main() {
	ctx, cancel := signal.NotifyContext(context.TODO(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-ctx.Done()
		cancel()
	}()
	
	// Connect to your WSS endpoint
	client, err := ethclient.Dial("wss://ethereum-mainnet.core.chainstack.com/52fe0d05347a608831b02990cf1de889")
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.NewDexConfig(ctx)
	if err != nil {
		log.Fatalf("failed config init: %v", err)
	}

	// database
	dbClient, err := mysql.NewMysqlClient(cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.Base)
	if err != nil {
		log.Fatalf("failed database init: %v", err)
	}
	defer dbClient.CloseConnect()

	

	// Create a channel for new headers
	headers := make(chan *types.Header)

	// Subscribe to new block headers
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting to monitor Uniswap swaps...")

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:
			dexRepositories := repository.NewDexRepositories(dbClient.DB)
			go processBlock(client, header, dexRepositories)
		}
	}
}

func processBlock(client *ethclient.Client, header *types.Header, dexRepositories *repository.DexRepositories) {
	block, err := client.BlockByHash(context.Background(), header.Hash())
	if err != nil {
		log.Printf("Error getting block: %v", err)
		return
	}

	fmt.Printf("Processing block: %d\n", block.Number().Uint64())

	// Process each transaction in the block
	for _, tx := range block.Transactions() {
		if tx.To() == nil {
			continue
		}

		toAddress := strings.ToLower(tx.To().Hex())

		// Check if transaction is to either Uniswap Router
		if toAddress == strings.ToLower(UniswapRouterAddress) ||
			toAddress == strings.ToLower(UniswapV3RouterAddress) {

			from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
			if err != nil {
				log.Printf("Error getting sender: %v", err)
				continue
			}

			ethValue := new(big.Float).Quo(
				new(big.Float).SetInt(tx.Value()),
				new(big.Float).SetFloat64(1e18),
			)

			version := "V2"
			if toAddress == strings.ToLower(UniswapV3RouterAddress) {
				version = "V3"
			}

			threshold := GetThresholdForChain(tx.ChainId().Uint64())

			if ethValue.Cmp(big.NewFloat(threshold)) >= 0 {
				fmt.Println("Threshold met")
				// Save to database
				dexRepositories.SaveSwap(context.Background(), tx, version)
			}

			fmt.Println("-----------------------------------------------------")
			fmt.Printf("Uniswap %s Transaction\n", version)
			fmt.Printf("Transaction Hash: %s\n", tx.Hash().Hex())
			fmt.Printf("From: %s\n", from.Hex())
			fmt.Printf("Value: %f ETH\n", ethValue)
			fmt.Println("-----------------------------------------------------")
		}
	}
}
