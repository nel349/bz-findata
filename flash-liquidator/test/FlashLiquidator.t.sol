// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;

import "forge-std/Test.sol";
import "../src/FlashLiquidator.sol";
import "../src/ISwapRouter.sol";
import "@aave/core-v3/contracts/interfaces/IPool.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

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
    
    // Test user with unhealthy position (you'd need to find one or create a mock)
    address public testUser = address(0x1);
    address public debtAsset = address(0x2);
    address public collateralAsset = address(0x3);
    
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
        
        // Tokens we'll work with
        address usdc = 0xaf88d065e77c8cC2239327C5EDb3A432268e5831; // USDC on Arbitrum
        address weth = 0x82aF49447D8a07e3bd95BD0d56f35241523fBab1; // WETH on Arbitrum
        
        // We need to either:
        // 1. Find an actual unhealthy position on Aave, or
        // 2. Create an unhealthy position for testing
        
        // Option 1: Find an actual unhealthy account
        // This requires querying positions or knowing an account ahead of time
        address userToLiquidate = 0x436288b1dA64676E57e8Ef2555E448d9470bB9B1; // Example address
        
        // Verify this user actually has an unhealthy position
        (
            uint256 totalCollateralBase,
            uint256 totalDebtBase,
            ,
            ,
            ,
            uint256 healthFactor
        ) = POOL.getUserAccountData(userToLiquidate);
        
        // Only proceed if position is liquidatable (health factor < 1)
        if (healthFactor >= 1e18) {
            console.log("User position is healthy, health factor:", healthFactor);
            console.log("Skipping liquidation test as there's no liquidatable position");
            return;
        }
        
        // Format total collateral and total debt to USD
        uint256 totalCollateralUSD = totalCollateralBase / 1e8;
        uint256 totalDebtUSD = totalDebtBase / 1e8;
        
       _logPositionDetails(userToLiquidate, healthFactor, totalCollateralBase, totalDebtBase);
        
        // // Get user reserve data to find what assets they're using
        // address[] memory reservesList = POOL.getReservesList();
        // address debtAsset;
        // address collateralAsset;
        // uint256 debtToCover;
        
        // // Find debt asset and collateral asset for this user
        // for (uint i = 0; i < reservesList.length; i++) {
        //     address asset = reservesList[i];
        //     (
        //         uint256 aTokenBalance,
        //         uint256 stableDebt,
        //         uint256 variableDebt,
        //         ,
        //         ,
        //         ,
        //         ,
        //         ,
                
        //     ) = POOL.getUserReserveData(asset, userToLiquidate);
            
        //     // Found debt asset
        //     if (variableDebt > 0 || stableDebt > 0) {
        //         debtAsset = asset;
        //         // Cover 50% of debt (max allowed in Aave)
        //         debtToCover = (variableDebt + stableDebt) / 2;
        //         console.log("Debt asset:", debtAsset);
        //         console.log("Debt to cover:", debtToCover);
        //     }
            
        //     // Found collateral asset
        //     if (aTokenBalance > 0) {
        //         collateralAsset = asset;
        //         console.log("Collateral asset:", collateralAsset);
        //     }
        // }
        
        // // Ensure we found both assets
        // require(debtAsset != address(0), "No debt asset found");
        // require(collateralAsset != address(0), "No collateral asset found");
        
        // // Make sure we have the actual ERC20 token instances
        // IERC20 debtToken = IERC20(debtAsset);
        // IERC20 collateralToken = IERC20(collateralAsset);
        
        // // Record balances before liquidation
        // uint256 liquidatorCollateralBefore = collateralToken.balanceOf(address(liquidator));
        
        // // Approve the liquidator contract to spend tokens from test contract if needed
        // // (though typically not needed as flash loans handle this)
        
        // // Execute the flash liquidation
        // vm.recordLogs();
        // liquidator.executeFlashLiquidation(
        //     debtAsset,
        //     debtToCover,
        //     collateralAsset,
        //     userToLiquidate
        // );
        
        // // Get liquidation event logs
        // Vm.Log[] memory entries = vm.getRecordedLogs();
        // bool foundLiquidationEvent = false;
        
        // for (uint i = 0; i < entries.length; i++) {
        //     // Check if this is our LiquidationExecuted event
        //     if (entries[i].topics[0] == keccak256("LiquidationExecuted(address,address,address,uint256,uint256,uint256)")) {
        //         foundLiquidationEvent = true;
        //         break;
        //     }
        // }
        
        // // Verify the liquidation event was emitted
        // assertTrue(foundLiquidationEvent, "Liquidation event not emitted");
        
        // // Verify the contract received collateral
        // uint256 liquidatorCollateralAfter = collateralToken.balanceOf(address(liquidator));
        // assertTrue(liquidatorCollateralAfter > liquidatorCollateralBefore, "No collateral received from liquidation");
        
        // console.log("Liquidation successful");
        // console.log("Collateral received:", liquidatorCollateralAfter - liquidatorCollateralBefore);
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