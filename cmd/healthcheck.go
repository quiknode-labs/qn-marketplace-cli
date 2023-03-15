/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/quiknode-labs/qn-marketplace-cli/marketplace"
	"github.com/spf13/cobra"
)

// healthcheckCmd represents the healthcheck command
var healthcheckCmd = &cobra.Command{
	Use:   "healthcheck",
	Short: "Allows you to test your add-on's healthcheck implementation",
	Args:  cobra.OnlyValidArgs,
	Run: func(cmd *cobra.Command, args []string) {
		header := color.New(color.FgWhite, color.BgBlue).SprintFunc()
		fmt.Printf("%s\n\n", header("        Healthcheck        "))
		verbose := cmd.Flag("verbose").Value.String() == "true"
		url := cmd.Flag("url").Value.String()

		if verbose {
			color.Blue("→ GET %s:\n", url)
		}
		statusCode, responseBody, err := marketplace.Healthcheck(url)
		if err != nil {
			color.Red("%s", err)
			os.Exit(1)
		}

		if verbose {
			fmt.Println("\nResponse body:")
			fmt.Printf("%s\n\n", responseBody)
		}

		color.Green("  ✓ Healthcheck was successful with HTTP response code: %d", statusCode)
	},
}

func init() {
	rootCmd.AddCommand(healthcheckCmd)

	healthcheckCmd.PersistentFlags().StringP("url", "u", "", "The URL of the add-on's healthcheck URL to test")
}
