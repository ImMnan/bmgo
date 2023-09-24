/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package find

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Find test details",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		testId, _ := cmd.Flags().GetInt("tid")
		findTest(testId)
	},
}

func init() {
	FindCmd.AddCommand(testCmd)
	testCmd.Flags().Int("tid", 0, "Provide a test Id to find a test")
}

type FindTestsResponse struct {
	Result findTestsResult `json:"result"`
}
type findTestsResult struct {
	Name        string `json:"name"`
	Id          int    `json:"id"`
	LastRunTime int    `json:"lastRunTime"`
}

func findTest(testId int) {
	testIdStr := strconv.Itoa(testId)
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/tests/"+testIdStr, nil)
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
	var responseObjectTest FindTestsResponse
	json.Unmarshal(bodyText, &responseObjectTest)
	fmt.Printf("%-10s %-20s %-15s\n", "TEST ID", "LAST RUN", "TEST NAME")
	testName := responseObjectTest.Result.Name
	testLastRunEp1 := responseObjectTest.Result.LastRunTime
	testLastRunEp := int64(responseObjectTest.Result.LastRunTime)
	if testLastRunEp1 != 0 {
		testLastRun := time.Unix(testLastRunEp, 0)
		testLastRunSp := fmt.Sprint(testLastRun)
		fmt.Printf("\n%-10v %-20s %-15s", testId, testLastRunSp[0:16], testName)
	} else {
		testLastRun := testLastRunEp1
		fmt.Printf("\n%-10v %-20v %-15s", testId, testLastRun, testName)
	}
}
