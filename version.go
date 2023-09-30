package main

import "fmt"

var appVersion = "No Version Provided"
var appBuild = "No Build Provided"
var appDate = "No Build Date"

func printver() {
	fmt.Printf("GoShield %s\n", appVersion)
	fmt.Printf("Build Number: %s\n", appBuild)
	fmt.Printf("Build Date: %s\n", appDate)
}
