package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"howett.net/plist"
)

type ProvisioningProfile struct {
	Filename     string
	Name         string
	Created      time.Time
	Expires      time.Time
	AppID        string
	TeamID       string
	Entitlements map[string]interface{}
}

func sortProfilesByCreationDateAndName(profiles []ProvisioningProfile) {
	sort.Slice(profiles, func(i, j int) bool {
		return profiles[i].Created.After(profiles[j].Created) || (profiles[i].Created.Equal(profiles[j].Created) && profiles[i].Name < profiles[j].Name)
	})
}

func printProfiles(profiles []ProvisioningProfile) {
	if len(profiles) == 0 {
		fmt.Println("No provisioning profiles found")
		return
	}

	fmt.Println("Found", len(profiles), "provisioning profile"+(plural(len(profiles)))+":")
	for _, profile := range profiles {
		fmt.Printf("  %s (%s)", profile.Name, profile.TeamID)

		// Print in red "EXPIRED" if the profile is expired
		if profile.Expires.Before(time.Now()) {
			fmt.Printf("\033[31m%s\033[0m", " !EXPIRED!")
		}

		fmt.Println()
	}
}

// Convert the swift code to Go code
func getProfiles() []ProvisioningProfile {
	var output []ProvisioningProfile

	// Get Provisioning Profiles from the file system
	// these are located in ~/Library/MobileDevice/Provisioning Profiles
	if libraryDirectory, err := os.UserHomeDir(); err == nil {
		provisioningProfilesPath := filepath.Join(libraryDirectory, "Library/MobileDevice/Provisioning Profiles")
		if files, err := ioutil.ReadDir(provisioningProfilesPath); err == nil {
			for _, file := range files {
				if filepath.Ext(file.Name()) == ".mobileprovision" {
					profileFilename := filepath.Join(provisioningProfilesPath, file.Name())
					profile, err := CreateProvisioningProfile(profileFilename)
					if err != nil {
						fmt.Println(err)
					} else {
						output = append(output, profile)
					}
				}
			}
		}
	}

	// Sort the profiles by creation date
	sortProfilesByCreationDateAndName(output)

	// Remove duplicates
	var newProfiles []ProvisioningProfile
	var names []string
	for _, profile := range output {
		if !contains(names, profile.Name+profile.AppID) {
			newProfiles = append(newProfiles, profile)
			names = append(names, profile.Name+profile.AppID)
		}
	}

	return newProfiles
}

type MobileProvision struct {
	ExpirationDate time.Time              `xml:"ExpirationDate"`
	CreationDate   time.Time              `xml:"CreationDate"`
	Name           string                 `xml:"Name"`
	Entitlements   map[string]interface{} `xml:"Entitlements"`
	_              map[string]interface{} `xml:",any"`
}

func CreateProvisioningProfile(filename string) (ProvisioningProfile, error) {
	var provisioningProfile ProvisioningProfile
	fmt.Println("Creating provisioning profile from file: " + filename)

	// Create the security command
	securityArgs := []string{"/usr/bin/security", "cms", "-D", "-i", filename}

	// Execute the security command
	bytes, status, err := executeProcess(securityArgs)
	if err != nil {
		return ProvisioningProfile{}, err
	}

	if status == 0 {
		// Get the output of the security command
		output := string(bytes)

		// Get the xml start tag
		xmlIndex := strings.Index(string(output), "<?xml")

		// Get the raw xml
		rawXML := string(output)[xmlIndex:]

		// Parse the plist
		var mobileProvision MobileProvision
		_, err := plist.Unmarshal([]byte(rawXML), &mobileProvision)
		if err != nil {
			return ProvisioningProfile{}, err
		}

		// Fill the provisioning profile information
		appID := mobileProvision.Entitlements["application-identifier"].(string)
		periodIndex := strings.Index(appID, ".")
		provisioningProfile.AppID = appID[periodIndex+1:]
		provisioningProfile.TeamID = appID[:periodIndex]

		provisioningProfile.Filename = filename
		provisioningProfile.Expires = mobileProvision.ExpirationDate
		provisioningProfile.Created = mobileProvision.CreationDate
		provisioningProfile.Name = mobileProvision.Name
		provisioningProfile.Entitlements = mobileProvision.Entitlements
	}

	return provisioningProfile, nil
}

func (profile ProvisioningProfile) removeGetTaskAllow() {
	delete(profile.Entitlements, "get-task-allow")
}

func (profile ProvisioningProfile) update(trueAppID string) {
	if _, ok := profile.Entitlements["application-identifier"].(string); ok {
		newIdentifier := profile.TeamID + "." + trueAppID
		profile.Entitlements["application-identifier"] = newIdentifier
	}
}
