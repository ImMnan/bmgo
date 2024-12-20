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

	"github.com/jsuar/go-cron-descriptor/pkg/crondescriptor"
	"github.com/spf13/cobra"
)

// schedulesCmd represents the schedules command
var schedulesCmd = &cobra.Command{
	Use:   "schedules",
	Short: "Get a list of schedules in the account or workspace",
	Long: `Use the command to list Schedules within a specified workspace or account. Tests can be scheduled to run on frequencies up to every minute. One or more schedules can be configured per test. The output includes Schedule ID, Test ID, Cron, etc.

	For example: [bmgo get -w <workspace id> schedules] OR
	             [bmgo get -a <account id> schedules]
	For default: [bmgo get --ws schedules] OR 
	             [bmgo get --ac schedules]`,
	Run: func(cmd *cobra.Command, args []string) {
		ac, _ := cmd.Flags().GetBool("ac")
		ws, _ := cmd.Flags().GetBool("ws")
		var accountId, workspaceId int
		if ac {
			accountId = defaultAccount()
		} else {
			accountId, _ = cmd.Flags().GetInt("accountid")
		}
		if ws {
			workspaceId = defaultWorkspace()
		} else {
			workspaceId, _ = cmd.Flags().GetInt("workspaceid")
		}
		rawOutput, _ := cmd.Flags().GetBool("raw")
		switch {
		case (accountId != 0) && (workspaceId == 0):
			getShedulesA(accountId, rawOutput)
		case (accountId == 0) && (workspaceId != 0):
			getShedulesWs(workspaceId, rawOutput)
		default:
			cmd.Help()
		}
	},
}

func init() {
	GetCmd.AddCommand(schedulesCmd)
}

type shedulesResponse struct {
	Result []scheduleResult `json:"result"`
	Error  errorResult      `json:"error"`
}
type scheduleResult struct {
	Id          string `json:"id"`
	TestId      int    `json:"testId"`
	NextRun     int    `json:"nextRun"`
	Cron        string `json:"cron"`
	CreatedById int    `json:"createdById"`
	Enabled     bool   `json:"enabled"`
}

func getShedulesA(accountId int, rawOutput bool) {
	apiId, apiSecret := Getapikeys()
	accountIdStr := strconv.Itoa(accountId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/schedules?accountId="+accountIdStr+"&limit=1000", nil)
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
		var responseBodyShedules shedulesResponse
		json.Unmarshal(bodyText, &responseBodyShedules)
		if responseBodyShedules.Error.Code == 0 {
			//	fmt.Printf("\n%-10s %-10s %-8s %-28s %-50s \n", "TEST", "OWNER", "ENABLED", "SCHEDULE ID", "CRON")
			tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			// Print headers
			fmt.Fprintln(tabWriter, "TEST\tOWNER\tENABLED\tSCHEDULE_ID\tCRON")
			for i := 0; i < len(responseBodyShedules.Result); i++ {
				sheduleId := responseBodyShedules.Result[i].Id
				scheduleTest := responseBodyShedules.Result[i].TestId
				sheduleOwn := responseBodyShedules.Result[i].CreatedById
				sheduleCron := responseBodyShedules.Result[i].Cron
				sheduleEnabled := responseBodyShedules.Result[i].Enabled
				//	sheduleNextRun := responseBodyShedules.Result[i].Next
				cd, _ := crondescriptor.NewCronDescriptor(sheduleCron)
				fullDescription, _ := cd.GetDescription(crondescriptor.Full)
				//	fmt.Printf("\n%-10v %-10v %-8v %-28s %-50s ", scheduleTest, sheduleOwn, sheduleEnabled, sheduleId, *fullDescription)
				fmt.Fprintf(tabWriter, "%d\t%d\t%t\t%s\t%s\n", scheduleTest, sheduleOwn, sheduleEnabled, sheduleId, *fullDescription)
			}
			tabWriter.Flush()
			fmt.Println("-")
		} else {
			errorCode := responseBodyShedules.Error.Code
			errorMessage := responseBodyShedules.Error.Message
			fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
		}
	}
}

func getShedulesWs(workspaceId int, rawOutput bool) {
	apiId, apiSecret := Getapikeys()
	wprkspaceIdStr := strconv.Itoa(workspaceId)

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/schedules?workspaceId="+wprkspaceIdStr+"&limit=1000", nil)
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
		var responseBodyShedulesWs shedulesResponse
		json.Unmarshal(bodyText, &responseBodyShedulesWs)
		if responseBodyShedulesWs.Error.Code == 0 {
			//	fmt.Printf("\n%-10s %-10s %-8s %-28s %-50s \n", "TEST", "OWNER", "ENABLED", "SCHEDULE ID", "CRON")
			tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			// Print headers
			fmt.Fprintln(tabWriter, "TEST\tOWNER\tENABLED\tSCHEDULE_ID\tCRON")
			for i := 0; i < len(responseBodyShedulesWs.Result); i++ {
				sheduleId := responseBodyShedulesWs.Result[i].Id
				scheduleTest := responseBodyShedulesWs.Result[i].TestId
				sheduleOwn := responseBodyShedulesWs.Result[i].CreatedById
				sheduleCron := responseBodyShedulesWs.Result[i].Cron
				sheduleEnabled := responseBodyShedulesWs.Result[i].Enabled
				//	sheduleNextRun := responseBodyShedules.Result[i].Next
				cd, _ := crondescriptor.NewCronDescriptor(sheduleCron)
				fullDescription, _ := cd.GetDescription(crondescriptor.Full)
				//	fmt.Printf("\n%-10v %-10v %-8v %-28s %-50s ", scheduleTest, sheduleOwn, sheduleEnabled, sheduleId, *fullDescription)
				fmt.Fprintf(tabWriter, "%d\t%d\t%t\t%s\t%s\n", scheduleTest, sheduleOwn, sheduleEnabled, sheduleId, *fullDescription)
			}
			tabWriter.Flush()
			fmt.Println("-")
		} else {
			errorCode := responseBodyShedulesWs.Error.Code
			errorMessage := responseBodyShedulesWs.Error.Message
			fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
		}
	}
}
