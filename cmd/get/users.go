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

// usersCmd represents the users command
var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "[>]Get a list of users part of the account",
	Long: `Use the command to list Users within a specified account, team (API monitoring), or a workspace. The output includes User ID, Name, Roles, Email, etc. The output can be further filtered by switching disabled flag as true to only display disabled users, --disabled.
	
	For example: [bmgo get -w <workspace id> users] OR
	             [bmgo get -a <account id> users] OR
		     [bmgo get -t <team id> users] OR
			     [bmgo get -t <team id> users --pages <page_number>] OR
	             [bmgo get -w <workspace id> users --disabled] OR
	             [bmgo get -a <account id> users --disabled]

    For default: [bmgo get --ws users] OR
	             [bmgo get --ac users] OR 
	             [bmgo get --tm users] OR
				 [bmgo get --tm users --pages <page_number>] OR
	             [bmgo get --ws users --disabled] OR
	             [bmgo get --ac users --disabled]`,
	Run: func(cmd *cobra.Command, args []string) {
		ac, _ := cmd.Flags().GetBool("ac")
		ws, _ := cmd.Flags().GetBool("ws")
		tm, _ := cmd.Flags().GetBool("tm")
		var accountId, workspaceId int
		var teamId string
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
		if tm {
			teamId = defaultTeam()
		} else {
			teamId, _ = cmd.Flags().GetString("teamid")
		}
		rawOutput, _ := cmd.Flags().GetBool("raw")
		disabledUsers, _ := cmd.Flags().GetBool("disabled")
		csvOutput, _ := cmd.Flags().GetBool("csv")
		pages, _ := cmd.Flags().GetInt("pages")
		switch {
		case accountId != 0 && workspaceId == 0 && teamId == "":
			getUsersA(accountId, disabledUsers, rawOutput, csvOutput)
		case accountId == 0 && workspaceId != 0 && teamId == "":
			getUsersWS(workspaceId, disabledUsers, rawOutput, csvOutput)
		case accountId == 0 && workspaceId == 0 && teamId != "":
			getUsersTm(teamId, rawOutput, csvOutput, pages)
		default:
			cmd.Help()
		}
	},
}

func init() {
	GetCmd.AddCommand(usersCmd)
	usersCmd.Flags().Bool("disabled", false, "[Optional] will show disabled users only")
	usersCmd.Flags().Bool("csv", false, "This will output in csv format")
	usersCmd.Flags().IntP("pages", "p", 1, "Total pages of output, 1 page only contains 200 max entries for this")
}

type usersResponse struct {
	Result []usersResult `json:"result"`
	Data   []usersData   `json:"data"`
	Error  errorResult   `json:"error"`
}

type usersResult struct {
	Id          int      `json:"id"`
	Email       string   `json:"email"`
	DisplayName string   `json:"displayName"`
	Enabled     bool     `json:"enabled"`
	Roles       []string `json:"roles"`
	LastAccess  int      `json:"lastAccess"`
}
type usersData struct {
	Uuid          string `json:"uuid"`
	Email         string `json:"email"`
	Role_name     string `json:"role_name"`
	Created_at    string `json:"created_at"`
	Name          string `json:"name"`
	Last_login_at string `json:"last_login_at"`
}

