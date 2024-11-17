package util

import (
	"strconv"
	"time"

	"math/rand"
)

func GenerateRandomNumber() string {
	rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	randomNumber := rand.Intn(900000) + 100000
	return strconv.Itoa(randomNumber)
}
