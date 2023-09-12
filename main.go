package main

import (
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func main() {

	inFile, err := filepath.Abs(flagFile)
	if err != nil {
		log.Fatalf("cipher err: %v", err.Error())
	}
	workDir, fileName := filepath.Split(inFile)

	log.Debug("Base directory:", inFile)
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
