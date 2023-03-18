package main

import "fmt"

func contains(slice []string, item string) bool {
	for _, element := range slice {
		if element == item {
			return true
		}
	}
	return false
}

func plural(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}

func askForConfirmation(s string) bool {
	var response string
	for {
		fmt.Printf("%s [y/n]: ", s)
		_, err := fmt.Scanln(&response)
		if err != nil {
			return false
		}
		if response == "y" || response == "n" {
			break
		}
	}
	return response == "y"
}
