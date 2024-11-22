package util

import (
	"strconv"
	"strings"
	"time"

	"math/rand"
)

func GenerateRandomNumber() string {
	rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	randomNumber := rand.Intn(900000) + 100000
	return strconv.Itoa(randomNumber)
}

// generateString generates a random alphanumeric string of length n
func GenerateRandomString(n int) (string, error) {
	const lettersDigits = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var result strings.Builder
	result.Grow(n)

	for i := 0; i < n; i++ {
		if err := result.WriteByte(lettersDigits[rand.Intn(len(lettersDigits))]); err != nil {
			return "", err
		}
	}

	return result.String(), nil
}
