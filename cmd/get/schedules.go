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

	"github.com/jsuar/go-cron-descriptor/pkg/crondescriptor"
	"github.com/spf13/cobra"
)

// schedulesCmd represents the schedules command
var schedulesCmd = &cobra.Command{
	Use:   "schedules",
	Short: "Get a list of schedules in the account",
	Long:  ``,
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
		case (accountId != 0) && (workspaceId == 0) && rawOutput:
			getShedulesAraw(accountId)
		case (accountId == 0) && (workspaceId != 0) && rawOutput:
			getShedulesWsraw(workspaceId)
		case (accountId != 0) && (workspaceId == 0) && !rawOutput:
			getShedulesA(accountId)
		case (accountId == 0) && (workspaceId != 0) && !rawOutput:
			getShedulesWs(workspaceId)
		default:
			fmt.Println("Please provide a correct workspace Id or Account Id to get the info")
			fmt.Println("[bmgo get -a <account_id>...] OR [bmgo get -w <workspace_id>...]\n-")
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

func getShedulesA(accountId int) {
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
	var responseBodyShedules shedulesResponse
	json.Unmarshal(bodyText, &responseBodyShedules)
	if responseBodyShedules.Error.Code == 0 {
		fmt.Printf("\n%-10s %-10s %-8s %-28s %-50s \n", "TEST", "OWNER", "ENABLED", "SCHEDULE ID", "CRON")
		for i := 0; i < len(responseBodyShedules.Result); i++ {
			sheduleId := responseBodyShedules.Result[i].Id
			scheduleTest := responseBodyShedules.Result[i].TestId
			sheduleOwn := responseBodyShedules.Result[i].CreatedById
			sheduleCron := responseBodyShedules.Result[i].Cron
			sheduleEnabled := responseBodyShedules.Result[i].Enabled
			//	sheduleNextRun := responseBodyShedules.Result[i].Next
			cd, _ := crondescriptor.NewCronDescriptor(sheduleCron)
			fullDescription, _ := cd.GetDescription(crondescriptor.Full)
			fmt.Printf("\n%-10v %-10v %-8v %-28s %-50s ", scheduleTest, sheduleOwn, sheduleEnabled, sheduleId, *fullDescription)
		}
		fmt.Println("\n-")
	} else {
		errorCode := responseBodyShedules.Error.Code
		errorMessage := responseBodyShedules.Error.Message
		fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
	}

}
func getShedulesAraw(accountId int) {
	apiId, apiSecret := Getapikeys()

	accountIdStr := strconv.Itoa(accountId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/schedules?accountId="+accountIdStr+"&limit=500", nil)
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
	fmt.Printf("%s\n", bodyText)
}

func getShedulesWs(workspaceId int) {
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
	var responseBodyShedulesWs shedulesResponse
	json.Unmarshal(bodyText, &responseBodyShedulesWs)
	if responseBodyShedulesWs.Error.Code == 0 {
		fmt.Printf("\n%-10s %-10s %-8s %-28s %-50s \n", "TEST", "OWNER", "ENABLED", "SCHEDULE ID", "CRON")
		for i := 0; i < len(responseBodyShedulesWs.Result); i++ {
			sheduleId := responseBodyShedulesWs.Result[i].Id
			scheduleTest := responseBodyShedulesWs.Result[i].TestId
			sheduleOwn := responseBodyShedulesWs.Result[i].CreatedById
			sheduleCron := responseBodyShedulesWs.Result[i].Cron
			sheduleEnabled := responseBodyShedulesWs.Result[i].Enabled
			//	sheduleNextRun := responseBodyShedules.Result[i].Next
			cd, _ := crondescriptor.NewCronDescriptor(sheduleCron)
			fullDescription, _ := cd.GetDescription(crondescriptor.Full)
			fmt.Printf("\n%-10v %-10v %-8v %-28s %-50s ", scheduleTest, sheduleOwn, sheduleEnabled, sheduleId, *fullDescription)
		}
		fmt.Println("\n-")
	} else {
		errorCode := responseBodyShedulesWs.Error.Code
		errorMessage := responseBodyShedulesWs.Error.Message
		fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
	}
}

func getShedulesWsraw(workspaceId int) {
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
	fmt.Printf("%s\n", bodyText)
}
