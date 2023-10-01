package main

import (
	log "github.com/sirupsen/logrus"
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
