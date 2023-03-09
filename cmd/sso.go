/*
Copyright © 2023 QuickNode, Inc.

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/fatih/color"
	"github.com/quiknode-labs/qn-marketplace-cli/marketplace"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
)

// ssoCmd represents the sso command
var ssoCmd = &cobra.Command{
	Use:   "sso",
	Short: "Allows you to test your add-on's SSO implementation",
	Long: `Use this command to make sure your add-on's SSO implementation is working as expected.
	
Learn more at https://www.quicknode.com/guides/quicknode-products/marketplace/how-sso-works-for-marketplace-partners/
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("*** SSO ***\n\n")

		provisionURL := cmd.Flag("url").Value.String()
		if provisionURL == "" {
			fmt.Print("Please provide a URL for the provision API via the --url flag\n")
			os.Exit(1)
		}

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

		color.Magenta("→ POST %s:\n", provisionURL)
		requestJson, _ := json.MarshalIndent(request, "", "  ")
		fmt.Printf("%s\n", requestJson)

		provisionResponse, err := marketplace.Provision(provisionURL, request, cmd.Flag("basic-auth").Value.String())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("\nProvision was successful:\n")
		fmt.Printf("  Status:     %s\n", provisionResponse.Status)
		fmt.Printf("  Dashboard URL:     %s\n", provisionResponse.DashboardURL)
		fmt.Printf("  Access URL:     %s\n", provisionResponse.AccessURL)

		dashboardURL := provisionResponse.DashboardURL
		if dashboardURL == "" {
			color.Red("The server did not return a dashboard-url. Please make sure your provision endpoint is returning the correct response.\n")
			os.Exit(1)
		}

		user := marketplace.User{
			QuicknodeID:      cmd.Flag("quicknode-id").Value.String(),
			Name:             cmd.Flag("name").Value.String(),
			Email:            cmd.Flag("email").Value.String(),
			OrganizationName: cmd.Flag("org").Value.String(),
		}
		color.Magenta("\n\n→ SSO into %s:\n", dashboardURL)
		userJson, _ := json.MarshalIndent(user, "", "  ")
		fmt.Printf("%s\n", userJson)

		jwtSecret := cmd.Flag("jwt-secret").Value.String()
		jwtToken, err := marketplace.GetJWT(jwtSecret, user)
		if err != nil {
			color.Red("Could not generate JWT: %s", err)
			os.Exit(1)
		}

		fmt.Printf("JWT Token: %s", jwtToken)

		// # Open the browser
		openbrowser(fmt.Sprintf("%s?jwt=%s", dashboardURL, jwtToken))
	},
}

func init() {
	rootCmd.AddCommand(ssoCmd)

	ssoCmd.PersistentFlags().StringP("url", "u", "", "The URL of the add-on's provision endpoint")

	// Note: basic auth defaults to username = Aladdin and password = open sesame
	ssoCmd.PersistentFlags().String("basic-auth", "QWxhZGRpbjpvcGVuIHNlc2FtZQ==", "The basic auth credentials for the add-on. Defaults to username = Aladdin and password = open sesame")

	ssoCmd.PersistentFlags().StringP("quicknode-id", "q", uuid.NewV4().String(), "The QuickNode ID to provision the add-on for (optional)")
	ssoCmd.PersistentFlags().StringP("endpoint-id", "e", uuid.NewV4().String(), "The endpoint ID to provision the add-on for (optional)")
	ssoCmd.PersistentFlags().StringP("chain", "c", "ethereum", "The chain to provision the add-on for")
	ssoCmd.PersistentFlags().StringP("network", "n", "mainnet", "The network to provision the add-on for")
	ssoCmd.PersistentFlags().StringP("plan", "p", "discover", "The plan to provision the add-on for")

	ssoCmd.PersistentFlags().StringP("jwt-secret", "j", "", "The JWT secret for the add-on")
	ssoCmd.PersistentFlags().String("name", "", "The name of the user trying to SSO into the add-on")
	ssoCmd.PersistentFlags().String("email", "", "The email of the user trying to SSO into the add-on")
	ssoCmd.PersistentFlags().String("org", "", "The organization name for the user trying to SSO into the add-on")
}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}
