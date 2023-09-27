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
	"strconv"
	"strings"
	"time"

	"github.com/jsuar/go-cron-descriptor/pkg/crondescriptor"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// scheduleCmd represents the schedule command
var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Enable or Disable Schedule",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		scheduleId, _ := cmd.Flags().GetString("sid")
		rawOutput, _ := cmd.Flags().GetBool("raw")
		if scheduleId != "" && rawOutput {
			updateScheduleraw(scheduleId)
		} else if scheduleId != "" {
			updateSchedule(scheduleId)
		} else {
			fmt.Println("\nPlease provide a correct Schedule ID to update the Schedule")
			fmt.Println("[bmgo update schedule --sid <schedule ID>")
		}
	},
}

func init() {
	UpdateCmd.AddCommand(scheduleCmd)
	scheduleCmd.Flags().String("sid", "", "Provide the schedule ID to modify")
}

func isEnabledPromt() bool {
	prompt1 := promptui.Select{
		Label:        "Enable:",
		Items:        []bool{true, false},
		HideSelected: true,
	}
	_, attachAuto, err := prompt1.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}
	boolVal, _ := strconv.ParseBool(attachAuto)
	return boolVal
}

type updateshedulesResponse struct {
	Result updatescheduleResult `json:"result"`
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
	var responseBodyUpdateShedules updateshedulesResponse
	json.Unmarshal(bodyText, &responseBodyUpdateShedules)
	fmt.Printf("\n%-10s %-10s %-50s %-20s", "TEST", "ENABLED", "CRON", "CREATED ON")
	scheduleTest := responseBodyUpdateShedules.Result.TestId
	sheduleEnabled := responseBodyUpdateShedules.Result.Enabled
	sheduleCron := responseBodyUpdateShedules.Result.Cron
	cd, _ := crondescriptor.NewCronDescriptor(sheduleCron)
	sheduleCronStr, _ := cd.GetDescription(crondescriptor.Full)

	sheduleCreatedEp := int64(responseBodyUpdateShedules.Result.Created)
	sheduleCreated := time.Unix(sheduleCreatedEp, 0)
	sheduleCreatedStr := fmt.Sprint(sheduleCreated)
	fmt.Printf("\n%-10v %-10t %-50s %-20s\n\n", scheduleTest, sheduleEnabled, *sheduleCronStr, sheduleCreatedStr[0:16])
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
