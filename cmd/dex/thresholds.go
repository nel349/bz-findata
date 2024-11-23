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
    EthereumSwapThreshold = 0   // ETH  // original 10.0 , and 0 for testing
    BaseSwapThreshold     = 10.0   // ETH
    ArbitrumSwapThreshold = 10.0   // ETH
    PolygonSwapThreshold  = 10000  // MATIC
    AvalancheSwapThreshold = 500  // AVAX
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