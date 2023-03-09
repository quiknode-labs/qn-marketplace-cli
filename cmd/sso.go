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
		dashboardURL := cmd.Flag("dashboard-url").Value.String()
		if dashboardURL == "" {
			fmt.Print("Please provide a dashboard URL for the provision API via the --dashboard-url flag\n")
			os.Exit(1)
		}

		user := marketplace.User{
			QuicknodeID:      cmd.Flag("quicknode-id").Value.String(),
			Name:             cmd.Flag("name").Value.String(),
			Email:            cmd.Flag("email").Value.String(),
			OrganizationName: cmd.Flag("org").Value.String(),
		}
		color.Magenta("→ SSO into %s:\n", dashboardURL)
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

	ssoCmd.PersistentFlags().StringP("dashboard-url", "u", "", "The Dashbaord URL of the add-on you want to SSO into")

	ssoCmd.PersistentFlags().StringP("jwt-secret", "j", "", "The JWT secret for the add-on")
	ssoCmd.PersistentFlags().String("name", "", "The name of the user trying to SSO into the add-on")
	ssoCmd.PersistentFlags().String("email", "", "The email of the user trying to SSO into the add-on")
	ssoCmd.PersistentFlags().String("org", "", "The organization name for the user trying to SSO into the add-on")
	ssoCmd.PersistentFlags().StringP("quicknode-id", "q", uuid.NewV4().String(), "The QuickNode ID to provision the add-on for (optional)")
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
