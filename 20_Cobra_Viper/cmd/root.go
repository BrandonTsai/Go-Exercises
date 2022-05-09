/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var descend bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cim",
	Short: "CLI for Container Images Management",
	Long: `CLI for Container Images Management.
	It shows a list sorted container images' name by default
	`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

		// Handle Local Flag `-v, --version`
		version, err := cmd.Flags().GetBool("version")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if version {
			fmt.Println("Version: ", "0.1.0")
			os.Exit(0)
		}

		// Call function to perform basic actions
		images, _ := listImages()
		// Print the output
		for _, value := range images {
			fmt.Println(value)
		}
	},
}

func listImages() ([]string, error) {

	// Handle the settings in config file via Viper
	baseTool := "docker"
	if viper.GetBool("podman") {
		baseTool = "podman"
	}

	// Get images via command
	cmd := exec.Command(baseTool, "images", "--format", "{{.Repository}}:{{.Tag}}")
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Sort images

	splitFn := func(c rune) bool {
		return c == '\n'
	}
	imgs := strings.FieldsFunc(string(stdout), splitFn)

	//imgs := strings.Split(string(stdout), "\n")
	if descend {
		sort.Sort(sort.Reverse(sort.StringSlice(imgs)))
	} else {
		sort.Strings(imgs)
	}

	return imgs, nil

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
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cim.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&descend, "descend", "d", false, "show descend sorted images")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("version", "v", false, "Show current version.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cim" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cim")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
