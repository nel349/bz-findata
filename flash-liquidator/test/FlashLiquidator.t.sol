// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;

import "forge-std/Test.sol";
import "../src/FlashLiquidator.sol";
import "../src/ISwapRouter.sol";
import "@aave/core-v3/contracts/interfaces/IPool.sol";
import "@aave/core-v3/contracts/flashloan/base/FlashLoanSimpleReceiverBase.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "aave-v3-core/contracts/misc/AaveProtocolDataProvider.sol";

contract FlashLiquidatorTest is Test {
    // Add event declaration to match the contract's event
    event LiquidationExecuted(
        address indexed user,
        address indexed debtAsset,
        address indexed collateralAsset,
        uint256 debtAmount,
        uint256 collateralReceived,
        uint256 profit
    );
    
    FlashLiquidator public liquidator;
    address public constant AAVE_ADDRESS_PROVIDER = 0xa97684ead0e402dC232d5A977953DF7ECBaB3CDb; // Arbitrum
    address public constant UNISWAP_ROUTER = 0xE592427A0AEce92De3Edee1F18E0157C05861564; // Uniswap v3 Router on Arbitrum
    IPool public constant POOL = IPool(0x794a61358D6845594F94dc1DB02A252b5b4814aD); // Aave V3 Pool on Arbitrum
    AaveProtocolDataProvider public constant DATA_PROVIDER = AaveProtocolDataProvider(0x69FA688f1Dc47d4B5d8029D5a35FB7a548310654);

    // Tokens we'll work with
    address _usdc = 0xaf88d065e77c8cC2239327C5EDb3A432268e5831; // USDC on Arbitrum
    address _weth = 0x82aF49447D8a07e3bd95BD0d56f35241523fBab1; // WETH on Arbitrum


    // Test user with unhealthy position
    address public testUser = address(0x436288b1dA64676E57e8Ef2555E448d9470bB9B1);
    // address public debtAsset = address(0x2);
    address public collateralAsset = _weth;
    
    // Fork Arbitrum mainnet
    function setUp() public {
        // Use vm.envString to get RPC URL from environment variables
        string memory arbitrumRpcUrl = vm.envOr("ARBITRUM_RPC_URL", string(""));
        
        // If no RPC URL is provided, skip the test
        if (bytes(arbitrumRpcUrl).length == 0) {
            console.log("No ARBITRUM_RPC_URL environment variable found. Skipping test.");
            return;
        }
        
        // Fork Arbitrum at a more recent block or use the 'latest' keyword
        vm.createSelectFork(arbitrumRpcUrl);
        
        // Deploy the liquidator
        liquidator = new FlashLiquidator(AAVE_ADDRESS_PROVIDER, UNISWAP_ROUTER);
    }
    
    // Test successful deployment
    function testDeployment() public {
        // Skip the test if no RPC URL was provided
        string memory arbitrumRpcUrl = vm.envOr("ARBITRUM_RPC_URL", string(""));
        if (bytes(arbitrumRpcUrl).length == 0) {
            return;
        }
        
        assertEq(liquidator.owner(), address(this));
    }
    
    function testFlashLoanAndLiquidation() public {
        // Skip the test if no RPC URL was provided
        string memory arbitrumRpcUrl = vm.envOr("ARBITRUM_RPC_URL", string(""));
        if (bytes(arbitrumRpcUrl).length == 0) {
            return;
        }
        
        
        // We need to either:
        // 1. Find an actual unhealthy position on Aave, or
        // 2. Create an unhealthy position for testing
        
        // Option 1: Find an actual unhealthy account
        // This requires querying positions or knowing an account ahead of time
        // address userToLiquidate = 0x436288b1dA64676E57e8Ef2555E448d9470bB9B1; // Example address 
        
        // Verify this user actually has an unhealthy position
        (
            uint256 totalCollateralBase,
            uint256 totalDebtBase,
            ,
            ,
            ,
            uint256 healthFactor
        ) = POOL.getUserAccountData(testUser);
        
        // Only proceed if position is liquidatable (health factor < 1)
        if (healthFactor >= 1e18) {
            console.log("User position is healthy, health factor:", healthFactor);
            console.log("Skipping liquidation test as there's no liquidatable position");
            return;
        }
        
        // Format total collateral and total debt to USD
        // uint256 totalCollateralUSD = totalCollateralBase / 1e8;
        // uint256 totalDebtUSD = totalDebtBase / 1e8;
        
       _logPositionDetails(testUser, healthFactor, totalCollateralBase, totalDebtBase);
        
        // // Get token addresses from the protocol data provider
        // (
        //     address aTokenAddress,
        //     ,  // stableDebtTokenAddress (not needed for this test)
        //     address variableDebtTokenAddress
        // ) = DATA_PROVIDER.getReserveTokensAddresses(debtAsset);
        
        // For the collateral asset (WETH):
        (
            address collateralATokenAddress,
            address collateralStableDebtTokenAddress,  // Keep this as a named variable
            address collateralVariableDebtTokenAddress  // Keep this as a named variable
        ) = DATA_PROVIDER.getReserveTokensAddresses(collateralAsset);
        
        // Log token addresses to verify
        console.log("stableDebtTokenAddress:", collateralStableDebtTokenAddress);
        console.log("variableDebtTokenAddress:", collateralVariableDebtTokenAddress);
        console.log("Collateral aToken address:", collateralATokenAddress);


        
         // Get all reserves from Aave pool
        address[] memory reservesList = POOL.getReservesList();
        
        // Variables to track the debt and collateral assets
        address debtAsset;
        uint256 maxDebt = 0;
        address collateralAsset;
        uint256 maxCollateral = 0;
        
        console.log("Checking reserves for user:", testUser);
        
        // Iterate through all reserves to find debt and collateral
        for (uint i = 0; i < reservesList.length; i++) {
            address asset = reservesList[i];
            
            // Get reserve token addresses
            (
                address aTokenAddress,
                address stableDebtTokenAddress,
                address variableDebtTokenAddress
            ) = DATA_PROVIDER.getReserveTokensAddresses(asset);
            
            // Check if user has collateral in this asset
            uint256 aTokenBalance = IERC20(aTokenAddress).balanceOf(testUser);
            
            // Check if user has debt in this asset
            uint256 stableDebt = IERC20(stableDebtTokenAddress).balanceOf(testUser);
            uint256 variableDebt = IERC20(variableDebtTokenAddress).balanceOf(testUser);
            uint256 totalDebt = stableDebt + variableDebt;
            
            // Log asset details
            if (aTokenBalance > 0 || totalDebt > 0) {
                console.log("Asset:", asset);
                console.log("  aToken balance:", aTokenBalance);
                console.log("  Total debt:", totalDebt);
            }
            
            // Track the asset with the highest debt
            if (totalDebt > maxDebt) {
                maxDebt = totalDebt;
                debtAsset = asset;
            }
            
            // Track the asset with the highest collateral
            if (aTokenBalance > maxCollateral) {
                    maxCollateral = aTokenBalance;
                    collateralAsset = asset;
            }
        }

        console.log("Found debt asset:", debtAsset);
        console.log("Debt amount:", maxDebt);
        console.log("Found collateral asset:", collateralAsset);
        console.log("Collateral amount:", maxCollateral);



        
        // Ensure we found both assets
        require(debtAsset != address(0), "No debt asset found");
        require(collateralAsset != address(0), "No collateral asset found");
        
        // // Make sure we have the actual ERC20 token instances
        // IERC20 debtToken = IERC20(debtAsset);
        IERC20 collateralToken = IERC20(collateralAsset);
        
        // Record balances before liquidation
        uint256 liquidatorCollateralBefore = collateralToken.balanceOf(address(liquidator));
        uint256 liquidatorDebtBefore = IERC20(debtAsset).balanceOf(address(liquidator));
        
        console.log("===== BEFORE LIQUIDATION =====");
        console.log("Liquidator collateral balance:", liquidatorCollateralBefore);
        console.log("Liquidator debt token balance:", liquidatorDebtBefore);
        
        // Before executing the liquidation, we need to ensure our contract has enough ETH for gas
        // and possibly for any required approvals or interactions
        vm.deal(address(liquidator), 1 ether);  // Provide some ETH for gas
        
        // We'll try to liquidate up to 50% of the debt (Aave's max liquidation close factor)
        uint256 debtToCover = maxDebt / 2;
        
        console.log("Attempting to liquidate:");
        console.log("User:", testUser);
        console.log("Debt asset:", debtAsset);
        console.log("Debt to cover:", debtToCover);
        console.log("Collateral asset:", collateralAsset);
        
        // Record logs for event analysis
        vm.recordLogs();
        
        // For a real live test, we'll just call the executeFlashLoan function
        // which will trigger the flash loan and then execute liquidation in the callback
        liquidator.executeFlashLoan(
            debtAsset,
            debtToCover,  // Use half the debt as a safer amount
            collateralAsset,
            testUser
        );

        console.log("Flash loan executed");
        
        // Get liquidation event logs
        Vm.Log[] memory entries = vm.getRecordedLogs();
        bool foundLiquidationEvent = false;
        
        console.log("===== EVENT LOGS =====");
        for (uint i = 0; i < entries.length; i++) {
            bytes32 topic0 = entries[i].topics[0];
            
            // Check if this is our LiquidationExecuted event
            if (topic0 == keccak256("LiquidationExecuted(address,address,address,uint256,uint256,uint256)")) {
                foundLiquidationEvent = true;
                console.log("Found LiquidationExecuted event!");
                
                // Decode the event data for more detailed information
                (address user, address debt, address collateral, uint256 debtAmount, uint256 collateralReceived, uint256 profit) = 
                    abi.decode(entries[i].data, (address, address, address, uint256, uint256, uint256));
                
                console.log("  User liquidated:", user);
                console.log("  Debt asset:", debt);
                console.log("  Collateral asset:", collateral);
                console.log("  Debt amount covered:", debtAmount);
                console.log("  Collateral received:", collateralReceived);
                console.log("  Profit:", profit);
            }
            // Log other significant events from Aave
            else if (topic0 == keccak256("LiquidationCall(address,address,address,uint256,uint256,address,bool)")) {
                console.log("Aave LiquidationCall event detected");
            }
            else if (topic0 == keccak256("FlashLoan(address,address,uint256,uint256,uint16)")) {
                console.log("Aave FlashLoan event detected");
            }
        }
        
        // Verify the liquidation event was emitted
        assertTrue(foundLiquidationEvent, "Liquidation event not emitted");
        
        // Get final balances to verify the outcome
        uint256 liquidatorCollateralAfter = collateralToken.balanceOf(address(liquidator));
        uint256 liquidatorDebtAfter = IERC20(debtAsset).balanceOf(address(liquidator));
        
        console.log("===== AFTER LIQUIDATION =====");
        console.log("Liquidator collateral balance:", liquidatorCollateralAfter);
        console.log("Liquidator debt token balance:", liquidatorDebtAfter);
        console.log("Collateral gained:", liquidatorCollateralAfter - liquidatorCollateralBefore);
        console.log("Debt token gained:", liquidatorDebtAfter - liquidatorDebtBefore);
        
        // Verify the contract received collateral
        assertTrue(liquidatorCollateralAfter > liquidatorCollateralBefore, "No collateral received from liquidation");
        
        // Check if there's profit (in either collateral or debt tokens)
        if (liquidatorDebtAfter > liquidatorDebtBefore) {
            console.log("Profit made in debt tokens:", liquidatorDebtAfter - liquidatorDebtBefore);
        }
        
        console.log("Liquidation test completed successfully");
        
        // Optional: withdraw the profits to the test contract
        vm.prank(address(this));
        liquidator.rescueTokens(collateralAsset);
    }




    // Helper function to log financial position details
    function _logPositionDetails(
        address user,
        uint256 healthFactor,
        uint256 totalCollateralBase,
        uint256 totalDebtBase
    ) internal {
        // Format values for display - dividing by 1e8 for USD values
        uint256 totalCollateralUSD = totalCollateralBase / 1e8;
        uint256 totalDebtUSD = totalDebtBase / 1e8;
        
        // Display health factor in its raw form (will be x * 10^18 for a factor of x)
        console.log("===== Position Details =====");
        console.log("User address:", user);
        console.log("Health factor (raw):", healthFactor);
        console.log("Health factor (normalized):", healthFactor / 1e18, ".", (healthFactor % 1e18) / 1e14);
        console.log("Total collateral (USD):", totalCollateralUSD);
        console.log("Total debt (USD):", totalDebtUSD);
        console.log("===========================");
    }
}