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
		header := color.New(color.FgWhite, color.BgBlue).SprintFunc()
		fmt.Printf("%s\n\n", header("        PUDD        "))
		verbose := cmd.Flag("verbose").Value.String() == "true"
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

		// Check that it is protected by basic auth
		isProtectedByBasicAuth, err := marketplace.RequiresBasicAuth(provisionUrl, "POST")
		if err != nil {
			color.Red("%s", err)
			os.Exit(1)
		}
		if !isProtectedByBasicAuth {
			color.Red("  ✘ The provision API is not protected by basic auth.")
			os.Exit(1)
		} else {
			color.Green("  ✓ Provision API is protected by basic auth.")
		}

		if verbose {
			color.Blue("→ POST %s:\n", provisionUrl)
		}
		requestJson, _ := json.MarshalIndent(request, "", "  ")
		if verbose {
			fmt.Printf("%s\n", requestJson)
		}

		provisionResponse, err := marketplace.Provision(provisionUrl, request, cmd.Flag("basic-auth").Value.String())
		if err != nil {
			color.Red("%s", err)
			os.Exit(1)
		}

		if verbose {
			fmt.Printf("\nProvision was successful:\n")
			fmt.Printf("  Status:     %s\n", provisionResponse.Status)
			fmt.Printf("  Dashboard URL:     %s\n", provisionResponse.DashboardURL)
			fmt.Printf("  Access URL:     %s\n\n", provisionResponse.AccessURL)
		}

		color.Green("  ✓ Provision #1 was successful")

		// Then Provision again to test for idempotent provisions
		if verbose {
			color.Blue("\n\n→ POST %s (again to test idempotent provisions):\n", provisionUrl)
			fmt.Printf("%s\n", requestJson)
		}
		provisionResponseTwo, err := marketplace.Provision(provisionUrl, request, cmd.Flag("basic-auth").Value.String())
		if err != nil {
			color.Red("%s", err)
			os.Exit(1)
		}

		if verbose {
			fmt.Printf("\nSecond Provision was successful:\n")
			fmt.Printf("  Status:     %s\n", provisionResponseTwo.Status)
			fmt.Printf("  Dashboard URL:     %s\n", provisionResponseTwo.DashboardURL)
			fmt.Printf("  Access URL:     %s\n\n", provisionResponseTwo.AccessURL)
		}

		color.Green("  ✓ Provision #2 was successful")

		// Now, let's Update
		updateUrl := baseUrl + "/update"
		if verbose {
			color.Blue("\n\n→ PUT %s:\n", updateUrl)
		}

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

		// Check that it is protected by basic auth
		updateIsProtectedByBasicAuth, err := marketplace.RequiresBasicAuth(updateUrl, "PUT")
		if err != nil {
			color.Red("%s", err)
			os.Exit(1)
		}
		if !updateIsProtectedByBasicAuth {
			color.Red("  ✘ The update API is not protected by basic auth.")
			os.Exit(1)
		} else {
			color.Green("  ✓ Update API is protected by basic auth.")
		}

		updateRequestJson, _ := json.MarshalIndent(updateRequest, "", "  ")
		if verbose {
			fmt.Printf("%s\n", updateRequestJson)
		}

		updateResponse, err := marketplace.Update(updateUrl, updateRequest, cmd.Flag("basic-auth").Value.String())
		if err != nil {
			color.Red("%s", err)
			os.Exit(1)
		}
		if verbose {
			fmt.Printf("\nUpdate was successful:\n")
			fmt.Printf("  Status:     %s\n\n", updateResponse.Status)
		}

		color.Green("  ✓ Update was successful")

		// Let's deactivate the endpoint
		deactivateUrl := baseUrl + "/deactivate_endpoint"
		if verbose {
			color.Blue("\n\n→ DELETE %s:\n", deactivateUrl)
		}
		deactivateRequest := marketplace.DeactivateRequest{
			QuickNodeId:  cmd.Flag("quicknode-id").Value.String(),
			EndpointId:   cmd.Flag("endpoint-id").Value.String(),
			Chain:        cmd.Flag("chain").Value.String(),
			Network:      cmd.Flag("network").Value.String(),
			DeactivateAt: time.Now().Format(time.RFC3339),
		}

		// Check that it is protected by basic auth
		deactivateIsProtectedByBasicAuth, err := marketplace.RequiresBasicAuth(deactivateUrl, "DELETE")
		if err != nil {
			color.Red("%s", err)
			os.Exit(1)
		}
		if !deactivateIsProtectedByBasicAuth {
			color.Red("  ✘ The deactivate_endpoint API is not protected by basic auth.")
			os.Exit(1)
		} else {
			color.Green("  ✓ Deactivate Endpoint API is protected by basic auth.")
		}

		deactivateRequestJson, _ := json.MarshalIndent(deactivateRequest, "", "  ")
		if verbose {
			fmt.Printf("%s\n", deactivateRequestJson)
		}

		deactivateResponse, err := marketplace.Deactivate(deactivateUrl, deactivateRequest, cmd.Flag("basic-auth").Value.String())
		if err != nil {
			color.Red("%s", err)
			os.Exit(1)
		}
		if verbose {
			fmt.Printf("\nDeactivate Endpoint was successful:\n")
			fmt.Printf("  Status:     %s\n\n", deactivateResponse.Status)
		}

		color.Green("  ✓ Deactivate Endpoint was successful")

		// Finally, deprovision
		deprovisionUrl := baseUrl + "/deprovision"
		deprovisionRequest := marketplace.DeprovisionRequest{
			QuickNodeId:   cmd.Flag("quicknode-id").Value.String(),
			EndpointId:    cmd.Flag("endpoint-id").Value.String(),
			Chain:         cmd.Flag("chain").Value.String(),
			Network:       cmd.Flag("network").Value.String(),
			DeprovisionAt: time.Now().Format(time.RFC3339),
		}

		// Check that it is protected by basic auth
		deprovisionIsProtectedByBasicAuth, err := marketplace.RequiresBasicAuth(deprovisionUrl, "DELETE")
		if err != nil {
			color.Red("%s", err)
			os.Exit(1)
		}
		if !deprovisionIsProtectedByBasicAuth {
			color.Red("  ✘ The deprovision API is not protected by basic auth.")
			os.Exit(1)
		} else {
			color.Green("  ✓ Deprovision API is protected by basic auth.")
		}

		if verbose {
			color.Blue("\n\n→ DELETE %s:\n", deprovisionUrl)
		}
		deprovisionRequestJson, _ := json.MarshalIndent(deprovisionRequest, "", "  ")
		if verbose {
			fmt.Printf("%s\n", deprovisionRequestJson)
		}

		deprovisionResponse, err := marketplace.Deprovision(deprovisionUrl, deprovisionRequest, cmd.Flag("basic-auth").Value.String())
		if err != nil {
			color.Red("%s", err)
			os.Exit(1)
		}
		if verbose {
			fmt.Printf("\nDeprovision was successful:\n")
			fmt.Printf("\tStatus: \t\t%s\n\n", deprovisionResponse.Status)
		}

		color.Green("  ✓ Deprovision was successful")
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
