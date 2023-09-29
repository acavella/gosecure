package main

import (
	"flag"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

var (
	CryptPw    string
	flagFile   string
	decryptPtr bool
	encryptPtr bool
	verbose    bool
	debugPtr   bool
	inFile     string
	fileName   string
	workDir    string
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

func getfiles() (inFile, fileName, workDir string) {
	inFile, err := filepath.Abs(flagFile)
	if err != nil {
		log.Fatalf("file error: %v", err.Error())
	}
	workDir, fileName = filepath.Split(inFile)
	return
}
