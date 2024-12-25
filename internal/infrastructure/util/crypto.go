// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package util

import (
	"crypto/rand"
	"crypto/sha256"
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

func GenerateTextHash(text string) string {
	hashText := text + "None"
	hash := sha256.New()
	hash.Write([]byte(hashText))
	return hex.EncodeToString(hash.Sum(nil))
}
