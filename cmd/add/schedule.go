/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package add

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jsuar/go-cron-descriptor/pkg/crondescriptor"
	"github.com/spf13/cobra"
)

// scheduleCmd represents the schedule command
var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Add new schedule",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		testId, _ := cmd.Flags().GetInt("tid")
		rawOutput, _ := cmd.Flags().GetBool("raw")
		if testId != 0 && rawOutput {
			addScheduleraw(testId)
		} else if testId != 0 {
			addSchedule(testId)
		} else {
			fmt.Println("\nPlease provide a correct Test ID to add a schedule to")
			fmt.Println("[bmgo add schedule --tid <test id>]")
		}
	},
}

func init() {
	AddCmd.AddCommand(scheduleCmd)
	scheduleCmd.Flags().Int("tid", 0, "Provide a test id to create a schedule for")
}

type addShedulesResponse struct {
	Result addScheduleResult `json:"result"`
	Error  errorResult       `json:"error"`
}
type addScheduleResult struct {
	TestId      int    `json:"testId"`
	Id          string `json:"id"`
	NextRun     []int  `json:"nextRun"`
	Cron        string `json:"cron"`
	CreatedById int    `json:"createdById"`
	Created     int    `json:"created"`
	Enabled     bool   `json:"enabled"`
}

func addSchedule(testId int) {
	apiId, apiSecret := Getapikeys()
	cronExpression := cronPrompt()
	client := &http.Client{}
	data := fmt.Sprintf(`{"cron":"%s","testId":%d}`, cronExpression, testId)
	reqBodydata := strings.NewReader(data)
	req, err := http.NewRequest("POST", "https://a.blazemeter.com/api/v4/schedules", reqBodydata)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
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
	//fmt.Printf("%s\n", bodyText)
	var responseBodyAddShedules addShedulesResponse
	if responseBodyAddShedules.Error.Code == 0 {
		json.Unmarshal(bodyText, &responseBodyAddShedules)
		fmt.Printf("\n%-28s %-10s %-20s %-40s", "SCHEDULE ID", "ENABLED", "NEXT RUN", "CRON")
		scheduleId := responseBodyAddShedules.Result.Id
		//sheduleOwn := responseBodyAddShedules.Result.CreatedById
		sheduleCron := responseBodyAddShedules.Result.Cron
		cd, _ := crondescriptor.NewCronDescriptor(sheduleCron)
		sheduleCronStr, _ := cd.GetDescription(crondescriptor.Full)
		scheduleEnabled := responseBodyAddShedules.Result.Enabled
		sheduleNextrunEP := int64(responseBodyAddShedules.Result.Created)
		sheduleNextRun := time.Unix(sheduleNextrunEP, 0)
		sheduleNextRunStr := fmt.Sprint(sheduleNextRun)
		fmt.Printf("\n%-28s %-10t %-20s %-40s\n\n", scheduleId, scheduleEnabled, sheduleNextRunStr[0:16], *sheduleCronStr)
	} else {
		errorCode := responseBodyAddShedules.Error.Code
		errorMessage := responseBodyAddShedules.Error.Message
		fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
	}
}
func addScheduleraw(testId int) {
	apiId, apiSecret := Getapikeys()
	cronExpression := cronPrompt()
	client := &http.Client{}
	data := fmt.Sprintf(`{"cron":"%s","testId":%d}`, cronExpression, testId)
	reqBodydata := strings.NewReader(data)
	req, err := http.NewRequest("POST", "https://a.blazemeter.com/api/v4/schedules", reqBodydata)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
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
