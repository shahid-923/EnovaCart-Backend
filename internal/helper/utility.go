package helper

import (
	"crypto/rand"
	"math/big"
)

func RandomNumbers(length int) (int, error) {

	min := 10000000
	max := 99999999

	n, err := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	if err != nil {
		return 0, err
	}

	return int(n.Int64()) + min, nil
}