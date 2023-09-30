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
	"strconv"
	"strings"

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
		case (workspaceId != 0) && (accountId == 0) && !rawOutput:
			addUserByUidWs(userId, workspaceId)
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
