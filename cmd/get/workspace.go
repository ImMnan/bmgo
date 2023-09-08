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

	"github.com/spf13/cobra"
)

// workspaceCmd represents the workspace pallete command
var workspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Get details about the workspace, use with other sub-commands to get specific/detailed info",
	Long:  `.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Getting workspace details...")
		workspaceId, _ := cmd.Flags().GetInt("id")
		rawOutput, _ := cmd.Flags().GetBool("raw")
		if rawOutput {
			getWorkspaceRaw(workspaceId)
		} else {
			getWorkspace(workspaceId)
		}
	},
}

type responseBodyWS struct {
	Result resultWS
}

type resultWS struct {
	Name         string `json:"name"`
	Enabled      bool   `json:"enabled"`
	MembersCount int    `json:"membersCount"`
	AccountId    int    `json:"accountId"`
}

func getWorkspace(workspaceId int) {
	apiId, apiSecret := Getapikeys()

	workspaceIdStr := strconv.Itoa(workspaceId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/workspaces/"+workspaceIdStr, nil)
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
	var responseObjectWS responseBodyWS
	json.Unmarshal(bodyText, &responseObjectWS)

	workspaceName := responseObjectWS.Result.Name
	members := responseObjectWS.Result.MembersCount
	accountId := responseObjectWS.Result.AccountId
	enabled := responseObjectWS.Result.Enabled

	fmt.Printf("\n%-20s %-10s %-10s %-10s\n", "NAME", "ACCOUNT", "MEMBERS", "ENABLED")
	fmt.Printf("%-20s %-10d %-10d %-10t\n\n", workspaceName, accountId, members, enabled)
}

func getWorkspaceRaw(workspaceId int) {
	apiId, apiSecret := Getapikeys()

	workspaceIdStr := strconv.Itoa(workspaceId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/workspaces/"+workspaceIdStr, nil)
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

func init() {
	GetCmd.AddCommand(workspaceCmd)
	workspaceCmd.PersistentFlags().Int("id", 0, "Confirm the workspace id")
	workspaceCmd.MarkPersistentFlagRequired("id")
	workspaceCmd.Flags().BoolP("raw", "r", false, "[Optional] If set, the output will be raw json")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// accountCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// accountCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
