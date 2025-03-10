// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;

import "forge-std/Script.sol";
import "../src/FlashLiquidator.sol";

contract DeployFlashLiquidator is Script {
    function run() external {
        // Get the private key from environment variable
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        
        // Addresses for Arbitrum
        address aaveAddressProvider = 0xa97684ead0e402dC232d5A977953DF7ECBaB3CDb;
        address uniswapRouter = 0xE592427A0AEce92De3Edee1F18E0157C05861564;
        
        vm.startBroadcast(deployerPrivateKey);
        
        // Deploy the liquidator
        FlashLiquidator liquidator = new FlashLiquidator(
            aaveAddressProvider,
            uniswapRouter
        );
        
        vm.stopBroadcast();
        
        // Output the contract address
        console.log("FlashLiquidator deployed at:", address(liquidator));
    }
}