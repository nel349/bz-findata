package main

const (
    // Chain IDs
    ChainEthereum = 1
    ChainBase     = 8453
    ChainArbitrum = 42161
    ChainPolygon  = 137
    ChainAvalanche = 43114

    // Base thresholds in native token decimals
    // These would need to be adjusted based on current market prices
    EthereumSwapThreshold = 0.01   // ETH
    BaseSwapThreshold     = 0.5   // ETH
    ArbitrumSwapThreshold = 3.0   // ETH
    PolygonSwapThreshold  = 5000  // MATIC
    AvalancheSwapThreshold = 100  // AVAX
)

func GetThresholdForChain(chainId uint64) float64 {
    switch chainId {
    case ChainEthereum:
        return EthereumSwapThreshold
    case ChainBase:
        return BaseSwapThreshold
    case ChainArbitrum:
        return ArbitrumSwapThreshold
    case ChainPolygon:
        return PolygonSwapThreshold
    case ChainAvalanche:
        return AvalancheSwapThreshold
    default:
        return 0 // Log all trades for unknown chains
    }
}