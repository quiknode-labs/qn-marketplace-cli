/*
Copyright © 2023 QuickNode, Inc.

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/quiknode-labs/qn-marketplace-cli/marketplace"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
)

// provisionCmd represents the provision command
var provisionCmd = &cobra.Command{
	Use:   "provision",
	Short: "Allows you to test your add-on's PROVISION implementation",
	Long: `Use this command to make sure your API implementation for PROVISION is working as expected.

Learn more at https://www.quicknode.com/guides/quicknode-products/marketplace/how-provisioning-works-for-marketplace-partners/`,
	Args: cobra.OnlyValidArgs,
	Run: func(cmd *cobra.Command, args []string) {
		header := color.New(color.FgWhite, color.BgGreen).SprintFunc()
		fmt.Printf("%s\n\n", header("        PROVISION        "))
		url := cmd.Flag("url").Value.String()
		if url == "" {
			fmt.Print("Please provide a URL for the provision API via the --url flag\n")
			os.Exit(1)
		}
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

		color.Magenta("→ POST %s:\n", url)
		requestJson, _ := json.MarshalIndent(request, "", "  ")
		fmt.Printf("%s\n", requestJson)

		response, err := marketplace.Provision(url, request, cmd.Flag("basic-auth").Value.String())
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

	provisionCmd.PersistentFlags().StringP("url", "u", "", "The URL of the add-on's provision endpoint")

	// Note: basic auth defaults to username = Aladdin and password = open sesame
	provisionCmd.PersistentFlags().String("basic-auth", "QWxhZGRpbjpvcGVuIHNlc2FtZQ==", "The basic auth credentials for the add-on. Defaults to username = Aladdin and password = open sesame")

	provisionCmd.PersistentFlags().StringP("quicknode-id", "q", uuid.NewV4().String(), "The QuickNode ID to provision the add-on for (optional)")
	provisionCmd.PersistentFlags().StringP("endpoint-id", "e", uuid.NewV4().String(), "The endpoint ID to provision the add-on for (optional)")
	provisionCmd.PersistentFlags().StringP("chain", "c", "ethereum", "The chain to provision the add-on for")
	provisionCmd.PersistentFlags().StringP("network", "n", "mainnet", "The network to provision the add-on for")
	provisionCmd.PersistentFlags().StringP("plan", "p", "discover", "The plan to provision the add-on for")
}
