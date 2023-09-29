package main

import (
	log "github.com/sirupsen/logrus"
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
