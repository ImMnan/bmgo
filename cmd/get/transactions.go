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

// transactionsCmd represents the transactions command
var transactionsCmd = &cobra.Command{
	Use:   "transactions",
	Short: "Get all transactions in a workspace > within a service",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ws, _ := cmd.Flags().GetBool("ws")
		var workspaceId int
		if ws {
			workspaceId = defaultWorkspace()
		} else {
			workspaceId, _ = cmd.Flags().GetInt("workspaceid")
		}
		rawOutput, _ := cmd.Flags().GetBool("raw")
		if workspaceId != 0 && rawOutput {
			getTransactionsWsraw(workspaceId)
		} else if workspaceId != 0 {
			getTransactionsWs(workspaceId)
		} else {
			fmt.Println("\nPlease provide a valid Workspace ID to get list of tests")
			fmt.Println("[bmgo get -w <workspace id>...")
		}
	},
}

type transactionsResponse struct {
	Result []transactionsResult `json:"result"`
	Error  errorResult          `json:"error"`
}
type transactionsResult struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	ServiceId   int    `json:"serviceId"`
	ServiceName string `json:"serviceName"`
}

func init() {
	GetCmd.AddCommand(transactionsCmd)
}
func getTransactionsWs(workspaceId int) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	workspaceIdStr := strconv.Itoa(workspaceId)
	//	serviceId := serviceIdPrompt()
	req, err := http.NewRequest("GET", "https://mock.blazemeter.com/api/v1/workspaces/"+workspaceIdStr+"/transactions", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("accept", "*/*")
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
	//	fmt.Printf("%s\n", bodyText)
	var responseBodyTransactions transactionsResponse
	json.Unmarshal(bodyText, &responseBodyTransactions)
	if len(responseBodyTransactions.Result) >= 1 {
		fmt.Printf("\n%-10s %-35s %-10s %-15s\n", "ID", "NAME", "SERVICE", "SERVICE NAME")
		for i := 0; i < len(responseBodyTransactions.Result); i++ {
			transId := responseBodyTransactions.Result[i].Id
			transName := responseBodyTransactions.Result[i].Name
			serviceId := responseBodyTransactions.Result[i].ServiceId
			serviceName := responseBodyTransactions.Result[i].ServiceName
			fmt.Printf("\n%-10d %-35s %-10d %-15s", transId, transName, serviceId, serviceName)
		}
		fmt.Println("\n-")
	} else {
		errorCode := 400
		fmt.Printf("\nError code: %v\nError Message: No Transactions found in workspace %v, provide correct workspace Id.\n\n", errorCode, workspaceId)
	}
}
func getTransactionsWsraw(workspaceId int) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	workspaceIdStr := strconv.Itoa(workspaceId)
	//	serviceId := serviceIdPrompt()
	req, err := http.NewRequest("GET", "https://mock.blazemeter.com/api/v1/workspaces/"+workspaceIdStr+"/transactions", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("accept", "*/*")
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
