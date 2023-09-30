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
	"strings"

	"github.com/spf13/cobra"
)

// invitationsCmd represents the invitations command
var invitationsCmd = &cobra.Command{
	Use:   "invitations",
	Short: "Get a list of pending invitations within your account",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ac, _ := cmd.Flags().GetBool("ac")
		var accountId int
		if ac {
			accountId = defaultAccount()
		} else {
			accountId, _ = cmd.Flags().GetInt("accountid")
		}
		getInvitations(accountId)
	},
}

func init() {
	GetCmd.AddCommand(invitationsCmd)
}

type invitesResponse struct {
	Result []invitesResult `json:"result"`
	Error  errorResult     `json:"error"`
}
type invitesResult struct {
	Id              string   `json:"id"`
	InviteeEmail    string   `json:"inviteeEmail"`
	AccountName     string   `json:"accountName"`
	WorkspaceNames  []string `json:"workspaceNames"`
	AccountRoles    []string `json:"accountRoles"`
	WorkspacesRoles []string `json:"workspacesRoles"`
}

func getInvitations(accountId int) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	accountIdStr := strconv.Itoa(accountId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/accounts/"+accountIdStr+"/invitations?limit=300", nil)
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
	var responseBodyInvites invitesResponse
	json.Unmarshal(bodyText, &responseBodyInvites)
	if responseBodyInvites.Error.Code == 0 {
		totalWsNames := []string{}
		totalARoles := []string{}
		//totalWsRoles := []string{}
		fmt.Printf("\n%-25s %-20s %-10s %-5s\n", "INVITEE EMAIL", "ACCOUNT", "AC_ROLE", "WORKSPACE/S & ROLES")
		for i := 0; i < len(responseBodyInvites.Result); i++ {
			accountName := responseBodyInvites.Result[i].AccountName
			userEmail := responseBodyInvites.Result[i].InviteeEmail

			for w := 0; w < len(responseBodyInvites.Result[i].WorkspaceNames); w++ {
				arr := responseBodyInvites.Result[i].WorkspaceNames[w]
				totalWsNames = append(totalWsNames, arr)
			}
			for ar := 0; ar < len(responseBodyInvites.Result[i].AccountRoles); ar++ {
				arr1 := responseBodyInvites.Result[i].AccountRoles[ar]
				totalARoles = append(totalARoles, arr1)
			}
			for wr := 0; wr < len(responseBodyInvites.Result[i].WorkspacesRoles); wr++ {
				arr2 := responseBodyInvites.Result[i].WorkspacesRoles[wr]
				totalWsNames = append(totalWsNames, arr2)
			}
			result1 := strings.Join(totalARoles, ",")
			fmt.Printf("\n%-25s %-20s %-10s %-5s\n", userEmail, accountName, result1, totalWsNames)
		}
		fmt.Println("\n-")
	} else {
		errorCode := responseBodyInvites.Error.Code
		errorMessage := responseBodyInvites.Error.Message
		fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
	}

}
