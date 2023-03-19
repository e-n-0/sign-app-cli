package utils

import (
	"fmt"
	"io/ioutil"
)

func Contains(slice []string, item string) bool {
	for _, element := range slice {
		if element == item {
			return true
		}
	}
	return false
}

func Plural(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}

func AskForConfirmation(s string) bool {
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

func CopyFile(src string, dst string) error {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dst, input, 0644)
	if err != nil {
		return err
	}

	return nil
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
