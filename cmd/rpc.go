/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/quiknode-labs/qn-marketplace-cli/marketplace"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
)

// rpcCmd represents the rpc command
var rpcCmd = &cobra.Command{
	Use:   "rpc",
	Short: "Allows you to test your add-on's RPC methods",
	Args:  cobra.OnlyValidArgs,
	Run: func(cmd *cobra.Command, args []string) {
		header := color.New(color.FgWhite, color.BgBlue).SprintFunc()
		fmt.Printf("%s\n\n", header("        RPC        "))
		verbose := cmd.Flag("verbose").Value.String() == "true"
		provisionURL := cmd.Flag("url").Value.String()
		if provisionURL == "" {
			fmt.Print("Please provide a URL for the provision API via the --url flag\n")
			os.Exit(1)
		}

		rpcURL := cmd.Flag("rpc-url").Value.String()
		if rpcURL == "" {
			fmt.Print("Please provide a URL for the RPC API via the --rpc-url flag\n")
			os.Exit(1)
		}

		rpcMethod := cmd.Flag("rpc-method").Value.String()
		if rpcMethod == "" {
			color.Red("Please provide an RPC Method for the provision API via the --rpc-method flag\n")
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

		// Now we can make the RPC call
		// First, Create an RPC request object
		var params []interface{}
		var paramsFlag = cmd.Flag("rpc-params").Value.String()

		if paramsFlag != "" {
			err := json.Unmarshal([]byte(paramsFlag), &params)
			if err != nil {
				color.Red("Error parsing params: %s", err)
				os.Exit(1)
			}
		} else {
			params = []interface{}{}
		}

		req := marketplace.RPCRequest{
			Method: cmd.Flag("rpc-method").Value.String(),
			Params: params,
			ID:     uuid.NewV4().String(),
		}

		// Encode the request object into a JSON string
		reqBody, err := json.Marshal(req)
		if err != nil {
			color.Red("Error encoding JSON: %", err)
			os.Exit(1)
		}

		reqBodyIndented, err := json.MarshalIndent(req, "", "  ")
		if err != nil {
			color.Red("Error encoding JSON: %", err)
			os.Exit(1)
		}
		// Create an HTTP request with the JSON-RPC request body
		if verbose {
			color.Blue("\n→ POST %s:\n", cmd.Flag("rpc-url").Value.String())
			fmt.Printf("%s\n", reqBodyIndented)
		}

		httpReq, err := http.NewRequest("POST", cmd.Flag("rpc-url").Value.String(), bytes.NewBuffer(reqBody))
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

		color.Green("  ✓ RPC call was successful and returned:")
		responseJson, _ := json.MarshalIndent(respBody, "", "  ")
		color.White("\n%s\n", responseJson)
	},
}

func init() {
	rootCmd.AddCommand(rpcCmd)

	rpcCmd.PersistentFlags().StringP("url", "u", "", "The URL of the add-on's provision endpoint")

	// Note: basic auth defaults to username = Aladdin and password = open sesame
	rpcCmd.PersistentFlags().String("basic-auth", "QWxhZGRpbjpvcGVuIHNlc2FtZQ==", "The basic auth credentials for the add-on. Defaults to username = Aladdin and password = open sesame")

	rpcCmd.PersistentFlags().StringP("quicknode-id", "q", uuid.NewV4().String(), "The QuickNode ID to provision the add-on for (optional)")
	rpcCmd.PersistentFlags().StringP("endpoint-id", "e", uuid.NewV4().String(), "The endpoint ID to provision the add-on for (optional)")
	rpcCmd.PersistentFlags().StringP("chain", "c", "ethereum", "The chain to provision the add-on for")
	rpcCmd.PersistentFlags().StringP("network", "n", "mainnet", "The network to provision the add-on for")
	rpcCmd.PersistentFlags().StringP("plan", "p", "discover", "The plan to provision the add-on for")

	rpcCmd.PersistentFlags().String("rpc-url", "", "The URL to make the RPC calls to")
	rpcCmd.PersistentFlags().String("rpc-method", "", "The RPC Method to call")
	rpcCmd.PersistentFlags().String("rpc-params", "", "The RPC Params to call the RPC Method with in JSON format")
}
