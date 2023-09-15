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

// deprovisionCmd represents the deprovision command
var deprovisionCmd = &cobra.Command{
	Use:   "deprovision",
	Short: "Allows you to test your add-on's DEPROVISION implementation",
	Long: `Use this command to make sure your API implementation for DEPROVISION is working as expected.

Learn more at https://www.quicknode.com/guides/quicknode-products/marketplace/how-provisioning-works-for-marketplace-partners/`,
	Args: cobra.OnlyValidArgs,
	Run: func(cmd *cobra.Command, args []string) {
		header := color.New(color.FgWhite, color.BgBlue).SprintFunc()
		fmt.Printf("%s\n\n", header("        DEPROVISION        "))
		verbose := cmd.Flag("verbose").Value.String() == "true"
		url := cmd.Flag("url").Value.String()
		if url == "" {
			fmt.Print("Please provide a URL for the deprovision API via the --url flag\n")
			os.Exit(1)
		}

		request := marketplace.DeprovisionRequest{
			QuickNodeId:   cmd.Flag("quicknode-id").Value.String(),
			EndpointId:    cmd.Flag("endpoint-id").Value.String(),
			Chain:         cmd.Flag("chain").Value.String(),
			Network:       cmd.Flag("network").Value.String(),
			DeprovisionAt: time.Now().Unix(),
		}

		// Check that it is protected by basic auth
		isProtectedByBasicAuth, err := marketplace.RequiresBasicAuth(url, "DELETE")
		if err != nil {
			color.Red("%s", err)
			os.Exit(1)
		}
		if !isProtectedByBasicAuth {
			color.Red("  ✘ The deprovision API is not protected by basic auth.")
			os.Exit(1)
		}

		if verbose {
			color.Blue("→ DELETE %s:\n", url)
		}
		requestJson, _ := json.MarshalIndent(request, "", "  ")
		if verbose {
			fmt.Printf("%s\n", requestJson)
		}

		response, err := marketplace.Deprovision(url, request, cmd.Flag("basic-auth").Value.String())
		if err != nil {
			color.Red("%s", err)
			os.Exit(1)
		}

		if verbose {
			fmt.Printf("\nDeprovision was successful:\n")
			fmt.Printf("\tStatus: \t\t%s\n\n", response.Status)
		}

		color.Green("  ✓ Deprovision was successful")
	},
}

func init() {
	rootCmd.AddCommand(deprovisionCmd)

	deprovisionCmd.PersistentFlags().StringP("url", "u", "", "The URL of the add-on's provision endpoint")

	// Note: basic auth defaults to username = Aladdin and password = open sesame
	deprovisionCmd.PersistentFlags().String("basic-auth", "QWxhZGRpbjpvcGVuIHNlc2FtZQ==", "The basic auth credentials for the add-on. Defaults to username = Aladdin and password = open sesame")

	deprovisionCmd.PersistentFlags().StringP("quicknode-id", "q", uuid.NewV4().String(), "The QuickNode ID to provision the add-on for (optional)")
	deprovisionCmd.PersistentFlags().StringP("endpoint-id", "e", uuid.NewV4().String(), "The endpoint ID to provision the add-on for (optional)")
	deprovisionCmd.PersistentFlags().StringP("chain", "c", "ethereum", "The chain to provision the add-on for")
	deprovisionCmd.PersistentFlags().StringP("network", "n", "mainnet", "The network to provision the add-on for")
}
