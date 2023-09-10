/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package get

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
)

// usersCmd represents the users command
var usersWSCmd = &cobra.Command{
	Use:   "users",
	Short: "Get a list of users part of the account",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("users called")
		workspaceId, _ := cmd.Flags().GetInt("id")
		rawOutput, _ := cmd.Flags().GetBool("raw")
		enabledUsers, _ := usersCmd.Flags().GetBool("enabled")
		if rawOutput && enabledUsers {
			getUsersWSraw(workspaceId)
		} else if rawOutput {
			getUsersWSrawDis(workspaceId)
		} else if enabledUsers {
			getUsersWSDis(workspaceId)
		} else {
			getUsersWS(workspaceId)
		}
	},
}

func init() {
	workspaceCmd.AddCommand(usersWSCmd)
	usersCmd.Flags().Bool("enabled", true, "[Optional] will show enabled users only")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// usersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// usersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type usersResponse struct {
	Result []usersResult `json:"result"`
}

type usersResult struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
	Enabled     bool   `json:"enabled"`
}

func getUsersWS(workspaceId int) {
	apiId, apiSecret := Getapikeys()

	client := &http.Client{}
	workspaceIdStr := strconv.Itoa(workspaceId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/workspaces/"+workspaceIdStr+"/users?limit=1000&enabled=true", nil)
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
	var responseBodyWsUsers usersResponse
	json.Unmarshal(bodyText, &responseBodyWsUsers)
	fmt.Printf("\n%-10s %-25s %-25s %-10s\n", "ID", "NAME", "EMAIL", "ENABLED")
	for i := 0; i < len(responseBodyWsUsers.Result); i++ {
		fmt.Printf("\n%-10v %-25s %-25s %-10t", (responseBodyWsUsers.Result[i].Id), (responseBodyWsUsers.Result[i].DisplayName), (responseBodyWsUsers.Result[i].Email), (responseBodyWsUsers.Result[i].Enabled))
	}
	fmt.Println("\n")
}

func getUsersWSraw(workspaceId int) {
	apiId, apiSecret := Getapikeys()

	client := &http.Client{}
	workspaceIdStr := strconv.Itoa(workspaceId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/workspaces/"+workspaceIdStr+"/users?limit=1000&enabled=true", nil)
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

func getUsersWSDis(workspaceId int) {
	apiId, apiSecret := Getapikeys()

	client := &http.Client{}
	workspaceIdStr := strconv.Itoa(workspaceId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/workspaces/"+workspaceIdStr+"/users?limit=1000&enabled=true", nil)
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
	var responseBodyWsUsers usersResponse
	json.Unmarshal(bodyText, &responseBodyWsUsers)
	fmt.Printf("\n%-10s %-25s %-25s %-10s\n", "ID", "NAME", "EMAIL", "ENABLED")
	for i := 0; i < len(responseBodyWsUsers.Result); i++ {
		fmt.Printf("\n%-10v %-25s %-25s %-10t", (responseBodyWsUsers.Result[i].Id), (responseBodyWsUsers.Result[i].DisplayName), (responseBodyWsUsers.Result[i].Email), (responseBodyWsUsers.Result[i].Enabled))
	}
	fmt.Println("\n")
}

func getUsersWSrawDis(workspaceId int) {
	apiId, apiSecret := Getapikeys()

	client := &http.Client{}
	workspaceIdStr := strconv.Itoa(workspaceId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/workspaces/"+workspaceIdStr+"/users?limit=1000&enabled=true", nil)
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
