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
	"time"

	"github.com/jsuar/go-cron-descriptor/pkg/crondescriptor"
	"github.com/spf13/cobra"
)

// userCmd represents the user command
var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Find details about the specific schedule",
	Long:  `.`,
	Run: func(cmd *cobra.Command, args []string) {
		scheduleId, _ := cmd.Flags().GetString("sid")
		rawOutput, _ := cmd.Flags().GetBool("raw")
		if scheduleId != "" && rawOutput {
			findScheduleraw(scheduleId)
		} else if scheduleId != "" {
			findSchedule(scheduleId)
		} else {
			fmt.Println("\nPlease provide a correct Schedule ID to find the Schedule")
			fmt.Println("[bmgo find schedule --sid <schedule ID>")
		}
	},
}

func init() {
	FindCmd.AddCommand(scheduleCmd)
	scheduleCmd.Flags().String("sid", "", "Provide the Schedule ID")
	scheduleCmd.MarkFlagRequired("sid")
}

type findshedulesResponse struct {
	Result findscheduleResult `json:"result"`
}
type findscheduleResult struct {
	TestId         int    `json:"testId"`
	NextExecutions []int  `json:"nextExecutions"`
	Cron           string `json:"cron"`
	CreatedById    int    `json:"createdById"`
	Created        int    `json:"created"`
}

func findSchedule(scheduleId string) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/schedules/"+scheduleId, nil)
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
	//fmt.Printf("%s\n", bodyText)
	var responseBodyFindShedules findshedulesResponse
	json.Unmarshal(bodyText, &responseBodyFindShedules)
	fmt.Printf("\n%-10s %-10s %-50s %-20s", "TEST", "OWNER", "CRON", "CREATED ON")
	scheduleTest := responseBodyFindShedules.Result.TestId
	sheduleOwn := responseBodyFindShedules.Result.CreatedById
	sheduleCron := responseBodyFindShedules.Result.Cron
	cd, _ := crondescriptor.NewCronDescriptor(sheduleCron)
	sheduleCronStr, _ := cd.GetDescription(crondescriptor.Full)

	sheduleCreatedEp := int64(responseBodyFindShedules.Result.Created)
	sheduleCreated := time.Unix(sheduleCreatedEp, 0)
	sheduleCreatedStr := fmt.Sprint(sheduleCreated)
	fmt.Printf("\n%-10v %-10v %-50s %-20s\n", scheduleTest, sheduleOwn, *sheduleCronStr, sheduleCreatedStr[0:16])
	fmt.Println("\n---------------------------------------------------------------------------------------------")
	fmt.Println("List of upcomming test runs\n-")
	for i := 0; i < len(responseBodyFindShedules.Result.NextExecutions); i++ {
		nextRunsEp := int64(responseBodyFindShedules.Result.NextExecutions[i])
		nextRun := time.Unix(nextRunsEp, 0)
		fmt.Println(nextRun)
	}
	fmt.Println("\n-")
}
func findScheduleraw(scheduleId string) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/schedules/"+scheduleId, nil)
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
