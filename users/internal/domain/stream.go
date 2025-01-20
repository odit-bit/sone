package domain

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/segmentio/ksuid"
)

type Stream struct {
	Key   string
	ID    string
	Title string

	CreatedAt       time.Time
	LastStartStream time.Time
	LastEndStream   time.Time
	IsStreaming     bool
}

func NewStream(title string) Stream {
	key := generateKey()

	id := ksuid.New()
	if title == "" {
		title = fmt.Sprintf("video-%s", generateKey())
	}

	s := Stream{
		Key:       key,
		ID:        id.String(),
		Title:     title,
		CreatedAt: time.Now(),
	}
	return s
}

// type Key string

func GenerateKey() string {
	return generateKey()
}

func generateKey() string {
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

	return b
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
