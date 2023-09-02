/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package find

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Get details about the user",
	Long:  `.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("user find called")
		emailIdUser, _ := cmd.Flags().GetString("email")
		getUserByEmail(emailIdUser)
	},
}

func init() {
	FindCmd.AddCommand(userCmd)
	userCmd.Flags().StringP("email", "e", "", "Source directory to read from")
	userCmd.MarkFlagRequired("email")
}

type responseBody struct {
	Result []userResult `json:"result"`
}

type userResult struct {
	Id             int      `json:"id"`
	FirstName      string   `json:"firstName"`
	LastName       string   `json:"LastName"`
	DefaultProject struct{} `json:"defaultProject"`
}

type DefaultProject struct {
	AccountId   int `json:"accountId"`
	WorkspaceId int `json:"workspaceId"`
}

func getUserByEmail(emailIdUser string) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/admin/users?email="+emailIdUser, nil)
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
	var responseObject responseBody
	json.Unmarshal(bodyText, &responseObject)
	fmt.Println("Total users found: ", len(responseObject.Result))
	fmt.Println("userId: ", responseObject.Result[0].Id)
	fmt.Println("firstName: ", responseObject.Result[0].FirstName)
	fmt.Println("lastName: ", responseObject.Result[0].LastName)

}

// Here you will define your flags and configuration settings.

// Cobra supports Persistent Flags which will work for this command
// and all subcommands, e.g.:
// userCmd.PersistentFlags().String("foo", "", "A help for foo")

// Cobra supports local flags which will only run when this command
// is called directly, e.g.:
// userCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
