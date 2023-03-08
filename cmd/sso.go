/*
Copyright © 2023 QuickNode, Inc.

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ssoCmd represents the sso command
var ssoCmd = &cobra.Command{
	Use:   "sso",
	Short: "Allows you to test your add-on's SSO implementation",
	Long: `Use this command to make sure your add-on's SSO implementation is working as expected.
	
Learn more at https://www.quicknode.com/guides/quicknode-products/marketplace/how-sso-works-for-marketplace-partners/
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sso called")
	},
}

func init() {
	rootCmd.AddCommand(ssoCmd)

	ssoCmd.PersistentFlags().String("url", "", "The SSO URL for the add-on")
  ssoCmd.PersistentFlags().String("basic-auth", "", "The basic auth credentials for the add-on")
}
