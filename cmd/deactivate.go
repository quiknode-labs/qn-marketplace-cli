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

// deactivateCmd represents the deactivate command
var deactivateCmd = &cobra.Command{
	Use:   "deactivate",
	Short: "Allows you to test your add-on's DEACTIVATE_ENDPOINT implementation",
	Long: `Use this command to make sure your API implementation for DEACTIVATE_ENDPOINT is working as expected.

Learn more at https://www.quicknode.com/guides/quicknode-products/marketplace/how-provisioning-works-for-marketplace-partners/`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("*** DEACTIVATE ***\n\n")
		url := cmd.Flag("url").Value.String()
		request := marketplace.DeactivateRequest{
			QuickNodeId:  cmd.Flag("quicknode-id").Value.String(),
			EndpointId:   cmd.Flag("endpoint-id").Value.String(),
			Chain:        cmd.Flag("chain").Value.String(),
			Network:      cmd.Flag("network").Value.String(),
			DeactivateAt: time.Now().Format(time.RFC3339),
		}

		fmt.Printf("DELETE %s:\n", url)
		requestJson, _ := json.MarshalIndent(request, "", "  ")
		fmt.Printf("%s\n", requestJson)

		response, err := marketplace.Deactivate(url, request)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("\nDeactivate Endpoint was successful:\n")
		fmt.Printf("\tStatus: \t\t%s\n", response.Status)
	},
}

func init() {
	rootCmd.AddCommand(deactivateCmd)

	deactivateCmd.PersistentFlags().String("url", "", "The URL of the add-on's update endpoint")
	deactivateCmd.PersistentFlags().String("basic-auth", "", "The basic auth credentials for the add-on")
	deactivateCmd.PersistentFlags().String("quicknode-id", "", "The Quicknode ID for the endpoint's account")
	deactivateCmd.PersistentFlags().String("endpoint-id", "", "The endpoint ID for the endpoint you want to deactivate")
	deactivateCmd.PersistentFlags().String("chain", "", "The chain to provision the add-on for")
	deactivateCmd.PersistentFlags().String("network", "", "The network to provision the add-on for")
}
