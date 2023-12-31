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

// projectsCmd represents the projects command
var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Get the list of Projects under workspace",
	Long: `Use the command to list Projects within a specified workspace or account. Projects are designed to organize tests and reports and track usage within a Workspace. The output includes Project NAME, ID, Test count, etc.

	For example: [bmgo get -w <workspace id> projects] OR 
	             [bmgo get -a <account id> projects]
	For default: [bmgo get --ac projects] OR 
	             [bmgo get --ws projects]`,
	Run: func(cmd *cobra.Command, args []string) {
		ac, _ := cmd.Flags().GetBool("ac")
		ws, _ := cmd.Flags().GetBool("ws")
		var accountId, workspaceId int
		if ac {
			accountId = defaultAccount()
		} else {
			accountId, _ = cmd.Flags().GetInt("accountid")
		}
		if ws {
			workspaceId = defaultWorkspace()
		} else {
			workspaceId, _ = cmd.Flags().GetInt("workspaceid")
		}
		rawOutput, _ := cmd.Flags().GetBool("raw")
		switch {
		case workspaceId == 0 && accountId != 0 && rawOutput:
			getProjectsAraw(accountId)
		case workspaceId != 0 && accountId == 0 && rawOutput:
			getProjectsWsraw(workspaceId)
		case workspaceId == 0 && accountId != 0 && !rawOutput:
			getProjectsA(accountId)
		case workspaceId != 0 && accountId == 0 && !rawOutput:
			getProjectsWs(workspaceId)
		default:
			cmd.Help()
		}
	},
}

func init() {
	GetCmd.AddCommand(projectsCmd)
}

type projectsResponse struct {
	Result []projectsResult `json:"result"`
}
type projectsResult struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	TestsCount int    `json:"testsCount"`
	Created    int    `json:"created"`
}

func getProjectsWs(workspaceId int) {
	apiId, apiSecret := Getapikeys()
	workspaceIdstr := strconv.Itoa(workspaceId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/projects?workspaceId="+workspaceIdstr+"&limit=0", nil)
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
	var responseBodyProjectsWs projectsResponse
	json.Unmarshal(bodyText, &responseBodyProjectsWs)
	fmt.Printf("\n%-10s %-25s %-8s %-20s\n", "ID", "NAME", "TESTS", "CREATED")
	for i := 0; i < len(responseBodyProjectsWs.Result); i++ {
		projectId := responseBodyProjectsWs.Result[i].Id
		projectName := responseBodyProjectsWs.Result[i].Name
		projectTests := responseBodyProjectsWs.Result[i].TestsCount
		pCreatedEpoch := int64(responseBodyProjectsWs.Result[i].Created)
		projectCreated := fmt.Sprint(time.Unix(pCreatedEpoch, 0))
		fmt.Printf("\n%-10v %-25s %-8v %-20v", projectId, projectName, projectTests, projectCreated[0:10])
	}
	fmt.Println("\n-")
}
func getProjectsWsraw(workspaceId int) {
	apiId, apiSecret := Getapikeys()
	workspaceIdstr := strconv.Itoa(workspaceId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/projects?workspaceId="+workspaceIdstr+"&limit=0", nil)
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

func getProjectsA(accountId int) {
	apiId, apiSecret := Getapikeys()
	accountIdstr := strconv.Itoa(accountId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/projects?accountId="+accountIdstr+"&limit=0", nil)
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
	var responseBodyProjectsA projectsResponse
	json.Unmarshal(bodyText, &responseBodyProjectsA)
	fmt.Printf("\n%-10s %-25s %-8s %-10s\n", "ID", "NAME", "TESTS", "CREATED")
	for i := 0; i < len(responseBodyProjectsA.Result); i++ {
		projectId := responseBodyProjectsA.Result[i].Id
		projectName := responseBodyProjectsA.Result[i].Name
		projectTests := responseBodyProjectsA.Result[i].TestsCount
		pCreatedEpoch := int64(responseBodyProjectsA.Result[i].Created)
		projectCreated := fmt.Sprint(time.Unix(pCreatedEpoch, 0))
		fmt.Printf("\n%-10v %-25s %-8v %-10v", projectId, projectName, projectTests, projectCreated[0:10])
	}
	fmt.Println("\n-")
}
func getProjectsAraw(accountId int) {
	apiId, apiSecret := Getapikeys()
	accountIdstr := strconv.Itoa(accountId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/projects?accountId="+accountIdstr+"&limit=0", nil)
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
