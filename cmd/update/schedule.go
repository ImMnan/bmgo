/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package update

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
	Short: "Enable or Disable Schedule",
	Long: `Use the command to update Schedule for the test, we can either enable or disable the schedule for the test using this command. To update the schedule, you will need to know the schedule Id of the schedule. Use the flag --sid followed by the schedule Id to update it.

	For example: [bmgo update schedule --sid <schedule ID>] `,
	Run: func(cmd *cobra.Command, args []string) {
		scheduleId, _ := cmd.Flags().GetString("sid")
		rawOutput, _ := cmd.Flags().GetBool("raw")
		if scheduleId != "" && rawOutput {
			updateScheduleraw(scheduleId)
		} else if scheduleId != "" {
			updateSchedule(scheduleId)
		} else {
			cmd.Help()
		}
	},
}

func init() {
	UpdateCmd.AddCommand(scheduleCmd)
	scheduleCmd.Flags().String("scheduleid", "", "Provide the Schedule ID")
}

type updateschedulesResponse struct {
	Result updatescheduleResult `json:"result"`
	Error  errorResult          `json:"error"`
}
type updatescheduleResult struct {
	TestId  int    `json:"testId"`
	Cron    string `json:"cron"`
	Created int    `json:"created"`
	Enabled bool   `json:"enabled"`
}

func updateSchedule(scheduleId string) {
	apiId, apiSecret := Getapikeys()
	status := isEnabledPromt()
	client := &http.Client{}
	data := fmt.Sprintf(`{"enabled":%t}`, status)
	reqBodydata := strings.NewReader(data)
	req, err := http.NewRequest("PATCH", "https://a.blazemeter.com/api/v4/schedules/"+scheduleId, reqBodydata)
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
	var responseBodyUpdateSchedules updateschedulesResponse
	json.Unmarshal(bodyText, &responseBodyUpdateSchedules)
	fmt.Printf("\n%-10s %-10s %-20s %-40s", "TEST", "ENABLED", "CREATED ON", "CRON")
	if responseBodyUpdateSchedules.Error.Code == 0 {
		scheduleTest := responseBodyUpdateSchedules.Result.TestId
		sheduleEnabled := responseBodyUpdateSchedules.Result.Enabled
		sheduleCron := responseBodyUpdateSchedules.Result.Cron
		cd, _ := crondescriptor.NewCronDescriptor(sheduleCron)
		sheduleCronStr, _ := cd.GetDescription(crondescriptor.Full)

		sheduleCreatedEp := int64(responseBodyUpdateSchedules.Result.Created)
		sheduleCreated := time.Unix(sheduleCreatedEp, 0)
		sheduleCreatedStr := fmt.Sprint(sheduleCreated)
		fmt.Printf("\n%-10v %-10t %-20s %-40s\n\n", scheduleTest, sheduleEnabled, sheduleCreatedStr[0:16], *sheduleCronStr)
	} else {
		errorCode := responseBodyUpdateSchedules.Error.Code
		errorMessage := responseBodyUpdateSchedules.Error.Message
		fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
	}
}

func updateScheduleraw(scheduleId string) {
	apiId, apiSecret := Getapikeys()
	status := isEnabledPromt()
	client := &http.Client{}
	data := fmt.Sprintf(`{"enabled":%t}`, status)
	reqBodydata := strings.NewReader(data)
	req, err := http.NewRequest("PATCH", "https://a.blazemeter.com/api/v4/schedules/"+scheduleId, reqBodydata)
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
