/*
Copyright © 2023 QuickNode, Inc.

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// puddCmd represents the pudd command
var puddCmd = &cobra.Command{
	Use:   "pudd",
	Short: "Allows you to test your add-on's entire provisioning workflows (all four actions)",
	Long: `Use this command to make sure your API implementation for provisioning workflows works across the board.

This only works if your API URLs ends with:
  - /provision
  - /update
  - /deactivate_endpoint
  - /deprovision

The tool will use the base-url you pass to it and append these to the base URL to call your API.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pudd called")
	},
}

func init() {
	rootCmd.AddCommand(puddCmd)

	puddCmd.PersistentFlags().String("base-url", "", "The base URL of the add-on's provisioning API")
  puddCmd.PersistentFlags().String("basic-auth", "", "The basic auth credentials for the add-on")
	puddCmd.PersistentFlags().String("chain", "", "The chain to provision the add-on for")
	puddCmd.PersistentFlags().String("network", "", "The network to provision the add-on for")
	puddCmd.PersistentFlags().String("plan", "", "The plan to provision the add-on for")
}
