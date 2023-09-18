/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package add

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Get details about the user",
	Long:  `.`,
	Run: func(cmd *cobra.Command, args []string) {
		//	fmt.Println("user add called")
		userId, _ := cmd.Flags().GetInt("uid")
		emailId, _ := cmd.Flags().GetString("email")
		workspaceId, _ := cmd.Flags().GetInt("workspaceid")
		accountId, _ := cmd.Flags().GetInt("accountid")
		rawOutput, _ := cmd.Flags().GetBool("raw")
		switch {
		case (workspaceId == 0) && (accountId != 0) && (emailId != "") && rawOutput:
			addUserByEmailAraw(emailId, accountId)
		case (workspaceId == 0) && (accountId != 0) && (emailId != "") && !rawOutput:
			addUserByEmailA(emailId, accountId)
		case (workspaceId != 0) && (accountId == 0) && rawOutput:
			addUserByUidWsraw(userId, workspaceId)
		case (workspaceId == 0) && (accountId != 0) && rawOutput:
			addUserByUidA(userId, accountId)
		case (workspaceId != 0) && (accountId == 0) && !rawOutput:
			addUserByUidWs(userId, workspaceId)
		case (workspaceId == 0) && (accountId != 0) && !rawOutput:
			addUserByUidA(userId, accountId)
		default:
			fmt.Println("\nPlease provide a correct workspace Id or Account Id to add user")
			fmt.Println("[bmgo add -a <account_id>...] OR [bmgo add -w <workspace_id>...]")
		}
	},
}

func init() {
	AddCmd.AddCommand(userCmd)
	userCmd.Flags().Int("uid", 0, "User ID for the user")
	userCmd.Flags().String("email", "", "Enter the Email ID of the user invited!")
}
func userRoleSelectorA() (string, bool) {
	prompt := promptui.Select{
		Label:        "Select Account Role",
		Items:        []string{"admin", "standard", "user_manager", "billing"},
		HideSelected: true,
	}
	prompt1 := promptui.Select{
		Label:        "attachAutomatically",
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
func userRoleSelectorWs() string {
	prompt := promptui.Select{
		Label:        "Select Workspace Role",
		Items:        []string{"tester", "manager", "viewer"},
		HideSelected: true,
	}
	_, roleSelected, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}
	return roleSelected
}
func workspaceIdPrompt() string {
	validate := func(input string) error {
		_, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return errors.New("invalid workspace")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:       "Provide Workspace/s-[Array supported]",
		HideEntered: true,
		Validate:    validate,
	}
	resultWsId, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}
	return resultWsId
}

type addUsersResponse struct {
	Result []addusersResult `json:"result"`
}
type addusersResult struct {
	Id           string   `json:"id"`
	InviteeEmail string   `json:"inviteeEmail"`
	AccountId    int      `json:"accountId"`
	WorkspacesId []int    `json:"workspacesId"`
	DisplayName  string   `json:"displayName"`
	Email        string   `json:"email"`
	Roles        []string `json:"roles"`
}

