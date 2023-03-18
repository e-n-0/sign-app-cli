package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {
	// Run this program only on Macos or iOS (darwin architecture)
	if runtime.GOOS != "darwin" && runtime.GOOS != "ios" {
		fmt.Println("This program only runs on Macos or iOS (darwin architecture).")

		// exit the program with error code 1
		os.Exit(1)
	}

	profiles := getProfiles()

	// print how many items
	printProfiles(profiles)

	cs, _ := getCodesigningCerts()
	printCodesigningCerts(cs)
}
