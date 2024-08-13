/*
Copyright © 2024 Manan Patel - github.com/immnan
*/
package update

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "[!]Update users in Account or Workspace",
	Long: `Use the command to update user entry, we can either enable or disable the user or make changes to the user-roles within a workspace or account level. To update the user, you will need the user Id & use the flag --accountid or --workspaceid to make changes to specific level.

	For example: [bmgo update user --uid <user Id> -a <account_id>] OR 
                 [bmgo update user --uid <user Id> -w <workspace_id>] 
	For default: [bmgo update user --uid <user Id> --ac  OR 
                 [bmgo update user --uid <user Id> --ws]`,
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
		case (workspaceId != 0) && (accountId == 0):
			updateUserWs(userId, workspaceId, rawOutput)
		case (accountId != 0) && (workspaceId == 0):
			updateUserA(userId, accountId, rawOutput)
		default:
			cmd.Help()
		}
	},
}

func init() {
	UpdateCmd.AddCommand(userCmd)
	userCmd.Flags().Int("uid", 0, "Enter the User ID")
	userCmd.Flags().IntP("accountid", "a", 0, "Account ID of the resource expected to being updated")
	userCmd.Flags().IntP("workspaceid", "w", 0, "Workspace ID of the resource expected to being updated")
	userCmd.Flags().Bool("ac", false, "Use default account Id (bmConfig.yaml)")
	userCmd.Flags().Bool("ws", false, "Use default workspace Id (bmConfig.yaml)")
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

func updateUserA(userId, accountId int, rawOutput bool) {
	apiId, apiSecret := Getapikeys()
	accountIdStr := strconv.Itoa(accountId)
	client := &http.Client{}
	userIdStr := strconv.Itoa(userId)
	enableA := isEnabledPromt()
	var bodyText []byte

	if !enableA {
		//var statusData string
		statusData := fmt.Sprintf(`{"enabled": %t}`, enableA)
		var reqBodyDataA = strings.NewReader(statusData)
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
		bodyText, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		roleA := updateUserRolesA()
		var data string
		if roleA == "admin + owner" {
			data = fmt.Sprintf(`{"roles": ["admin", "owner"], "enabled": %t}`, enableA)
		} else {
			data = fmt.Sprintf(`{"roles": ["%s"], "enabled": %t}`, roleA, enableA)
		}
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
		bodyText, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
	}
	if rawOutput {
		fmt.Printf("%s\n", bodyText)
	} else {
		var responseBodyUpdateUserA updateUserResponse
		var userRoleA string
		json.Unmarshal(bodyText, &responseBodyUpdateUserA)
		if responseBodyUpdateUserA.Error.Code == 0 {

			userEmailA := responseBodyUpdateUserA.Result.Email
			userTypeA := responseBodyUpdateUserA.Result.Type
			userStatusA := responseBodyUpdateUserA.Result.Enabled
			//fmt.Printf("\n%-25s %-12s %-10s %-10s", "EMAIL", "TYPE", "ENABLE", "ROLE")
			tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			// Print headers
			fmt.Fprintln(tabWriter, "EMAIL\tTYPE\tENABLE\tROLE")
			for i := 0; i < len(responseBodyUpdateUserA.Result.Roles); i++ {
				userRoleA = responseBodyUpdateUserA.Result.Roles[i]
			}
			//	fmt.Printf("\n%-25s %-12s %-10t %-10s\n\n", userEmailA, userTypeA, userStatusA, userRoleA)
			fmt.Fprintf(tabWriter, "%s\t%s\t%t\t%s\n", userEmailA, userTypeA, userStatusA, userRoleA)
			tabWriter.Flush()
			fmt.Println("-")
		} else {
			errorCode := responseBodyUpdateUserA.Error.Code
			errorMessage := responseBodyUpdateUserA.Error.Message
			fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
		}
	}
}

func updateUserWs(userId, workspaceId int, rawOutput bool) {
	workspaceIdStr := strconv.Itoa(workspaceId)
	apiId, apiSecret := Getapikeys()
	userIdStr := strconv.Itoa(userId)
	var bodyText []byte
	client := &http.Client{}
	enableWs := isEnabledPromt()
	var data string
	if !enableWs {
		data = fmt.Sprintf(`{"enabled": %t}`, enableWs)
	} else {
		roleWs := updateUserRolesWs()
		data = fmt.Sprintf(`{"roles": ["%s"], "enabled": %t}`, roleWs, enableWs)
	}
	var reqBodyDataA = strings.NewReader(data)
	req, err := http.NewRequest("PUT", "https://a.blazemeter.com/api/v4/workspaces/"+workspaceIdStr+"/users/"+userIdStr, reqBodyDataA)
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
	bodyText, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if rawOutput {
		fmt.Printf("%s\n", bodyText)
	} else {
		var responseBodyUpdateUserWs updateUserResponse
		var userRoleWs string
		json.Unmarshal(bodyText, &responseBodyUpdateUserWs)
		if responseBodyUpdateUserWs.Error.Code == 0 {

			userEmailWs := responseBodyUpdateUserWs.Result.Email
			userTypeWs := responseBodyUpdateUserWs.Result.Type
			userStatusWs := responseBodyUpdateUserWs.Result.Enabled
			//	fmt.Printf("\n%-25s %-12s %-10s %-10s", "EMAIL", "TYPE", "ENABLE", "ROLE")
			tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			// Print headers
			fmt.Fprintln(tabWriter, "EMAIL\tTYPE\tENABLE\tROLE")
			for i := 0; i < len(responseBodyUpdateUserWs.Result.Roles); i++ {
				userRoleWs = responseBodyUpdateUserWs.Result.Roles[i]
			}
			//fmt.Printf("\n%-25s %-12s %-10t %-10s\n\n", userEmailWs, userTypeWs, userStatusWs, userRoleWs)
			fmt.Fprintf(tabWriter, "%s\t%s\t%t\t%s\n", userEmailWs, userTypeWs, userStatusWs, userRoleWs)
			tabWriter.Flush()
			fmt.Println("-")
		} else {
			errorCode := responseBodyUpdateUserWs.Error.Code
			errorMessage := responseBodyUpdateUserWs.Error.Message
			fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
		}
	}
}
