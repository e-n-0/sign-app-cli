/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/e-n-0/sign-app-cli/codesigning"
	"github.com/spf13/cobra"
)

// listCodesigningCertsCmd represents the listCodesigningCerts command
var listCodesigningCertsCmd = &cobra.Command{
	Use:   "listCodesigningCerts",
	Short: "List all codesigning certificates available in your keychain",
	Long: `
This command will list all codesigning certificates available in your keychain.
You can use this command to find the name of the certificate you want to use.
Then you can use the name to sign your app.
For example:
$ sign-app-cli sign [...] --certificate <name>`,
	Run: func(cmd *cobra.Command, args []string) {
		certs, err := codesigning.GetCodesigningCerts()
		if err != nil {
			panic(err)
		}

		codesigning.PrintCodesigningCerts(certs)
	},
}

func init() {
	rootCmd.AddCommand(listCodesigningCertsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCodesigningCertsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCodesigningCertsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
