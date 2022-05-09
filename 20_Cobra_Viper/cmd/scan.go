/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Get scan result from Trivy",
	Long: `Get scan result from Trivy
	You need to setting the connection of Trivy in config file.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		trivy := viper.GetStringMapString("trivy")
		if trivy == nil {
			fmt.Println("Please specify trivy server info in config file")
			os.Exit(1)
		}

		// Get images list
		images, _ := listImages()

		// for each image, scan it via trivy fmt.Println(trivy)
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Image", "Result"})

		for _, img := range images {

			cmd := "trivy client --remote http://" + trivy["server"] + ":" + trivy["port"] + " " + img + " | grep 'Total'"
			out, err := exec.Command("bash", "-c", cmd).Output()

			if err != nil {
				// fmt.Println("Failed to execute command:", cmd)
				t.AppendRow([]interface{}{img, "Unsupported"})
			} else {
				t.AppendRow([]interface{}{img, string(out)})
			}
			t.AppendSeparator()
		}
		t.Render()

	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}
