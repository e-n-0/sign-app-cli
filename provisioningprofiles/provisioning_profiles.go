package provisioningprofiles

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/e-n-0/sign-app-cli/utils"
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
	Path         string
}

func sortProfilesByCreationDateAndName(profiles []ProvisioningProfile) {
	sort.Slice(profiles, func(i, j int) bool {
		return profiles[i].Created.After(profiles[j].Created) || (profiles[i].Created.Equal(profiles[j].Created) && profiles[i].Name < profiles[j].Name)
	})
}

func PrintProfiles(profiles []ProvisioningProfile) {
	if len(profiles) == 0 {
		fmt.Println("No provisioning profiles found")
		return
	}

	fmt.Println("Found", len(profiles), "provisioning profile"+(utils.Plural(len(profiles)))+":")
	for _, profile := range profiles {
		fmt.Printf("  %s (%s)", profile.Name, profile.TeamID)

		// Print in red "EXPIRED" if the profile is expired
		if profile.Expires.Before(time.Now()) {
			fmt.Printf("\033[31m%s\033[0m", " !EXPIRED!")
		}

		fmt.Println()
	}
}

func GetProfile(name string) (ProvisioningProfile, error) {
	profiles := GetProfiles()
	for _, profile := range profiles {
		if fmt.Sprintf("%s (%s)", profile.Name, profile.TeamID) == name {
			return profile, nil
		}
	}
	return ProvisioningProfile{}, fmt.Errorf("failed to find provisioning profile with name: %s", name)
}

// Convert the swift code to Go code
func GetProfiles() []ProvisioningProfile {
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
						profile.Path = filepath.Join(provisioningProfilesPath, file.Name())
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
		if !utils.Contains(names, profile.Name+profile.AppID) {
			newProfiles = append(newProfiles, profile)
			names = append(names, profile.Name+profile.AppID)
		}
	}

	return newProfiles
}

type mobileProvision struct {
	ExpirationDate time.Time              `xml:"ExpirationDate"`
	CreationDate   time.Time              `xml:"CreationDate"`
	Name           string                 `xml:"Name"`
	Entitlements   map[string]interface{} `xml:"Entitlements"`
	_              map[string]interface{} `xml:",any"`
}

func CreateProvisioningProfile(filename string) (ProvisioningProfile, error) {
	var provisioningProfile ProvisioningProfile

	// Execute the security command
	bytes, status, err := utils.ExecuteProcess("/usr/bin/security", "cms", "-D", "-i", filename)
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
		var mobileProvision mobileProvision
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

func (profile ProvisioningProfile) RemoveGetTaskAllow() {
	delete(profile.Entitlements, "get-task-allow")
}

func (profile ProvisioningProfile) Update(trueAppID string) {
	if _, ok := profile.Entitlements["application-identifier"].(string); ok {
		newIdentifier := profile.TeamID + "." + trueAppID
		profile.Entitlements["application-identifier"] = newIdentifier
	}
}

func (profile ProvisioningProfile) GetEntitlements() map[string]interface{} {
	return profile.Entitlements
}
