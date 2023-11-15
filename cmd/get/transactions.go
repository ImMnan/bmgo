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
	Long: `Use the command to list Transactions within a specified workspace. A Transaction is a request/response pair that is associated with a given Service. The output includes transactions ID, Name, Service Id, etc. The output can be further filtered by specifying a service id by using the --serviceid flag.

	For example: [bmgo get -w <workspace id> transactions] OR
	             [bmgo get -w <workspace id> transactions --serviceid <service id>]
	For default: [bmgo get --ws transactions] OR 
	             [bmgo get --ws transactions --svc <service id>]`,
	Run: func(cmd *cobra.Command, args []string) {
		ws, _ := cmd.Flags().GetBool("ws")
		var workspaceId int
		if ws {
			workspaceId = defaultWorkspace()
		} else {
			workspaceId, _ = cmd.Flags().GetInt("workspaceid")
		}
		rawOutput, _ := cmd.Flags().GetBool("raw")
		serviceId, _ := cmd.Flags().GetInt("svc")
		switch {
		case (workspaceId != 0 || serviceId != 0) && rawOutput:
			getTransactionsWsraw(workspaceId, serviceId)
		case (workspaceId != 0 || serviceId != 0) && !rawOutput:
			getTransactionsWs(workspaceId, serviceId)
		default:
			cmd.Help()
		}
	},
}

func init() {
	GetCmd.AddCommand(transactionsCmd)
	transactionsCmd.Flags().Int("svc", 0, "Provide a service Id to filter the results")
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

func getTransactionsWs(workspaceId, serviceId int) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	workspaceIdStr := strconv.Itoa(workspaceId)
	var req *http.Request
	var err error
	if serviceId != 0 {
		serviceIdStr := strconv.Itoa(serviceId)
		req, err = http.NewRequest("GET", "https://mock.blazemeter.com/api/v1/workspaces/"+workspaceIdStr+"/transactions?serviceId="+serviceIdStr+"&limit=-1", nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		req, err = http.NewRequest("GET", "https://mock.blazemeter.com/api/v1/workspaces/"+workspaceIdStr+"/transactions?limit=-1", nil)
		if err != nil {
			log.Fatal(err)
		}
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
func getTransactionsWsraw(workspaceId, serviceId int) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	workspaceIdStr := strconv.Itoa(workspaceId)
	var req *http.Request
	var err error
	if serviceId != 0 {
		serviceIdStr := strconv.Itoa(serviceId)
		req, err = http.NewRequest("GET", "https://mock.blazemeter.com/api/v1/workspaces/"+workspaceIdStr+"/transactions?serviceId="+serviceIdStr+"&limit=-1", nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		req, err = http.NewRequest("GET", "https://mock.blazemeter.com/api/v1/workspaces/"+workspaceIdStr+"/transactions?limit=-1", nil)
		if err != nil {
			log.Fatal(err)
		}
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
