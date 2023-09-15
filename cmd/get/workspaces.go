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

// workspacesCmd represents the workspaces command
var workspacesCmd = &cobra.Command{
	Use:   "workspaces",
	Short: "Get a list of workspaces in the account",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("workspaces called")
		accountId, _ := cmd.Flags().GetInt("accountid")
		rawOutput, _ := cmd.Flags().GetBool("raw")
		if accountId != 0 && rawOutput {
			getWorkspacesraw(accountId)
		} else if accountId != 0 {
			getWorkspaces(accountId)
		} else {
			fmt.Println("\nPlease provide a correct workspace Id")
			fmt.Println("[bmgo get -w <workspace_id>...]")
			cmd.Help()
		}
	},
}

func init() {
	GetCmd.AddCommand(workspacesCmd)
}

type workspacesResponse struct {
	Result []wsResult `json:"result"`
}
type wsResult struct {
	Name         string `json:"name"`
	Enabled      bool   `json:"enabled"`
	MembersCount int    `json:"membersCount"`
	AccountId    int    `json:"accountId"`
	Created      int    `json:"created"`
}

func getWorkspaces(accountId int) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	accountIdStr := strconv.Itoa(accountId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/workspaces?accountId="+accountIdStr+"&limit=200", nil)
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
	var responseObjectWS workspacesResponse
	json.Unmarshal(bodyText, &responseObjectWS)
	fmt.Printf("\n%-30s %-10s %-10s %-30s\n", "NAME", "MEMBERS", "ENABLED", "CREATED")
	for i := 0; i < len(responseObjectWS.Result); i++ {
		workspaceName := responseObjectWS.Result[i].Name
		members := responseObjectWS.Result[i].MembersCount
		created := int64(responseObjectWS.Result[i].Created)
		enabled := responseObjectWS.Result[i].Enabled
		createdepoch := time.Unix(created, 0)
		fmt.Printf("\n%-30s %-10d %-10t %-30v", workspaceName, members, enabled, createdepoch)
	}
	fmt.Println("\n-")
}

func getWorkspacesraw(accountId int) {
	apiId, apiSecret := Getapikeys()

	client := &http.Client{}
	accountIdStr := strconv.Itoa(accountId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/workspaces?accountId="+accountIdStr+"&limit=200", nil)
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