func addUserByUidWs(userId, workspaceId int) {
	roleWs := userRoleSelectorWs()
	apiId, apiSecret := Getapikeys()
	workspaceIdStr := strconv.Itoa(workspaceId)
	client := &http.Client{}
	//	var data = strings.NewReader(`{"usersIds":[%v],"roles": ["manager"]}`)
	data := fmt.Sprintf(`{"usersIds":[%v],"roles": ["%s"]}`, userId, roleWs)
	var reqBodyData = strings.NewReader(data)
	req, err := http.NewRequest("POST", "https://a.blazemeter.com/api/v4/workspaces/"+workspaceIdStr+"/users", reqBodyData)
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
	var responseBodyWsAddUser addUsersResponse
	json.Unmarshal(bodyText, &responseBodyWsAddUser)
	totalRoles := []string{}
	fmt.Printf("\n%-20s %-30s %-5s", "NAME", "EMAIL", "ROLES")
	for i := 0; i < len(responseBodyWsAddUser.Result); i++ {
		userName := responseBodyWsAddUser.Result[i].DisplayName
		userEmail := responseBodyWsAddUser.Result[i].Email
		for r := 0; r < len(responseBodyWsAddUser.Result[i].Roles); r++ {
			arr := responseBodyWsAddUser.Result[i].Roles[r]
			totalRoles = append(totalRoles, arr)
		}
		fmt.Printf("\n%-20s %-30s %-5s\n", userName, userEmail, totalRoles)
	}
}
func addUserByUidWsraw(userId, workspaceId int) {
	roleWs := userRoleSelectorWs()
	apiId, apiSecret := Getapikeys()
	workspaceIdStr := strconv.Itoa(workspaceId)
	client := &http.Client{}
	data := fmt.Sprintf(`{"usersIds":[%v],"roles": ["%s"]}`, userId, roleWs)
	var reqBodyData = strings.NewReader(data)
	req, err := http.NewRequest("POST", "https://a.blazemeter.com/api/v4/workspaces/"+workspaceIdStr+"/users", reqBodyData)
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

func addUserByEmailA(emailId string, accountId int) {
	apiId, apiSecret := Getapikeys()
	accountIdStr := strconv.Itoa(accountId)
	rolesA, boolVal := userRoleSelectorA()
	resultWsId := workspaceIdPrompt()
	roleWs := userRoleSelectorWs()
	client := &http.Client{}
	data := fmt.Sprintf(`{"invitations":[{"inviteeEmail":"%s","attachAutomatically":%t,"accountRoles":["%s"], "workspacesId":[%v],"workspacesRoles":["%s"]}]}`, emailId, boolVal, rolesA, resultWsId, roleWs)
	var reqBodyData = strings.NewReader(data)
	req, err := http.NewRequest("POST", "https://a.blazemeter.com/api/v4/accounts/"+accountIdStr+"/invitations", reqBodyData)
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
	var responseBodyInviteUser addUsersResponse
	json.Unmarshal(bodyText, &responseBodyInviteUser)
	totalWsinvited := []int{}
	fmt.Printf("\n%-25s %-30s %-15s %-5s\n", "INVITE_ID", "INVITEE", "ACCOUNT", "WORKSPACE")
	for i := 0; i < len(responseBodyInviteUser.Result); i++ {
		inviteId := responseBodyInviteUser.Result[i].Id
		invitee := responseBodyInviteUser.Result[i].InviteeEmail
		invitedAccount := responseBodyInviteUser.Result[i].AccountId

		for w := 0; w < len(responseBodyInviteUser.Result[i].WorkspacesId); w++ {
			wsIdarr := responseBodyInviteUser.Result[i].WorkspacesId[w]
			totalWsinvited = append(totalWsinvited, wsIdarr)
		}
		fmt.Printf("%-25s %-30s %-15v %-5v\n", inviteId, invitee, invitedAccount, totalWsinvited)
	}
}
func addUserByEmailAraw(emailId string, accountId int) {
	apiId, apiSecret := Getapikeys()
	accountIdStr := strconv.Itoa(accountId)
	rolesA, boolVal := userRoleSelectorA()
	resultWsId := workspaceIdPrompt()
	roleWs := userRoleSelectorWs()
	client := &http.Client{}
	data := fmt.Sprintf(`{"invitations":[{"inviteeEmail":"%s","attachAutomatically":%t,"accountRoles":["%s"],
	"workspacesId":[%v],"workspacesRoles":["%s"]}]}`, emailId, boolVal, rolesA, resultWsId, roleWs)
	var reqBodyData = strings.NewReader(data)
	req, err := http.NewRequest("POST", "https://a.blazemeter.com/api/v4/accounts/"+accountIdStr+"/invitations", reqBodyData)
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

// This below is an Admin level command
func addUserByUidA(userId, accountId int) {
	apiId, apiSecret := Getapikeys()
	accountIdStr := strconv.Itoa(accountId)
	client := &http.Client{}
	data := fmt.Sprintf(`{ "accountId": %v, "id": %v }`, accountId, userId)
	var reqBodyData = strings.NewReader(data)
	req, err := http.NewRequest("POST", "https://a.blazemeter.com/api/v4/"+accountIdStr+"/{s}/users", reqBodyData)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("accept", "application/json")
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
