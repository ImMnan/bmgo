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
		rawOutput, _ := cmd.Flags().GetBool("raw")
		if testId != 0 && rawOutput {
			findTestraw(testId)
		} else if testId != 0 {
			findTest(testId)
		} else {
			fmt.Println("\nPlease provide a valid Test ID to find the Test")
			fmt.Println("[bmgo find test --tid <Test id>")
		}
	},
}

func init() {
	FindCmd.AddCommand(testCmd)
	testCmd.Flags().Int("tid", 0, "Provide a test Id to find a test")
	testCmd.MarkFlagRequired("tid")
}

type FindTestsResponse struct {
	Result findTestsResult `json:"result"`
}
type findTestsResult struct {
	Name               string                    `json:"name"`
	Id                 int                       `json:"id"`
	LastRunTime        int                       `json:"lastRunTime"`
	OverrideExecutions []overrideExecutionsArray `json:"overrideExecutions"`
	ProjectId          int                       `json:"projectId"`
}
type overrideExecutionsArray struct {
	Executor string `json:"executor"`
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
	//fmt.Printf("%s\n", bodyText)
	var responseObjectTest FindTestsResponse
	json.Unmarshal(bodyText, &responseObjectTest)
	testName := responseObjectTest.Result.Name
	fmt.Printf("\nTEST NAME: %s\n", testName)
	testLastRunEp1 := responseObjectTest.Result.LastRunTime
	testProjectId := responseObjectTest.Result.ProjectId
	testLastRunEp := int64(responseObjectTest.Result.LastRunTime)
	var testExecutor string

	fmt.Println("\n---------------------------------------------------------------------------------------------")
	fmt.Printf("%-10v %-20s %-10s %-10s\n", "TEST ID", "LAST RUN", "PROJECT", "EXECUTOR")
	for i := 0; i < len(responseObjectTest.Result.OverrideExecutions); i++ {
		testExecutor = responseObjectTest.Result.OverrideExecutions[i].Executor
	}
	if testLastRunEp1 != 0 {
		testLastRun := time.Unix(testLastRunEp, 0)
		testLastRunSp := fmt.Sprint(testLastRun)
		fmt.Printf("%-10v %-20s %-10d %-10s", testId, testLastRunSp[0:16], testProjectId, testExecutor)
	} else {
		testLastRun := testLastRunEp1
		fmt.Printf("%-10v %-20v %-10d %-10s\n", testId, testLastRun, testProjectId, testExecutor)
	}
}
func findTestraw(testId int) {
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
}