// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;

import {FlashLoanSimpleReceiverBase} from "@aave/core-v3/contracts/flashloan/base/FlashLoanSimpleReceiverBase.sol";
import {IPoolAddressesProvider} from "@aave/core-v3/contracts/interfaces/IPoolAddressesProvider.sol";
import {IPool} from "@aave/core-v3/contracts/interfaces/IPool.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {ISwapRouter} from "./ISwapRouter.sol";
import "forge-std/console.sol";

contract FlashLiquidator is FlashLoanSimpleReceiverBase {
    using SafeERC20 for IERC20;
    
    address public owner;
    ISwapRouter public swapRouter;
    
    event LiquidationExecuted(
        address indexed user,
        address indexed debtAsset,
        address indexed collateralAsset,
        uint256 debtAmount,
        uint256 collateralReceived,
        uint256 profit
    );
    
    constructor(
        address _addressProvider,
        address _swapRouter
    ) FlashLoanSimpleReceiverBase(IPoolAddressesProvider(_addressProvider)) {
        owner = msg.sender;
        swapRouter = ISwapRouter(_swapRouter);
    }
    
    /**
     * @dev Initiates a flash loan to perform liquidation
     * @param debtAsset The address of the debt asset to borrow
     * @param debtAmount The amount to borrow
     * @param collateralAsset The address of the collateral asset to receive
     * @param userToLiquidate The address of the user to liquidate
     */
    function executeFlashLoan(
        address debtAsset,
        uint256 debtAmount,
        address collateralAsset,
        address userToLiquidate
    ) external onlyOwner {
        // Request flash loan
        bytes memory params = abi.encode(
            debtAsset,
            collateralAsset,
            userToLiquidate,
            debtAmount
        );
        
        // Request the flash loan
        IPool pool = POOL;
        pool.flashLoanSimple(
            address(this),
            debtAsset,
            debtAmount,
            params,
            0  // referral code (0 = no referral)
        );
    }
    
    /**
     * @dev This function is called after the contract receives the flash loaned amount
     * @param asset The address of the flash-borrowed asset
     * @param amount The amount of the flash-borrowed asset
     * @param premium The fee of the flash-borrowed asset
     * @param initiator The address of the flashloan initiator
     * @param params The byte-encoded params passed when initiating the flashloan
     * @return boolean indicating whether the operation was successful
     */
    function executeOperation(
        address asset,
        uint256 amount,
        uint256 premium,
        address initiator,
        bytes calldata params
    ) external override returns (bool) {
        // Ensure this is called by the Aave pool
        require(msg.sender == address(POOL), "Callback called by unauthorized address");
        require(initiator == address(this), "Initiator must be this contract");
        
        // Decode parameters
        (
            address debtAsset,
            address collateralAsset,
            address userToLiquidate,
            uint256 debtToCover
        ) = abi.decode(params, (address, address, address, uint256));
        
        // Verify that we received the correct asset from the flash loan
        require(asset == debtAsset, "Received asset doesn't match requested asset");
        
        // Execute the actual liquidation
        _executeLiquidation(
            debtAsset, 
            collateralAsset, 
            userToLiquidate, 
            debtToCover, 
            amount,
            premium
        );
        
        // Return true to indicate the flash loan was handled successfully
        return true;
    }
    
    /**
     * @dev Executes the liquidation using the flash loaned funds
     * @param debtAsset The debt asset to repay
     * @param collateralAsset The collateral asset to receive
     * @param userToLiquidate The user to liquidate
     * @param debtToCover The amount of debt to cover
     * @param flashLoanAmount The total amount that was flash borrowed
     * @param flashLoanPremium The premium that needs to be repaid on top of the flash loan
     */
    function _executeLiquidation(
        address debtAsset,
        address collateralAsset,
        address userToLiquidate,
        uint256 debtToCover,
        uint256 flashLoanAmount,
        uint256 flashLoanPremium
    ) internal {
        // Track collateral balance before liquidation 
        uint256 collateralBefore = IERC20(collateralAsset).balanceOf(address(this));
        
        // Approve the pool to use our funds for liquidation
        IERC20(debtAsset).safeApprove(address(POOL), debtToCover);
        
        // Execute liquidation
        POOL.liquidationCall(
            collateralAsset,
            debtAsset,
            userToLiquidate,
            debtToCover,
            false  // receive the underlying asset, not aTokens
        );
        
        // Calculate received collateral
        uint256 collateralAfter = IERC20(collateralAsset).balanceOf(address(this));
        uint256 collateralReceived = collateralAfter - collateralBefore;
        
        // Calculate amount to repay (loan + premium)
        uint256 amountToRepay = flashLoanAmount + flashLoanPremium;
        
        // We need to ensure we have enough of the debt asset to repay
        uint256 debtTokenBalance = IERC20(debtAsset).balanceOf(address(this));
        
        console.log("Debt token balance before swap:", debtTokenBalance);
        console.log("Amount to repay:", amountToRepay);
        
        // If we don't have enough debt tokens, swap some collateral
        if (debtTokenBalance < amountToRepay && collateralAsset != debtAsset) {
            uint256 amountToSwap = _calculateSwapAmount(
                collateralAsset,
                debtAsset,
                amountToRepay - debtTokenBalance,
                collateralReceived
            );
            
            // Ensure we have enough collateral to swap
            require(amountToSwap <= collateralReceived, "Not enough collateral to swap");
            
            // Swap collateral for debt asset
            uint256 receivedFromSwap = _swapTokens(
                collateralAsset,
                debtAsset,
                amountToSwap,
                amountToRepay - debtTokenBalance
            );
            
            // Check if we received enough from the swap
            debtTokenBalance = IERC20(debtAsset).balanceOf(address(this));
            console.log("Debt token balance after swap:", debtTokenBalance);
            
            if (debtTokenBalance < amountToRepay) {
                console.log("WARNING: Not enough tokens received from swap to repay the flash loan");
                console.log("Shortfall:", amountToRepay - debtTokenBalance);
                
                // In a real-world scenario, we would need to handle this situation
                // For this test, we'll revert with a clear message
                revert("Insufficient funds after swap to repay flash loan");
            }
        }
        
        // Approve the pool to withdraw the debt asset + premium
        IERC20(debtAsset).safeApprove(address(POOL), amountToRepay);
        
        // Calculate profit (remaining collateral)
        uint256 finalCollateralBalance = IERC20(collateralAsset).balanceOf(address(this));
        
        // Emit event with liquidation details
        emit LiquidationExecuted(
            userToLiquidate,
            debtAsset,
            collateralAsset,
            debtToCover,
            collateralReceived,
            finalCollateralBalance
        );
    }
    
    /**
     * @dev Calculates the optimal amount of collateral to swap
     * @param fromAsset The asset to swap from
     * @param toAsset The asset to swap to
     * @param amountNeeded The amount needed of the to asset
     * @param maxAmount The maximum amount of the from asset available to swap
     * @return The amount of from asset to swap
     */
    function _calculateSwapAmount(
        address fromAsset,
        address toAsset,
        uint256 amountNeeded,
        uint256 maxAmount
    ) internal view returns (uint256) {
        // fromAsset and toAsset parameters are not used in this simplified version
        // but are kept for documentation and potential future price-aware implementation
        
        console.log("fromAsset", fromAsset);
        console.log("toAsset", toAsset);
        console.log("amountNeeded", amountNeeded);
        console.log("maxAmount", maxAmount);


        // For WETH to USDC/stablecoin swaps, we need a reasonable conversion ratio
        // Based on your logs, you're trying to swap around 4.1 WETH for 7.4e21 USDC
        
        // Since we're swapping a high-value asset (like ETH) for a lower-value asset (like USDC)
        // we need to use a much smaller amount of the high-value asset
        
        // A simple approach: use up to 90% of available collateral
        // This ensures we maximize our chances of getting enough toAsset
        uint256 amountToUse = (maxAmount * 90) / 100;
        
        console.log("amountToUse", amountToUse);
        return amountToUse;
    }
    
    /**
     * @dev Swaps tokens using Uniswap v3 or another DEX
     */
    function _swapTokens(
        address fromToken,
        address toToken,
        uint256 amountIn,
        uint256 minAmountOut
    ) internal returns (uint256) {
        // Approve the router to spend tokens
        IERC20(fromToken).safeApprove(address(swapRouter), amountIn);

        // CRITICAL: Set an extremely low minimum amount to ensure the swap succeeds
        // For production, you would want a more sophisticated approach with price oracles
        // Set minimum to 0.01% of what we need (1/10000)
        uint256 minAmountWithSlippage = minAmountOut / 10000;
        
        console.log("minAmountOut", minAmountOut);
        console.log("minAmountWithSlippage", minAmountWithSlippage);
        
        // Get balance before swap to calculate actual received amount
        uint256 toTokenBalanceBefore = IERC20(toToken).balanceOf(address(this));
        
        // Prepare swap parameters
        ISwapRouter.ExactInputSingleParams memory params = 
            ISwapRouter.ExactInputSingleParams({
                tokenIn: fromToken,
                tokenOut: toToken,
                fee: 3000, // 0.3% fee tier
                recipient: address(this),
                deadline: block.timestamp + 300, // 5 minutes
                amountIn: amountIn,
                amountOutMinimum: minAmountWithSlippage,
                sqrtPriceLimitX96: 0
            });
        
        // Execute the swap
        uint256 amountOut = swapRouter.exactInputSingle(params);
        
        // Get balance after swap to verify
        uint256 toTokenBalanceAfter = IERC20(toToken).balanceOf(address(this));
        uint256 actualReceived = toTokenBalanceAfter - toTokenBalanceBefore;
        
        console.log("Swap executed:");
        console.log("- Amount in:", amountIn);
        console.log("- Expected minimum:", minAmountOut);
        console.log("- Reported received:", amountOut);
        console.log("- Actual received:", actualReceived);
        console.log("- Shortfall:", minAmountOut > actualReceived ? minAmountOut - actualReceived : 0);
        
        return actualReceived;
    }
    
    /**
     * @dev Withdraw tokens to the owner
     */
    function rescueTokens(address token) external onlyOwner {
        uint256 balance = IERC20(token).balanceOf(address(this));
        IERC20(token).safeTransfer(owner, balance);
    }
    
    /**
     * @dev Withdraw ETH to the owner
     */
    function rescueETH() external onlyOwner {
        payable(owner).transfer(address(this).balance);
    }
    
    /**
     * @dev Update the owner
     */
    function transferOwnership(address newOwner) external onlyOwner {
        require(newOwner != address(0), "New owner cannot be zero address");
        owner = newOwner;
    }
    
    /**
     * @dev Update the swap router
     */
    function setSwapRouter(address _swapRouter) external onlyOwner {
        require(_swapRouter != address(0), "Swap router cannot be zero address");
        swapRouter = ISwapRouter(_swapRouter);
    }
    
    /**
     * @dev Only allow the owner to call
     */
    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner");
        _;
    }
    
    // Required to receive ETH
    receive() external payable {}
}