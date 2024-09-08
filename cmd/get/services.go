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

// servicesCmd represents the services command
var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "Get services [for Mock service] within workspace",
	Long: `Use the command to list Services within a specified workspace. Within BlazeMeter, a Service is a logical grouping of Transactions. These Transactions can be anything, but typically, a Service is a grouping of Transactions that are related to a specific live service. The output includes service ID, Name, etc. 

	For example: [bmgo get -w <workspace id> services] OR
	For default: [bmgo get --ws services]`,
	Run: func(cmd *cobra.Command, args []string) {
		ws, _ := cmd.Flags().GetBool("ws")
		var workspaceId int
		if ws {
			workspaceId = defaultWorkspace()
		} else {
			workspaceId, _ = cmd.Flags().GetInt("workspaceid")
		}
		rawOutput, _ := cmd.Flags().GetBool("raw")
		if workspaceId != 0 {
			getServicesWs(workspaceId, rawOutput)
		} else {
			cmd.Help()
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

func getServicesWs(workspaceId int, rawOutput bool) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	workspaceIdStr := strconv.Itoa(workspaceId)
	req, err := http.NewRequest("GET", "https://mock.blazemeter.com/api/v1/workspaces/"+workspaceIdStr+"/services?limit=-1&skip=0&sort=name", nil)
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
	if rawOutput {
		fmt.Printf("%s\n", bodyText)
	} else {
		var responseBodyServices servicesResponse
		json.Unmarshal(bodyText, &responseBodyServices)
		if len(responseBodyServices.Result) >= 1 {
			//	fmt.Printf("\n%-10s %-30s %-30s\n", "ID", "NAME", "DESCRIPTION")
			tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			// Print headers
			fmt.Fprintln(tabWriter, "ID\tNAME\tDESCRIPTION")
			for i := 0; i < len(responseBodyServices.Result); i++ {
				serviceId := responseBodyServices.Result[i].Id
				serviceName := responseBodyServices.Result[i].Name
				serviceDescr := responseBodyServices.Result[i].Description
				//fmt.Printf("\n%-10d %-30s %-30s", serviceId, serviceName, serviceDescr)
				fmt.Fprintf(tabWriter, "%d\t%s\t%s\n", serviceId, serviceName, serviceDescr)
			}
			tabWriter.Flush()
			fmt.Println("-")
		} else if responseBodyServices.Error != "" {
			errorCode := 404
			fmt.Printf("\nError Code: %v\nError Message: %v\n-", errorCode, responseBodyServices.Error)
		} else {
			errorCode := 404
			fmt.Printf("\nError code: %v\nError Message:No Transactions found in workspace %v\nPlease provide correct workspace Id.\n-", errorCode, workspaceId)
		}
	}
}
