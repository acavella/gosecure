package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/argon2"
)

var (
	CryptPw     string
	flagFile    string
	decryptPtr  bool
	encryptPtr  bool
	verbose     bool
	debugPtr    bool
	absPath     string
	fileName    string
	workingPath string
)

func init() {

	flag.BoolVar(&encryptPtr, "e", false, "Encrypt the input data.")
	flag.BoolVar(&decryptPtr, "d", false, "Decrypt the input data.")
	flag.StringVar(&flagFile, "in", "", "The input filename, standard input by default.")
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

	getfiles()
	printver()

}

func main() {

	log.Trace("Password:", CryptPw)

	if encryptPtr {
		log.Info("Encrypting file:", fileName)
		encryptFile()
	} else if decryptPtr {
		log.Info("Decrypting file:", fileName)
		decryptFile()
	}
}

func getfiles() (absPath, workingPath, fileName string) {
	absPath, err := filepath.Abs(flagFile)
	if err != nil {
		log.Fatalf("file error: %v", err.Error())
	}
	workingPath, fileName = filepath.Split(absPath)
	return
}

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

func extractciphertext() {
	// Reading ciphertext file
	inFile, err := filepath.Abs(flagFile)
	if err != nil {
		log.Fatalf("filepath error: %v", err.Error())
	}

	fin, err := os.Open(inFile)
	if err != nil {
		panic(err)
	}
	defer fin.Close()

	fout, err := os.Create("dest.txt")
	if err != nil {
		panic(err)
	}
	defer fout.Close()

	// Offset is the number of bytes you want to exclude
	_, err = fin.Seek(32, io.SeekStart)
	if err != nil {
		panic(err)
	}

	n, err := io.Copy(fout, fin)
	fmt.Printf("Copied %d bytes, err: %v", n, err)
}

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

	extractciphertext()

	cipherText, err := os.ReadFile("dest.txt")
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
