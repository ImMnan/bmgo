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
	Long: `Use the command to list Users within a specified account, team (API monitoring), or a workspace. The output includes User ID, Name, Roles, Email, etc. The output can be further filtered by switching disabled flag as true to only display disabled users, --disabled.
	
	For example: [bmgo get -w <workspace id> users] OR
	             [bmgo get -a <account id> users] OR
		     [bmgo get -t <team id> users] OR
	             [bmgo get -w <workspace id> users --disabled] OR
	             [bmgo get -a <account id> users --disabled]

    For default: [bmgo get --ws users] OR
	             [bmgo get --ac users] OR 
	             [bmgo get --tm users] OR
	             [bmgo get --ws users --disabled] OR
	             [bmgo get --ac users --disabled]`,
	Run: func(cmd *cobra.Command, args []string) {
		ac, _ := cmd.Flags().GetBool("ac")
		ws, _ := cmd.Flags().GetBool("ws")
		tm, _ := cmd.Flags().GetBool("tm")
		var accountId, workspaceId int
		var teamId string
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
		if tm {
			teamId = defaultTeam()
		} else {
			teamId, _ = cmd.Flags().GetString("teamid")
		}
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
		} else if (accountId == 0) && (workspaceId == 0) && (teamId != "") && rawOutput {
			getUsersTmraw(teamId)
		} else if (accountId == 0) && (workspaceId == 0) && (teamId != "") && !rawOutput {
			getUsersTm(teamId)
		} else if (workspaceId != 0) && (accountId == 0) && disabledUsers {
			getUsersWSDis(workspaceId)
		} else if (accountId != 0) && (workspaceId == 0) && disabledUsers {
			getUsersADis(accountId)
		} else if (workspaceId != 0) && (accountId == 0) {
			getUsersWS(workspaceId)
		} else if (accountId != 0) && (workspaceId == 0) {
			getUsersA(accountId)
		} else {
			cmd.Help()
		}
	},
}

func init() {
	GetCmd.AddCommand(usersCmd)
	usersCmd.Flags().Bool("disabled", false, "[Optional] will show enabled users only")
}

type usersResponse struct {
	Result []usersResult `json:"result"`
	Data   []usersData   `json:"data"`
	Error  errorResult   `json:"error"`
}

type usersResult struct {
	Id          int      `json:"id"`
	Email       string   `json:"email"`
	DisplayName string   `json:"displayName"`
	Enabled     bool     `json:"enabled"`
	Roles       []string `json:"roles"`
}
type usersData struct {
	Uuid       string `json:"uuid"`
	Email      string `json:"email"`
	Role_name  string `json:"role_name"`
	Created_at string `json:"created_at"`
	Name       string `json:"name"`
}

