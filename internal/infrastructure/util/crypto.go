// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package util

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
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

func EncodeFloat32ToBase64(vectors []float32) (string, error) {
	buf := new(bytes.Buffer)

	for _, vector := range vectors {
		if err := binary.Write(buf, binary.LittleEndian, vector); err != nil {
			return "", err
		}
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func DecodeBase64ToFloat32(encodedString string) ([]float32, error) {
	data, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewReader(data)
	var vector []float32
	for buf.Len() > 0 {
		var val float32
		if err := binary.Read(buf, binary.LittleEndian, &val); err != nil {
			return nil, err
		}

		if errors.Is(err, io.EOF) {
			break
		}
		vector = append(vector, val)
	}

	return vector, nil
}
