/*
Copyright © 2023 Flavien Darche 'en0'
*/
package cmd

import (
	"github.com/e-n-0/sign-app-cli/provisioningprofiles"
	"github.com/spf13/cobra"
)

// listProvisioningProfilesCmd represents the listProvisioningProfiles command
var listProvisioningProfilesCmd = &cobra.Command{
	Use:   "listProvisioningProfiles",
	Short: "List all provisioning profiles available in your keychain",
	Long: `
This command will list all provisioning profiles available in your keychain.
You can use this command to find the name of the provisioning profile you want to use.
Then you can use the name to sign your app.
For example:
$ sign-app-cli sign [...] --provisioning-profile <name>`,
	Run: func(cmd *cobra.Command, args []string) {
		profiles := provisioningprofiles.GetProfiles()
		provisioningprofiles.PrintProfiles(profiles)
	},
}

func init() {
	rootCmd.AddCommand(listProvisioningProfilesCmd)
}
