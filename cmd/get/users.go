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

// usersCmd represents the users command
var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Get a list of users part of the account",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("users called")
		accountId, _ := cmd.Flags().GetInt("accountid")
		workspaceId, _ := cmd.Flags().GetInt("workspaceid")
		rawOutput, _ := cmd.Flags().GetBool("raw")
		disabledUsers, _ := cmd.Flags().GetBool("disabled")

		if (workspaceId != 0) && (accountId == 0) && rawOutput && disabledUsers {
			getUsersWSrawDis(workspaceId)
		} else if (accountId != 0) && (workspaceId == 0) && rawOutput && disabledUsers {
			getUsersArawDis(accountId)
		} else if (workspaceId != 0) && (accountId == 0) && rawOutput {
			getUsersWSraw(workspaceId)
		} else if (accountId != 0) && (workspaceId == 0) && rawOutput {
			getUsersAraw(accountId)
		} else if (workspaceId != 0) && (accountId == 0) && disabledUsers {
			getUsersWSDis(workspaceId)
		} else if (accountId != 0) && (workspaceId == 0) && disabledUsers {
			getUsersADis(accountId)
		} else if (workspaceId != 0) && (accountId == 0) {
			getUsersWS(workspaceId)
		} else if (accountId != 0) && (workspaceId == 0) {
			getUsersA(accountId)
		} else {
			fmt.Println("\nPlease provide a correct workspace Id or Account Id to get the info")
			fmt.Println("[bmgo get -a <account_id>...] OR [bmgo get -w <workspace_id>...]")
		}
	},
}

func init() {
	GetCmd.AddCommand(usersCmd)
	usersCmd.Flags().Bool("disabled", false, "[Optional] will show enabled users only")
}

type usersResponse struct {
	Result []usersResult `json:"result"`
}

type usersResult struct {
	Id          int      `json:"id"`
	Email       string   `json:"email"`
	DisplayName string   `json:"displayName"`
	Enabled     bool     `json:"enabled"`
	Roles       []string `json:"roles"`
}

