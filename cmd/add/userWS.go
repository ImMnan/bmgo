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
var userWsCmd = &cobra.Command{
	Use:   "user",
	Short: "Get details about the user",
	Long:  `.`,
	Run: func(cmd *cobra.Command, args []string) {
		//	fmt.Println("user add called")
		userId, _ := cmd.Flags().GetInt("uid")
		workspaceId, _ := cmd.Flags().GetInt("workspaceid")
		if workspaceId == 0 {
			fmt.Println("\nPlease provide a Workspace ID or an Account ID")
			fmt.Println("[bmgo add -w <workspace_ID>...] OR [bmgo add -a <account_ID>...]")
		} else {
			addUserByUid(userId, workspaceId)
		}
		//rawOutput, _ := cmd.Flags().GetBool("raw")
		//		if rawOutput {
		//	addUserByUidRaw(userId, workspaceId)
		//	} else {

		//	}
	},
}

func init() {
	AddCmd.AddCommand(userWsCmd)
	userWsCmd.Flags().Int("uid", 0, "User ID for the user")
	userWsCmd.MarkFlagRequired("uid")
	//	userWsCmd.Flags().BoolP("raw", "r", false, "[Optional] If set, the output will be raw json")
}

func userRoleSelector() string {
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

func addUserByUid(userId, workspaceId int) {
	role := userRoleSelector()
	apiId, apiSecret := Getapikeys()
	workspaceIdStr := strconv.Itoa(workspaceId)
	client := &http.Client{}
	//	var data = strings.NewReader(`{"usersIds":[%v],"roles": ["manager"]}`)
	data := fmt.Sprintf(`{"usersIds":[%v],"roles": ["%s"]}`, userId, role)
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
