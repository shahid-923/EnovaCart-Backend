package helper

import (
	"crypto/rand"
	"math/big"
)

func GenerateOTP() (int, error) {
	min := 100000
	max := 999999

	n, err := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	if err != nil {
		return 0, err
	}

	return int(n.Int64()) + min, nil
}

func GenerateOrderID() (int64, error) {
	min := int64(1000000000)
	max := int64(9999999999)

	n, err := rand.Int(rand.Reader, big.NewInt(max-min+1))
	if err != nil {
		return 0, err
	}

	return n.Int64() + min, nil
}
