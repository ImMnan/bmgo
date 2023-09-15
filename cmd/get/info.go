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

// accountCmd represents the account command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get details about the account, use with other sub-commands to get specific/detailed info",
	Long:  `.`,
	Run: func(cmd *cobra.Command, args []string) {
		accountId, _ := cmd.Flags().GetInt("accountid")
		workspaceId, _ := cmd.Flags().GetInt("workspaceid")
		rawOutput, _ := cmd.Flags().GetBool("raw")

		if (workspaceId != 0) && (accountId == 0) && rawOutput {
			getWorkspaceRaw(workspaceId)
		} else if (accountId != 0) && (workspaceId == 0) && rawOutput {
			getAccountIdRaw(accountId)
		} else if (workspaceId != 0) && (accountId == 0) {
			getWorkspace(workspaceId)
		} else if (accountId != 0) && (workspaceId == 0) {
			getAccountId(accountId)
		} else {
			fmt.Println("\nPlease provide a correct workspace Id or Account Id to get the info")
			fmt.Println("[bmgo get -a <account_id>...] OR [bmgo get -w <workspace_id>...]")
		}
	},
}

func init() {
	GetCmd.AddCommand(infoCmd)
}

type inforesponseA struct {
	Result result
}
type result struct {
	Name            string `json:"name"`
	Owner           owner
	MembersCount    int `json:"membersCount"`
	WorkspacesCount int `json:"workspacesCount"`
	Plan            plan
	CloudProviders  []string
	Credits         int `json:"credits"`
	Expiration      int `json:"expiration"`
}

type owner struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type plan struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	ReportRetention  int    `json:"reportRetention"`
	ThreadsPerEngine int    `json:"threadsPerEngine"`
}

func getAccountId(accountId int) {
	apiId, apiSecret := Getapikeys()

	accountIdStr := strconv.Itoa(accountId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/accounts/"+accountIdStr, nil)
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
	var responseObject inforesponseA
	json.Unmarshal(bodyText, &responseObject)

	accountName := responseObject.Result.Name
	ownerEmail := responseObject.Result.Owner.Email
	workspaceCount := responseObject.Result.WorkspacesCount
	memberCount := responseObject.Result.MembersCount

	accountPlanId := responseObject.Result.Plan.Id
	accountPlanName := responseObject.Result.Plan.Name
	accountReportRet := responseObject.Result.Plan.ReportRetention
	accountThreadsPE := responseObject.Result.Plan.ThreadsPerEngine
	accountCredits := responseObject.Result.Credits
	accountExpiration := int64(responseObject.Result.Expiration)
	mytimeExpiration := time.Unix(accountExpiration, 0)
	fmt.Printf("\n%-20s %-30s %-10s %-10s %-20s\n", "NAME", "OWNER", "WORKSPACES", "USERS", "PLAN NAME")
	fmt.Printf("%-20s %-30s %-10d %-10d %-20s\n", accountName, ownerEmail, workspaceCount, memberCount, accountPlanName)

	fmt.Printf("\n------------------------------------------------------------------------------------------------------------")

	fmt.Printf("\n%-20s %-10s %-10s %-10s %-20s\n", "PLAN ID", "CREDITS", "REP RET.", "TPE", "EXPIRATION")
	fmt.Printf("%-20s %-10v %-10d %-10d %-20v\n", accountPlanId, accountCredits, accountReportRet, accountThreadsPE, mytimeExpiration)

	cloudProviders := []string{}
	for i := 0; i < len(responseObject.Result.CloudProviders); i++ {
		cloudProlist := responseObject.Result.CloudProviders[i]
		cloudProviders = append(cloudProviders, cloudProlist)
	}
	fmt.Printf("\n------------------------------------------------------------------------------------------------------------")
	fmt.Printf("\nSupported cloud providers: %v \n\n", cloudProviders)
}

func getAccountIdRaw(accountId int) {
	apiId, apiSecret := Getapikeys()

	accountIdStr := strconv.Itoa(accountId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/accounts/"+accountIdStr, nil)
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

type inforesponseBodyWS struct {
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
	var responseObjectWS inforesponseBodyWS
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
