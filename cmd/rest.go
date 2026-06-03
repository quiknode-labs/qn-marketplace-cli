/*
Copyright © 2023 QuikNode Inc.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/quiknode-labs/qn-marketplace-cli/marketplace"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
)

var restCmd = &cobra.Command{
	Use:   "rest",
	Short: "Allows you to test your add-on's REST paths",
	Args:  cobra.OnlyValidArgs,
	Run: func(cmd *cobra.Command, args []string) {
		header := color.New(color.FgWhite, color.BgBlue).SprintFunc()
		fmt.Printf("%s\n\n", header("        REST        "))
		verbose := cmd.Flag("verbose").Value.String() == "true"

		restURL := cmd.Flag("rest-url").Value.String()
		if restURL == "" {
			fmt.Print("Please provide a URL for the REST API via the --rest-url flag\n")
			os.Exit(1)
		}

		restVerb := cmd.Flag("rest-verb").Value.String()
		if restVerb == "" {
			color.Red("Please provide a REST HTTP Verb (e.g. GET or POST) via the --rest-verb flag\n")
			os.Exit(1)
		}

		provisionURL := cmd.Flag("url").Value.String()
		customHeaders, err := cmd.Flags().GetStringArray("header")
		if err != nil {
			color.Red("Error reading --header flag: %s", err)
			os.Exit(1)
		}
		mode := detectAuthMode(provisionURL, customHeaders, cmd.Flag("basic-auth").Changed)

		quicknodeID := cmd.Flag("quicknode-id").Value.String()
		endpointID := cmd.Flag("endpoint-id").Value.String()

		if mode == authModeProvisioning {
			request := marketplace.ProvisionRequest{
				QuickNodeId:       quicknodeID,
				EndpointId:        endpointID,
				Chain:             cmd.Flag("chain").Value.String(),
				Network:           cmd.Flag("network").Value.String(),
				Plan:              cmd.Flag("plan").Value.String(),
				WSSURL:            cmd.Flag("wss-url").Value.String(),
				HTTPURL:           cmd.Flag("endpoint-url").Value.String(),
				Referers:          []string{"https://quicknode.com"},
				ContractAddresses: []string{"0x4d224452801ACEd8B2F0aebE155379bb5D594381"},
				AddOnSlug:         cmd.Flag("add-on-slug").Value.String(),
				AddOnId:           cmd.Flag("add-on-id").Value.String(),
			}

			if verbose {
				color.Blue("→ POST %s:\n", provisionURL)
				requestJson, _ := json.MarshalIndent(request, "", "  ")
				fmt.Printf("%s\n", requestJson)
			}

			provisionResponse, err := marketplace.Provision(provisionURL, request, cmd.Flag("basic-auth").Value.String())
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
		}

		requestBody := cmd.Flag("rest-body").Value.String()

		if verbose {
			color.Blue("\n→ %s %s:\n", restVerb, restURL)
			fmt.Printf("%s\n", requestBody)
		}

		httpReq, err := http.NewRequest(restVerb, restURL, strings.NewReader(requestBody))
		if err != nil {
			color.Red("Error creating HTTP request: %s", err)
			os.Exit(1)
		}

		for _, w := range applyAuthHeaders(httpReq, mode, customHeaders, cmd.Flag("basic-auth").Value.String()) {
			color.Yellow("Warning: %s\n", w)
		}

		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("X-QUICKNODE-ID", quicknodeID)
		httpReq.Header.Set("X-INSTANCE-ID", endpointID)
		httpReq.Header.Set("X-QN-CHAIN", cmd.Flag("chain").Value.String())
		httpReq.Header.Set("X-QN-NETWORK", cmd.Flag("network").Value.String())
		httpReq.Header.Add("X-QN-TESTING", "true")

		client := http.Client{}
		resp, err := client.Do(httpReq)
		if err != nil {
			color.Red("Error sending HTTP request: %s", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			color.Red("Error reading response body: %s", err)
			os.Exit(1)
		}

		printBody := func(b []byte) {
			if len(b) == 0 {
				return
			}
			var respBody interface{}
			if jsonErr := json.Unmarshal(b, &respBody); jsonErr == nil {
				responseJson, _ := json.MarshalIndent(respBody, "", "  ")
				color.White("\n%s\n", responseJson)
			} else {
				color.White("\n%s\n", string(b))
			}
		}

		if resp.StatusCode == 200 {
			color.Green("  ✓ REST call was successful and returned:")
			printBody(body)
		} else {
			color.Red("  ✘ REST call failed:     %s\n\n", resp.Status)
			printBody(body)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(restCmd)

	restCmd.PersistentFlags().StringP("url", "u", "", "The URL of the add-on's provision endpoint (optional, enables provisioning mode)")
	restCmd.PersistentFlags().String("basic-auth", "QWxhZGRpbjpvcGVuIHNlc2FtZQ==", "The basic auth credentials for the add-on")
	restCmd.PersistentFlags().StringP("quicknode-id", "q", uuid.NewV4().String(), "The QuickNode ID (optional)")
	restCmd.PersistentFlags().StringP("endpoint-id", "e", uuid.NewV4().String(), "The endpoint ID (optional)")
	restCmd.PersistentFlags().StringP("endpoint-url", "l", "https://long-late-firefly.quiknode.pro/4bb1e6b2dec8294938b6fdfdb7cf0cf70c4e97a2/", "The endpoint HTTP URL (optional)")
	restCmd.PersistentFlags().StringP("wss-url", "w", "wss://long-late-firefly.quiknode.pro/4bb1e6b2dec8294938b6fdfdb7cf0cf70c4e97a2/", "The endpoint WSS URL (optional)")
	restCmd.PersistentFlags().StringP("chain", "c", "ethereum", "The chain")
	restCmd.PersistentFlags().StringP("network", "n", "mainnet", "The network")
	restCmd.PersistentFlags().StringP("plan", "p", "discover", "The plan slug")
	restCmd.PersistentFlags().StringP("add-on-id", "i", "33", "The add-on ID")
	restCmd.PersistentFlags().StringP("add-on-slug", "s", "myslug", "The add-on slug")
	restCmd.PersistentFlags().String("rest-url", "", "The URL to make the REST calls to")
	restCmd.PersistentFlags().String("rest-verb", "", "The REST HTTP Method or verb to use (e.g. GET or POST)")
	restCmd.PersistentFlags().String("rest-body", "", "The REST request body")
	restCmd.PersistentFlags().StringArray("header", []string{}, "HTTP header in \"Key: Value\" format (repeatable, enables header auth mode)")
}
