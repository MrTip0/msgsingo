package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"math/big"
)

func rightSizePass(password *big.Int) []byte {
	pass := password.Bytes()
	if len(pass) >= 32 {
		pass = pass[:32]
	} else if len(pass) >= 16 {
		pass = pass[:16]
	} else {
		for len(pass) < 16 {
			pass = append(pass, '\x00')
		}
	}

	return pass
}

func gcmEncrypter(key []byte, data string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nonce, nonce, []byte(data), nil)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func gcmDecrypter(key []byte, data string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	text, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	nonceSize := aesgcm.NonceSize()
	if len(text) < nonceSize {
		return "", fmt.Errorf("corrupted data")
	}

	nonce, ciphertext := text[:nonceSize], text[nonceSize:]

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
