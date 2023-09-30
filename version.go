package main

import "fmt"

var appVersion = "v0.0.0"
var appBuild = "0000000"
var appDate = "0/0/00"

func printver() {
	fmt.Printf("GoShield %s\n", appVersion)
	fmt.Printf("Build Number: %s\n", appBuild)
	fmt.Printf("Build Date: %s\n", appDate)
}
