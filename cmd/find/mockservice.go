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
	"time"

	"github.com/spf13/cobra"
)

// testCmd represents the test command
var mockserviceCmd = &cobra.Command{
	Use:   "mockservice",
	Short: "Find mockservice details",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		mockId, _ := cmd.Flags().GetInt("mockid")
		workspaceId, _ := cmd.Flags().GetInt("workspaceid")
		rawOutput, _ := cmd.Flags().GetBool("raw")
		if mockId != 0 && rawOutput {
			findMockraw(mockId, workspaceId)
		} else if mockId != 0 {
			findMock(mockId, workspaceId)
		} else {
			fmt.Println("\nPlease provide a valid Mock ID to find the Test")
			fmt.Println("[bmgo find mockservice --mockid <Mock Service id> OR bmgo find -w <workspace Id> mockservice --mockid <Mock Service id>")
		}
	},
}

func init() {
	FindCmd.AddCommand(mockserviceCmd)
	mockserviceCmd.Flags().Int("mockid", 0, "Provide a mock Id to find a test")
	mockserviceCmd.MarkFlagRequired("mockid")
}

type findmockResponse struct {
	Result mockResult `json:"result"`
	Error  string     `json:"error"`
}
type mockResult struct {
	Name          string             `json:"name"`
	ServiceId     int                `json:"serviceId"`
	ServiceName   string             `json:"serviceName"`
	Status        string             `json:"status"`
	Location      string             `json:"locationName"`
	ShipId        string             `json:"shipId"`
	HttpEndpoint  string             `json:"httpEndpoint"`
	HttpsEndpoint string             `json:"httpsEndpoint"`
	Created       int                `json:"created"`
	Badges        []processingAction `json:"badges"`
}

type processingAction struct {
	Webhook     string `json:"WEBHOOK"`
	HttpCall    string `json:"HTTP_CALL"`
	StateUpdate string `json:"STATE_UPDATE"`
}

func findMock(mockId, workspaceId int) {
	mockIdStr := strconv.Itoa(mockId)
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	var workspaceIdStr string
	if workspaceId == 0 {
		workspaceIdStr = workspaceIdPrompt()
	} else {
		workspaceIdStr = strconv.Itoa(workspaceId)
	}
	req, err := http.NewRequest("GET", "https://mock.blazemeter.com/api/v1/workspaces/"+workspaceIdStr+"/service-mocks/"+mockIdStr, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("accept", "application/json;charset=UTF-8")
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
	var responseObjectMockservice findmockResponse
	json.Unmarshal(bodyText, &responseObjectMockservice)
	if responseObjectMockservice.Error == "" {
		mockName := responseObjectMockservice.Result.Name
		serviceName := responseObjectMockservice.Result.ServiceName
		serviceId := responseObjectMockservice.Result.ServiceId
		mockStatus := responseObjectMockservice.Result.Status
		mockCreated := int64(responseObjectMockservice.Result.Created)
		mockCreatedStr := fmt.Sprint(time.Unix(mockCreated, 0))
		mockHttpEndpoint := responseObjectMockservice.Result.HttpEndpoint
		mockHttpsEndpoint := responseObjectMockservice.Result.HttpsEndpoint
		mockShipId := responseObjectMockservice.Result.ShipId
		mockLocation := responseObjectMockservice.Result.Location
		//	fmt.Printf("Name: %s\nService Name: %s   Service ID: %d\n", mockName, serviceName, serviceId)
		fmt.Printf("\n%-30s %-10s %-18s %-10s %-10s\n", "MOCK NAME", "STATUS", "CREATED", "SERVICE", "SERVICE NAME")
		fmt.Printf("%-30s %-10s %-18s %-10d %-10s\n", mockName, mockStatus, mockCreatedStr[0:16], serviceId, serviceName)
		fmt.Printf("\nHTTP ENDPOINT: %s\nHTTPS ENDPOINT: %s", mockHttpEndpoint, mockHttpsEndpoint)
		fmt.Printf("\nLOCATION: %s\nAGENT: %s", mockLocation, mockShipId)
		fmt.Println("\n-")
	} else {
		fmt.Printf("Error code: %v\n", responseObjectMockservice.Error)
	}
}

func findMockraw(mockId, workspaceId int) {
	mockIdStr := strconv.Itoa(mockId)
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	var workspaceIdStr string
	if workspaceId == 0 {
		workspaceIdStr = workspaceIdPrompt()
	} else {
		workspaceIdStr = strconv.Itoa(workspaceId)
	}
	req, err := http.NewRequest("GET", "https://mock.blazemeter.com/api/v1/workspaces/"+workspaceIdStr+"/service-mocks/"+mockIdStr, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("accept", "application/json;charset=UTF-8")
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
