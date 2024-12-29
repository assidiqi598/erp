package utils

import (
	"crypto/rand"
	"math/big"
)

func GenerateSecureRandomNumber(digits int) int {
	if digits <= 0 {
		return 0
	}

	// Calculate the minimum and maximum values based on the number of digits
	min := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(digits-1)), nil) // 10^(digits-1)
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(digits)), nil)   // 10^digits

	// Generate a random number in the range [min, max)
	num, _ := rand.Int(rand.Reader, new(big.Int).Sub(max, min)) // max - min is the range

	// Add min to ensure the number is in the correct range
	num.Add(num, min)

	// Convert to int (assuming the number fits in an int)
	return int(num.Int64())
}