func getUsersA(accountId int, disabledUsers, rawOutput, csvOutput bool) {
	apiId, apiSecret := Getapikeys()

	client := &http.Client{}
	accountIdStr := strconv.Itoa(accountId)
	var req *http.Request
	var err error
	if disabledUsers {
		req, err = http.NewRequest("GET", "https://a.blazemeter.com/api/v4/accounts/"+accountIdStr+"/users?limit=-1&enabled=false", nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		req, err = http.NewRequest("GET", "https://a.blazemeter.com/api/v4/accounts/"+accountIdStr+"/users?limit=-1&enabled=true", nil)
		if err != nil {
			log.Fatal(err)
		}
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
		var responseBodyAUsers usersResponse
		json.Unmarshal(bodyText, &responseBodyAUsers)
		if responseBodyAUsers.Error.Code == 0 {
			//	fmt.Printf("\n%-10s %-25s %-8s %-5s        %-35s\n", "ID", "DISPLAY NAME", "ENABLED", "ROLES", "EMAIL")
			tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			// Print headers
			if csvOutput {
				fmt.Fprintln(tabWriter, "ID,DISPLAY_NAME,ENABLED,ROLES,EMAIL,LAST_ACCESS")
			} else {
				fmt.Fprintln(tabWriter, "ID\tDISPLAY_NAME\tENABLED\tROLES\tEMAIL\tLAST_ACCESS")
			}
			//fmt.Fprintln(tabWriter, "ID\tDISPLAY_NAME\tENABLED\tROLES\tEMAIL\tLAST_ACCESS")
			//fmt.Printf("\n%-28s %-8s %-18s %-10s\n", "SHIP ID", "STATE", "LA
			//	rolesListTotal := []string{}
			for i := 0; i < len(responseBodyAUsers.Result); i++ {
				userIdWS := responseBodyAUsers.Result[i].Id
				displayNameWS := responseBodyAUsers.Result[i].DisplayName
				emailIdWS := responseBodyAUsers.Result[i].Email
				enabledUserWS := responseBodyAUsers.Result[i].Enabled
				var totalRoles []string
				for r := 0; r < len(responseBodyAUsers.Result[i].Roles); r++ {
					rolesArr := responseBodyAUsers.Result[i].Roles[r]
					totalRoles = append(totalRoles, rolesArr)
				}
				var lastAccessStr string
				lastAccessEp := int64(responseBodyAUsers.Result[i].LastAccess)
				// This is because there are epoch values as "0", it converts to a time on 1970s, so we want to condition that here:
				if lastAccessEp != 0 {
					lastAccess := time.Unix(lastAccessEp, 0)
					lastAccessStr = fmt.Sprint(lastAccess)
				} else {
					lastAccessStr = "NA"
				}
				if csvOutput {
					fmt.Printf("%d,%s,%t,%s,%s,%s\n", userIdWS, displayNameWS, enabledUserWS, totalRoles, emailIdWS, lastAccessStr)
				} else {
					fmt.Fprintf(tabWriter, "%d\t%s\t%t\t%s\t%s\t%s\n", userIdWS, displayNameWS, enabledUserWS, totalRoles, emailIdWS, lastAccessStr)
				}
			}
			tabWriter.Flush()
			fmt.Println("-")
		} else {
			errorCode := responseBodyAUsers.Error.Code
			errorMessage := responseBodyAUsers.Error.Message
			fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
		}
	}
}

func getUsersWS(workspaceId int, disabledUsers, rawOutput, csvOutput bool) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	workspaceIdStr := strconv.Itoa(workspaceId)
	var req *http.Request
	var err error
	if disabledUsers {
		req, err = http.NewRequest("GET", "https://a.blazemeter.com/api/v4/workspaces/"+workspaceIdStr+"/users?limit=-1&enabled=false", nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		req, err = http.NewRequest("GET", "https://a.blazemeter.com/api/v4/workspaces/"+workspaceIdStr+"/users?limit=-1&enabled=true", nil)
		if err != nil {
			log.Fatal(err)
		}
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
		var responseBodyWsUsers usersResponse
		json.Unmarshal(bodyText, &responseBodyWsUsers)
		if responseBodyWsUsers.Error.Code == 0 {
			//	fmt.Printf("\n%-10s %-25s %-12s %-10s %-10s\n", "ID", "DISPLAY NAME", "ROLES", "ENABLED", "EMAIL")
			tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			// Print headers
			if csvOutput {
				fmt.Fprintln(tabWriter, "ID,DISPLAY_NAME,ROLES,ENABLED,EMAIL,LAST_ACCESS")
			} else {
				fmt.Fprintln(tabWriter, "ID\tDISPLAY_NAME\tROLES\tENABLED\tEMAIL\tLAST_ACCESS")
			}
			//fmt.Fprintln(tabWriter, "ID\tDISPLAY_NAME\tROLES\tENABLED\tEMAIL\tLAST_ACCESS")
			for i := 0; i < len(responseBodyWsUsers.Result); i++ {
				userIdWS := responseBodyWsUsers.Result[i].Id
				displayNameWS := responseBodyWsUsers.Result[i].DisplayName
				emailIdWS := responseBodyWsUsers.Result[i].Email
				enabledUserWS := responseBodyWsUsers.Result[i].Enabled
				var lastAccessStr string
				lastAccessEp := int64(responseBodyWsUsers.Result[i].LastAccess)
				// This is because there are epoch values as "0", it converts to a time on 1970s, so we want to condition that here:
				if lastAccessEp != 0 {
					lastAccess := time.Unix(lastAccessEp, 0)
					lastAccessStr = fmt.Sprint(lastAccess)
				} else {
					lastAccessStr = "NA"
				}
				if csvOutput {
					fmt.Printf("%d,%s,%s,%t,%s,%s\n", userIdWS, displayNameWS, responseBodyWsUsers.Result[i].Roles[0], enabledUserWS, emailIdWS, lastAccessStr)
				} else {
					fmt.Fprintf(tabWriter, "%d\t%s\t%s\t%t\t%s\t%s\n", userIdWS, displayNameWS, responseBodyWsUsers.Result[i].Roles[0], enabledUserWS, emailIdWS, lastAccessStr)
				}
				//fmt.Printf("\n%-10v %-25s %-12s %-10t %-10s", userIdWS, displayNameWS, responseBodyWsUsers.Result[i].Roles[0], enabledUserWS, emailIdWS)
			}
			tabWriter.Flush()
			fmt.Println("-")
		} else {
			errorCode := responseBodyWsUsers.Error.Code
			errorMessage := responseBodyWsUsers.Error.Message
			fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
		}
	}
}

func getUsersTm(teamId string, rawOutput, csvOutput bool, pages int) {
	Bearer := fmt.Sprintf("Bearer %v", GetPersonalAccessToken())
	client := &http.Client{}
	tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	// Print headers
	if csvOutput {
		fmt.Fprintln(tabWriter, "UUID,ROLES,NAME,EMAIL,LAST_ACCESS")
	} else {
		fmt.Fprintln(tabWriter, "UUID\tROLES\tNAME\tEMAIL\tLAST_ACCESS")
	}
	for i := 1; i <= pages; i++ {
		req, err := http.NewRequest("GET", "https://api.runscope.com/teams/"+teamId+"/people?count=200&page="+strconv.Itoa(i), nil)
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
			var responseBodyTmUsers usersResponse
			json.Unmarshal(bodyText, &responseBodyTmUsers)
			if responseBodyTmUsers.Error.Status == 0 {

				for i := 0; i < len(responseBodyTmUsers.Data); i++ {
					userIdTm := responseBodyTmUsers.Data[i].Uuid
					userNameTm := responseBodyTmUsers.Data[i].Name
					userEmailTm := responseBodyTmUsers.Data[i].Email
					userRoleTm := responseBodyTmUsers.Data[i].Role_name
					lastAccessedTm := responseBodyTmUsers.Data[i].Last_login_at
					//	userCreatedTm := responseBodyTmUsers.Data[i].Created_at
					//fmt.Printf("\n%-38s %-14s %-28s %-10s", userIdTm, userRoleTm, userNameTm, userEmailTm)
					if csvOutput {
						fmt.Printf("%s,%s,%s,%s,%s\n", userIdTm, userRoleTm, userNameTm, userEmailTm, lastAccessedTm)
					} else {
						fmt.Fprintf(tabWriter, "%s\t%s\t%s\t%s\t%s\n", userIdTm, userRoleTm, userNameTm, userEmailTm, lastAccessedTm)
					}
				}
				tabWriter.Flush()
			} else {
				errorCode := responseBodyTmUsers.Error.Status
				errorMessage := responseBodyTmUsers.Error.Message
				fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
			}
		}
	}
}
