/*
Copyright © 2024 Manan Patel - github.com/immnan
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
	Short: "[!]Add user to the account or workspace using User ID or invite email",
	Long: `Add users to workspace using the User Id, considering the user is already in the same account. Also add users to an account by inviting the user via email address, in this case user is not part of Blazemeter. 
	
	For example: [bmgo add -a <account_id> user <user email>] OR [bmgo add -w <workspace_id> user --uid <user id>]
	For default: [bmgo add --ac user <user email>] OR [bmgo add --ws user --uid <user id>]`,

	Run: func(cmd *cobra.Command, args []string) {
		//	fmt.Println("user add called")
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
		userId, _ := cmd.Flags().GetInt("uid")
		emailId, _ := cmd.Flags().GetString("email")
		rawOutput, _ := cmd.Flags().GetBool("raw")
		switch {
		case (workspaceId == 0) && (accountId != 0) && (emailId != ""):
			addUserByEmailA(emailId, accountId, rawOutput)
		case (workspaceId != 0) && (accountId == 0) && (userId != 0):
			addUserByUidWs(userId, workspaceId, rawOutput)
		default:
			cmd.Help()
		}
	},
}

func init() {
	AddCmd.AddCommand(userCmd)
	userCmd.Flags().Int("uid", 0, "User ID for the user")
	userCmd.Flags().String("email", "", "Enter the Email ID of the user invited!")
	userCmd.Flags().IntP("accountid", "a", 0, " Provide Account ID to add a resource to")
	userCmd.Flags().Bool("ac", false, "Use default account Id (bmConfig.yaml)")
	userCmd.Flags().IntP("workspaceid", "w", 0, " Provide Workspace ID to add a resource to")
	userCmd.Flags().Bool("ws", false, "Use default workspace Id (bmConfig.yaml)")
}

type addUsersResponse struct {
	Result []addusersResult `json:"result"`
	Error  errorResult      `json:"error"`
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

func addUserByUidWs(userId, workspaceId int, rawOutput bool) {
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
	if rawOutput {
		fmt.Printf("%s\n", bodyText)
	} else {
		var responseBodyWsAddUser addUsersResponse
		json.Unmarshal(bodyText, &responseBodyWsAddUser)
		if responseBodyWsAddUser.Error.Code == 0 {
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
		} else {
			errorCode := responseBodyWsAddUser.Error.Code
			errorMessage := responseBodyWsAddUser.Error.Message
			fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
		}
	}
}

func addUserByEmailA(emailId string, accountId int, rawOutput bool) {
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
	if rawOutput {
		fmt.Printf("%s\n", bodyText)
	} else {
		var responseBodyInviteUser addUsersResponse
		json.Unmarshal(bodyText, &responseBodyInviteUser)
		if responseBodyInviteUser.Error.Code == 0 {
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
		} else {
			errorCode := responseBodyInviteUser.Error.Code
			errorMessage := responseBodyInviteUser.Error.Message
			fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
		}
	}
}
