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

// puddCmd represents the pudd command
var puddCmd = &cobra.Command{
	Use:   "pudd",
	Short: "Allows you to test your add-on's entire provisioning workflows (all four actions)",
	Long: `Use this command to make sure your API implementation for provisioning workflows works across the board.

This only works if your API URLs ends with:
  - /provision
  - /update
  - /deactivate_endpoint
  - /deprovision

The tool will use the base-url you pass to it and append these to the base URL to call your API.
`,
	Args: cobra.OnlyValidArgs,
	Run: func(cmd *cobra.Command, args []string) {
		header := color.New(color.FgWhite, color.BgGreen).SprintFunc()
		fmt.Printf("%s\n\n", header("        PUDD        "))
		baseUrl := cmd.Flag("base-url").Value.String()
		if baseUrl == "" {
			fmt.Print("Please provide a base URL for the provisioning API via the --base-url flag\n")
			os.Exit(1)
		}

		// First Provision
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

		provisionUrl := baseUrl + "/provision"
		color.Magenta("→ POST %s:\n", provisionUrl)
		requestJson, _ := json.MarshalIndent(request, "", "  ")
		fmt.Printf("%s\n", requestJson)

		provisionResponse, err := marketplace.Provision(provisionUrl, request, cmd.Flag("basic-auth").Value.String())
		if err != nil {
			color.Red("%s", err)
			os.Exit(1)
		}
		fmt.Printf("\nProvision was successful:\n")
		fmt.Printf("  Status:     %s\n", provisionResponse.Status)
		fmt.Printf("  Dashboard URL:     %s\n", provisionResponse.DashboardURL)
		fmt.Printf("  Access URL:     %s\n", provisionResponse.AccessURL)

		// Then Provision again to test for idempotent provisions
		color.Magenta("\n\n→ POST %s (again to test idempotent provisions):\n", provisionUrl)
		fmt.Printf("%s\n", requestJson)
		provisionResponseTwo, err := marketplace.Provision(provisionUrl, request, cmd.Flag("basic-auth").Value.String())
		if err != nil {
			color.Red("%s", err)
			os.Exit(1)
		}
		fmt.Printf("\nSecond Provision was successful:\n")
		fmt.Printf("  Status:     %s\n", provisionResponseTwo.Status)
		fmt.Printf("  Dashboard URL:     %s\n", provisionResponseTwo.DashboardURL)
		fmt.Printf("  Access URL:     %s\n", provisionResponseTwo.AccessURL)

		// Now, let's Update
		updateUrl := baseUrl + "/update"
		color.Magenta("\n\n→ PUT %s:\n", updateUrl)

		updateRequest := marketplace.UpdateRequest{
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

		updateRequestJson, _ := json.MarshalIndent(updateRequest, "", "  ")
		fmt.Printf("%s\n", updateRequestJson)

		updateResponse, err := marketplace.Update(updateUrl, updateRequest, cmd.Flag("basic-auth").Value.String())
		if err != nil {
			color.Red("%s", err)
			os.Exit(1)
		}
		fmt.Printf("\nUpdate was successful:\n")
		fmt.Printf("  Status:     %s\n", updateResponse.Status)

		// Let's deactivate the endpoint
		deactivateUrl := baseUrl + "/deactivate_endpoint"
		color.Magenta("\n\n→ DELETE %s:\n", deactivateUrl)
		deactivateRequest := marketplace.DeactivateRequest{
			QuickNodeId:  cmd.Flag("quicknode-id").Value.String(),
			EndpointId:   cmd.Flag("endpoint-id").Value.String(),
			Chain:        cmd.Flag("chain").Value.String(),
			Network:      cmd.Flag("network").Value.String(),
			DeactivateAt: time.Now().Format(time.RFC3339),
		}

		deactivateRequestJson, _ := json.MarshalIndent(deactivateRequest, "", "  ")
		fmt.Printf("%s\n", deactivateRequestJson)

		deactivateResponse, err := marketplace.Deactivate(deactivateUrl, deactivateRequest, cmd.Flag("basic-auth").Value.String())
		if err != nil {
			color.Red("%s", err)
			os.Exit(1)
		}
		fmt.Printf("\nDeactivate Endpoint was successful:\n")
		fmt.Printf("  Status:     %s\n", deactivateResponse.Status)

		// Finally, deprovision
		deprovisionUrl := baseUrl + "/deprovision"
		deprovisionRequest := marketplace.DeprovisionRequest{
			QuickNodeId:   cmd.Flag("quicknode-id").Value.String(),
			EndpointId:    cmd.Flag("endpoint-id").Value.String(),
			Chain:         cmd.Flag("chain").Value.String(),
			Network:       cmd.Flag("network").Value.String(),
			DeprovisionAt: time.Now().Format(time.RFC3339),
		}

		color.Magenta("\n\n→ DELETE %s:\n", deprovisionUrl)
		deprovisionRequestJson, _ := json.MarshalIndent(deprovisionRequest, "", "  ")
		fmt.Printf("%s\n", deprovisionRequestJson)

		deprovisionResponse, err := marketplace.Deprovision(deprovisionUrl, deprovisionRequest, cmd.Flag("basic-auth").Value.String())
		if err != nil {
			color.Red("%s", err)
			os.Exit(1)
		}
		fmt.Printf("\nDeprovision was successful:\n")
		fmt.Printf("\tStatus: \t\t%s\n", deprovisionResponse.Status)
	},
}

func init() {
	rootCmd.AddCommand(puddCmd)

	puddCmd.PersistentFlags().StringP("base-url", "u", "", "The base URL of the add-on's provisioning API")

	// Note: basic auth defaults to username = Aladdin and password = open sesame
	puddCmd.PersistentFlags().String("basic-auth", "QWxhZGRpbjpvcGVuIHNlc2FtZQ==", "The basic auth credentials for the add-on. Defaults to username = Aladdin and password = open sesame")

	puddCmd.PersistentFlags().StringP("quicknode-id", "q", uuid.NewV4().String(), "The QuickNode ID to provision the add-on for (optional)")
	puddCmd.PersistentFlags().StringP("endpoint-id", "e", uuid.NewV4().String(), "The endpoint ID to provision the add-on for (optional)")
	puddCmd.PersistentFlags().StringP("chain", "c", "ethereum", "The chain to provision the add-on for")
	puddCmd.PersistentFlags().StringP("network", "n", "mainnet", "The network to provision the add-on for")
	puddCmd.PersistentFlags().StringP("plan", "p", "discover", "The plan to provision the add-on for")
}
