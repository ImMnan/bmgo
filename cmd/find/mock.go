/*
Copyright © 2024 Manan Patel - github.com/immnan
*/
package find

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

// mockserviceCmd represents the test command
var mockCmd = &cobra.Command{
	Use:   "mock",
	Short: "Find mockservice details",
	Long: ` Use this command to find details about the specofied mock service (--mockid) along with the workspace id it belongs to (--ws Or -w). Global Flag --raw can be used for raw Json output. 
	For example: [bmgo find mock --mockid <Mock Service id> --workspaceid <workspace Id>] OR 
	For default: [bmgo find mock --mockid <Mock Service id> --ws]
	To Download: [[bmgo find mock --mockid <Mock Service id> --ws --download] OR 
	[bmgo find mock --mockid <Mock Service id> --ws --download --filename <filename>]`,
	Run: func(cmd *cobra.Command, args []string) {
		mockId, _ := cmd.Flags().GetInt("mockid")
		download, _ := cmd.Flags().GetBool("download")
		fileName, _ := cmd.Flags().GetString("filename")
		ws, _ := cmd.Flags().GetBool("ws")
		var workspaceId int
		if ws {
			workspaceId = defaultWorkspace()
		} else {
			workspaceId, _ = cmd.Flags().GetInt("workspaceid")
		}
		rawOutput, _ := cmd.Flags().GetBool("raw")
		if mockId != 0 && workspaceId != 0 {
			findMock(mockId, workspaceId, rawOutput, download, fileName)
		} else if workspaceId == 0 {
			fmt.Println("Workspace Id is required, please specify a workspace Id using -w flag or use --ws flag to use default workspace ID")
		} else {
			cmd.Help()
		}
	},
}

func init() {
	FindCmd.AddCommand(mockCmd)
	mockCmd.Flags().Int("mockid", 0, "Provide a mock service Id to find a mockservice test")
	mockCmd.MarkFlagRequired("mockid")
	mockCmd.Flags().IntP("workspaceid", "w", 0, "Provide Workspace ID to look for a mockservice test")
	mockCmd.Flags().BoolP("download", "d", false, "Specify if you want to download the mockservice transactions")
	mockCmd.Flags().Bool("ws", false, "Use default workspace Id (bmConfig.yaml)")
	mockCmd.Flags().StringP("filename", "f", "mockservice", "Specify the filename to download the mockservice transactions to")
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
	HarbourId     string             `json:"harborId"`
	CreatedBy     string             `json:"createdBy"`
}

type processingAction struct {
	Webhook     string `json:"WEBHOOK"`
	HttpCall    string `json:"HTTP_CALL"`
	StateUpdate string `json:"STATE_UPDATE"`
}

func findMock(mockId, workspaceId int, rawOutput, download bool, fileName string) {
	mockIdStr := strconv.Itoa(mockId)
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	var workspaceIdStr string
	if workspaceId == 0 {
		workspaceIdStr = workspaceIdPrompt()
	} else {
		workspaceIdStr = strconv.Itoa(workspaceId)
	}
	downloadFile := fileName + ".txs.json"

	var req *http.Request
	var err error
	if download {
		req, err = http.NewRequest("GET", "https://mock.blazemeter.com/api/v1/workspaces/"+workspaceIdStr+"/transactions?serviceMockId="+mockIdStr+"&sort=-id&limit=-1&type=GENERIC_DSL&exportAsFile=true", nil)
	} else {
		req, err = http.NewRequest("GET", "https://mock.blazemeter.com/api/v1/workspaces/"+workspaceIdStr+"/service-mocks/"+mockIdStr, nil)
	}

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

	if download {
		// Handle the downloaded file
		err = os.WriteFile(downloadFile, bodyText, 0644)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Mock downloaded to %s\n", downloadFile)
	} else {
		if rawOutput {
			fmt.Printf("%s\n", bodyText)
		} else {
			var responseObjectMockservice findmockResponse
			json.Unmarshal(bodyText, &responseObjectMockservice)
			if responseObjectMockservice.Error == "" {
				mockCreatedBy := responseObjectMockservice.Result.CreatedBy
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
				mockHarbour := responseObjectMockservice.Result.HarbourId
				//	fmt.Printf("Name: %s\nService Name: %s   Service ID: %d\n", mockName, serviceName, serviceId)
				fmt.Printf("NAME: %v\n\n", mockName)
				//fmt.Printf("\n%-35s %-10s %-18s %-10s %-10s\n", "CREATED BY", "STATUS", "CREATED", "SERVICE", "SERVICE NAME")
				tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
				// Print headers
				fmt.Fprintln(tabWriter, "CREATED_BY\tSTATUS\tCREATED\tSERVICE\tSERVICE_NAME")
				//fmt.Printf("%-35s %-10s %-18s %-10d %-10s\n", mockCreatedBy, mockStatus, mockCreatedStr[0:16], serviceId, serviceName)
				fmt.Fprintf(tabWriter, "%s\t%s\t%s\t%d\t%s\n", mockCreatedBy, mockStatus, mockCreatedStr[0:16], serviceId, serviceName)
				tabWriter.Flush()
				fmt.Println("-")
				fmt.Printf("HTTP ENDPOINT: %s\nHTTPS ENDPOINT: %s", mockHttpEndpoint, mockHttpsEndpoint)
				fmt.Printf("\nLOCATION: %s\nHARBOUR: %s\nAGENT: %s\n", mockLocation, mockHarbour, mockShipId)
			} else {
				fmt.Printf("Error code: %v\n", responseObjectMockservice.Error)
			}
		}
	}
}
