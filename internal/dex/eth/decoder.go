package eth

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/nel349/bz-findata/pkg/entity"
)

func DecodeSwapExactTokensForTokens(data []byte) (*entity.SwapTransaction, error) {

	// Debug prints
	fmt.Println("Raw data length:", len(data))
	fmt.Printf("Raw data hex: 0x%x\n", data)
	fmt.Println("First 4 bytes (method signature):", fmt.Sprintf("0x%x", data[:4]))

    // The data format appears to be:
    // 4 bytes: method signature
    // 32 bytes: amountIn
    // 32 bytes: unknown (possibly array offset)
    // 32 bytes: to address
    // 32 bytes: deadline
    // 32 bytes: array length (for path)
    // 32 bytes: first token address
    // 32 bytes: second token address

	// Convert the byte data into string array, each element representing a 32-byte parameter
    if len(data) < 4+32*7 {
        return nil, errors.New("invalid input length")
    }

    // Convert amountIn (first parameter)
    amountIn := new(big.Int)
    amountIn.SetString(fmt.Sprintf("%x", data[4:36]), 16)
    amountInDecimal := new(big.Float).Quo(
        new(big.Float).SetInt(amountIn),
        new(big.Float).SetFloat64(1e18),
    )

    // Get to address (third parameter)
    toAddress := fmt.Sprintf("%x", data[68:100])[24:] // Take last 20 bytes

    // Get deadline (fourth parameter)
    deadline := new(big.Int)
    deadline.SetString(fmt.Sprintf("%x", data[100:132]), 16)
    // deadlineUnix := time.Unix(deadline.Int64(), 0)

    // Get token addresses (last two parameters)
    tokenPathFrom := fmt.Sprintf("%x", data[164:196])[24:] // Take last 20 bytes
    tokenPathTo := fmt.Sprintf("%x", data[196:228])[24:]   // Take last 20 bytes

	// Debug prints	
	fmt.Println("swap transaction-----------------------------------------------------")
	fmt.Printf("Amount In: %v tokens\n", amountInDecimal)
	fmt.Printf("To Address: 0x%s\n", toAddress)
	// fmt.Printf("Deadline: %v (%v)\n", deadline, deadlineUnix)
	fmt.Printf("Token Path:\n")
	fmt.Printf("  From: 0x%s\n", tokenPathFrom)
	fmt.Printf("  To: 0x%s\n", tokenPathTo)
	fmt.Println("-----------------------------------------------------")
	return &entity.SwapTransaction{
		AmountIn:      amountInDecimal.String(),
		ToAddress:     toAddress,
		// Deadline:      deadlineUnix.Format(time.RFC3339),
		TokenPathFrom: tokenPathFrom,
		TokenPathTo:   tokenPathTo,
	}, nil
}
