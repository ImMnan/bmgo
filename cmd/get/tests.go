/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package get

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

// testsCmd represents the tests command
var testsCmd = &cobra.Command{
	Use:   "tests",
	Short: "Get list of tests",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		workspacesId, _ := cmd.Flags().GetInt("workspaceid")
		rawOutput, _ := cmd.Flags().GetBool("raw")
		if workspacesId != 0 && rawOutput {
			listTestsWSraw(workspacesId)
		} else if workspacesId != 0 {
			listTestsWS(workspacesId)
		} else {
			fmt.Println("\nPlease provide a valid Workspace ID to get list of tests")
			fmt.Println("[bmgo get -w <workspace id>...")
		}
	},
}

func init() {
	GetCmd.AddCommand(testsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type ListTestsResponse struct {
	Result []listTestsResult `json:"result"`
}
type listTestsResult struct {
	Name        string `json:"name"`
	Id          int    `json:"id"`
	LastRunTime int    `json:"lastRunTime"`
}

func listTestsWS(workspaceId int) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	workspaceIdStr := strconv.Itoa(workspaceId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/tests?workspaceId="+workspaceIdStr, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
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
	var responseObjectListTests ListTestsResponse
	json.Unmarshal(bodyText, &responseObjectListTests)
	fmt.Printf("%-10s %-30s %-15s\n", "TEST ID", "LAST RUN", "TEST NAME")
	for i := 0; i < len(responseObjectListTests.Result); i++ {
		testName := responseObjectListTests.Result[i].Name
		testId := responseObjectListTests.Result[i].Id
		testLastRunEp1 := responseObjectListTests.Result[i].LastRunTime
		testLastRunEp := int64(responseObjectListTests.Result[i].LastRunTime)
		if testLastRunEp1 != 0 {
			testLastRun := time.Unix(testLastRunEp, 0)
			fmt.Printf("\n%-10v %-30v %-15s", testId, testLastRun, testName)
		} else {
			testLastRun := testLastRunEp1
			fmt.Printf("\n%-10v %-30v %-15s", testId, testLastRun, testName)
		}
	}
	fmt.Println("\n-")
}

func listTestsWSraw(workspaceId int) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	workspaceIdStr := strconv.Itoa(workspaceId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/tests?workspaceId="+workspaceIdStr, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
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
