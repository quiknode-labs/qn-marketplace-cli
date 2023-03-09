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

// provisionCmd represents the provision command
var provisionCmd = &cobra.Command{
	Use:   "provision",
	Short: "Allows you to test your add-on's PROVISION implementation",
	Long: `Use this command to make sure your API implementation for PROVISION is working as expected.

Learn more at https://www.quicknode.com/guides/quicknode-products/marketplace/how-provisioning-works-for-marketplace-partners/`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("*** PROVISION ***\n\n")
		url := cmd.Flag("url").Value.String()
		request := marketplace.ProvisionRequest{
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

		fmt.Printf("POST %s:\n", url)
		requestJson, _ := json.MarshalIndent(request, "", "  ")
		fmt.Printf("%s\n", requestJson)

		response, err := marketplace.Provision(url, request)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("\nProvision was successful:\n")
		fmt.Printf("\tStatus: \t\t%s\n", response.Status)
		fmt.Printf("\tDashboard URL: \t\t%s\n", response.DashboardURL)
		fmt.Printf("\tAccess URL: \t\t%s\n", response.AccessURL)
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
