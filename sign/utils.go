package sign

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/e-n-0/sign-app-cli/provisioningprofiles"
	"howett.net/plist"
)

func locateAppFolder(inputFolder string) (string, error) {
	// Check if the input folder is a folder
	fileInfo, err := os.Stat(inputFolder)
	if err != nil {
		return "", fmt.Errorf("failed to stat the input folder: %s", err)
	}

	if !fileInfo.IsDir() {
		return "", fmt.Errorf("the input folder is not a folder")
	}

	appFolder := ""
	if files, err := ioutil.ReadDir(inputFolder); err == nil {
		for _, file := range files {
			if file.IsDir() {
				appFolder = filepath.Join(inputFolder, file.Name())
				break
			}
		}
	}

	if appFolder == "" {
		return "", fmt.Errorf("failed to find the app folder")
	}

	return appFolder, nil
}

func updateAppIdIfNeeded(appFolder string, profile provisioningprofiles.ProvisioningProfile) error {
	// Check if the app id is already set
	plistPath := filepath.Join(appFolder, "Info.plist")
	plistData, err := ioutil.ReadFile(plistPath)
	if err != nil {
		return fmt.Errorf("failed to read Info.plist file: %s", err)
	}

	var plistDataMap map[string]interface{}
	_, err = plist.Unmarshal(plistData, &plistDataMap)
	if err != nil {
		return fmt.Errorf("failed to unmarshal Info.plist file: %s", err)
	}

	// Check if the app id is already set
	if plistDataMap["CFBundleIdentifier"] == nil {
		return fmt.Errorf("the app id is not set in the Info.plist file")
	}

	bundleIdentifier := plistDataMap["CFBundleIdentifier"].(string)
	isWilcard := strings.HasSuffix(profile.AppID, "*")

	if !isWilcard && bundleIdentifier != profile.AppID {
		return fmt.Errorf("the app id in the Info.plist file (%s) is not the same as the one in the provisioning profile (%s)", bundleIdentifier, profile.AppID)
	}

	if isWilcard && !strings.HasPrefix(bundleIdentifier, strings.TrimSuffix(profile.AppID, "*")) {
		return fmt.Errorf("the app id in the Info.plist file (%s) is not the same as the one in the provisioning profile (%s)", bundleIdentifier, profile.AppID)
	}

	if isWilcard {
		// if bundleIdentifier != newApplicationID
		// TODO: Custom data via cmd

		profile.Update(bundleIdentifier)
	}

	return nil
}
