/*
Copyright © 2023 QuickNode, Inc.

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/quiknode-labs/qn-marketplace-cli/marketplace"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Allows you to test your add-on's UPDATE implementation",
	Long: `Use this command to make sure your API implementation for UPDATE is working as expected.

Learn more at https://www.quicknode.com/guides/quicknode-products/marketplace/how-provisioning-works-for-marketplace-partners/`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("*** UPDATE ***\n\n")
		url := cmd.Flag("url").Value.String()
		request := marketplace.UpdateRequest{
			QuickNodeId:       cmd.Flag("quicknode-id").Value.String(),
			EndpointId:        cmd.Flag("endpoint-id").Value.String(),
			Chain:             cmd.Flag("chain").Value.String(),
			Network:           cmd.Flag("network").Value.String(),
			Plan:              cmd.Flag("plan").Value.String(),
			WSSURL:            "wss://long-late-firefly.quiknode.pro/4bb1e6b2dec8294938b6fdfdb7cf0cf70c4e97a2/",
			HTTPURL:           "https://long-late-firefly.quiknode.pro/4bb1e6b2dec8294938b6fdfdb7cf0cf70c4e97a2/",
			Referers:          []string{"https://quicknode.com"},
			ContractAddresses: []string{"0x4d224452801ACEd8B2F0aebE155379bb5D594381"},
		}

		fmt.Printf("PUT %s:\n", url)
		requestJson, _ := json.MarshalIndent(request, "", "  ")
		fmt.Printf("%s\n", requestJson)

		response, err := marketplace.Update(url, request)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("\nUpdate was successful:\n")
		fmt.Printf("\tStatus: \t\t%s\n", response.Status)
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
