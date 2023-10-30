/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package update

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Update users in Account or Workspace",
	Long: `Use the command to update user entry, we can either enable or disable the user or make changes to the user-roles within a workspace or account level. To update the user, you will need the user Id & use the flag --accountid or --workspaceid to make changes to specific level.

	For example: [bmgo update -a <account_id> user --uid <user Id>] OR 
                 [bmgo update -w <workspace_id> user --uid <user Id>] 
	For default: [bmgo update --ac user --uid <user Id>] OR 
                 [bmgo update --ws user --uid <user Id>] `,
	Run: func(cmd *cobra.Command, args []string) {
		userId, _ := cmd.Flags().GetInt("uid")
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
		switch {
		case (workspaceId != 0) && (accountId == 0) && rawOutput:
			updateUserWsraw(userId, workspaceId)
		case (workspaceId == 0) && (accountId != 0) && rawOutput:
			updateUserAraw(userId, accountId)
		case (workspaceId != 0) && (accountId == 0) && !rawOutput:
			updateUserWs(userId, workspaceId)
		case (workspaceId == 0) && (accountId != 0) && !rawOutput:
			updateUserA(userId, accountId)
		default:
			cmd.Help()
		}
	},
}

func init() {
	UpdateCmd.AddCommand(userCmd)
	userCmd.Flags().Int("uid", 0, "Enter the User ID")
	userCmd.MarkFlagRequired("uid")
}

type updateUserResponse struct {
	Result updateUserResult `json:"result"`
	Error  errorResult      `json:"error"`
}
type updateUserResult struct {
	Email   string   `json:"email"`
	Enabled bool     `json:"enabled"`
	Type    string   `json:"type"`
	Roles   []string `json:"roles"`
}

func updateUserA(userId, accountId int) {
	apiId, apiSecret := Getapikeys()
	roleA := updateUserRolesA()
	enableA := isEnabledPromt()
	accountIdStr := strconv.Itoa(accountId)
	client := &http.Client{}
	userIdStr := strconv.Itoa(userId)
	// var data = strings.NewReader(`{"roles":["user_manager"],"enabled": false}`)
	data := fmt.Sprintf(`{"roles": ["%s"], "enabled": %t}`, roleA, enableA)
	var reqBodyDataA = strings.NewReader(data)
	req, err := http.NewRequest("PUT", "https://a.blazemeter.com/api/v4/accounts/"+accountIdStr+"/users/"+userIdStr, reqBodyDataA)
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
	var responseBodyUpdateUserA updateUserResponse
	var userRoleA string
	json.Unmarshal(bodyText, &responseBodyUpdateUserA)
	if responseBodyUpdateUserA.Error.Code == 0 {

		userEmailA := responseBodyUpdateUserA.Result.Email
		userTypeA := responseBodyUpdateUserA.Result.Type
		userStatusA := responseBodyUpdateUserA.Result.Enabled
		fmt.Printf("\n%-25s %-12s %-10s %-10s", "EMAIL", "TYPE", "ENABLE", "ROLE")
		for i := 0; i < len(responseBodyUpdateUserA.Result.Roles); i++ {
			userRoleA = responseBodyUpdateUserA.Result.Roles[i]
		}
		fmt.Printf("\n%-25s %-12s %-10t %-10s\n\n", userEmailA, userTypeA, userStatusA, userRoleA)
	} else {
		errorCode := responseBodyUpdateUserA.Error.Code
		errorMessage := responseBodyUpdateUserA.Error.Message
		fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
	}
}
func updateUserAraw(userId, accountId int) {
	apiId, apiSecret := Getapikeys()
	roleA := updateUserRolesA()
	enableA := isEnabledPromt()
	accountIdStr := strconv.Itoa(accountId)
	client := &http.Client{}
	userIdStr := strconv.Itoa(userId)
	// var data = strings.NewReader(`{"roles":["user_manager"],"enabled": false}`)
	data := fmt.Sprintf(`{"roles": ["%s"], "enabled": %t}`, roleA, enableA)
	var reqBodyDataA = strings.NewReader(data)
	req, err := http.NewRequest("PUT", "https://a.blazemeter.com/api/v4/accounts/"+accountIdStr+"/users/"+userIdStr, reqBodyDataA)
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

func updateUserWs(userId, workspaceId int) {
	workspaceIdStr := strconv.Itoa(workspaceId)
	apiId, apiSecret := Getapikeys()
	roleWs := updateUserRolesWs()
	enableWs := isEnabledPromt()
	userIdStr := strconv.Itoa(userId)
	client := &http.Client{}
	data := fmt.Sprintf(`{"roles":["%s"],"enabled": %t}`, roleWs, enableWs)
	var reqBodyDataWs = strings.NewReader(data)
	req, err := http.NewRequest("PUT", "https://a.blazemeter.com/api/v4/workspaces/"+workspaceIdStr+"/users/"+userIdStr, reqBodyDataWs)
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
	var responseBodyUpdateUserWs updateUserResponse
	var userRoleWs string
	json.Unmarshal(bodyText, &responseBodyUpdateUserWs)
	if responseBodyUpdateUserWs.Error.Code == 0 {

		userEmailWs := responseBodyUpdateUserWs.Result.Email
		userTypeWs := responseBodyUpdateUserWs.Result.Type
		userStatusWs := responseBodyUpdateUserWs.Result.Enabled
		fmt.Printf("\n%-25s %-12s %-10s %-10s", "EMAIL", "TYPE", "ENABLE", "ROLE")
		for i := 0; i < len(responseBodyUpdateUserWs.Result.Roles); i++ {
			userRoleWs = responseBodyUpdateUserWs.Result.Roles[i]
		}
		fmt.Printf("\n%-25s %-12s %-10t %-10s\n\n", userEmailWs, userTypeWs, userStatusWs, userRoleWs)
	} else {
		errorCode := responseBodyUpdateUserWs.Error.Code
		errorMessage := responseBodyUpdateUserWs.Error.Message
		fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
	}
}

func updateUserWsraw(userId, workspaceId int) {
	workspaceIdStr := strconv.Itoa(workspaceId)
	apiId, apiSecret := Getapikeys()
	roleWs := updateUserRolesWs()
	enableWs := isEnabledPromt()
	userIdStr := strconv.Itoa(userId)
	client := &http.Client{}
	data := fmt.Sprintf(`{"roles":["%s"],"enabled": %t}`, roleWs, enableWs)
	var reqBodyDataWs = strings.NewReader(data)
	req, err := http.NewRequest("PUT", "https://a.blazemeter.com/api/v4/workspaces/"+workspaceIdStr+"/users/"+userIdStr, reqBodyDataWs)
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
