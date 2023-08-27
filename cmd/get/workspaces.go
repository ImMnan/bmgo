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

// workspacesCmd represents the workspaces command
var workspacesCmd = &cobra.Command{
	Use:   "workspaces",
	Short: "Get a list of workspaces in the account",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("workspaces called")
		accountId, _ := cmd.Flags().GetInt("id")
		getWorkspaces(accountId)
	},
}

func init() {
	accountCmd.AddCommand(workspacesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// workspacesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// workspacesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getWorkspaces(accountId int) {
	apiId, apiSecret := Getapikeys()

	client := &http.Client{}
	accountIdStr := strconv.Itoa(accountId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/workspaces?accountId="+accountIdStr+"&limit=200", nil)
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
