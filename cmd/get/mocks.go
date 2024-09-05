/*
Copyright © 2024 Manan Patel - github.com/immnan
*/
package get

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/tabwriter"

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
		case (workspaceId != 0 || serviceId != 0):
			getMocksWs(workspaceId, serviceId, rawOutput)
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

func getMocksWs(workspaceId, serviceId int, rawOutput bool) {
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
	if rawOutput {
		fmt.Printf("%s\n", bodyText)
	} else {
		var responseBodyMocks mocksResponse
		json.Unmarshal(bodyText, &responseBodyMocks)
		if len(responseBodyMocks.Result) >= 1 {
			//fmt.Printf("\n%-10s %-35s %-10s %-15s\n", "ID", "NAME", "SERVICE", "SERVICE NAME")
			tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			// Print headers
			fmt.Fprintln(tabWriter, "ID\tNAME\tSERVICE\tSERVICE_NAME")
			for i := 0; i < len(responseBodyMocks.Result); i++ {
				mockId := responseBodyMocks.Result[i].Id
				mockName := responseBodyMocks.Result[i].Name
				serviceId := responseBodyMocks.Result[i].ServiceId
				serviceName := responseBodyMocks.Result[i].ServiceName
				//		fmt.Printf("\n%-10d %-35s %-10d %-15s", mockId, mockName, serviceId, serviceName)
				fmt.Fprintf(tabWriter, "%d\t%s\t%d\t%s\n", mockId, mockName, serviceId, serviceName)
			}
			tabWriter.Flush()
			fmt.Println("-")
		} else {
			errorCode := 404
			fmt.Printf("\nError code: %v\nError Message: No Transactions found in workspace %v, provide correct workspace Id.\n\n", errorCode, workspaceId)
		}
	}
}
