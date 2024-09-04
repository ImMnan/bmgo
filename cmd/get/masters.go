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
	"time"

	"github.com/spf13/cobra"
)

// mastersCmd represents the masters command
var mastersCmd = &cobra.Command{
	Use:   "masters",
	Short: "Get masters for a test",
	Long: `Use the command to list masters/test runs for a specific test, use the test ID to list the masters for the test. The command outputs STATUS, START TIME, END TIME, etc. for the lister masters.
	
	For example: [bmgo get masters --tid <test id>]`,
	Run: func(cmd *cobra.Command, args []string) {
		testId, _ := cmd.Flags().GetInt("tid")
		rawOutput, _ := cmd.Flags().GetBool("raw")
		switch {
		case testId != 0:
			getMasters(testId, rawOutput)
		default:
			cmd.Help()
		}
	},
}

func init() {
	GetCmd.AddCommand(mastersCmd)
	mastersCmd.Flags().Int("tid", 0, "Provide the test ID to list masters")
}

type mastersResponse struct {
	Result []mastersResult `json:"result"`
	Error  errorResult     `json:"error"`
}

type mastersResult struct {
	Id           int      `json:"id"`
	Status       string   `json:"reportStatus"`
	Created      int      `json:"created"`
	Ended        int      `json:"ended"`
	Locations    []string `json:"locations"`
	SessionId    []string `json:"sessionsId"`
	ProjectId    int      `json:"projectId"`
	RunnerUserId int      `json:"runnerUserId"`
}

func getMasters(testId int, rawOutput bool) {
	apiId, apiSecret := Getapikeys()
	testIdStr := strconv.Itoa(testId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/masters?testId="+testIdStr+"&limit=0", nil)
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
	if rawOutput {
		fmt.Printf("%s\n", bodyText)
	} else {
		var responseObjectMasters mastersResponse
		json.Unmarshal(bodyText, &responseObjectMasters)
		if responseObjectMasters.Error.Code == 0 {
			//fmt.Printf("\n%-10s %-8s %-10s %-10s %-15s %-15s\n", "MASTER", "STATUS", "RUN_BY", "PROJECT", "START", "END")
			tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			// Print headers
			fmt.Fprintln(tabWriter, "MASTER\tSTATUS\tRUN_BY\tPROJECT\tSTART\tEND")

			for i := 0; i < len(responseObjectMasters.Result); i++ {
				masterId := responseObjectMasters.Result[i].Id
				masterStatus := responseObjectMasters.Result[i].Status
				masterProject := responseObjectMasters.Result[i].ProjectId
				masterCreatedEp := int64(responseObjectMasters.Result[i].Created)
				masterEndEp := int64(responseObjectMasters.Result[i].Ended)
				masterRunner := responseObjectMasters.Result[i].RunnerUserId
				if masterCreatedEp != 0 && masterEndEp != 0 {
					masterCreatedStr := fmt.Sprint(time.Unix(masterCreatedEp, 0))
					masterEndStr := fmt.Sprint(time.Unix(masterEndEp, 0))
					//	fmt.Printf("\n%-10d %-8s %-10d %-10d %-15s %-15s", masterId, masterStatus, masterRunner, masterProject, masterCreatedStr[2:16], masterEndStr[5:16])
					fmt.Fprintf(tabWriter, "%d\t%s\t%d\t%d\t%s\t%s\n", masterId, masterStatus, masterRunner, masterProject, masterCreatedStr[2:16], masterEndStr[5:16])
				} else {
					fmt.Fprintf(tabWriter, "%d\t%s\t%d\t%d\t%d\t%d\n", masterId, masterStatus, masterRunner, masterProject, masterCreatedEp, masterEndEp)
				}
			}
			tabWriter.Flush()
			fmt.Println("\n-")
			for i := 0; i < len(responseObjectMasters.Result); i++ {
				masterId := responseObjectMasters.Result[i].Id
				totalLocations := []string{}
				for l := 0; l < len(responseObjectMasters.Result[i].Locations); l++ {
					locations := responseObjectMasters.Result[i].Locations[l]
					totalLocations = append(totalLocations, locations)
				}
				totalSessions := []string{}
				for rv := 0; rv < len(responseObjectMasters.Result[i].SessionId); rv++ {
					sessions := responseObjectMasters.Result[i].SessionId[rv]
					totalSessions = append(totalSessions, sessions)
				}
				fmt.Printf("MASTER: %d\nLOCATIONS: %s\nSESSIONS:  %s\n\n", masterId, totalLocations, totalSessions)
			}
		} else {
			errorCode := responseObjectMasters.Error.Code
			errorMessage := responseObjectMasters.Error.Message
			fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
		}
	}

}
