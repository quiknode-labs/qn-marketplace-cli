/*
Copyright © 2023 QuickNode, Inc.

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "qn-marketplace-cli",
	Short: "A command line interface for the QuickNode Marketplace",
	Long: `
The qn-marketplace-cli is a Command Line Interface for the QuickNode Marketplace:

Among other things, it provides:
  - A set of commands to test an add-on's provisioning implementation
  - A command to test your an add-on's SSO implementation
  - A command to test your an add-on's RPC methods
	
For more information, visit https://www.quicknode.com/marketplace`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().Bool("verbose", false, "Verbose output")
}
