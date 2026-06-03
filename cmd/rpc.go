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

var rpcCmd = &cobra.Command{
	Use:   "rpc",
	Short: "Allows you to test your add-on's RPC methods",
	Args:  cobra.OnlyValidArgs,
	Run: func(cmd *cobra.Command, args []string) {
		header := color.New(color.FgWhite, color.BgBlue).SprintFunc()
		fmt.Printf("%s\n\n", header("        RPC        "))
		verbose := cmd.Flag("verbose").Value.String() == "true"

		rpcURL := cmd.Flag("rpc-url").Value.String()
		if rpcURL == "" {
			fmt.Print("Please provide a URL for the RPC API via the --rpc-url flag\n")
			os.Exit(1)
		}

		rpcMethod := cmd.Flag("rpc-method").Value.String()
		if rpcMethod == "" {
			color.Red("Please provide an RPC Method via the --rpc-method flag\n")
			os.Exit(1)
		}

		provisionURL := cmd.Flag("url").Value.String()
		customHeaders, _ := cmd.Flags().GetStringArray("header")
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

		var params []interface{}
		if paramsFlag := cmd.Flag("rpc-params").Value.String(); paramsFlag != "" {
			if err := json.Unmarshal([]byte(paramsFlag), &params); err != nil {
				color.Red("Error parsing params: %s", err)
				os.Exit(1)
			}
		} else {
			params = []interface{}{}
		}

		req := marketplace.RPCRequest{
			Method: rpcMethod,
			Params: params,
			ID:     uuid.NewV4().String(),
		}

		reqBody, err := json.Marshal(req)
		if err != nil {
			color.Red("Error encoding JSON: %s", err)
			os.Exit(1)
		}

		if verbose {
			reqBodyIndented, _ := json.MarshalIndent(req, "", "  ")
			color.Blue("\n→ POST %s:\n", rpcURL)
			fmt.Printf("%s\n", reqBodyIndented)
		}

		httpReq, err := http.NewRequest("POST", rpcURL, bytes.NewBuffer(reqBody))
		if err != nil {
			color.Red("Error creating HTTP request: %s", err)
			os.Exit(1)
		}

		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("X-QUICKNODE-ID", quicknodeID)
		httpReq.Header.Set("X-INSTANCE-ID", endpointID)
		httpReq.Header.Set("X-QN-CHAIN", cmd.Flag("chain").Value.String())
		httpReq.Header.Set("X-QN-NETWORK", cmd.Flag("network").Value.String())
		httpReq.Header.Add("X-QN-TESTING", "true")

		for _, w := range applyAuthHeaders(httpReq, mode, customHeaders, cmd.Flag("basic-auth").Value.String()) {
			color.Yellow("Warning: %s\n", w)
		}

		client := http.Client{}
		resp, err := client.Do(httpReq)
		if err != nil {
			color.Red("Error sending HTTP request: %s", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		var respBody interface{}
		if err = json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
			color.Red("Error decoding JSON: %s", err)
			os.Exit(1)
		}

		responseJson, _ := json.MarshalIndent(respBody, "", "  ")
		if resp.StatusCode == 200 {
			color.Green("  ✓ RPC call was successful and returned:")
			color.White("\n%s\n", responseJson)
		} else {
			color.Red("  ✘ RPC call failed:     %s\n\n", resp.Status)
			color.White("\n%s\n", responseJson)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(rpcCmd)

	rpcCmd.PersistentFlags().StringP("url", "u", "", "The URL of the add-on's provision endpoint (optional, enables provisioning mode)")
	rpcCmd.PersistentFlags().String("basic-auth", "QWxhZGRpbjpvcGVuIHNlc2FtZQ==", "The basic auth credentials for the add-on")
	rpcCmd.PersistentFlags().StringP("quicknode-id", "q", uuid.NewV4().String(), "The QuickNode ID (optional)")
	rpcCmd.PersistentFlags().StringP("endpoint-id", "e", uuid.NewV4().String(), "The endpoint ID (optional)")
	rpcCmd.PersistentFlags().StringP("endpoint-url", "l", "https://long-late-firefly.quiknode.pro/4bb1e6b2dec8294938b6fdfdb7cf0cf70c4e97a2/", "The endpoint HTTP URL (optional)")
	rpcCmd.PersistentFlags().StringP("wss-url", "w", "wss://long-late-firefly.quiknode.pro/4bb1e6b2dec8294938b6fdfdb7cf0cf70c4e97a2/", "The endpoint WSS URL (optional)")
	rpcCmd.PersistentFlags().StringP("chain", "c", "ethereum", "The chain")
	rpcCmd.PersistentFlags().StringP("network", "n", "mainnet", "The network")
	rpcCmd.PersistentFlags().StringP("plan", "p", "discover", "The plan slug")
	rpcCmd.PersistentFlags().StringP("add-on-id", "i", "33", "The add-on ID")
	rpcCmd.PersistentFlags().StringP("add-on-slug", "s", "myslug", "The add-on slug")
	rpcCmd.PersistentFlags().String("rpc-url", "", "The URL to make the RPC calls to")
	rpcCmd.PersistentFlags().String("rpc-method", "", "The RPC Method to call")
	rpcCmd.PersistentFlags().String("rpc-params", "", "The RPC params in JSON format")
	rpcCmd.PersistentFlags().StringArray("header", []string{}, "HTTP header in \"Key: Value\" format (repeatable, enables header auth mode)")
}
