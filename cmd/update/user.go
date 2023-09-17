/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package update

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
	Short: "Update users in Account or Workspace",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("user called")
	},
}

func init() {
	UpdateCmd.AddCommand(userCmd)
	userCmd.Flags().Int("uid", 0, "Enter the User ID")
	userCmd.MarkFlagRequired("uid")
}

func userRoleSelectorA() string {
	prompt := promptui.Select{
		Label: "Select Role",
		Items: []string{"admin", "standard", "user_manager", "billing"},
	}
	_, roleSelected, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}
	fmt.Printf("You choose %q\n", roleSelected)
	return roleSelected
}

func updateUserByUidA(userId, accountId int) {
	roleA := userRoleSelectorA()
	apiId, apiSecret := Getapikeys()
	accountIdStr := strconv.Itoa(accountId)
	client := &http.Client{}
	userIdStr := strconv.Itoa(userId)
	// var data = strings.NewReader(`{"roles":["user_manager"],"enabled": false}`)
	data := fmt.Sprintf(`{"roles": ["%s"], "enabled": true}`, roleA)
	fmt.Println(data)
	var reqBodyData = strings.NewReader(data)
	fmt.Println(reqBodyData)
	req, err := http.NewRequest("PUT", "https://a.blazemeter.com/api/v4/accounts/"+accountIdStr+"/users/"+userIdStr, reqBodyData)
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
