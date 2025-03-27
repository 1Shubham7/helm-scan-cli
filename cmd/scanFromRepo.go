/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/1shubham7/helm-scan-cli/cmd/repo"
)

// scanFromRepoCmd represents the scanFromRepo command
var scanFromRepoCmd = &cobra.Command{
	Use:   "scanFromRepo",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		
		repo, err := cmd.Flags().GetString("repo")
		if err != nil || repo == "" {
			fmt.Println("Error: Repository name is required")
			fmt.Println("Please use the -r or --repo flag to specify the repository")
			os.Exit(1)
		}

		chart, err := cmd.Flags().GetString("chart")
		if err != nil || chart == "" {
			fmt.Println("Error: Chart name is required")
			fmt.Println("Please use the -c or --chart flag to specify the chart name")
			os.Exit(1)
		}

		c, err := repository.DownloadChart(repo, chart)
		if err != nil {
			fmt.Println("Error in downloading chart...")
			os.Exit(1)
		}

		resp, err := repository.FindImages(c)
		if err != nil {
			fmt.Println("Error in finding images...")
			os.Exit(1)
		}

		for i, img := range resp {
			fmt.Printf("%d. Image scanned:\n", i+1)
			fmt.Printf("   Image:  %s\n", img.Name)
			fmt.Printf("   Size:   %s\n", img.Size)
			fmt.Printf("   Layers: %d\n\n", img.Layers)
		}			
	},
}

func init() {
	rootCmd.AddCommand(scanFromRepoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanFromRepoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	scanFromRepoCmd.Flags().StringP("repo", "r", "", "Repository Name")
	scanFromRepoCmd.Flags().StringP("chart", "c", "", "Chart Name")
}
