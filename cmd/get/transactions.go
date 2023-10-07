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
		serviceId, _ := cmd.Flags().GetInt("serviceid")
		switch {
		case workspaceId != 0 && serviceId == 0 && rawOutput:
			getTransactionsWsraw(workspaceId)
		case workspaceId != 0 && serviceId != 0 && rawOutput:
			getTransactionsServiceraw(workspaceId, serviceId)
		case workspaceId != 0 && serviceId == 0:
			getTransactionsWs(workspaceId)
		case workspaceId != 0 && serviceId != 0:
			getTransactionsService(workspaceId, serviceId)
		default:
			fmt.Println("\nPlease provide a valid Workspace ID &+OR Service id to get list of transactions")
			fmt.Println("[bmgo get -w <workspace id>...")
		}
	},
}

func init() {
	GetCmd.AddCommand(transactionsCmd)
	transactionsCmd.Flags().Int("serviceid", 0, "Provide a service Id to filter the results")
}

type transactionsResponse struct {
	Result []transactionsResult `json:"result"`
	Error  string               `json:"error"`
}
type transactionsResult struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	ServiceId   int    `json:"serviceId"`
	ServiceName string `json:"serviceName"`
}

func getTransactionsWs(workspaceId int) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	workspaceIdStr := strconv.Itoa(workspaceId)
	//	serviceId := serviceIdPrompt()
	req, err := http.NewRequest("GET", "https://mock.blazemeter.com/api/v1/workspaces/"+workspaceIdStr+"/transactions?limit=-1", nil)
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
		errorCode := 404
		fmt.Printf("\nError code: %v\nError Message: No Transactions found in workspace %v, provide correct workspace Id.\n\n", errorCode, workspaceId)
	}
}
func getTransactionsWsraw(workspaceId int) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	workspaceIdStr := strconv.Itoa(workspaceId)
	//	serviceId := serviceIdPrompt()
	req, err := http.NewRequest("GET", "https://mock.blazemeter.com/api/v1/workspaces/"+workspaceIdStr+"/transactions?limit=-1", nil)
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
func getTransactionsService(workspaceId, serviceId int) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	workspaceIdStr := strconv.Itoa(workspaceId)
	serviceIdStr := strconv.Itoa(serviceId)
	req, err := http.NewRequest("GET", "https://mock.blazemeter.com/api/v1/workspaces/"+workspaceIdStr+"/transactions?serviceId="+serviceIdStr+"&limit=-1", nil)
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
	} else if responseBodyTransactions.Error != "" {
		errorCode := 404
		fmt.Printf("\nError Code: %v\nError Message: %v\n-", errorCode, responseBodyTransactions.Error)
	} else {
		errorCode := 404
		fmt.Printf("\nError code: %v\nError Message:No Transactions found in workspace %v or service %v.\nPlease provide correct workspace Id or Service Id.\n\n-", errorCode, workspaceId, serviceId)
	}
}
func getTransactionsServiceraw(workspaceId, serviceId int) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	workspaceIdStr := strconv.Itoa(workspaceId)
	serviceIdStr := strconv.Itoa(serviceId)
	req, err := http.NewRequest("GET", "https://mock.blazemeter.com/api/v1/workspaces/"+workspaceIdStr+"/transactions?serviceId="+serviceIdStr+"&limit=-1", nil)
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
