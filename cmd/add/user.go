/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package add

import (
	"fmt"
	"io"
	"log"
	"net/http"
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
		workspaceId, _ := cmd.Flags().GetInt("workspaceid")
		accountId, _ := cmd.Flags().GetInt("accountid")
		rawOutput, _ := cmd.Flags().GetBool("raw")

		if (workspaceId != 0) && (accountId == 0) && rawOutput {
			addUserByUidWs(userId, workspaceId)
		} else if (accountId != 0) && (workspaceId == 0) && rawOutput {
			addUserByUidA(userId, accountId)
		} else if (workspaceId != 0) && (accountId == 0) {
			addUserByUidWs(userId, workspaceId)
		} else if (accountId != 0) && (workspaceId == 0) {
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
	userCmd.MarkFlagRequired("uid")
	//	userWsCmd.Flags().BoolP("raw", "r", false, "[Optional] If set, the output will be raw json")
}

func userRoleSelectorWs() string {
	prompt := promptui.Select{
		Label: "Select Role",
		Items: []string{"tester", "manager", "viewer"},
	}
	_, roleSelected, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}
	fmt.Printf("You choose %q\n", roleSelected)
	return roleSelected
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

func addUserByUidA(userId, accountId int) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	data := fmt.Sprintf(`{ "accountId": %v, "id": %v }`, accountId, userId)
	fmt.Println(data)
	var reqBodyData = strings.NewReader(data)
	req, err := http.NewRequest("POST", "https://a.blazemeter.com/api/v4/accounts/{s}/users", reqBodyData)
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
