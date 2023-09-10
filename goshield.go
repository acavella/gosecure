package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"flag"
	"io"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"golang.org/x/crypto/argon2"
)

var (
	CryptPw    string
	FilePath   string
	decryptPtr bool
	encryptPtr bool
	verbose    bool
	debugPtr   bool
)

func init() {

	flag.BoolVar(&encryptPtr, "e", false, "Encrypt the input data.")
	flag.BoolVar(&decryptPtr, "d", false, "Decrypt the input data.")
	flag.StringVar(&FilePath, "in", "", "The input filename, standard input by default.")
	flag.StringVar(&CryptPw, "k", "", "The password to derive the key from.")
	flag.BoolVar(&verbose, "v", false, "Enables verbosity to default logger")
	flag.BoolVar(&debugPtr, "debug", false, "Enables debug output to default logger")
	flag.Parse()

	if verbose {
		log.SetLevel(log.InfoLevel)
	} else if debugPtr {
		log.SetLevel(log.TraceLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

}

func main() {

	basePath, err := filepath.Abs(FilePath)
	if err != nil {
		log.Fatalf("cipher err: %v", err.Error())
	}
	workDir, fileName := filepath.Split(basePath)

	log.Debug("Base directory:", basePath)
	log.Debug("The file dir is:", workDir)
	log.Debug("The file name is:", fileName)
	log.Trace("Password:", CryptPw)

	if encryptPtr {
		log.Info("Encrypting file:", fileName)
		encryptFile()
	} else if decryptPtr {
		log.Info("Decrypting file:", fileName)
		decryptFile()
	}
}

func encryptFile() {
	// Reading plaintext file
	absPath, err := filepath.Abs(FilePath)
	if err != nil {
		log.Fatalf("filepath error: %v", err.Error())
	}

	// Generate random salt
	salt := make([]byte, 32)
	log.Trace("Salt:", salt)

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
	decFile := strings.TrimSuffix(absPath, filepath.Ext(absPath))
	err = os.WriteFile(decFile, plainText, 0777)
	if err != nil {
		log.Fatalf("write file err: %v", err.Error())
	} else {
		log.Info("Writing decrypted file:", decFile)
	}
}
