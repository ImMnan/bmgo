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

// servicesCmd represents the services command
var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "Get services [for Mock service] within workspace",
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
			getServicesWsraw(workspaceId)
		} else if workspaceId != 0 {
			getServicesWs(workspaceId)
		} else {
			fmt.Println("Please provide a valid Workspace ID to get list of tests")
			fmt.Println("[bmgo get -w <workspace id>...\n-")
		}
	},
}

func init() {
	GetCmd.AddCommand(servicesCmd)
}

type servicesResponse struct {
	Result []servicesResult `json:"result"`
	Error  string           `json:"error"`
}
type servicesResult struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func getServicesWs(workspaceId int) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	workspaceIdStr := strconv.Itoa(workspaceId)
	req, err := http.NewRequest("GET", "https://mock.blazemeter.com/api/v1/workspaces/"+workspaceIdStr+"/services?limit=500&skip=0&sort=name", nil)
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
	var responseBodyServices servicesResponse
	json.Unmarshal(bodyText, &responseBodyServices)
	if len(responseBodyServices.Result) >= 1 {
		fmt.Printf("\n%-10s %-30s %-30s\n", "ID", "NAME", "DESCRIPTION")
		for i := 0; i < len(responseBodyServices.Result); i++ {
			serviceId := responseBodyServices.Result[i].Id
			serviceName := responseBodyServices.Result[i].Name
			serviceDescr := responseBodyServices.Result[i].Description
			fmt.Printf("\n%-10d %-30s %-30s", serviceId, serviceName, serviceDescr)
		}
		fmt.Println("\n-")
	} else if responseBodyServices.Error != "" {
		errorCode := 404
		fmt.Printf("\nError Code: %v\nError Message: %v\n-", errorCode, responseBodyServices.Error)
	} else {
		errorCode := 404
		fmt.Printf("\nError code: %v\nError Message:No Transactions found in workspace %v\nPlease provide correct workspace Id.\n-", errorCode, workspaceId)
	}
}
func getServicesWsraw(workspaceId int) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	workspaceIdStr := strconv.Itoa(workspaceId)
	req, err := http.NewRequest("GET", "https://mock.blazemeter.com/api/v1/workspaces/"+workspaceIdStr+"/services?limit=500&skip=0&sort=name", nil)
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
