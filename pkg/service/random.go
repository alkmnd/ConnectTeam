package service

import (
	"crypto/rand"
	"math/big"
)

func randomCrypto() int {
	randomDigit, err := rand.Int(rand.Reader, big.NewInt(10))
	if err != nil {
		return 0
	}

	return int(randomDigit.Int64())
}