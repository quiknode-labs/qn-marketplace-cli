/*
Copyright © 2023 QuickNode, Inc.

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// deprovisionCmd represents the deprovision command
var deprovisionCmd = &cobra.Command{
	Use:   "deprovision",
	Short: "Allows you to test your add-on's DEPROVISION implementation",
	Long: `Use this command to make sure your API implementation for DEPROVISION is working as expected.

Learn more at https://www.quicknode.com/guides/quicknode-products/marketplace/how-provisioning-works-for-marketplace-partners/`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("deprovision called")
	},
}

func init() {
	rootCmd.AddCommand(deprovisionCmd)

	deprovisionCmd.PersistentFlags().String("url", "", "The URL of the add-on's update endpoint")
  deprovisionCmd.PersistentFlags().String("basic-auth", "", "The basic auth credentials for the add-on")
  deprovisionCmd.PersistentFlags().String("quicknode-id", "", "The QuickNode ID to deprovision the add-on for")
}
