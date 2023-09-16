/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package add

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os/user"
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

		if (workspaceId != 0) && (accountId == 0) && (emailId == "") {
			addUserByUidWs(userId, workspaceId)
		} else if (workspaceId == 0) && (accountId != 0) && (emailId != "") {
			addUserByEmailA(emailId, accountId)
		} else if (workspaceId != 0) && (accountId == 0) {
			addUserByUidWs(userId, workspaceId)
		} else if (accountId != 0) && (workspaceId == 0) {
			addUserByUidA(userId, accountId)
		} else if (workspaceId != 0) && (accountId == 0) && rawOutput && emailId != "" {
			addUserByUidA(userId, accountId)
		} else {
			fmt.Println("\nPlease provide a correct workspace Id or Account Id to get the info")
			fmt.Println("[bmgo get -a <account_id>...] OR [bmgo get -w <workspace_id>...]")
		}
	},
}

func init() {
	AddCmd.AddCommand(userCmd)
	userCmd.Flags().Int("uid", 0, "User ID for the user")
	//	userCmd.MarkFlagRequired("uid")
	userCmd.Flags().String("email", "", "Enter the Email ID of the user invited!")
	//	userWsCmd.Flags().BoolP("raw", "r", false, "[Optional] If set, the output will be raw json")
}
func userRoleSelectorA() (string, bool) {
	prompt := promptui.Select{
		Label: "Select Account Role",
		Items: []string{"admin", "standard", "user_manager", "billing"},
	}
	prompt1 := promptui.Select{
		Label: "attachAutomatically",
		Items: []bool{true, false},
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
	fmt.Printf("You choose %q and %t\n", roleSelected, boolVal)
	return roleSelected, boolVal
}
func userRoleSelectorWs() string {
	prompt := promptui.Select{
		Label: "Select Workspace Role",
		Items: []string{"tester", "manager", "viewer"},
	}
	_, roleSelected, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}
	fmt.Printf("You choose %q\n", roleSelected)
	return roleSelected
}

func workspaceIdPrompt() string {
	var username string
	u, err := user.Current()
	if err == nil {
		username = u.Username
	}
	prompt := promptui.Prompt{
		Label:   "Provide Workspace ID to add user into",
		Default: username,
	}
	resultWsId, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}
	return resultWsId
}

func addUserByUidWs(userId, workspaceId int) {
	roleWs := userRoleSelectorWs()
	apiId, apiSecret := Getapikeys()
	workspaceIdStr := strconv.Itoa(workspaceId)
	client := &http.Client{}
	//	var data = strings.NewReader(`{"usersIds":[%v],"roles": ["manager"]}`)
	data := fmt.Sprintf(`{"usersIds":[%v],"roles": ["%s"]}`, userId, roleWs)
	var reqBodyData = strings.NewReader(data)
	fmt.Println(reqBodyData)
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

// This below is an Admin level command
func addUserByUidA(userId, accountId int) {
	apiId, apiSecret := Getapikeys()
	accountIdStr := strconv.Itoa(accountId)
	client := &http.Client{}
	data := fmt.Sprintf(`{ "accountId": %v, "id": %v }`, accountId, userId)
	fmt.Println(data)
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

func addUserByEmailA(emailId string, accountId int) {
	apiId, apiSecret := Getapikeys()
	accountIdStr := strconv.Itoa(accountId)
	roleWs := userRoleSelectorWs()
	rolesA, boolVal := userRoleSelectorA()
	resultWsId := workspaceIdPrompt()
	client := &http.Client{}
	data := fmt.Sprintf(`{"invitations":[{"inviteeEmail":"%s","attachAutomatically":%t,"accountRoles":["%s"],
	"workspacesId":[%v],"workspacesRoles":["%s"]}]}`, emailId, boolVal, rolesA, resultWsId, roleWs)
	fmt.Println(data)
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
