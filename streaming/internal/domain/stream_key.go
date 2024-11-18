package domain

import (
	"crypto/rand"
	"math/big"
)

type Key string

func CreateKey() Key {
	var b string
	count := 0
	for {
		v, err := generateRandomString(12)
		if err != nil {
			if count < 3 {
				count++
				continue
			} else {
				panic("failed generate random string more than 3 times")
			}
		}
		b = v
		break
	}

	return Key(b)
}

func generateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[randomIndex.Int64()]
	}
	return string(b), nil
}
