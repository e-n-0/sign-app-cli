package sign

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/e-n-0/sign-app-cli/codesigning"
	"github.com/e-n-0/sign-app-cli/utils"
)

func tryCodeSign(codesignCertificate string, tmpFolder string) error {
	// Copy own binary to tmp folder
	ownBinary := os.Args[0]
	testBinaryPath := filepath.Join(tmpFolder, "test-sign-file")
	err := utils.CopyFile(ownBinary, testBinaryPath)
	if err != nil {
		return err
	}

	// Try to sign the binary
	codeSign(testBinaryPath, codesignCertificate, "", "")

	// Check if the binary is signed
	_, status, err := utils.ExecuteProcess("codesign", "-v", testBinaryPath)
	os.Remove(testBinaryPath)
	if err != nil || status != 0 {
		if err == nil {
			err = fmt.Errorf("codesign failed with status code %d", status)
		}

		return err
	}

	return nil
}

// Try to sign an arbitrary file to test if the certificate is valid
func trySignCodeFail(tmpFolder string, codesignCertificate string) error {
	testTmpFolder := filepath.Join(tmpFolder, "test-codesign")
	err := os.Mkdir(testTmpFolder, 0755)
	if err != nil {
		return err
	}

	if err := tryCodeSign(codesignCertificate, testTmpFolder); err != nil {
		codesigning.FixSigningError()

		// Try again
		if err := tryCodeSign(codesignCertificate, testTmpFolder); err != nil {
			return fmt.Errorf("failed to resolve the codesigning issue: %s", err)
		}

		fmt.Println("Codesigning issue resolved")
	}
	os.RemoveAll(testTmpFolder)

	return nil
}
