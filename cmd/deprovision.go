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
	"github.com/spf13/cobra"
)

// deprovisionCmd represents the deprovision command
var deprovisionCmd = &cobra.Command{
	Use:   "deprovision",
	Short: "Allows you to test your add-on's DEPROVISION implementation",
	Long: `Use this command to make sure your API implementation for DEPROVISION is working as expected.

Learn more at https://www.quicknode.com/guides/quicknode-products/marketplace/how-provisioning-works-for-marketplace-partners/`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("*** DEPROVISION ***\n\n")
		url := cmd.Flag("url").Value.String()
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

		response, err := marketplace.Deprovision(url, request)
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

	deprovisionCmd.PersistentFlags().String("url", "", "The URL of the add-on's update endpoint")
	deprovisionCmd.PersistentFlags().String("basic-auth", "", "The basic auth credentials for the add-on")
	deprovisionCmd.PersistentFlags().String("quicknode-id", "", "The QuickNode ID to deprovision the add-on for")
	deprovisionCmd.PersistentFlags().String("endpoint-id", "", "The endpoint ID for the endpoint you want to deactivate")
	deprovisionCmd.PersistentFlags().String("chain", "", "The chain to provision the add-on for")
	deprovisionCmd.PersistentFlags().String("network", "", "The network to provision the add-on for")
}
