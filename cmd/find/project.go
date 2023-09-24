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

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Find Project using Project ID",
	Long:  `.`,
	Run: func(cmd *cobra.Command, args []string) {
		projectId, _ := cmd.Flags().GetInt("pid")
		rawOutput, _ := cmd.Flags().GetBool("raw")
		if projectId != 0 && rawOutput {
			findProjectraw(projectId)
		} else if projectId != 0 {
			findProject(projectId)
		} else {
			fmt.Println("\nPlease provide a valid Project ID to find the Project")
			fmt.Println("[bmgo find project --pid <Project id>")
		}
	},
}

func init() {
	FindCmd.AddCommand(projectCmd)
	projectCmd.Flags().Int("pid", 0, "Provide the project id to find")
	projectCmd.MarkFlagRequired("pid")
}

type ProjectResponse struct {
	Result projectResult `json:"result"`
}
type projectResult struct {
	Name        string `json:"name"`
	WorkspaceId int    `json:"workspaceId"`
	Created     int    `json:"created"`
}
type ListTestsResponse struct {
	Result []listTestsResult `json:"result"`
}
type listTestsResult struct {
	Name        string `json:"name"`
	Id          int    `json:"id"`
	LastRunTime int    `json:"lastRunTime"`
}

func findProject(projectId int) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	projectIdStr := strconv.Itoa(projectId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/projects/"+projectIdStr, nil)
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
	var responseObjectProject ProjectResponse
	json.Unmarshal(bodyText, &responseObjectProject)
	fmt.Printf("\n%-20s %-15s %-15s\n", "PROJECT NAME", "WORKSPACE", "CREATED")
	projectName := responseObjectProject.Result.Name
	projectWorkspace := responseObjectProject.Result.WorkspaceId
	projectCreatedEp := int64(responseObjectProject.Result.Created)
	projectCreated := time.Unix(projectCreatedEp, 0)
	fmt.Printf("%-20s %-15v %-15v\n", projectName, projectWorkspace, projectCreated)
	listTestsProject(projectId)
}
func listTestsProject(projectId int) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	projectIdStr := strconv.Itoa(projectId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/tests?projectId="+projectIdStr, nil)
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
	fmt.Println("\n---------------------------------------------------------------------------------------------")
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

func findProjectraw(projectId int) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	projectIdStr := strconv.Itoa(projectId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/projects/"+projectIdStr, nil)
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
