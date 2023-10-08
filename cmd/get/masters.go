/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package get

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
)

// mastersCmd represents the masters command
var mastersCmd = &cobra.Command{
	Use:   "masters",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("masters called")
		testId, _ := cmd.Flags().GetInt("tid")
		getMasters(testId)
	},
}

func init() {
	GetCmd.AddCommand(mastersCmd)
	mastersCmd.Flags().Int("tid", 0, "Provide the test ID to list masters")
}

func getMasters(testId int) {
	apiId, apiSecret := Getapikeys()
	testIdStr := strconv.Itoa(testId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/masters?testId="+testIdStr+"&limit=0", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(apiId, apiSecret)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", bodyText)
}
