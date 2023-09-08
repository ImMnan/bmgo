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
	"strconv"

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
		rawOutput, _ := cmd.Flags().GetBool("raw")

		if rawOutput {
			getUserByEmailRaw(emailIdUser)
		} else {
			getUserByEmail(emailIdUser)
		}
	},
}

func init() {
	FindCmd.AddCommand(userCmd)
	userCmd.Flags().StringP("email", "e", "", "Source directory to read from")
	userCmd.MarkFlagRequired("email")
	userCmd.Flags().BoolP("raw", "r", false, "[Optional] If set, the output will be raw json")
}

type responseBody struct {
	Result []userResult `json:"result"`
}

type userResult struct {
	Id             int    `json:"id"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"LastName"`
	DefaultProject defaultProject
	Roles          []string
}

type defaultProject struct {
	AccountId     int    `json:"accountId"`
	WorkspaceId   int    `json:"workspaceId"`
	AccountName   string `json:"accountName"`
	WorkspaceName string `json:"workspaceName"`
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
	//fmt.Printf("%s\n", bodyText)
	rolesListTotal := []string{}
	var responseObject responseBody
	json.Unmarshal(bodyText, &responseObject)

	userId := responseObject.Result[0].Id
	firstName := responseObject.Result[0].FirstName
	lastName := responseObject.Result[0].LastName
	accountIdstr := strconv.Itoa(responseObject.Result[0].DefaultProject.AccountId)
	workspaceIdstr := strconv.Itoa(responseObject.Result[0].DefaultProject.WorkspaceId)

	for i := 0; i < len(responseObject.Result[0].Roles); i++ {
		rolesList := responseObject.Result[0].Roles[i]
		rolesListTotal = append(rolesListTotal, rolesList)
	}

	fmt.Printf("\n%-15s %-15s %-15s %-5s\n", "USERID", "FIRSTNAME", "LASTNAME", "ROLES")
	fmt.Printf("%-15d %-15s %-15s %-5s \n\n", userId, firstName, lastName, rolesListTotal)

	fmt.Println("Navigate to Blazemeter User account: ", "https://a.blazemeter.com/app/#/accounts/"+accountIdstr+"/workspaces/"+workspaceIdstr+"/dashboard")
}

func getUserByEmailRaw(emailIdUser string) {
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
}

// Here you will define your flags and configuration settings.

// Cobra supports Persistent Flags which will work for this command
// and all subcommands, e.g.:
// userCmd.PersistentFlags().String("foo", "", "A help for foo")

// Cobra supports local flags which will only run when this command
// is called directly, e.g.:
// userCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
