/*
Copyright © 2023 QuickNode, Inc.

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// provisionCmd represents the provision command
var provisionCmd = &cobra.Command{
	Use:   "provision",
	Short: "Allows you to test your add-on's PROVISION implementation",
	Long: `Use this command to make sure your API implementation for PROVISION is working as expected.

Learn more at https://www.quicknode.com/guides/quicknode-products/marketplace/how-provisioning-works-for-marketplace-partners/`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("provision called")
	},
}

func init() {
	rootCmd.AddCommand(provisionCmd)

	provisionCmd.PersistentFlags().String("url", "", "The URL of the add-on's provision endpoint")
  provisionCmd.PersistentFlags().String("basic-auth", "", "The basic auth credentials for the add-on")
	provisionCmd.PersistentFlags().String("chain", "", "The chain to provision the add-on for")
	provisionCmd.PersistentFlags().String("network", "", "The network to provision the add-on for")
	provisionCmd.PersistentFlags().String("plan", "", "The plan to provision the add-on for")
	provisionCmd.PersistentFlags().String("quicknode-id", "", "The QuickNode ID to provision the add-on for (optional)")
  provisionCmd.PersistentFlags().String("endpoint-id", "", "The endpoint ID to provision the add-on for (optional)")
}
