/*
Copyright © 2023 QuickNode, Inc.

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

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
		fmt.Printf("*** DEPROVISION ***\n\n")
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
			DeprovisionAt: time.Now().Format(time.RFC3339),
		}

		fmt.Printf("DELETE %s:\n", url)
		requestJson, _ := json.MarshalIndent(request, "", "  ")
		fmt.Printf("%s\n", requestJson)

		response, err := marketplace.Deprovision(url, request, cmd.Flag("basic-auth").Value.String())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("\nDeprovision was successful:\n")
		fmt.Printf("\tStatus: \t\t%s\n", response.Status)
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
