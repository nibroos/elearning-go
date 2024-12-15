package scheduler

import (
	"log"
	"math/rand"
	"time"
)

func GenerateRandomString() {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 10
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	log.Printf("Generated Random String: %s", string(b))
}

func GenerateRandomNumber() {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	number := seededRand.Intn(100)
	log.Printf("Generated Random Number: %d", number)
}
