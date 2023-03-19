/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/e-n-0/sign-app-cli/codesigning"
	"github.com/e-n-0/sign-app-cli/provisioningprofiles"
	"github.com/e-n-0/sign-app-cli/sign"
	"github.com/e-n-0/sign-app-cli/utils"
	"github.com/spf13/cobra"
)

var (
	provisioningProfileName string
	provisioningProfilePath string
	codesigningCertName     string

	inputFile  string
	outputFile string

	entitlementsFile string
)

// signCmd represents the sign command
var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Sign the provided file",
	Long: `
Sign the provided file with the provided provisioning profile and codesigning certificate.
The file can be an .ipa or a .app.
`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if the input file exists
		if !utils.FileExists(inputFile) {
			panic("The input file does not exist")
		}

		// Check provisioning profile name
		var provisioningProfile provisioningprofiles.ProvisioningProfile
		if provisioningProfileName != "" {
			// Check if the provisioning profile exists
			p, err := provisioningprofiles.GetProfile(provisioningProfileName)
			if err != nil {
				end(err)
			}
			provisioningProfile = p
		} else if provisioningProfilePath != "" {
			// Check if the provisioning profile exists
			if !utils.FileExists(provisioningProfilePath) {
				end(fmt.Errorf("the provisioning profile does not exist"))
			}

			p, err := provisioningprofiles.CreateProvisioningProfile(provisioningProfilePath)
			if err != nil {
				end(err)
			}

			provisioningProfile = p
		} else {
			end(fmt.Errorf("you must provide a provisioning profile"))
		}

		// Check if the codesigning certificate exists
		codesignCert, err := codesigning.GetCodesigningCert(codesigningCertName)
		if err != nil {
			end(err)
		}

		// Check if the entitlements file exists
		if entitlementsFile != "" && !utils.FileExists(entitlementsFile) {
			end(fmt.Errorf("the entitlements file does not exist"))
		}

		err = sign.Sign(sign.SignerParams{
			ProvisioninngProfile: provisioningProfile,
			CodesignCertificate:  codesignCert,
			InputFile:            inputFile,
			OutputFile:           outputFile,
			EntitlementsFile:     entitlementsFile,
		})

		if err != nil {
			panic(err)
		}
	},
}

func end(err error) {
	fmt.Println("error:", err)
	os.Exit(1)
}

func init() {
	rootCmd.AddCommand(signCmd)

	// Add cobra command
	signCmd.Flags().StringVarP(&provisioningProfileName, "profile", "p", "", "The name of the provisioning profile to use installed on the machine (list with 'sign-app-cli listProvisioningProfiles')")
	signCmd.Flags().StringVarP(&provisioningProfilePath, "profilePath", "P", "", "The path of the provisioning profile to use")
	signCmd.Flags().StringVarP(&codesigningCertName, "certificate", "c", "", "The name of the codesigning certificate to use installed on the machine (list with 'sign-app-cli listCodesigningCerts')")
	signCmd.Flags().StringVarP(&inputFile, "input", "i", "", "The path of the file to sign")
	signCmd.Flags().StringVarP(&outputFile, "output", "o", "", "The path of the signed file")
	signCmd.Flags().StringVarP(&entitlementsFile, "entitlements", "e", "", "The path of the entitlements file to use")

	signCmd.MarkFlagFilename("profilePath")
	signCmd.MarkFlagFilename("input")
	signCmd.MarkFlagFilename("output")
	signCmd.MarkFlagFilename("entitlements")

	signCmd.MarkFlagRequired("certificate")
	signCmd.MarkFlagRequired("input")
	signCmd.MarkFlagRequired("output")
	//signCmd.MarkFlagRequired("profile")
	//signCmd.MarkFlagRequired("profilePath")

	signCmd.MarkFlagsMutuallyExclusive("profile", "profilePath")
}
