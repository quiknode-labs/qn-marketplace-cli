/*
Copyright © 2023 QuickNode, Inc.

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Allows you to test your add-on's UPDATE implementation",
	Long: `Use this command to make sure your API implementation for UPDATE is working as expected.

Learn more at https://www.quicknode.com/guides/quicknode-products/marketplace/how-provisioning-works-for-marketplace-partners/`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("update called")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.PersistentFlags().String("url", "", "The URL of the add-on's update endpoint")
  updateCmd.PersistentFlags().String("basic-auth", "", "The basic auth credentials for the add-on")
	updateCmd.PersistentFlags().String("chain", "", "The chain for the installation to update")
	updateCmd.PersistentFlags().String("network", "", "The network for the installation to update")
	updateCmd.PersistentFlags().String("plan", "", "The plan for the installation to update")
	updateCmd.PersistentFlags().String("quicknode-id", "", "The QuickNode ID for the account you want to update (optional)")
  updateCmd.PersistentFlags().String("endpoint-id", "", "The endpoint ID for the endpoint you want to update (optional)")
}
