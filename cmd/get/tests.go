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
		ws, _ := cmd.Flags().GetBool("ws")
		var workspaceId int
		if ws {
			workspaceId = defaultWorkspace()
		} else {
			workspaceId, _ = cmd.Flags().GetInt("workspaceid")
		}
		rawOutput, _ := cmd.Flags().GetBool("raw")
		projectId, _ := cmd.Flags().GetInt("pid")

		switch {
		case workspaceId != 0 && projectId != 0 && rawOutput:
			listTestsWSProjectraw(workspaceId, projectId)
		case workspaceId != 0 && projectId != 0:
			listTestsWSProject(workspaceId, projectId)
		case workspaceId != 0 && rawOutput:
			listTestsWSraw(workspaceId)
		case workspaceId != 0:
			listTestsWS(workspaceId)
		default:
			fmt.Println("\nPlease provide a valid Workspace ID to get list of tests")
			fmt.Println("[bmgo get -w <workspace id>...")
		}
	},
}

func init() {
	GetCmd.AddCommand(testsCmd)
	testsCmd.Flags().Int("pid", 0, "Provide the project ID to filter the results")
}

type ListTestsResponse struct {
	Result []listTestsResult `json:"result"`
	Error  errorResult       `json:"error"`
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
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/tests?workspaceId="+workspaceIdStr+"&limit=0", nil)
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
	var responseObjectListTests ListTestsResponse
	json.Unmarshal(bodyText, &responseObjectListTests)
	if responseObjectListTests.Error.Code == 0 {
		fmt.Printf("\n%-10s %-20s %-15s\n", "TEST ID", "LAST RUN", "TEST NAME")
		for i := 0; i < len(responseObjectListTests.Result); i++ {
			testName := responseObjectListTests.Result[i].Name
			testId := responseObjectListTests.Result[i].Id
			testLastRunEp1 := responseObjectListTests.Result[i].LastRunTime
			testLastRunEp := int64(responseObjectListTests.Result[i].LastRunTime)
			if testLastRunEp1 != 0 {
				testLastRun := time.Unix(testLastRunEp, 0)
				testLastRunSp := fmt.Sprint(testLastRun)
				fmt.Printf("\n%-10v %-20s %-15s", testId, testLastRunSp[0:16], testName)
			} else {
				testLastRun := testLastRunEp1
				fmt.Printf("\n%-10v %-20v %-15s", testId, testLastRun, testName)
			}
		}
		fmt.Println("\n-")
	} else {
		errorCode := responseObjectListTests.Error.Code
		errorMessage := responseObjectListTests.Error.Message
		fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
	}

}

func listTestsWSraw(workspaceId int) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	workspaceIdStr := strconv.Itoa(workspaceId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/tests?workspaceId="+workspaceIdStr+"&limit=0", nil)
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

func listTestsWSProject(workspaceId, projectId int) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	workspaceIdStr := strconv.Itoa(workspaceId)
	projectIdStr := strconv.Itoa(projectId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/tests?workspaceId="+workspaceIdStr+"&projectId="+projectIdStr+"&limit=0", nil)
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
	var responseObjectListTests ListTestsResponse
	json.Unmarshal(bodyText, &responseObjectListTests)
	if responseObjectListTests.Error.Code == 0 {
		fmt.Printf("\n%-10s %-20s %-15s\n", "TEST ID", "LAST RUN", "TEST NAME")
		for i := 0; i < len(responseObjectListTests.Result); i++ {
			testName := responseObjectListTests.Result[i].Name
			testId := responseObjectListTests.Result[i].Id
			testLastRunEp1 := responseObjectListTests.Result[i].LastRunTime
			testLastRunEp := int64(responseObjectListTests.Result[i].LastRunTime)
			if testLastRunEp1 != 0 {
				testLastRun := time.Unix(testLastRunEp, 0)
				testLastRunSp := fmt.Sprint(testLastRun)
				fmt.Printf("\n%-10v %-20s %-15s", testId, testLastRunSp[0:16], testName)
			} else {
				testLastRun := testLastRunEp1
				fmt.Printf("\n%-10v %-20v %-15s", testId, testLastRun, testName)
			}
		}
		fmt.Println("\n-")
	} else {
		errorCode := responseObjectListTests.Error.Code
		errorMessage := responseObjectListTests.Error.Message
		fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
	}

}
func listTestsWSProjectraw(workspaceId, projectId int) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	workspaceIdStr := strconv.Itoa(workspaceId)
	projectIdStr := strconv.Itoa(projectId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/tests?workspaceId="+workspaceIdStr+"&projectId="+projectIdStr+"&limit=0", nil)
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
