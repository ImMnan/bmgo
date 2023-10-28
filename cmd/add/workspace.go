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

// workspaceCmd represents the workspace command
var workspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Add workspace to an account",
	Long: `Add a new workspace into your existing account using this command, just specify the name of the workspace. Workspace Id for the newly created workspace is returned in the output.
	
	For example: [bmgo add -a <account_id> workspace --name <workspace name>]
	For default: [bmgo add --ac workspace --name <workspace name>]`,
	Run: func(cmd *cobra.Command, args []string) {
		ac, _ := cmd.Flags().GetBool("ac")
		var accountId int
		if ac {
			accountId = defaultAccount()
		} else {
			accountId, _ = cmd.Flags().GetInt("accountid")
		}
		workspaceName, _ := cmd.Flags().GetString("name")
		rawOutput, _ := cmd.Flags().GetBool("raw")
		if accountId != 0 && rawOutput {
			addWorkspaceraw(workspaceName, accountId)
		} else if accountId != 0 {
			addWorkspace(workspaceName, accountId)
		} else {
			cmd.Help()
		}
	},
}

func init() {
	AddCmd.AddCommand(workspaceCmd)
	workspaceCmd.Flags().String("name", "", "Name your workspace")
	workspaceCmd.MarkFlagRequired("name")
}

type addWorkspaceResponse struct {
	Result addWorkspaceResult `json:"result"`
	Error  errorResult        `json:"error"`
}
type addWorkspaceResult struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

func addWorkspace(workspaceName string, accountId int) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	data := fmt.Sprintf(`{"dedicatedIpsEnabled": true, "enabled": true, "name": "%s", "privateLocationsEnabled": true, "accountId": %v}`, workspaceName, accountId)
	reqBodydata := strings.NewReader(data)
	req, err := http.NewRequest("POST", "https://a.blazemeter.com/api/v4/workspaces", reqBodydata)
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
	var responseBodyAddWorkspace addWorkspaceResponse
	json.Unmarshal(bodyText, &responseBodyAddWorkspace)
	if responseBodyAddWorkspace.Error.Code == 0 {
		wsName := responseBodyAddWorkspace.Result.Name
		wsEnabled := responseBodyAddWorkspace.Result.Enabled
		wsId := responseBodyAddWorkspace.Result.Id

		fmt.Printf("\n%-10s %-20s %-10s", "ID", "NAME", "ENABLED")
		fmt.Printf("\n%-10v %-20s %-10t\n\n", wsId, wsName, wsEnabled)
	} else {
		errorCode := responseBodyAddWorkspace.Error.Code
		errorMessage := responseBodyAddWorkspace.Error.Message
		fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
	}
}
func addWorkspaceraw(workspaceName string, accountId int) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	data := fmt.Sprintf(`{"dedicatedIpsEnabled": true, "enabled": true, 
	"name": "%s", "privateLocationsEnabled": true, "accountId": %v}`, workspaceName, accountId)
	reqBodydata := strings.NewReader(data)
	req, err := http.NewRequest("POST", "https://a.blazemeter.com/api/v4/workspaces", reqBodydata)
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
