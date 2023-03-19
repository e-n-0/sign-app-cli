package sign

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/e-n-0/sign-app-cli/provisioningprofiles"
	"github.com/e-n-0/sign-app-cli/utils"

	"howett.net/plist"
)

type SignerParams struct {
	ProvisioninngProfile provisioningprofiles.ProvisioningProfile
	CodesignCertificate  string
	InputFile            string
	OutputFile           string
	EntitlementsFile     string
}

var validBinariesExtensions = []string{".app", ".framework", ".dylib", ".appex", ".so", "0", ".vis", ".pvr"}

// Sign all the files in the folder recursively
func signPath(folderOrFolderPath string, params SignerParams) error {

	// 1 - Recursively sign the folder content
	if utils.IsFolder(folderOrFolderPath) {
		entries, err := os.ReadDir(folderOrFolderPath)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			if entry.IsDir() {
				err := signPath(filepath.Join(folderOrFolderPath, entry.Name()), params)
				if err != nil {
					return err
				}
			}
		}
	}

	// 2 - Check the file extension
	fileExtension := filepath.Ext(folderOrFolderPath)
	if !utils.StringInSlice(fileExtension, validBinariesExtensions) {
		return nil
	}

	// 3 - Sign the path (file or folder)
	// The folder executable (if bundle) must be signed after all its content is signed
	err := codeSign(folderOrFolderPath, params.CodesignCertificate, params.EntitlementsFile, params.ProvisioninngProfile.Path)
	if err != nil {
		return err
	}

	return nil
}

func Sign(params SignerParams) error {
	fmt.Println("Starting the signing process for file:", params.InputFile)

	tmpFolder, err := os.MkdirTemp("", "sign-app-cli-*")
	if err != nil {
		return fmt.Errorf("failed to create temporary folder: %s", err)
	}
	defer os.RemoveAll(tmpFolder)

	// Try to sign an arbitrary file to test if the certificate is valid
	err = trySignCodeFail(tmpFolder, params.CodesignCertificate)
	if err != nil {
		return err
	}

	filenameExt := filepath.Ext(params.InputFile)

	// Create a folder "work" inside the tmp folder
	workingTmpFolder := filepath.Join(tmpFolder, "work")
	err = os.Mkdir(workingTmpFolder, 0755)
	if err != nil {
		return err
	}

	switch filenameExt {
	case ".ipa":
		// Unzip the ipa file
		fmt.Println("Extracting ipa file...")
		err = utils.ExtractZip(params.InputFile, workingTmpFolder)
		if err != nil {
			return err
		}

		// Get the Payload folder
		payloadFolder := filepath.Join(workingTmpFolder, "Payload")

		// Get the app folder
		appFolder, err := locateAppFolder(payloadFolder)
		if err != nil {
			return err
		}

		// Retreive entitlements from the provisioning profile and save it to a file
		/*err = updateAppIdIfNeeded(appFolder, params.ProvisioninngProfile)
		if err != nil {
			return err
		}*/

		entitlements := params.ProvisioninngProfile.GetEntitlements()

		params.EntitlementsFile = filepath.Join(workingTmpFolder, "entitlements.plist")
		entitlementsFile, err := os.Create(params.EntitlementsFile)
		if err != nil {
			return err
		}

		// Encode the entitlements to xml
		encoder := plist.NewEncoder(entitlementsFile)
		encoder.Indent("\t")
		err = encoder.Encode(entitlements)
		if err != nil {
			return err
		}

		entitlementsFile.Close()
		////

		// Sign the app folder
		fmt.Println("Signing the app folder...")
		err = signPath(appFolder, params)
		if err != nil {
			return err
		}

		// Zip the Payload folder and save it to the output file
		err = utils.CreateZip(params.OutputFile, payloadFolder)
		if err != nil {
			return err
		}

		// Print in green
		fmt.Println("\033[32m" + "Successfully signed the ipa file" + "\033[0m")
		fmt.Println("The signed ipa file is available at:", params.OutputFile)
		return nil

	default:
		// Not supported file type
		return fmt.Errorf("unsupported file type: %s", filenameExt)
	}

}

func codeSign(inputFile string, codesignCertificate string, entitlementsFile string, mobileProvisionFile string) error {
	fmt.Println("Signing", inputFile, "...")

	fileExt := filepath.Ext(inputFile)
	filePath := inputFile
	useEntitlements := false

	// Get the executable file from different file types
	switch fileExt {
	case ".framework":
		fileName := filepath.Base(inputFile)
		fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))
		filePath = filepath.Join(inputFile, fileName)
	case ".app", ".appex":
		// Read executable file from Info.plist
		infoPlist := filepath.Join(inputFile, "Info.plist")
		plistBytes, err := os.ReadFile(infoPlist)
		if err != nil {
			return fmt.Errorf("failed to read Info.plist file, error: %s", err)
		}

		// Read the plist file into a map
		var plistData map[string]interface{}
		_, err = plist.Unmarshal(plistBytes, &plistData)
		if err != nil {
			return err
		}

		// Get the executable file
		executableFile := plistData["CFBundleExecutable"].(string)
		filePath = filepath.Join(inputFile, executableFile)

		// Check if the entitlements file exists
		if entitlementsFile != "" {
			_, err := os.Stat(entitlementsFile)
			if err != nil {
				return err
			}

			useEntitlements = true
		}

		mobileProvisionAppPath := filepath.Join(inputFile, "embedded.mobileprovision")

		// delete an existing embedded.mobileprovision file
		_ = os.Remove(filepath.Join(inputFile, "embedded.mobileprovision"))

		// copy the provisioning profile to the app folder
		err = utils.CopyFile(mobileProvisionFile, mobileProvisionAppPath)
		if err != nil {
			return fmt.Errorf("failed to copy the provisioning profile to the app folder: %s", err)
		}
	}

	// Create the arguments
	args := []string{"codesign", "-f", "-s", codesignCertificate, "--generate-entitlement-der"}
	if useEntitlements {
		args = append(args, "--entitlements", entitlementsFile)
	}
	args = append(args, filePath)

	// Execute the codesign command
	_, status, err := utils.ExecuteProcess(args...)
	if err != nil || status != 0 {
		if err == nil {
			err = fmt.Errorf("codesign failed with status code %d", status)
		}

		return err
	}

	return nil
}
