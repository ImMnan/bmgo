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
	Short: "Get details about the account or workspace",
	Long:  `.`,
	Run: func(cmd *cobra.Command, args []string) {
		rawOutput, _ := cmd.Flags().GetBool("raw")
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
	Error  errorResult `json:"error"`
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
	if responseObject.Error.Code == 0 {
		accountName := responseObject.Result.Name
		ownerEmail := responseObject.Result.Owner.Email
		workspaceCount := responseObject.Result.WorkspacesCount
		memberCount := responseObject.Result.MembersCount

		accountPlanId := responseObject.Result.Plan.Id
		//	accountPlanName := responseObject.Result.Plan.Name
		accountReportRet := responseObject.Result.Plan.ReportRetention
		accountThreadsPE := responseObject.Result.Plan.ThreadsPerEngine
		accountCredits := responseObject.Result.Credits
		accountExpiration := int64(responseObject.Result.Expiration)
		mytimeExpiration := time.Unix(accountExpiration, 0)
		expirationTimeStr := fmt.Sprint(mytimeExpiration)
		fmt.Printf("\n%-25s %-35s %-15s %-10s\n", "NAME", "OWNER", "WORKSPACES", "USERS")
		fmt.Printf("%-25s %-35s %-15d %-10d\n", accountName, ownerEmail, workspaceCount, memberCount)

		fmt.Printf("\n------------------------------------------------------------------------------------------------------------")

		fmt.Printf("\n%-35s %-10s %-10s %-10s %-20s\n", "PLAN ID", "CREDITS", "REP RET.", "TPE", "EXPIRATION")
		fmt.Printf("%-35s %-10v %-10d %-10d %-20v\n", accountPlanId, accountCredits, accountReportRet, accountThreadsPE, expirationTimeStr[0:16])

		cloudProviders := []string{}
		for i := 0; i < len(responseObject.Result.CloudProviders); i++ {
			cloudProlist := responseObject.Result.CloudProviders[i]
			cloudProviders = append(cloudProviders, cloudProlist)
		}
		fmt.Printf("\n------------------------------------------------------------------------------------------------------------")
		fmt.Printf("\nSupported cloud providers: %v", cloudProviders)
	} else {
		errorCode := responseObject.Error.Code
		errorMessage := responseObject.Error.Message
		fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
	}
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
	Error  errorResult `json:"error"`
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
	if responseObjectWS.Error.Code == 0 {
		workspaceName := responseObjectWS.Result.Name
		members := responseObjectWS.Result.MembersCount
		accountId := responseObjectWS.Result.AccountId
		enabled := responseObjectWS.Result.Enabled

		fmt.Printf("\n%-30s %-10s %-10s %-10s\n", "NAME", "ACCOUNT", "MEMBERS", "ENABLED")
		fmt.Printf("%-30s %-10d %-10d %-10t\n", workspaceName, accountId, members, enabled)
		fmt.Printf("\n------------------------------------------------------------------------------------------------------------")
	} else {
		errorCode := responseObjectWS.Error.Code
		errorMessage := responseObjectWS.Error.Message
		fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
	}
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
