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

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Update users in Account or Workspace",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		userId, _ := cmd.Flags().GetInt("uid")
		accountId, _ := cmd.Flags().GetInt("accountid")
		workspaceId, _ := cmd.Flags().GetInt("workspaceid")
		rawOutput, _ := cmd.Flags().GetBool("raw")
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
			fmt.Println("\nPlease provide a correct workspace Id or Account Id to get the info")
			fmt.Println("[bmgo update -a <account_id>...] OR [bmgo update -w <workspace_id>...]")
		}
	},
}

func init() {
	UpdateCmd.AddCommand(userCmd)
	userCmd.Flags().Int("uid", 0, "Enter the User ID")
	userCmd.MarkFlagRequired("uid")
}

func updateUserSelectorA() (string, bool) {
	prompt := promptui.Select{
		Label:        "Select Account Role:",
		Items:        []string{"admin", "standard", "user_manager", "billing"},
		HideSelected: true,
	}
	prompt1 := promptui.Select{
		Label:        "Enable:",
		Items:        []bool{true, false},
		HideSelected: true,
	}
	_, roleSelected, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}
	_, attachAuto, err := prompt1.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}
	boolVal, _ := strconv.ParseBool(attachAuto)
	return roleSelected, boolVal
}
func updateUserSelectorWs() (string, bool) {
	prompt := promptui.Select{
		Label:        "Select Account Role:",
		Items:        []string{"tester", "manager", "viewer"},
		HideSelected: true,
	}
	prompt1 := promptui.Select{
		Label:        "Enable:",
		Items:        []bool{true, false},
		HideSelected: true,
	}
	_, roleSelected, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}
	_, attachAuto, err := prompt1.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}
	boolVal, _ := strconv.ParseBool(attachAuto)
	return roleSelected, boolVal
}

type updateUserResponse struct {
	Result updateUserResult `json:"result"`
}
type updateUserResult struct {
	Email   string   `json:"email"`
	Enabled bool     `json:"enabled"`
	Type    string   `json:"type"`
	Roles   []string `json:"roles"`
}

func updateUserA(userId, accountId int) {
	roleA, enableA := updateUserSelectorA()
	apiId, apiSecret := Getapikeys()
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
	userEmailA := responseBodyUpdateUserA.Result.Email
	userTypeA := responseBodyUpdateUserA.Result.Type
	userStatusA := responseBodyUpdateUserA.Result.Enabled
	fmt.Printf("\n%-25s %-12s %-10s %-10s", "EMAIL", "TYPE", "ENABLE", "ROLE")
	for i := 0; i < len(responseBodyUpdateUserA.Result.Roles); i++ {
		userRoleA = responseBodyUpdateUserA.Result.Roles[i]
	}
	fmt.Printf("\n%-25s %-12s %-10t %-10s\n\n", userEmailA, userTypeA, userStatusA, userRoleA)
}
func updateUserAraw(userId, accountId int) {
	roleA, enableA := updateUserSelectorA()
	apiId, apiSecret := Getapikeys()
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
	roleWs, enableWs := updateUserSelectorWs()
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
	userEmailWs := responseBodyUpdateUserWs.Result.Email
	userTypeWs := responseBodyUpdateUserWs.Result.Type
	userStatusWs := responseBodyUpdateUserWs.Result.Enabled
	fmt.Printf("\n%-25s %-12s %-10s %-10s", "EMAIL", "TYPE", "ENABLE", "ROLE")
	for i := 0; i < len(responseBodyUpdateUserWs.Result.Roles); i++ {
		userRoleWs = responseBodyUpdateUserWs.Result.Roles[i]
	}
	fmt.Printf("\n%-25s %-12s %-10t %-10s\n\n", userEmailWs, userTypeWs, userStatusWs, userRoleWs)
}
func updateUserWsraw(userId, workspaceId int) {
	workspaceIdStr := strconv.Itoa(workspaceId)
	apiId, apiSecret := Getapikeys()
	roleWs, enableWs := updateUserSelectorWs()
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
