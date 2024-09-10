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
	"time"

	"github.com/spf13/cobra"
)

// workspacesCmd represents the workspaces command
var workspacesCmd = &cobra.Command{
	Use:   "workspaces",
	Short: "Get a list of workspaces in the account",
	Long: `Use the command to list workspaced within a specified account. The output includes workspace ID, Name, Members Count, etc.
	For example: [bmgo get -a <account id> workspaces]
	For default: [bmgo get --ac workspaces]`,
	Run: func(cmd *cobra.Command, args []string) {
		ac, _ := cmd.Flags().GetBool("ac")
		var accountId int
		if ac {
			accountId = defaultAccount()
		} else {
			accountId, _ = cmd.Flags().GetInt("accountid")
		}
		rawOutput, _ := cmd.Flags().GetBool("raw")
		if accountId != 0 {
			getWorkspaces(accountId, rawOutput)
		} else {
			cmd.Help()
		}
	},
}

func init() {
	GetCmd.AddCommand(workspacesCmd)
}

type workspacesResponse struct {
	Result []wsResult  `json:"result"`
	Error  errorResult `json:"error"`
}
type wsResult struct {
	Name         string `json:"name"`
	Enabled      bool   `json:"enabled"`
	MembersCount int    `json:"membersCount"`
	AccountId    int    `json:"accountId"`
	Created      int    `json:"created"`
	Id           int    `json:"id"`
}

func getWorkspaces(accountId int, rawOutput bool) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	accountIdStr := strconv.Itoa(accountId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/workspaces?accountId="+accountIdStr+"&limit=200", nil)
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
		var responseObjectWS workspacesResponse
		json.Unmarshal(bodyText, &responseObjectWS)
		if responseObjectWS.Error.Code == 0 {
			//	fmt.Printf("\n%-10s %-35s %-10s %-10s %-30s\n", "ID", "NAME", "MEMBERS", "ENABLED", "CREATED")
			tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			// Print headers
			fmt.Fprintln(tabWriter, "ID\tNAME\tMEMBERS\tENABLED\tCREATED")
			for i := 0; i < len(responseObjectWS.Result); i++ {
				workspaceId := responseObjectWS.Result[i].Id
				workspaceName := responseObjectWS.Result[i].Name
				members := responseObjectWS.Result[i].MembersCount
				createdepoch := int64(responseObjectWS.Result[i].Created)
				enabled := responseObjectWS.Result[i].Enabled
				created := time.Unix(createdepoch, 0)
				createdstr := fmt.Sprint(created)
				//fmt.Printf("\n% -10v %-35s %-10d %-10t %-30v", workspaceId, workspaceName, members, enabled, createdstr[0:16])
				fmt.Fprintf(tabWriter, "%d\t%s\t%d\t%t\t%s\n", workspaceId, workspaceName, members, enabled, createdstr[0:16])
			}
			tabWriter.Flush()
			fmt.Println("-")
		} else {
			errorCode := responseObjectWS.Error.Code
			errorMessage := responseObjectWS.Error.Message
			fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
		}
	}
}
