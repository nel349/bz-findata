// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;

import {FlashLoanSimpleReceiverBase} from "@aave/core-v3/contracts/flashloan/base/FlashLoanSimpleReceiverBase.sol";
import {IPoolAddressesProvider} from "@aave/core-v3/contracts/interfaces/IPoolAddressesProvider.sol";
import {IPool} from "@aave/core-v3/contracts/interfaces/IPool.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {ISwapRouter} from "./ISwapRouter.sol";

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
     * @param debtAsset The address of the debt asset to cover
     * @param debtAmount The amount of debt to cover
     * @param collateralAsset The address of the collateral asset to receive
     * @param userToLiquidate The address of the user to liquidate
     */
    function executeFlashLiquidation(
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
        uint256 amountToRepay = amount + premium;
        
        // We need to ensure we have enough of the debt asset to repay
        uint256 debtTokenBalance = IERC20(debtAsset).balanceOf(address(this));
        
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
            _swapTokens(
                collateralAsset,
                debtAsset,
                amountToSwap,
                amountToRepay - debtTokenBalance
            );
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
        
        // Return true to indicate the flash loan was handled successfully
        return true;
    }
    
    /**
     * @dev Calculates the optimal amount of collateral to swap
     */
    function _calculateSwapAmount(
        address fromToken,
        address toToken,
        uint256 amountNeeded,
        uint256 maxAmount
    ) internal view returns (uint256) {
        // Add a buffer for slippage (3%)
        uint256 amountWithBuffer = (amountNeeded * 103) / 100;
        return amountWithBuffer < maxAmount ? amountWithBuffer : maxAmount;
    }
    
    /**
     * @dev Swaps tokens using Uniswap v3 or another DEX
     */
    function _swapTokens(
        address fromToken,
        address toToken,
        uint256 amountIn,
        uint256 minAmountOut
    ) internal {
        // Approve the router to spend tokens
        IERC20(fromToken).safeApprove(address(swapRouter), amountIn);
        
        // Prepare swap parameters
        ISwapRouter.ExactInputSingleParams memory params = 
            ISwapRouter.ExactInputSingleParams({
                tokenIn: fromToken,
                tokenOut: toToken,
                fee: 3000, // 0.3% fee tier
                recipient: address(this),
                deadline: block.timestamp + 300, // 5 minutes
                amountIn: amountIn,
                amountOutMinimum: minAmountOut,
                sqrtPriceLimitX96: 0
            });
        
        // Execute the swap
        swapRouter.exactInputSingle(params);
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