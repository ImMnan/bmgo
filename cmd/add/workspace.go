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
	Short: "Add workspace to account",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		accountId, _ := cmd.Flags().GetInt("accountid")
		workspaceName, _ := cmd.Flags().GetString("name")
		rawOutput, _ := cmd.Flags().GetBool("raw")
		if accountId != 0 && rawOutput {
			addWorkspaceraw(workspaceName, accountId)
		} else if accountId != 0 {
			addWorkspace(workspaceName, accountId)
		} else {
			fmt.Println("\nPlease provide a correct Account Id to add workspace to")
			fmt.Println("[bmgo add -a <account_id>...]")
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
	wsName := responseBodyAddWorkspace.Result.Name
	wsEnabled := responseBodyAddWorkspace.Result.Enabled
	wsId := responseBodyAddWorkspace.Result.Id

	fmt.Printf("\n%-10s %-20s %-10s", "ID", "NAME", "ENABLED")
	fmt.Printf("\n%-10v %-20s %-10t\n\n", wsId, wsName, wsEnabled)
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