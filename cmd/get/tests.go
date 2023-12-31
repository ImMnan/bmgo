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
	Long: `Use the command to list Tests within a specified workspace. The output includes Test ID, Name, Project Id, etc. The output can be further filtered by specifying a project id by using the --pid flag.

	For example: [bmgo get -w <workspace id> tests] OR
	             [bmgo get -w <workspace id> tests --pid <project id>]
	For default: [bmgo get --ws tests] OR 
	             [bmgo get --ws tests --pid <project id>]`,
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
		case !rawOutput && (workspaceId != 0 || projectId != 0):
			listTestsWS(workspaceId, projectId)
		case rawOutput && (workspaceId != 0 || projectId != 0):
			listTestsWSraw(workspaceId, projectId)
		default:
			cmd.Help()
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
	Name          string             `json:"name"`
	Id            int                `json:"id"`
	LastRunTime   int                `json:"lastRunTime"`
	ProjectId     int                `json:"projectId"`
	Configuration config             `json:"configuration"`
	Executions    []executiondetails `json:"executions"`
}
type config struct {
	DedicatedIpsEnabled     bool `json:"dedicatedIpsEnabled"`
	EnableTestData          bool `json:"enableTestData"`
	EnableLoadConfiguration bool `json:"enableLoadConfiguration"`
}
type executiondetails struct {
	Locations string
}

func listTestsWS(workspaceId, projectId int) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	var req *http.Request
	var err error
	if projectId != 0 {
		projectIdStr := strconv.Itoa(projectId)
		req, err = http.NewRequest("GET", "https://a.blazemeter.com/api/v4/tests?projectId="+projectIdStr+"&limit=0", nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		workspaceIdStr := strconv.Itoa(workspaceId)
		req, err = http.NewRequest("GET", "https://a.blazemeter.com/api/v4/tests?workspaceId="+workspaceIdStr+"&limit=0", nil)
		if err != nil {
			log.Fatal(err)
		}
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
		fmt.Printf("\n%-10s %-20s %-10s %-15s\n", "TEST ID", "LAST RUN", "PROJECT", "TEST NAME")
		for i := 0; i < len(responseObjectListTests.Result); i++ {
			testName := responseObjectListTests.Result[i].Name
			testId := responseObjectListTests.Result[i].Id
			testLastRunEp1 := responseObjectListTests.Result[i].LastRunTime
			testProjectId := responseObjectListTests.Result[i].ProjectId
			testLastRunEp := int64(responseObjectListTests.Result[i].LastRunTime)
			// This is because there are epoch values as "0", it converts to a time on 1970s, so we want to condition that here:
			if testLastRunEp1 != 0 {
				testLastRun := time.Unix(testLastRunEp, 0)
				testLastRunSp := fmt.Sprint(testLastRun)
				fmt.Printf("\n%-10v %-20s %-10d %-15s", testId, testLastRunSp[0:16], testProjectId, testName)
			} else {
				fmt.Printf("\n%-10v %-20v %-10d %-15s", testId, testLastRunEp1, testProjectId, testName)
			}
		}
		fmt.Println("\n-")
	} else {
		errorCode := responseObjectListTests.Error.Code
		errorMessage := responseObjectListTests.Error.Message
		fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
	}
}
func listTestsWSraw(workspaceId, projectId int) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	var req *http.Request
	var err error
	if projectId != 0 {
		projectIdStr := strconv.Itoa(projectId)
		req, err = http.NewRequest("GET", "https://a.blazemeter.com/api/v4/tests?projectId="+projectIdStr+"&limit=0", nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		workspaceIdStr := strconv.Itoa(workspaceId)
		req, err = http.NewRequest("GET", "https://a.blazemeter.com/api/v4/tests?workspaceId="+workspaceIdStr+"&limit=0", nil)
		if err != nil {
			log.Fatal(err)
		}
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
