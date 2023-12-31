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

// mocksCmd represents the transactions command
var mocksCmd = &cobra.Command{
	Use:   "mocks",
	Short: "Get all mock services in a workspace Or within a service",
	Long: `Use the command to list the mock services within a specified workspace. A Mock Service is filled with a collection of Transactions, typically a subset of the Transactions in a Service. The output included Mock service ID, its NAME, SERVICE Name etc. The output can be further filtered by specifying a service by using the --serviceid flag.
	
	For example: [bmgo get -w <workspace id> mocks] OR 
	             [bmgo get -w <workspace id> mocks --serviceid <service id>
	For default: [bmgo get --ws mocks] OR 
	             [bmgo get --ws mocks --svc <service id>]`,
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
			getMocksWsraw(workspaceId, serviceId)
		case (workspaceId != 0 || serviceId != 0) && !rawOutput:
			getMocksWs(workspaceId, serviceId)
		default:
			cmd.Help()
		}
	},
}

func init() {
	GetCmd.AddCommand(mocksCmd)
	mocksCmd.Flags().Int("svc", 0, "Provide a service Id")
}

type mocksResponse struct {
	Result []mockResult `json:"result"`
	Error  string       `json:"error"`
}
type mockResult struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	ServiceId   int    `json:"serviceId"`
	ServiceName string `json:"serviceName"`
}

func getMocksWs(workspaceId, serviceId int) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	workspaceIdStr := strconv.Itoa(workspaceId)
	var req *http.Request
	var err error
	if serviceId != 0 {
		serviceIdStr := strconv.Itoa(serviceId)
		//	serviceId := serviceIdPrompt()
		req, err = http.NewRequest("GET", "https://mock.blazemeter.com/api/v1/workspaces/"+workspaceIdStr+"/service-mocks?serviceId="+serviceIdStr, nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		req, err = http.NewRequest("GET", "https://mock.blazemeter.com/api/v1/workspaces/"+workspaceIdStr+"/service-mocks", nil)
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
	var responseBodyMocks mocksResponse
	json.Unmarshal(bodyText, &responseBodyMocks)
	if len(responseBodyMocks.Result) >= 1 {
		fmt.Printf("\n%-10s %-35s %-10s %-15s\n", "ID", "NAME", "SERVICE", "SERVICE NAME")
		for i := 0; i < len(responseBodyMocks.Result); i++ {
			mockId := responseBodyMocks.Result[i].Id
			mockName := responseBodyMocks.Result[i].Name
			serviceId := responseBodyMocks.Result[i].ServiceId
			serviceName := responseBodyMocks.Result[i].ServiceName
			fmt.Printf("\n%-10d %-35s %-10d %-15s", mockId, mockName, serviceId, serviceName)
		}
		fmt.Println("\n-")
	} else {
		errorCode := 404
		fmt.Printf("\nError code: %v\nError Message: No Transactions found in workspace %v, provide correct workspace Id.\n\n", errorCode, workspaceId)
	}
}
func getMocksWsraw(workspaceId, serviceId int) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	workspaceIdStr := strconv.Itoa(workspaceId)
	var req *http.Request
	var err error
	if serviceId != 0 {
		serviceIdStr := strconv.Itoa(serviceId)
		//	serviceId := serviceIdPrompt()
		req, err = http.NewRequest("GET", "https://mock.blazemeter.com/api/v1/workspaces/"+workspaceIdStr+"/service-mocks?serviceId="+serviceIdStr, nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		req, err = http.NewRequest("GET", "https://mock.blazemeter.com/api/v1/workspaces/"+workspaceIdStr+"/service-mocks", nil)
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
