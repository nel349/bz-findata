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
    
    // Mainnet Arbitrum addresses
    address public constant MAINNET_AAVE_ADDRESS_PROVIDER = 0xa97684ead0e402dC232d5A977953DF7ECBaB3CDb;
    address public constant MAINNET_UNISWAP_ROUTER = 0xE592427A0AEce92De3Edee1F18E0157C05861564;
    IPool public constant MAINNET_POOL = IPool(0x794a61358D6845594F94dc1DB02A252b5b4814aD);
    AaveProtocolDataProvider public constant MAINNET_DATA_PROVIDER = AaveProtocolDataProvider(0x69FA688f1Dc47d4B5d8029D5a35FB7a548310654);
    
    // Testnet (Arbitrum Goerli) addresses
    address public constant TESTNET_AAVE_ADDRESS_PROVIDER = 0xf8aa90E66B8BaE13f2E4aDe6104ABAB8EeDebBBe;
    address public constant TESTNET_UNISWAP_ROUTER = 0xE592427A0AEce92De3Edee1F18E0157C05861564;
    IPool public constant TESTNET_POOL = IPool(0x8472e931d3d005Ed3113C2019343BfE242E5dA89);
    AaveProtocolDataProvider public constant TESTNET_DATA_PROVIDER = AaveProtocolDataProvider(0x84b7C502e1821b880AFf7c528748CDEfe2e7f5a8);

    // Dynamic references for the current network
    address public AAVE_ADDRESS_PROVIDER;
    address public UNISWAP_ROUTER;
    IPool public POOL;
    AaveProtocolDataProvider public DATA_PROVIDER;

    // Network selection flag
    bool public isMainnet;

    // Mainnet tokens 
    address public MAINNET_USDC = 0xaf88d065e77c8cC2239327C5EDb3A432268e5831;
    address public MAINNET_WETH = 0x82aF49447D8a07e3bd95BD0d56f35241523fBab1;
    address public MAINNET_LUSD = 0x93b346b6BC2548dA6A1E7d98E9a421B42541425b;

    // Testnet tokens (Arbitrum Goerli)
    address public TESTNET_USDC = 0x3D28771f2dfE34f57E1Db4e0dEA05d59f9FA92aE;
    address public TESTNET_WETH = 0xD513c4e3c0A45499cb8ad80bD4A4E887e38528aB;
    
    // Dynamic token references
    address public _usdc;
    address public _weth;

    // Test user with unhealthy position
    address public testUser = address(0x436288b1dA64676E57e8Ef2555E448d9470bB9B1); // Example address - update with your target
    address public collateralAsset;
    
    // Fork settings
    function setUp() public {
        // Determine which network to use
        string memory network = vm.envOr("AAVE_NETWORK", string("mainnet")); // Default to mainnet fork
        
        if (keccak256(abi.encodePacked(network)) == keccak256(abi.encodePacked("mainnet"))) {
            // Use mainnet settings
            isMainnet = true;
            AAVE_ADDRESS_PROVIDER = MAINNET_AAVE_ADDRESS_PROVIDER;
            UNISWAP_ROUTER = MAINNET_UNISWAP_ROUTER;
            POOL = MAINNET_POOL;
            DATA_PROVIDER = MAINNET_DATA_PROVIDER;
            _usdc = MAINNET_USDC;
            _weth = MAINNET_WETH;
            
            // Fork Arbitrum mainnet
            string memory arbitrumRpcUrl = vm.envOr("ARBITRUM_RPC_URL", string(""));
            
            // If no RPC URL is provided, skip the test
            if (bytes(arbitrumRpcUrl).length == 0) {
                console.log("No ARBITRUM_RPC_URL environment variable found. Skipping test.");
                return;
            }
            
            vm.createSelectFork(arbitrumRpcUrl);
            console.log("Running tests on Arbitrum mainnet fork");
        } else {
            // Use testnet settings
            isMainnet = false;
            AAVE_ADDRESS_PROVIDER = TESTNET_AAVE_ADDRESS_PROVIDER;
            UNISWAP_ROUTER = TESTNET_UNISWAP_ROUTER;
            POOL = TESTNET_POOL;
            DATA_PROVIDER = TESTNET_DATA_PROVIDER;
            _usdc = TESTNET_USDC;
            _weth = TESTNET_WETH;
            
            // Fork Arbitrum testnet (Goerli)
            string memory arbitrumGoerliRpcUrl = vm.envOr("ARBITRUM_GOERLI_RPC_URL", string(""));
            
            // If no RPC URL is provided, skip the test
            if (bytes(arbitrumGoerliRpcUrl).length == 0) {
                console.log("No ARBITRUM_GOERLI_RPC_URL environment variable found. Skipping test.");
                return;
            }
            
            vm.createSelectFork(arbitrumGoerliRpcUrl);
            console.log("Running tests on Arbitrum testnet fork");
            
            // For testnet, we may need to specify a testUser or create one
            string memory testUserInput = vm.envOr("TEST_USER", string(""));
            if (bytes(testUserInput).length > 0) {
                // Parse address from environment variable
                testUser = vm.parseAddress(testUserInput);
                console.log("Using test user from environment:", testUser);
            } else {
                console.log("Using default test user:", testUser);
            }
        }
        
        // Deploy the liquidator
        liquidator = new FlashLiquidator(AAVE_ADDRESS_PROVIDER, UNISWAP_ROUTER);
        console.log("Liquidator deployed at:", address(liquidator));
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
        address _localCollateralAsset;
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
                    _localCollateralAsset = asset;
            }
        }

        console.log("Found debt asset:", debtAsset);
        console.log("Debt amount:", maxDebt);
        console.log("Found collateral asset:", _localCollateralAsset);
        console.log("Collateral amount:", maxCollateral);



        
        // Ensure we found both assets
        require(debtAsset != address(0), "No debt asset found");
        require(_localCollateralAsset != address(0), "No collateral asset found");
        
        // // Make sure we have the actual ERC20 token instances
        // IERC20 debtToken = IERC20(debtAsset);
        IERC20 collateralToken = IERC20(_localCollateralAsset);
        
        // Record balances before liquidation
        uint256 liquidatorCollateralBefore = collateralToken.balanceOf(address(liquidator));
        uint256 liquidatorDebtBefore = IERC20(debtAsset).balanceOf(address(liquidator));
        
        console.log("===== BEFORE LIQUIDATION =====");
        console.log("Liquidator collateral balance:", liquidatorCollateralBefore);
        console.log("Liquidator debt token balance:", liquidatorDebtBefore);
        
        // Before executing the liquidation, we need to ensure our contract has enough ETH for gas
        // and possibly for any required approvals or interactions
        vm.deal(address(liquidator), 1 ether);  // Provide some ETH for gas
        
        // We'll liquidate the maximum allowed percentage (50%)
        // This maximizes our profit potential from the liquidation
        uint256 debtToCover = (maxDebt * 50) / 100;  // 50% of debt
        
        console.log("Executing flash loan for liquidation:");
        console.log("- User to liquidate:", testUser);
        console.log("- Debt asset:", debtAsset);
        console.log("- Debt amount:", maxDebt);
        console.log("- Debt to cover ( 50 %):", debtToCover);
        console.log("- Collateral asset:", _localCollateralAsset);
        
        // Since we're on mainnet, we don't need to fund the contract with tokens
        // The flash loan will provide the necessary tokens
        
        // Record logs for event analysis
        vm.recordLogs();
        
        // For a real live test, we'll just call the executeFlashLoan function
        // which will trigger the flash loan and then execute liquidation in the callback
        liquidator.executeFlashLoan(
            debtAsset,
            debtToCover,
            _localCollateralAsset,
            testUser,
            50  // Use 50% of debt to maximize profit
        );
        
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
        liquidator.rescueTokens(_localCollateralAsset);
    }
    
    /**
     * @dev Helper function to log financial position details
     * @param user The address of the user
     * @param healthFactor The health factor of the user's position
     * @param totalCollateralBase The total collateral in base units
     * @param totalDebtBase The total debt in base units
     */
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