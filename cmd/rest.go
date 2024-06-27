/*
Copyright © 2023 QuikNode Inc.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/quiknode-labs/qn-marketplace-cli/marketplace"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
)

// restCmd represents the rest command
var restCmd = &cobra.Command{
	Use:   "rest",
	Short: "Allows you to test your add-on's REST paths",
	Args:  cobra.OnlyValidArgs,
	Run: func(cmd *cobra.Command, args []string) {
		header := color.New(color.FgWhite, color.BgBlue).SprintFunc()
		fmt.Printf("%s\n\n", header("        REST        "))
		verbose := cmd.Flag("verbose").Value.String() == "true"
		provisionURL := cmd.Flag("url").Value.String()
		if provisionURL == "" {
			fmt.Print("Please provide a URL for the provision API via the --url flag\n")
			os.Exit(1)
		}

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

		// First Provision
		request := marketplace.ProvisionRequest{
			QuickNodeId:       cmd.Flag("quicknode-id").Value.String(),
			EndpointId:        cmd.Flag("endpoint-id").Value.String(),
			Chain:             cmd.Flag("chain").Value.String(),
			Network:           cmd.Flag("network").Value.String(),
			Plan:              cmd.Flag("plan").Value.String(),
			WSSURL:            cmd.Flag("wss-url").Value.String(),
			HTTPURL:           cmd.Flag("endpoint-url").Value.String(),
			Referers:          []string{"https://quicknode.com"},
			ContractAddresses: []string{"0x4d224452801ACEd8B2F0aebE155379bb5D594381"},
		}

		if verbose {
			color.Blue("→ POST %s:\n", provisionURL)
		}
		requestJson, _ := json.MarshalIndent(request, "", "  ")
		if verbose {
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

		// Now we can make the REST call
		var requestBody = cmd.Flag("rest-body").Value.String()

		// Create an HTTP request with the REST request
		if verbose {
			color.Blue("\n→ %s %s:\n", restVerb, cmd.Flag("rest-url").Value.String())
			fmt.Printf("%s\n", requestBody)
		}

		httpReq, err := http.NewRequest(restVerb, cmd.Flag("rest-url").Value.String(), strings.NewReader(requestBody))
		if err != nil {
			color.Red("Error creating HTTP request: %s", err)
			os.Exit(1)
		}

		// Set the HTTP request header to indicate that the request body is in JSON format
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("X-QUICKNODE-ID", cmd.Flag("quicknode-id").Value.String())
		httpReq.Header.Set("X-INSTANCE-ID", cmd.Flag("endpoint-id").Value.String())
		httpReq.Header.Set("X-QN-CHAIN", cmd.Flag("chain").Value.String())
		httpReq.Header.Set("X-QN-NETWORK", cmd.Flag("network").Value.String())
		httpReq.Header.Add("X-QN-TESTING", "true")

		// Send the HTTP request and capture the response
		client := http.Client{}
		resp, err := client.Do(httpReq)
		if err != nil {
			color.Red("Error sending HTTP request: %s", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		// Decode the response body into an interface{} object
		var respBody interface{}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		if err != nil {
			color.Red("Error decoding JSON:", err)
			os.Exit(1)
		}

		responseJson, _ := json.MarshalIndent(respBody, "", "  ")
		if resp.StatusCode == 200 {
			color.Green("  ✓ REST call was successful and returned:")
			color.White("\n%s\n", responseJson)
		} else {
			color.Red("  ✘ REST call failed:     %s\n\n", resp.Status)
			color.White("\n%s\n", responseJson)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(restCmd)

	restCmd.PersistentFlags().StringP("url", "u", "", "The URL of the add-on's provision endpoint")

	// Note: basic auth defaults to username = Aladdin and password = open sesame
	restCmd.PersistentFlags().String("basic-auth", "QWxhZGRpbjpvcGVuIHNlc2FtZQ==", "The basic auth credentials for the add-on. Defaults to username = Aladdin and password = open sesame")

	restCmd.PersistentFlags().StringP("quicknode-id", "q", uuid.NewV4().String(), "The QuickNode ID to provision the add-on for (optional)")
	restCmd.PersistentFlags().StringP("endpoint-id", "e", uuid.NewV4().String(), "The endpoint ID to provision the add-on for (optional)")
	restCmd.PersistentFlags().StringP("endpoint-url", "l", "https://long-late-firefly.quiknode.pro/4bb1e6b2dec8294938b6fdfdb7cf0cf70c4e97a2/", "The endpoint URL to provision the add-on for (optional - defaults to an ethereum mainnet endpoint")
	restCmd.PersistentFlags().StringP("wss-url", "w", "wss://long-late-firefly.quiknode.pro/4bb1e6b2dec8294938b6fdfdb7cf0cf70c4e97a2/", "The WSS URL to provision the add-on for (optional - defaults to an ethereum mainnet endpoint")
	restCmd.PersistentFlags().StringP("chain", "c", "ethereum", "The chain to provision the add-on for")
	restCmd.PersistentFlags().StringP("network", "n", "mainnet", "The network to provision the add-on for")
	restCmd.PersistentFlags().StringP("plan", "p", "discover", "The plan to provision the add-on for")

	restCmd.PersistentFlags().String("rest-url", "", "The URL to make the REST calls to")
	restCmd.PersistentFlags().String("rest-verb", "", "The REST HTTP Method or verb to use (e.g. GET or POST)")
	restCmd.PersistentFlags().String("rest-body", "", "The Rest Request Body")
}
