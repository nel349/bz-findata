package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	// "math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Contract ABIs as strings
const poolAbiJSON = `[
	{"inputs":[{"internalType":"address","name":"asset","type":"address"}],"name":"getReserveData","outputs":[{"components":[{"components":[{"internalType":"uint256","name":"data","type":"uint256"}],"internalType":"struct DataTypes.ReserveConfigurationMap","name":"configuration","type":"tuple"},{"internalType":"uint128","name":"liquidityIndex","type":"uint128"},{"internalType":"uint128","name":"currentLiquidityRate","type":"uint128"},{"internalType":"uint128","name":"variableBorrowIndex","type":"uint128"},{"internalType":"uint128","name":"currentVariableBorrowRate","type":"uint128"},{"internalType":"uint128","name":"currentStableBorrowRate","type":"uint128"},{"internalType":"uint40","name":"lastUpdateTimestamp","type":"uint40"},{"internalType":"uint16","name":"id","type":"uint16"},{"internalType":"address","name":"aTokenAddress","type":"address"},{"internalType":"address","name":"stableDebtTokenAddress","type":"address"},{"internalType":"address","name":"variableDebtTokenAddress","type":"address"},{"internalType":"address","name":"interestRateStrategyAddress","type":"address"},{"internalType":"uint128","name":"accruedToTreasury","type":"uint128"},{"internalType":"uint128","name":"unbacked","type":"uint128"},{"internalType":"uint128","name":"isolationModeTotalDebt","type":"uint128"}],"internalType":"struct DataTypes.ReserveData","name":"","type":"tuple"}],"stateMutability":"view","type":"function"},
	{"inputs":[],"name":"getReservesList","outputs":[{"internalType":"address[]","name":"","type":"address[]"}],"stateMutability":"view","type":"function"},
	{"inputs":[{"internalType":"address","name":"user","type":"address"}],"name":"getUserAccountData","outputs":[{"internalType":"uint256","name":"totalCollateralBase","type":"uint256"},{"internalType":"uint256","name":"totalDebtBase","type":"uint256"},{"internalType":"uint256","name":"availableBorrowsBase","type":"uint256"},{"internalType":"uint256","name":"currentLiquidationThreshold","type":"uint256"},{"internalType":"uint256","name":"ltv","type":"uint256"},{"internalType":"uint256","name":"healthFactor","type":"uint256"}],"stateMutability":"view","type":"function"}
]`

const aTokenAbiJSON = `[
	{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Transfer","type":"event"},
	{"inputs":[{"internalType":"address","name":"account","type":"address"}],"name":"balanceOf","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},
	{"inputs":[],"name":"totalSupply","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"}
]`

// Aave Pool Address on Arbitrum
const poolAddress = "0x794a61358D6845594F94dc1DB02A252b5b4814aD"

func main() {
	log.Println("Starting liquidator service")

	client, err := ethclient.Dial(os.Getenv("ARBITRUM_RPC_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to the Arbitrum network: %v", err)
	}
	defer client.Close()

	log.Println("Connected to Arbitrum network")

	err = monitorLiquidatablePositions(client)
	if err != nil {
		log.Fatalf("Error monitoring liquidatable positions: %v", err)
	}
}

type ReserveConfiguration struct {
	Data *big.Int `json:"data"`
}

type ReserveData struct {
	Configuration               ReserveConfiguration `json:"configuration"`
	LiquidityIndex              *big.Int             `json:"liquidityIndex"`
	CurrentLiquidityRate        *big.Int             `json:"currentLiquidityRate"`
	VariableBorrowIndex         *big.Int             `json:"variableBorrowIndex"`
	CurrentVariableBorrowRate   *big.Int             `json:"currentVariableBorrowRate"`
	CurrentStableBorrowRate     *big.Int             `json:"currentStableBorrowRate"`
	LastUpdateTimestamp         *big.Int             `json:"lastUpdateTimestamp"`
	ID                          uint16               `json:"id"`
	ATokenAddress               common.Address       `json:"aTokenAddress"`
	StableDebtTokenAddress      common.Address       `json:"stableDebtTokenAddress"`
	VariableDebtTokenAddress    common.Address       `json:"variableDebtTokenAddress"`
	InterestRateStrategyAddress common.Address       `json:"interestRateStrategyAddress"`
	AccruedToTreasury           *big.Int             `json:"accruedToTreasury"`
	Unbacked                    *big.Int             `json:"unbacked"`
	IsolationModeTotalDebt      *big.Int             `json:"isolationModeTotalDebt"`
}

