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

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Allows you to test your add-on's UPDATE implementation",
	Long: `Use this command to make sure your API implementation for UPDATE is working as expected.

Learn more at https://www.quicknode.com/guides/quicknode-products/marketplace/how-provisioning-works-for-marketplace-partners/`,
	Args: cobra.OnlyValidArgs,
	Run: func(cmd *cobra.Command, args []string) {
		header := color.New(color.FgWhite, color.BgBlue).SprintFunc()
		fmt.Printf("%s\n\n", header("        UPDATE        "))
		verbose := cmd.Flag("verbose").Value.String() == "true"
		url := cmd.Flag("url").Value.String()
		if url == "" {
			fmt.Print("Please provide a URL for the update API via the --url flag\n")
			os.Exit(1)
		}

		request := marketplace.UpdateRequest{
			QuickNodeId:       cmd.Flag("quicknode-id").Value.String(),
			EndpointId:        cmd.Flag("endpoint-id").Value.String(),
			Chain:             cmd.Flag("chain").Value.String(),
			Network:           cmd.Flag("network").Value.String(),
			Plan:              cmd.Flag("plan").Value.String(),
			WSSURL:            "wss://long-late-firefly.quiknode.pro/4bb1e6b2dec8294938b6fdfdb7cf0cf70c4e97a2/",
			HTTPURL:           cmd.Flag("endpoint-url").Value.String(),
			Referers:          []string{"https://quicknode.com"},
			ContractAddresses: []string{"0x4d224452801ACEd8B2F0aebE155379bb5D594381"},
		}

		// Check that it is protected by basic auth
		isProtectedByBasicAuth, err := marketplace.RequiresBasicAuth(url, "PUT")
		if err != nil {
			color.Red("%s", err)
			os.Exit(1)
		}
		if !isProtectedByBasicAuth {
			color.Red("  ✘ The update API is not protected by basic auth.")
			os.Exit(1)
		}

		if verbose {
			color.Blue("→ PUT %s:\n", url)
		}
		requestJson, _ := json.MarshalIndent(request, "", "  ")
		if verbose {
			fmt.Printf("%s\n", requestJson)
		}

		response, err := marketplace.Update(url, request, cmd.Flag("basic-auth").Value.String())
		if err != nil {
			color.Red("%s", err)
			os.Exit(1)
		}

		if verbose {
			fmt.Printf("\nUpdate was successful:\n")
			fmt.Printf("  Status:     %s\n\n", response.Status)
		}

		color.Green("  ✓ Update was successful")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.PersistentFlags().StringP("url", "u", "", "The URL of the add-on's provision endpoint")

	// Note: basic auth defaults to username = Aladdin and password = open sesame
	updateCmd.PersistentFlags().String("basic-auth", "QWxhZGRpbjpvcGVuIHNlc2FtZQ==", "The basic auth credentials for the add-on. Defaults to username = Aladdin and password = open sesame")

	updateCmd.PersistentFlags().StringP("quicknode-id", "q", uuid.NewV4().String(), "The QuickNode ID to provision the add-on for (optional)")
	updateCmd.PersistentFlags().StringP("endpoint-id", "e", uuid.NewV4().String(), "The endpoint ID to provision the add-on for (optional)")
	updateCmd.PersistentFlags().StringP("endpoint-url", "l", "https://long-late-firefly.quiknode.pro/4bb1e6b2dec8294938b6fdfdb7cf0cf70c4e97a2/", "The endpoint URL to provision the add-on for (optional - defaults to an ethereum mainnet endpoint")
	updateCmd.PersistentFlags().StringP("chain", "c", "ethereum", "The chain to provision the add-on for")
	updateCmd.PersistentFlags().StringP("network", "n", "mainnet", "The network to provision the add-on for")
	updateCmd.PersistentFlags().StringP("plan", "p", "discover", "The plan to provision the add-on for")
}
