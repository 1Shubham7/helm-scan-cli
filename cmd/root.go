/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const logo = `
	╔═══════════╗
 	║ helm-scan ║
 	╚═══════════╝
`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "helm-scan",
	Short: "A brief description of your application",
	Long: `
	
Helm-Scan is a command-line tool that scans a given Helm chart,
identifies container images used, and provides metadata such as image size
and the number of layers for each image.

Usage examples:
  helm-scan scan-from-link [CHART_LINK]
  helm-scan scan-from-repo -r [REPO_NAME] -c [CHART_NAME]
  helm-scan scan-from-file [FILE_PATH]

This tool helps developers and security teams analyze container images
embedded in Helm charts before deployment.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(logo)
		fmt.Println(cmd.Long)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


