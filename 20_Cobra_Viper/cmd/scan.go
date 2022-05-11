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

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	gojsonq "github.com/thedevsaddam/gojsonq/v2"
)

type vulResult struct {
	image   string
	high    int
	med     int
	low     int
	unknow  int
	support bool
}

// type vulResultList []vulResult

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
		if true {
			getJsonResult(images, trivy["server"], trivy["port"])
		} else {
			showBriefResult(images, trivy["server"], trivy["port"])
			showSortedResult(images, trivy["server"], trivy["port"])
		}
	},
}

func getJsonResult(images []string, server string, port string) {
	// for each image, scan it via trivy fmt.Println(trivy)
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Image", "High", "Medium", "Low", "Unknow"})

	for _, img := range images {

		cmd1 := "trivy client --format json --remote http://" + server + ":" + port + " " + img
		jsondata, err1 := exec.Command("bash", "-c", cmd1).Output()

		if err1 != nil {
			fmt.Println("Failed to execute command:", cmd1)
			os.Exit(1)
		}

		results, _ := gojsonq.New().FromString(string(jsondata)).From("Results.[].Vulnerabilities.[]").FindR("Severity")
		name, _ := results.String()
		fmt.Printf("%#v\n", name)
		//fmt.Printf("%#v\n", results.Get())
	}
}

func showSortedResult(images []string, server string, port string) {

	// for each image, scan it via trivy fmt.Println(trivy)

	imgVulResults := []vulResult{}
	for _, img := range images {

		// cmd := "trivy client --remote http://" + trivy["server"] + ":" + trivy["port"] + " " + img + " | grep 'Total'"
		// out, err := exec.Command("bash", "-c", cmd).Output()
		cmd1 := "trivy client --format json --remote http://" + server + ":" + port + " " + img
		parseResultCmd := cmd1 + " | jq -r \".Results[]\""
		_, err1 := exec.Command("bash", "-c", parseResultCmd).Output()

		if err1 != nil {
			fmt.Println("Failed to execute command:", cmd1)
			//t.AppendRow([]interface{}{img, -1, -1, -1, -1, "Unsupported"})
			imgVulResults = append(imgVulResults, vulResult{img, -1, -1, -1, -1, false})
			continue
		}

		cmd2 := cmd1 + " | jq -r \".Results[].Vulnerabilities[].Severity\""
		out2, err2 := exec.Command("bash", "-c", cmd2).Output()
		if err2 != nil {
			fmt.Println("Failed to parse result of command:", string(cmd1))
			//t.AppendRow([]interface{}{img, 0, 0, 0, 0})
			imgVulResults = append(imgVulResults, vulResult{img, 0, 0, 0, 0, true})
		} else {
			splitFn := func(c rune) bool {
				return c == '\n'
			}
			severity := strings.FieldsFunc(string(out2), splitFn)
			//fmt.Println(img, ":", severity)
			vulCount := countOccurence(severity)
			//t.AppendRow([]interface{}{img, vulCount["HIGH"], vulCount["MEDIUM"], vulCount["LOW"], vulCount["UNKNOWN"]})
			imgVulResults = append(imgVulResults, vulResult{img, vulCount["HIGH"], vulCount["MEDIUM"], vulCount["LOW"], vulCount["UNKNOWN"], true})
		}

		// t.AppendSeparator()
	}

	sort.Slice(imgVulResults, func(i, j int) bool {
		return imgVulResults[i].high < imgVulResults[j].high
	})

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Image", "High", "Medium", "Low", "Unknow"})
	for _, r := range imgVulResults {
		// fmt.Println(value.image, value.h, value.m, value.l)
		if r.support {
			t.AppendRow([]interface{}{r.image, r.high, r.med, r.low, r.unknow, ""})
		} else {
			t.AppendRow([]interface{}{r.image, r.high, r.med, r.low, r.unknow, "Unsupported"})
		}
		t.AppendSeparator()

	}
	t.Render()

}

func countOccurence(apps []string) map[string]int {
	dict := make(map[string]int)
	for _, v := range apps {
		dict[v]++
	}
	return dict
}

func showBriefResult(images []string, server string, port string) {
	// for each image, scan it via trivy fmt.Println(trivy)
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Image", "Result"})

	for _, img := range images {

		cmd := "trivy client --remote http://" + server + ":" + port + " " + img + " | grep 'Total'"
		out, err := exec.Command("bash", "-c", cmd).Output()

		if err != nil {
			//fmt.Println("Failed to execute command:", cmd)
			t.AppendRow([]interface{}{img, "Unsupported"})
		} else {
			t.AppendRow([]interface{}{img, string(out)})
		}
		t.AppendSeparator()
	}
	t.Render()

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