func getUsersA(accountId int) {
	apiId, apiSecret := Getapikeys()

	client := &http.Client{}
	accountIdStr := strconv.Itoa(accountId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/accounts/"+accountIdStr+"/users?limit=1500&enabled=true", nil)
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
	var responseBodyAUsers usersResponse
	json.Unmarshal(bodyText, &responseBodyAUsers)
	fmt.Printf("\n%-10s %-25s %-30s %-12s %-10s\n", "ID", "DISPLAY NAME", "EMAIL", "ROLES", "ENABLED")
	//	rolesListTotal := []string{}
	for i := 0; i < len(responseBodyAUsers.Result); i++ {
		userIdWS := responseBodyAUsers.Result[i].Id
		displayNameWS := responseBodyAUsers.Result[i].DisplayName
		emailIdWS := responseBodyAUsers.Result[i].Email
		enabledUserWS := responseBodyAUsers.Result[i].Enabled
		fmt.Printf("\n%-10v %-25s %-30s %-12s %-10t", userIdWS, displayNameWS, emailIdWS, responseBodyAUsers.Result[i].Roles[0], enabledUserWS)
	}
	fmt.Println("\n-")
}
func getUsersAraw(accountId int) {
	apiId, apiSecret := Getapikeys()

	client := &http.Client{}
	accountIdStr := strconv.Itoa(accountId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/accounts/"+accountIdStr+"/users?limit=1500&enabled=true", nil)
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

func getUsersADis(accountId int) {
	apiId, apiSecret := Getapikeys()

	client := &http.Client{}
	accountIdStr := strconv.Itoa(accountId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/accounts/"+accountIdStr+"/users?limit=1500&enabled=false", nil)
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
	var responseBodyAUsers usersResponse
	json.Unmarshal(bodyText, &responseBodyAUsers)
	fmt.Printf("\n%-10s %-25s %-30s %-12s %-10s\n", "ID", "DISPLAY NAME", "EMAIL", "ROLES", "ENABLED")
	//	rolesListTotal := []string{}
	for i := 0; i < len(responseBodyAUsers.Result); i++ {
		userIdWS := responseBodyAUsers.Result[i].Id
		displayNameWS := responseBodyAUsers.Result[i].DisplayName
		emailIdWS := responseBodyAUsers.Result[i].Email
		enabledUserWS := responseBodyAUsers.Result[i].Enabled
		fmt.Printf("\n%-10v %-25s %-30s %-12s %-10t", userIdWS, displayNameWS, emailIdWS, responseBodyAUsers.Result[i].Roles[0], enabledUserWS)
	}
	fmt.Println("\n-")
}
func getUsersArawDis(accountId int) {
	apiId, apiSecret := Getapikeys()

	client := &http.Client{}
	accountIdStr := strconv.Itoa(accountId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/accounts/"+accountIdStr+"/users?limit=1500&enabled=false", nil)
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

func getUsersWS(workspaceId int) {
	apiId, apiSecret := Getapikeys()

	client := &http.Client{}
	workspaceIdStr := strconv.Itoa(workspaceId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/workspaces/"+workspaceIdStr+"/users?limit=1000&enabled=true", nil)
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

	var responseBodyWsUsers usersResponse
	json.Unmarshal(bodyText, &responseBodyWsUsers)
	fmt.Printf("\n%-10s %-25s %-30s %-12s %-10s\n", "ID", "DISPLAY NAME", "EMAIL", "ROLES", "ENABLED")
	//	rolesListTotal := []string{}
	for i := 0; i < len(responseBodyWsUsers.Result); i++ {
		userIdWS := responseBodyWsUsers.Result[i].Id
		displayNameWS := responseBodyWsUsers.Result[i].DisplayName
		emailIdWS := responseBodyWsUsers.Result[i].Email
		enabledUserWS := responseBodyWsUsers.Result[i].Enabled
		fmt.Printf("\n%-10v %-25s %-30s %-12s %-10t", userIdWS, displayNameWS, emailIdWS, responseBodyWsUsers.Result[i].Roles[0], enabledUserWS)
	}
	fmt.Println("\n-")
}

func getUsersWSraw(workspaceId int) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	workspaceIdStr := strconv.Itoa(workspaceId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/workspaces/"+workspaceIdStr+"/users?limit=1000&enabled=true", nil)
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

func getUsersWSDis(workspaceId int) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	workspaceIdStr := strconv.Itoa(workspaceId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/workspaces/"+workspaceIdStr+"/users?limit=1000&enabled=false", nil)
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
	var responseBodyWsUsers usersResponse
	json.Unmarshal(bodyText, &responseBodyWsUsers)
	fmt.Printf("\n%-10s %-25s %-30s %-12s %-10s\n", "ID", "DISPLAY NAME", "EMAIL", "ROLES", "ENABLED")
	for i := 0; i < len(responseBodyWsUsers.Result); i++ {
		userIdWS := responseBodyWsUsers.Result[i].Id
		displayNameWS := responseBodyWsUsers.Result[i].DisplayName
		emailIdWS := responseBodyWsUsers.Result[i].Email
		enabledUserWS := responseBodyWsUsers.Result[i].Enabled
		fmt.Printf("\n%-10v %-25s %30s %-12s %-10t", userIdWS, displayNameWS, emailIdWS, responseBodyWsUsers.Result[i].Roles[0], enabledUserWS)
	}
	fmt.Println("\n-")
}

func getUsersWSrawDis(workspaceId int) {
	apiId, apiSecret := Getapikeys()

	client := &http.Client{}
	workspaceIdStr := strconv.Itoa(workspaceId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/workspaces/"+workspaceIdStr+"/users?limit=1000&enabled=false", nil)
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
