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
		accountId, _ := cmd.Flags().GetInt("accountid")
		getWorkspaces(accountId)
	},
}

func init() {
	GetCmd.AddCommand(workspacesCmd)
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