func monitorLiquidatablePositions(client *ethclient.Client) error {
	ctx := context.Background()

	// Parse the ABIs
	poolAbi, err := abi.JSON(strings.NewReader(poolAbiJSON))
	if err != nil {
		return fmt.Errorf("failed to parse pool ABI: %v", err)
	}

	aTokenAbi, err := abi.JSON(strings.NewReader(aTokenAbiJSON))
	if err != nil {
		return fmt.Errorf("failed to parse aToken ABI: %v", err)
	}

	// Create the pool contract address
	poolAddr := common.HexToAddress(poolAddress)

	// Step 1: Get all supported asset addresses
	callData, err := poolAbi.Pack("getReservesList")
	if err != nil {
		return fmt.Errorf("failed to pack getReservesList data: %v", err)
	}

	result, err := client.CallContract(ctx, ethereum.CallMsg{
		To:   &poolAddr,
		Data: callData,
	}, nil)
	if err != nil {
		return fmt.Errorf("failed to call getReservesList: %v", err)
	}

	var assetAddresses []common.Address
	err = poolAbi.UnpackIntoInterface(&assetAddresses, "getReservesList", result)
	if err != nil {
		return fmt.Errorf("failed to unpack getReservesList result: %v", err)
	}

	log.Printf("Found %d supported assets", len(assetAddresses))

	// For each asset
	for _, assetAddress := range assetAddresses {
		// Step 2: Get aToken for this asset
		callData, err := poolAbi.Pack("getReserveData", assetAddress)
		if err != nil {
			log.Printf("Error packing getReserveData for asset %s: %v", assetAddress.Hex(), err)
			continue
		}

		result, err := client.CallContract(ctx, ethereum.CallMsg{
			To:   &poolAddr,
			Data: callData,
		}, nil)
		if err != nil {
			log.Printf("Error calling getReserveData for asset %s: %v", assetAddress.Hex(), err)
			continue
		}

		// Let's print the raw data for the first asset to verify
		if assetAddress.Hex() == "0xDA10009cBd5D07dd0CeCc66161FC93D7c9000da1" {
			log.Printf("Raw data first 32 bytes: %x", result[:32])
			log.Printf("Raw data bytes 256-288: %x", result[256:288]) // 8*32 to 9*32
			log.Printf("Raw data bytes 288-320: %x", result[288:320]) // 9*32 to 10*32
		}

		// The correct offset for the aToken address is 8*32+12
		aTokenOffset := 8*32 + 12

		// Make sure we have enough data
		if len(result) < aTokenOffset+20 {
			log.Printf("Result data too short for asset %s", assetAddress.Hex())
			continue
		}

		// Extract the aToken address (20 bytes)
		aTokenAddressBytes := result[aTokenOffset : aTokenOffset+20]
		aTokenAddress := common.BytesToAddress(aTokenAddressBytes)

		log.Printf("Asset: %s", assetAddress.Hex())
		log.Printf("aToken: %s", aTokenAddress.Hex())

		// Step 3: Find holders of this aToken by monitoring Transfer events
		// We'll look back a certain number of blocks
		blockNumber, err := client.BlockNumber(ctx)
		if err != nil {
			log.Printf("Error getting latest block number: %v", err)
			continue
		}

		startBlock := blockNumber - 100000 // Last 100,000 blocks
		if startBlock < 0 {
			startBlock = 0
		}

		query := ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(startBlock)),
			ToBlock:   big.NewInt(int64(blockNumber)),
			Addresses: []common.Address{aTokenAddress},
			Topics: [][]common.Hash{{
				aTokenAbi.Events["Transfer"].ID,
			}},
		}

		logs, err := client.FilterLogs(ctx, query)
		if err != nil {
			log.Printf("Error filtering logs for aToken %s: %v", aTokenAddress.Hex(), err)
			continue
		}

		// Extract unique addresses
		uniqueAddresses := make(map[common.Address]bool)
		for _, vLog := range logs {
			// For Transfer events, the indexed parameters (from and to) are in the topics
			if len(vLog.Topics) >= 3 {
				fromAddr := common.BytesToAddress(vLog.Topics[1].Bytes())
				toAddr := common.BytesToAddress(vLog.Topics[2].Bytes())

				if fromAddr != (common.Address{}) { // Skip zero address
					uniqueAddresses[fromAddr] = true
				}
				if toAddr != (common.Address{}) { // Skip zero address
					uniqueAddresses[toAddr] = true
				}
			}
		}

		log.Printf("Found %d potential users for asset %s", len(uniqueAddresses), assetAddress.Hex())

		// Step 4: Check each user's health factor
		for userAddress := range uniqueAddresses {
			callData, err := poolAbi.Pack("getUserAccountData", userAddress)
			if err != nil {
				log.Printf("Error packing getUserAccountData for user %s: %v", userAddress.Hex(), err)
				continue
			}

			result, err := client.CallContract(ctx, ethereum.CallMsg{
				To:   &poolAddr,
				Data: callData,
			}, nil)
			if err != nil {
				log.Printf("Error calling getUserAccountData for user %s: %v", userAddress.Hex(), err)
				continue
			}

			var userData struct {
				TotalCollateralBase         *big.Int
				TotalDebtBase               *big.Int
				AvailableBorrowsBase        *big.Int
				CurrentLiquidationThreshold *big.Int
				Ltv                         *big.Int
				HealthFactor                *big.Int
			}

			err = poolAbi.UnpackIntoInterface(&userData, "getUserAccountData", result)
			if err != nil {
				log.Printf("Error unpacking getUserAccountData for user %s: %v", userAddress.Hex(), err)
				continue
			}

			// To compare with 1, we multiply by 10^18
			one := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)


			// 1.25 × 10^18 Less than 1.25
			// riskThreshold := new(big.Int).Mul(big.NewInt(125), new(big.Int).Exp(big.NewInt(10), big.NewInt(16), nil))

			// Less than 1.07
			// liquidationThreshold107 := new(big.Int).Mul(big.NewInt(107), new(big.Int).Exp(big.NewInt(10), big.NewInt(16), nil))
			if userData.HealthFactor.Cmp(one) < 0 {
				// Convert health factor to human-readable form (division by 10^18)
				healthFactorFloat := new(big.Float).Quo(
					new(big.Float).SetInt(userData.HealthFactor),
					new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),
				)

				log.Printf("LIQUIDATABLE: User %s has health factor %s",
					userAddress.Hex(),
					healthFactorFloat.Text('f', 6),
				)

				// display debt base in regular format	
				totalDebtBaseStr:= new(big.Float).Quo(
					new(big.Float).SetInt(userData.TotalDebtBase),
					new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(8), nil)),
				)

				userCollateralBaseStr:= new(big.Float).Quo(
					new(big.Float).SetInt(userData.TotalCollateralBase),
					new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(8), nil)),
				)

				// Lets print the user debt base and collateral base
				log.Printf("User debt base: %s", totalDebtBaseStr.Text('f', 2))
				log.Printf("User collateral base: %s", userCollateralBaseStr.Text('f', 2))

				// Here you would implement your liquidation logic
			}
		}
	}

	return nil
}
