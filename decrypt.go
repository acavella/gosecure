package main

import (
	"crypto/aes"
	"crypto/cipher"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/argon2"
)

func decryptFile() {
	// Reading ciphertext file
	inFile, err := filepath.Abs(flagFile)
	if err != nil {
		log.Fatalf("filepath error: %v", err.Error())
	}

	// Ensure input file is .enc type
	fileExt := filepath.Ext(inFile)
	if fileExt != ".enc" {
		log.Fatal("Input file expects .enc, provided:", fileExt)
	}

	// Read salt from file
	xfile, err := os.Open(inFile) // For read access.

	if err != nil {
		panic(err)
	}

	defer xfile.Close()

	xheadBytes := make([]byte, 32)
	m, err := xfile.Read(xheadBytes)
	if err != nil {
		panic(err)
	}

	salt := xheadBytes[:m]

	cipherText, err := os.ReadFile(inFile)
	if err != nil {
		log.Fatal(err)
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

	// Deattached nonce and decrypt
	nonce := cipherText[:gcm.NonceSize()]
	cipherText = cipherText[gcm.NonceSize():]
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		log.Fatalf("decrypt file err: %v", err.Error())
	}

	// Writing decryption content
	decFile := strings.TrimSuffix(inFile, filepath.Ext(inFile))
	err = os.WriteFile(decFile, plainText, 0777)
	if err != nil {
		log.Fatalf("write file err: %v", err.Error())
	} else {
		log.Info("Writing decrypted file:", decFile)
	}
}
