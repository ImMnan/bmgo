/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package get

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
)

// usersCmd represents the users command
var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Get a list of users part of the account",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("users called")
		accountId, _ := cmd.Flags().GetInt("id")
		getUsers(accountId)
	},
}

func init() {
	accountCmd.AddCommand(usersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// usersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// usersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getUsers(accountId int) {
	apiId, apiSecret := Getapikeys()

	client := &http.Client{}
	accountIdStr := strconv.Itoa(accountId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/accounts/"+accountIdStr+"/users?limit=1500", nil)
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
	fmt.Printf("%s\n", bodyText)
}
