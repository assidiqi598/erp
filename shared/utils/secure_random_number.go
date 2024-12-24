package util

import (
	"crypto/rand"
	"math/big"
)

func GenerateSecureRandomNumber(digits int) int {
	if digits <= 0 {
		return 0
	}

	// Calculate the range
	max := new(big.Int)
	max.Exp(big.NewInt(10), big.NewInt(int64(digits)), nil).Sub(max, big.NewInt(1))
	min := new(big.Int)
	min.Exp(big.NewInt(10), big.NewInt(int64(digits-1)), nil)

	// Generate a random number in the range
	num, _ := rand.Int(rand.Reader, new(big.Int).Sub(max, min).Add(max, big.NewInt(1)))
	return int(num.Add(num, min).Int64())
}
