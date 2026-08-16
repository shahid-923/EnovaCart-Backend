package helper

import (
	"crypto/rand"
	"errors"
	"math/big"
	"strconv"
)

func randomNumber(min, max int64) (int64, error) {
	if min > max {
		return 0, errors.New("invalid range")
	}

	n, err := rand.Int(
		rand.Reader,
		big.NewInt(max-min+1),
	)
	if err != nil {
		return 0, err
	}

	return n.Int64() + min, nil
}

func GenerateOTP() (int, error) {
	n, err := randomNumber(100000, 999999)
	if err != nil {
		return 0, err
	}

	return int(n), nil
}

func GenerateOrderID() (string, error) {
	n, err := randomNumber(1_000_000_000, 9_999_999_999)
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(n, 10), nil
}
