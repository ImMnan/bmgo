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

	"github.com/spf13/cobra"
)

// rolesCmd represents the agents command
var rolesCmd = &cobra.Command{
	Use:   "roles",
	Short: "[>]Get roles within an API Monitoring team",
	Long: `The command returns a list of available roles in an API Monitoring team (default & created), you will need to provide a team id to run the command. Outputs "ROLE NAME", "ROLE ID", etc.
	
	For example: [bmgo get -t <team id> roles ] OR 
	             [bmgo get --team <team id> roles]
	For default: [bmgo get --tm roles]]`,
	Run: func(cmd *cobra.Command, args []string) {
		tm, _ := cmd.Flags().GetBool("tm")
		var teamId string
		if tm {
			teamId = defaultTeam()
		} else {
			teamId, _ = cmd.Flags().GetString("teamid")
		}
		rawOutput, _ := cmd.Flags().GetBool("raw")
		switch {
		case teamId != "":
			getRolesTm(teamId, rawOutput)
		default:
			cmd.Help()
		}
	},
}

func init() {
	GetCmd.AddCommand(rolesCmd)
}

type UserGroupResponse struct {
	Meta struct {
		Status string `json:"status"`
	} `json:"meta"`

	Data []struct {
		UUID        string   `json:"uuid"`
		Name        string   `json:"name"`
		Permissions []string `json:"permissions"`
	} `json:"data"`

	Error errorResult `json:"error"`
}

func getRolesTm(teamId string, rawOutput bool) {
	Bearer := fmt.Sprintf("Bearer %v", GetPersonalAccessToken())
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.runscope.com/teams/"+teamId+"/roles", nil)
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
	if rawOutput {
		fmt.Printf("%s\n", bodyText)
	} else {
		var responseBodyTmRoles UserGroupResponse
		json.Unmarshal(bodyText, &responseBodyTmRoles)
		if responseBodyTmRoles.Error.Code == 0 {
			for _, group := range responseBodyTmRoles.Data {
				UUID := group.UUID
				name := group.Name
				//var permissions []string
				fmt.Printf("UUID: %s\nROLENAME: %s\nPERMISSIONS: ", UUID, name)
				for i := 0; i < len(group.Permissions); i++ {
					fmt.Println(group.Permissions[i])
				}
				fmt.Println("\n-")
			}
		} else {
			errorCode := responseBodyTmRoles.Error.Status
			errorMessage := responseBodyTmRoles.Error.Message
			fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
		}

	}
}
