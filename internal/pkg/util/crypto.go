package util

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// GenerateRefreshToken generates a secure random refresh token of the specified length
func GenerateRefreshToken(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("invalid length")
	}

	tokenBytes := make([]byte, length)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %v", err)
	}

	token := hex.EncodeToString(tokenBytes)

	return token, nil
}
