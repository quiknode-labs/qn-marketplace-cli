/*
Copyright © 2023 QuickNode, Inc.

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/quiknode-labs/qn-marketplace-cli/marketplace"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
)

// deactivateCmd represents the deactivate command
var deactivateCmd = &cobra.Command{
	Use:   "deactivate",
	Short: "Allows you to test your add-on's DEACTIVATE_ENDPOINT implementation",
	Long: `Use this command to make sure your API implementation for DEACTIVATE_ENDPOINT is working as expected.

Learn more at https://www.quicknode.com/guides/quicknode-products/marketplace/how-provisioning-works-for-marketplace-partners/`,
	Args: cobra.OnlyValidArgs,
	Run: func(cmd *cobra.Command, args []string) {
		header := color.New(color.FgWhite, color.BgBlue).SprintFunc()
		fmt.Printf("%s\n\n", header("        DEACTIVATE        "))
		verbose := cmd.Flag("verbose").Value.String() == "true"
		url := cmd.Flag("url").Value.String()
		if url == "" {
			fmt.Print("Please provide a URL for the deactivate API via the --url flag\n")
			os.Exit(1)
		}
		request := marketplace.DeactivateRequest{
			QuickNodeId:  cmd.Flag("quicknode-id").Value.String(),
			EndpointId:   cmd.Flag("endpoint-id").Value.String(),
			Chain:        cmd.Flag("chain").Value.String(),
			Network:      cmd.Flag("network").Value.String(),
			DeactivateAt: time.Now().Format(time.RFC3339),
		}

		if verbose {
			color.Blue("→ DELETE %s:\n", url)
		}
		requestJson, _ := json.MarshalIndent(request, "", "  ")
		if verbose {
			fmt.Printf("%s\n", requestJson)
		}

		response, err := marketplace.Deactivate(url, request, cmd.Flag("basic-auth").Value.String())
		if err != nil {
			color.Red("%s", err)
			os.Exit(1)
		}

		if verbose {
			fmt.Printf("\nDeactivate Endpoint was successful:\n")
			fmt.Printf("  Status:     %s\n\n", response.Status)
		}

		color.Green("  ✓ Deactivate Endpoint was successful")
	},
}

func init() {
	rootCmd.AddCommand(deactivateCmd)

	deactivateCmd.PersistentFlags().StringP("url", "u", "", "The URL of the add-on's provision endpoint")

	// Note: basic auth defaults to username = Aladdin and password = open sesame
	deactivateCmd.PersistentFlags().String("basic-auth", "QWxhZGRpbjpvcGVuIHNlc2FtZQ==", "The basic auth credentials for the add-on. Defaults to username = Aladdin and password = open sesame")

	deactivateCmd.PersistentFlags().StringP("quicknode-id", "q", uuid.NewV4().String(), "The QuickNode ID to provision the add-on for (optional)")
	deactivateCmd.PersistentFlags().StringP("endpoint-id", "e", uuid.NewV4().String(), "The endpoint ID to provision the add-on for (optional)")
	deactivateCmd.PersistentFlags().StringP("chain", "c", "ethereum", "The chain to provision the add-on for")
	deactivateCmd.PersistentFlags().StringP("network", "n", "mainnet", "The network to provision the add-on for")
}
