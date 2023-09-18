/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package add

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "A brief description of your command",
	Long:  `.`,
	Run: func(cmd *cobra.Command, args []string) {
		workspaceId, _ := cmd.Flags().GetInt("workspaceid")
		workspaceName, _ := cmd.Flags().GetString("name")
		rawOutput, _ := cmd.Flags().GetBool("raw")
		switch {
		case (workspaceId != 0) && (workspaceName != "") && rawOutput:
			addProjectraw(workspaceName, workspaceId)
		case (workspaceId != 0) && (workspaceName != ""):
			addProject(workspaceName, workspaceId)
		default:
			fmt.Println("\nPlease provide a correct Workspace Id & Project Name")
			fmt.Println("[bmgo add -w <workspace id> project --name <project name>]\n-")
		}
	},
}

func init() {
	AddCmd.AddCommand(projectCmd)
	projectCmd.Flags().String("name", "", "Name your Project")
	projectCmd.MarkFlagRequired("name")
}

type addProjectResponse struct {
	Result addProjectResult `json:"result"`
}
type addProjectResult struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func addProject(workspaceName string, workspaceId int) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	data := fmt.Sprintf(`{
		"name": "%s", 
		"description": "Project created through bmgo", 
		"workspaceId": %v}`, workspaceName, workspaceId)
	reqBodyData := strings.NewReader(data)
	req, err := http.NewRequest("POST", "https://a.blazemeter.com/api/v4/projects", reqBodyData)
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
	//fmt.Printf("%s\n", bodyText)
	var responseBodyAddProject addProjectResponse
	json.Unmarshal(bodyText, &responseBodyAddProject)
	projectId := responseBodyAddProject.Result.Id
	projectName := responseBodyAddProject.Result.Name
	fmt.Printf("\n%-15s %-15s", "PROJECT ID", "NAME")
	fmt.Printf("\n%-15v %-15s", projectId, projectName)
	fmt.Println("\n-")
}
func addProjectraw(workspaceName string, workspaceId int) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	data := fmt.Sprintf(`{
		"name": "%s", 
		"description": "Project created through bmgo", 
		"workspaceId": %v}`, workspaceName, workspaceId)
	reqBodyData := strings.NewReader(data)
	req, err := http.NewRequest("POST", "https://a.blazemeter.com/api/v4/projects", reqBodyData)
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