func getUsersA(accountId int) {
	apiId, apiSecret := Getapikeys()

	client := &http.Client{}
	accountIdStr := strconv.Itoa(accountId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/accounts/"+accountIdStr+"/users?limit=-1&enabled=true", nil)
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
	if responseBodyAUsers.Error.Code == 0 {
		fmt.Printf("\n%-10s %-25s %-12s %-10s %-10s\n", "ID", "DISPLAY NAME", "ROLES", "ENABLED", "EMAIL")
		//	rolesListTotal := []string{}
		for i := 0; i < len(responseBodyAUsers.Result); i++ {
			userIdWS := responseBodyAUsers.Result[i].Id
			displayNameWS := responseBodyAUsers.Result[i].DisplayName
			emailIdWS := responseBodyAUsers.Result[i].Email
			enabledUserWS := responseBodyAUsers.Result[i].Enabled
			fmt.Printf("\n%-10v %-25s %-12s %-10t %-10s", userIdWS, displayNameWS, responseBodyAUsers.Result[i].Roles[0], enabledUserWS, emailIdWS)
		}
		fmt.Println("\n-")
	} else {
		errorCode := responseBodyAUsers.Error.Code
		errorMessage := responseBodyAUsers.Error.Message
		fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
	}
}
func getUsersAraw(accountId int) {
	apiId, apiSecret := Getapikeys()

	client := &http.Client{}
	accountIdStr := strconv.Itoa(accountId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/accounts/"+accountIdStr+"/users?limit=-1&enabled=true", nil)
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
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/accounts/"+accountIdStr+"/users?limit=1000&enabled=false", nil)
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
	if responseBodyAUsers.Error.Code == 0 {
		fmt.Printf("\n%-10s %-25s %-12s %-10s %-10s\n", "ID", "DISPLAY NAME", "ROLES", "ENABLED", "EMAIL")
		//	rolesListTotal := []string{}
		for i := 0; i < len(responseBodyAUsers.Result); i++ {
			userIdWS := responseBodyAUsers.Result[i].Id
			displayNameWS := responseBodyAUsers.Result[i].DisplayName
			emailIdWS := responseBodyAUsers.Result[i].Email
			enabledUserWS := responseBodyAUsers.Result[i].Enabled
			fmt.Printf("\n%-10v %-25s %-12s %-10t %-10s", userIdWS, displayNameWS, responseBodyAUsers.Result[i].Roles[0], enabledUserWS, emailIdWS)
		}
		fmt.Println("\n-")
	} else {
		errorCode := responseBodyAUsers.Error.Code
		errorMessage := responseBodyAUsers.Error.Message
		fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
	}
}
func getUsersArawDis(accountId int) {
	apiId, apiSecret := Getapikeys()

	client := &http.Client{}
	accountIdStr := strconv.Itoa(accountId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/accounts/"+accountIdStr+"/users?limit=1000&enabled=false", nil)
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
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/workspaces/"+workspaceIdStr+"/users?limit=-1&enabled=true", nil)
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
	if responseBodyWsUsers.Error.Code == 0 {
		fmt.Printf("\n%-10s %-25s %-12s %-10s %-10s\n", "ID", "DISPLAY NAME", "ROLES", "ENABLED", "EMAIL")
		//	rolesListTotal := []string{}
		for i := 0; i < len(responseBodyWsUsers.Result); i++ {
			userIdWS := responseBodyWsUsers.Result[i].Id
			displayNameWS := responseBodyWsUsers.Result[i].DisplayName
			emailIdWS := responseBodyWsUsers.Result[i].Email
			enabledUserWS := responseBodyWsUsers.Result[i].Enabled
			fmt.Printf("\n%-10v %-25s %-12s %-10t %-10s", userIdWS, displayNameWS, responseBodyWsUsers.Result[i].Roles[0], enabledUserWS, emailIdWS)
		}
		fmt.Println("\n-")
	} else {
		errorCode := responseBodyWsUsers.Error.Code
		errorMessage := responseBodyWsUsers.Error.Message
		fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
	}
}
func getUsersWSraw(workspaceId int) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	workspaceIdStr := strconv.Itoa(workspaceId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/workspaces/"+workspaceIdStr+"/users?limit=-1&enabled=true", nil)
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
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/workspaces/"+workspaceIdStr+"/users?limit=-500&enabled=false", nil)
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
	if responseBodyWsUsers.Error.Code == 0 {
		fmt.Printf("\n%-10s %-25s %-12s %-10s %-10s\n", "ID", "DISPLAY NAME", "ROLES", "ENABLED", "EMAIL")
		for i := 0; i < len(responseBodyWsUsers.Result); i++ {
			userIdWS := responseBodyWsUsers.Result[i].Id
			displayNameWS := responseBodyWsUsers.Result[i].DisplayName
			emailIdWS := responseBodyWsUsers.Result[i].Email
			enabledUserWS := responseBodyWsUsers.Result[i].Enabled
			fmt.Printf("\n%-10v %-25s %-12s %-10t %-10s", userIdWS, displayNameWS, responseBodyWsUsers.Result[i].Roles[0], enabledUserWS, emailIdWS)
		}
		fmt.Println("\n-")
	} else {
		errorCode := responseBodyWsUsers.Error.Code
		errorMessage := responseBodyWsUsers.Error.Message
		fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
	}
}
func getUsersWSrawDis(workspaceId int) {
	apiId, apiSecret := Getapikeys()

	client := &http.Client{}
	workspaceIdStr := strconv.Itoa(workspaceId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/workspaces/"+workspaceIdStr+"/users?limit=500&enabled=false", nil)
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

func getUsersTm(teamId string) {
	Bearer := fmt.Sprintf("Bearer %v", GetPersonalAccessToken())
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.runscope.com/teams/"+teamId+"/people", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", Bearer)
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
	var responseBodyTmUsers usersResponse
	json.Unmarshal(bodyText, &responseBodyTmUsers)
	if responseBodyTmUsers.Error.Status == 0 {
		fmt.Printf("\n%-38s %-14s %-28s %-10s\n", "UUID", "ROLES", "NAME", "EMAIL")
		for i := 0; i < len(responseBodyTmUsers.Data); i++ {
			userIdTm := responseBodyTmUsers.Data[i].Uuid
			userNameTm := responseBodyTmUsers.Data[i].Name
			userEmailTm := responseBodyTmUsers.Data[i].Email
			userRoleTm := responseBodyTmUsers.Data[i].Role_name
			//	userCreatedTm := responseBodyTmUsers.Data[i].Created_at
			fmt.Printf("\n%-38s %-14s %-28s %-10s", userIdTm, userRoleTm, userNameTm, userEmailTm)
		}
		fmt.Println("\n-")
	} else {
		errorCode := responseBodyTmUsers.Error.Status
		errorMessage := responseBodyTmUsers.Error.Message
		fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
	}
}

func getUsersTmraw(teamId string) {
	Bearer := fmt.Sprintf("Bearer %v", GetPersonalAccessToken())
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.runscope.com/teams/"+teamId+"/people", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", Bearer)
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
