/*
Copyright © 2023 QuickNode, Inc.

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// deactivateCmd represents the deactivate command
var deactivateCmd = &cobra.Command{
	Use:   "deactivate",
	Short: "Allows you to test your add-on's DEACTIVATE_ENDPOINT implementation",
	Long: `Use this command to make sure your API implementation for DEACTIVATE_ENDPOINT is working as expected.

Learn more at https://www.quicknode.com/guides/quicknode-products/marketplace/how-provisioning-works-for-marketplace-partners/`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("deactivate called")
	},
}

func init() {
	rootCmd.AddCommand(deactivateCmd)

	deactivateCmd.PersistentFlags().String("url", "", "The URL of the add-on's update endpoint")
  deactivateCmd.PersistentFlags().String("basic-auth", "", "The basic auth credentials for the add-on")
  deactivateCmd.PersistentFlags().String("endpoint-id", "", "The endpoint ID for the endpoint you want to deactivate")
}
