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

// multitestsCmd represents the tests command
var multitestsCmd = &cobra.Command{
	Use:   "multitests",
	Short: "Get list of multi-tests",
	Long: `Use the command to list multi-tests within a specified workspace. The output includes Multi-test NAME, ID, SCENARIOS, etc. The output can be further filtered by specifying a project id by using the --pid flag.

	For example: [bmgo get -w <workspace id> multitests] OR 
	             [bmgo get -w <workspace id> multitests --pid <project id>]
	For default: [bmgo get --ws multitests] OR 
	             [bmgo get --ws multitests --pid <project id>]`,
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
			listMultiTestsWSProjectraw(projectId)
		case workspaceId != 0 && projectId != 0:
			listMultiTestsWSProject(projectId)
		case workspaceId != 0 && rawOutput:
			listMultiTestsWSraw(workspaceId)
		case workspaceId != 0:
			listMultiTestsWS(workspaceId)
		default:
			cmd.Help()
		}
	},
}

func init() {
	GetCmd.AddCommand(multitestsCmd)
	multitestsCmd.Flags().Int("pid", 0, "Provide the project ID to filter the results")
}

type ListMultiTestsResponse struct {
	Result []listMultiTestsResult `json:"result"`
	Error  errorResult            `json:"error"`
}
type listMultiTestsResult struct {
	Name               string      `json:"name"`
	Id                 int         `json:"id"`
	LastRunTime        int         `json:"lastRunTime"`
	TestsForExecutions []scenarios `json:"testsForExecutions"`
	ProjectId          int         `json:"projectId"`
}
type scenarios struct {
	TestId int `json:"testId"`
}

func listMultiTestsWS(workspaceId int) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	workspaceIdStr := strconv.Itoa(workspaceId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/multi-tests?workspaceId="+workspaceIdStr+"&limit=0", nil)
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
	var responseObjectListMultiTests ListMultiTestsResponse
	json.Unmarshal(bodyText, &responseObjectListMultiTests)
	if responseObjectListMultiTests.Error.Code == 0 {
		fmt.Printf("\n%-10s %-10s %-20s %-10s %-10s\n", "TEST ID", "SCENARIOS", "LAST RUN", "PROJECT", "TEST NAME")
		for i := 0; i < len(responseObjectListMultiTests.Result); i++ {
			testName := responseObjectListMultiTests.Result[i].Name
			testId := responseObjectListMultiTests.Result[i].Id
			testLastRunEp1 := responseObjectListMultiTests.Result[i].LastRunTime
			testLastRunEp := int64(responseObjectListMultiTests.Result[i].LastRunTime)
			testProjectId := responseObjectListMultiTests.Result[i].ProjectId
			totalscenarios := []int{}
			for s := 0; s < len(responseObjectListMultiTests.Result[i].TestsForExecutions); s++ {
				scenario := responseObjectListMultiTests.Result[i].TestsForExecutions[s].TestId
				totalscenarios = append(totalscenarios, scenario)
			}
			// This is because there are epoch values as "0", it converts to a time on 1970s, so we want to condition that here:
			if testLastRunEp1 != 0 {
				testLastRun := time.Unix(testLastRunEp, 0)
				testLastRunSp := fmt.Sprint(testLastRun)
				fmt.Printf("\n%-10v %-10v %-20s %-10d %-10s", testId, len(totalscenarios), testLastRunSp[0:16], testProjectId, testName)
			} else {
				testLastRun := testLastRunEp1
				fmt.Printf("\n%-10v %-10v %-20d %-10d %-10s", testId, len(totalscenarios), testLastRun, testProjectId, testName)
			}
		}
		fmt.Println("\n-")
	} else {
		errorCode := responseObjectListMultiTests.Error.Code
		errorMessage := responseObjectListMultiTests.Error.Message
		fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
	}
}
func listMultiTestsWSraw(workspaceId int) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	workspaceIdStr := strconv.Itoa(workspaceId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/multi-tests?workspaceId="+workspaceIdStr+"&limit=0", nil)
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

func listMultiTestsWSProject(projectId int) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	projectIdStr := strconv.Itoa(projectId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/multi-tests?projectId="+projectIdStr+"&limit=0", nil)
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
	var responseObjectListMultiTests ListMultiTestsResponse
	json.Unmarshal(bodyText, &responseObjectListMultiTests)
	if responseObjectListMultiTests.Error.Code == 0 {
		fmt.Printf("\n%-10s %-10s %-20s %-10s %-10s\n", "TEST ID", "SCENARIOS", "LAST RUN", "PROJECT", "TEST NAME")
		for i := 0; i < len(responseObjectListMultiTests.Result); i++ {
			testName := responseObjectListMultiTests.Result[i].Name
			testId := responseObjectListMultiTests.Result[i].Id
			testLastRunEp1 := responseObjectListMultiTests.Result[i].LastRunTime
			testLastRunEp := int64(responseObjectListMultiTests.Result[i].LastRunTime)
			testProjectId := responseObjectListMultiTests.Result[i].ProjectId
			totalscenarios := []int{}
			for s := 0; s < len(responseObjectListMultiTests.Result[i].TestsForExecutions); s++ {
				scenario := responseObjectListMultiTests.Result[i].TestsForExecutions[s].TestId
				totalscenarios = append(totalscenarios, scenario)
			}
			// This is because there are epoch values as "0", it converts to a time on 1970s, so we want to condition that here:
			if testLastRunEp1 != 0 {
				testLastRun := time.Unix(testLastRunEp, 0)
				testLastRunSp := fmt.Sprint(testLastRun)
				fmt.Printf("\n%-10v %-10v %-20s %-10d %-10s", testId, len(totalscenarios), testLastRunSp[0:16], testProjectId, testName)
			} else {
				fmt.Printf("\n%-10v %-10v %-20d %-10d %-10s", testId, len(totalscenarios), testLastRunEp1, testProjectId, testName)
			}
		}
		fmt.Println("\n-")
	} else {
		errorCode := responseObjectListMultiTests.Error.Code
		errorMessage := responseObjectListMultiTests.Error.Message
		fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
	}
}
func listMultiTestsWSProjectraw(projectId int) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	projectIdStr := strconv.Itoa(projectId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/multi-tests?projectId="+projectIdStr+"&limit=0", nil)
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
