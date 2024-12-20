/*
Copyright © 2024 Manan Patel - github.com/immnan
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
	Short: "Add Project into workspace",
	Long: `Projects are designed to organize tests and reports and track usage within a Workspace. A Test and all of its Reports belong to a Project. You can add a project using this command to an existing workspace.
	
	For example: [bmgo add -w <workspace id> project --name <project name>]
	For default: [bmgo add --ws project --name <project name>]`,
	Run: func(cmd *cobra.Command, args []string) {
		ws, _ := cmd.Flags().GetBool("ws")
		var workspaceId int
		if ws {
			workspaceId = defaultWorkspace()
		} else {
			workspaceId, _ = cmd.Flags().GetInt("workspaceid")
		}
		projectName, _ := cmd.Flags().GetString("name")
		rawOutput, _ := cmd.Flags().GetBool("raw")
		switch {
		case (workspaceId != 0) && (projectName != ""):
			addProject(projectName, workspaceId, rawOutput)
		default:
			cmd.Help()
		}
	},
}

func init() {
	AddCmd.AddCommand(projectCmd)
	projectCmd.Flags().String("name", "", "Name your Project")
	projectCmd.MarkFlagRequired("name")
	projectCmd.Flags().IntP("workspaceid", "w", 0, " Provide Workspace ID to add a resource to")
	projectCmd.Flags().Bool("ws", false, "Use default workspace Id (bmConfig.yaml)")
}

type addProjectResponse struct {
	Result addProjectResult `json:"result"`
	Error  errorResult      `json:"error"`
}
type addProjectResult struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func addProject(projectName string, workspaceId int, rawOutput bool) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	data := fmt.Sprintf(`{
		"name": "%s", 
		"description": "Project created through bmgo", 
		"workspaceId": %v}`, projectName, workspaceId)
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
	if rawOutput {
		fmt.Printf("%s\n", bodyText)
	} else {
		var responseBodyAddProject addProjectResponse
		json.Unmarshal(bodyText, &responseBodyAddProject)
		if responseBodyAddProject.Error.Code == 0 {
			projectIdres := responseBodyAddProject.Result.Id
			projectNameres := responseBodyAddProject.Result.Name
			fmt.Printf("\n%-15s %-15s", "PROJECT ID", "NAME")
			fmt.Printf("\n%-15v %-15s", projectIdres, projectNameres)
			fmt.Println("\n-")
		} else {
			errorCode := responseBodyAddProject.Error.Code
			errorMessage := responseBodyAddProject.Error.Message
			fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
		}
	}
}
