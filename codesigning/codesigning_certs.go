package codesigning

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/e-n-0/sign-app-cli/utils"
)

func GetCodesigningCerts() ([]string, error) {
	var output []string
	bytes, status, err := utils.ExecuteProcess([]string{"/usr/bin/security", "find-identity", "-v", "-p", "codesigning"})
	if err != nil || status != 0 {
		return output, err
	}

	securityResult := string(bytes)
	rawResult := strings.Split(securityResult, "\"")
	for index := 0; index <= len(rawResult)-2; index += 2 {
		if !(len(rawResult)-1 < index+1) {
			output = append(output, rawResult[index+1])
		}
	}
	sort.Strings(output)
	return output, nil
}

func PrintCodesigningCerts(certs []string) {
	if len(certs) == 0 {
		fmt.Println("No codesigning certificates found.")
		fixSigningError()
		return
	}

	fmt.Println("Found", len(certs), "codesigning certificate"+(utils.Plural(len(certs)))+":")
	for _, cert := range certs {
		fmt.Println("  ", cert)
	}
}

// if there is no codesign certificates available on the system
// it might be the cause of missing "Apple Worldwide Developer Relations Certification Authority" certificate
// this certificate can be installed from https://www.apple.com/certificateauthority/AppleWWDRCAG3.cer
func fixSigningError() {
	if installed, err := checkAppleCertInstalled(); err != nil {
		fmt.Println("Failed to check if Apple Worldwide Developer Relations Certification Authority certificate is installed:", err)
	} else if !installed {
		// Print in yellow
		fmt.Println("\033[33m", "An issue has been detected with your codesigning certificates.", "\033[0m")
		fmt.Println("\033[33m", "Do you want to try to fix this issue by installing the Apple Worldwide Developer Relations Certification Authority certificate?", "\033[0m")
		if utils.AskForConfirmation("Do you want to install the certificate now?") {
			fixSigningError()
		}

		fmt.Println("Please try to run this command again to list your codesigning certificates.")
	}
}

// Download a file from url and save it to a temporary file
func downloadFile(url string, filename string) (string, error) {
	// Create the file
	out, err := os.CreateTemp("", "*-"+filename)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return out.Name(), nil
}

func checkAppleCertInstalled() (bool, error) {
	// Check if the certificate is installed
	_, status, err := utils.ExecuteProcess([]string{"/usr/bin/security", "find-certificate", "-c", "Apple Worldwide Developer Relations Certification Authority", "-a"})
	if err != nil || status != 0 {
		return false, err
	}

	return true, nil
}

func installAppleCert() error {

	// Download the certificate to a temporary file
	filePathTemp, err := downloadFile("https://www.apple.com/certificateauthority/AppleWWDRCAG3.cer", "AppleWWDRCAG3.cer")
	if err != nil {
		return err
	}

	// Install the certificate
	_, status, err := utils.ExecuteProcess([]string{"sudo", "/usr/bin/security", "add-trusted-cert", "-d", "-r", "trustRoot", "-k", "/Library/Keychains/System.keychain", filePathTemp})
	if err != nil || status != 0 {
		return err
	}

	// Remove the downloaded certificate
	os.Remove(filePathTemp)

	fmt.Println("Apple Worldwide Developer Relations Certification Authority certificate installed.")
	return nil
}
