// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const PREFIX_LUNA_HYBRID = "LUNA_HYBRID:"

const RSA_PUBLIC_KEY = "LUNA RSA PRIVATE KEY"
const RSA_PRIVATE_KEY = "LUNA RSA PUBLIC KEY"

// Storage implementation
type Storage interface {
	Save(filepath string, data []byte) error
	Load(filepath string) ([]byte, error)
}

// FileStorage implements the Storage interface for file-based storage
type FileStorage struct{}

func (fs *FileStorage) Save(filePath string, data []byte) error {

	_, fullFilePath, _, ok := runtime.Caller(0)

	if !ok {
		return fmt.Errorf("fail to get runtime caller info")
	}
	fileLocateDir := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(fullFilePath))))

	filePath = fmt.Sprintf("%s/%s", fileLocateDir, filePath)

	if err := os.MkdirAll(filepath.Dir(filePath), 0700); err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0700)
}

func (fs *FileStorage) Load(filepath string) ([]byte, error) {
	return os.ReadFile(filepath)
}

// GenerateKeyPair generates RSA key pair and saves the private key
func GenerateKeyPair(storage Storage, tenantID string) (string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", err
	}
	publicKey := &privateKey.PublicKey

	privateKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  RSA_PRIVATE_KEY,
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	filepath := fmt.Sprintf("private_keys/%s/private.pem", tenantID)

	if err := storage.Save(filepath, privateKeyPem); err != nil {
		return "", err
	}

	publicKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  RSA_PUBLIC_KEY,
		Bytes: x509.MarshalPKCS1PublicKey(publicKey),
	})
	return string(publicKeyPem), nil
}

// Encrypt encrypts text using hybrid RSA and AES encryption
func Encrypt(text string, publicKeyPem string) (string, error) {
	block, _ := pem.Decode([]byte(publicKeyPem))
	if block == nil || block.Type != RSA_PUBLIC_KEY {
		return "", errors.New("invalid public key")
	}

	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	aesKey := make([]byte, 16)
	if _, err := rand.Read(aesKey); err != nil {
		return "", err
	}

	blockAES, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(blockAES)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nil, nonce, []byte(text), nil)

	encAESKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, aesKey, nil)
	if err != nil {
		return "", err
	}

	encryptedData := append(encAESKey, append(nonce, ciphertext...)...)

	return PREFIX_LUNA_HYBRID + base64.StdEncoding.EncodeToString(encryptedData), nil
}

// Decrypt decrypts the hybrid encrypted data using the tenant's private key
func Decrypt(encryptedText string, tenantID string, storage Storage) (string, error) {
	if !bytes.HasPrefix([]byte(encryptedText), []byte(PREFIX_LUNA_HYBRID)) {
		return "", errors.New("invalid encrypted text format")
	}

	encryptedText = encryptedText[len(PREFIX_LUNA_HYBRID):]
	encryptedData, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	// Load private key
	filepath := fmt.Sprintf("./private_keys/%s/private.pem", tenantID)
	privateKeyPem, err := storage.Load(filepath)
	if err != nil {
		return "", err
	}

	block, _ := pem.Decode(privateKeyPem)
	if block == nil || block.Type != RSA_PRIVATE_KEY {
		return "", errors.New("invalid private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// Decrypt AES key using RSA
	keySize := privateKey.Size()
	if len(encryptedData) < keySize {
		return "", errors.New("invalid encrypted data")
	}
	encAESKey := encryptedData[:keySize]
	aesKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, encAESKey, nil)
	if err != nil {
		return "", err
	}

	// Decrypt text using AES
	blockAES, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(blockAES)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(encryptedData[keySize:]) < nonceSize {
		return "", errors.New("invalid encrypted data")
	}

	nonce := encryptedData[keySize : keySize+nonceSize]
	ciphertext := encryptedData[keySize+nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// func main() {
// 	storage := &FileStorage{}

// 	publicKey, err := GenerateKeyPair(storage, "tenant123")
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("Public Key:", publicKey)

// 	encryptedText, err := Encrypt("Hello, world!", publicKey)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("Encrypted Text:", encryptedText)

// 	decryptedText, err := Decrypt(encryptedText, "tenant123", storage)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("Decrypted Text:", decryptedText)
// }
