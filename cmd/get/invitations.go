/*
Copyright © 2024 Manan Patel - github.com/immnan
*/
package get

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// invitationsCmd represents the invitations command
var invitationsCmd = &cobra.Command{
	Use:   "invitations",
	Short: "Get a list of pending invitations within your account",
	Long: `Use this command to list existing invitations for the account, only outputs the pending ones
	For example: [bmgo get -a <account id> invitations]
	For default: [bmgo get --ac invitations]`,
	Run: func(cmd *cobra.Command, args []string) {
		ac, _ := cmd.Flags().GetBool("ac")
		var accountId int
		if ac {
			accountId = defaultAccount()
		} else {
			accountId, _ = cmd.Flags().GetInt("accountid")
		}
		rawOutput, _ := cmd.Flags().GetBool("raw")
		getInvitations(accountId, rawOutput)
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

func getInvitations(accountId int, rawOutput bool) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	accountIdStr := strconv.Itoa(accountId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/accounts/"+accountIdStr+"/invitations?limit=0", nil)
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
	if rawOutput {
		fmt.Printf("%s\n", bodyText)
	} else {
		var responseBodyInvites invitesResponse
		json.Unmarshal(bodyText, &responseBodyInvites)
		if responseBodyInvites.Error.Code == 0 {
			//totalWsRoles := []string{}
			//fmt.Printf("\n%-30s %-20s %-20s %-10s %-5s\n", "INVITEE EMAIL", "ACCOUNT", "WORKSPACE", "AC_ROLE", "WORKSPACE ROLES")
			tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			// Print headers
			fmt.Fprintln(tabWriter, "INVITEE_EMAIL\tACCOUNT\tWORKSPACE\tAC_ROLE\tWORKSPACE_ROLES")
			for i := 0; i < len(responseBodyInvites.Result); i++ {
				accountName := responseBodyInvites.Result[i].AccountName
				userEmail := responseBodyInvites.Result[i].InviteeEmail
				var workspaceName, acRoles, wsRoles string
				for w := 0; w < len(responseBodyInvites.Result[i].WorkspaceNames); w++ {
					workspaceName = responseBodyInvites.Result[i].WorkspaceNames[w]
				}
				for ar := 0; ar < len(responseBodyInvites.Result[i].AccountRoles); ar++ {
					acRoles = responseBodyInvites.Result[i].AccountRoles[ar]
				}
				for wr := 0; wr < len(responseBodyInvites.Result[i].WorkspacesRoles); wr++ {
					wsRoles = responseBodyInvites.Result[i].WorkspacesRoles[wr]
				}

				fmt.Printf("\n%-30s %-20s %-20s %-10s %-5s", userEmail, accountName, workspaceName, acRoles, wsRoles)
				fmt.Fprintf(tabWriter, "%s\t%s\t%s\t%s\t%s\n", userEmail, accountName, workspaceName, acRoles, wsRoles)
			}
			tabWriter.Flush()
			fmt.Println("\n-")
		} else {
			errorCode := responseBodyInvites.Error.Code
			errorMessage := responseBodyInvites.Error.Message
			fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
		}
	}
}
