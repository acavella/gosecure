package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/argon2"
)

// Generate a random salt return the result as byte
func generateRandomSalt() []byte {
	var result = make([]byte, 32)

	_, err := rand.Read(result[:])

	if err != nil {
		log.Fatalf("error generating salt: %v", err.Error())
	}

	return result
}

func encryptFile() {

	// Generate salt
	salt := generateRandomSalt()
	log.Trace("Salt:", salt)

	absPath, err := filepath.Abs(flagFile)
	if err != nil {
		log.Fatalf("file error: %v", err.Error())
	}

	plainText, err := os.ReadFile(absPath)
	if err != nil {
		log.Fatalf("read file err: %v", err.Error())
	}

	// Generating derivative key
	dk := argon2.IDKey([]byte(CryptPw), salt, 3, 64*1024, 4, 32)
	log.Trace("Derived Key:", dk)

	// Creating block of algorithm
	block, err := aes.NewCipher(dk)
	if err != nil {
		log.Fatalf("cipher err: %v", err.Error())
	}

	// Creating GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatalf("cipher GCM err: %v", err.Error())
	}

	// Generating random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatalf("nonce  err: %v", err.Error())
	}
	log.Trace("Nonce:", nonce)

	// Decrypt file
	cipherText := gcm.Seal(nonce, nonce, plainText, nil)

	// Writing ciphertext file
	encFile := absPath + ".enc"

	// Writing IV to file
	err = os.WriteFile(encFile, salt, 0777)
	if err != nil {
		log.Fatalf("write file err: %v", err.Error())
	} else {
		log.Info("Writing salt to file:", encFile)
	}

	f, err := os.OpenFile(encFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		log.Fatalf("open file err: %v", err.Error())
	}

	defer f.Close()

	_, err2 := f.Write(cipherText)
	if err2 != nil {
		log.Fatalf("write file err: %v", err.Error())
	} else {
		log.Info("Writing ciphertext to file:", encFile)
	}

	/*
		err = os.WriteFile(encFile, cipherText, 0777)
		if err != nil {
			log.Fatalf("write file err: %v", err.Error())
		} else {
			log.Info("Writing encrypted file:", encFile)
		}
	*/
}
