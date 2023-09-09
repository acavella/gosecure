package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	CryptPw    string
	FilePath   string
	decryptPtr bool
	encryptPtr bool
)

func init() {

	flag.BoolVar(&encryptPtr, "e", false, "Flag sets mode to encrypt")
	flag.BoolVar(&encryptPtr, "encrypt", false, "Flag sets mode to encrypt")
	flag.BoolVar(&decryptPtr, "d", false, "Flag sets mode to decrypt")
	flag.BoolVar(&decryptPtr, "decrypt", false, "Flag sets mode to decrypt")
	flag.StringVar(&FilePath, "p", "", "Path to file")
	flag.StringVar(&FilePath, "path", "", "Path to file")
	flag.StringVar(&CryptPw, "k", "", "Password used to encrypt/decrypt")
	flag.StringVar(&CryptPw, "pass", "", "Password used to encrypt/decrypt")
	flag.Parse()

}

func main() {

	basePath, err := filepath.Abs(FilePath)
	if err != nil {
		log.Fatalf("cipher err: %v", err.Error())
	}
	workDir, fileName := filepath.Split(basePath)

	fmt.Println("Base directory:", basePath)
	fmt.Println("The file dir is:", workDir)
	fmt.Println("The file name is:", fileName)
	fmt.Println("Password:", CryptPw)

	if encryptPtr {
		log.Println("Encrypting file:", fileName)
		encryptFile()
	} else if decryptPtr {
		log.Println("Decrypting file:", fileName)
		decryptFile()
	}
}

func encryptFile() {
	// Reading plaintext file
	absPath, err := filepath.Abs(FilePath)
	if err != nil {
		log.Fatalf("filepath error: %v", err.Error())
	}

	plainText, err := os.ReadFile(absPath)
	if err != nil {
		log.Fatalf("read file err: %v", err.Error())
	}

	// Generating derivative key
	dk := argon2.IDKey([]byte(CryptPw), []byte("c123bdb6574e817ac0a5f8b2e097b986"), 3, 64*1024, 4, 32)
	fmt.Println("Derived Key:", dk)

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
	fmt.Println("Nonce:", nonce)

	// Decrypt file
	cipherText := gcm.Seal(nonce, nonce, plainText, nil)

	// Writing ciphertext file
	encFile := absPath + ".enc"
	err = os.WriteFile(encFile, cipherText, 0777)
	if err != nil {
		log.Fatalf("write file err: %v", err.Error())
	}

}

func decryptFile() {
	// Reading ciphertext file
	absPath, err := filepath.Abs(FilePath)
	if err != nil {
		log.Fatalf("filepath error: %v", err.Error())
	}

	cipherText, err := os.ReadFile(absPath)
	if err != nil {
		log.Fatal(err)
	}

	// Generating derivative key
	dk := argon2.IDKey([]byte(CryptPw), []byte("c123bdb6574e817ac0a5f8b2e097b986"), 3, 64*1024, 4, 32)
	fmt.Println("Derived Key:", dk)

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
	decFile := strings.TrimSuffix(absPath, filepath.Ext(absPath))
	err = os.WriteFile(decFile, plainText, 0777)
	if err != nil {
		log.Fatalf("write file err: %v", err.Error())
	}
}
