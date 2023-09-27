/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package add

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jsuar/go-cron-descriptor/pkg/crondescriptor"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// scheduleCmd represents the schedule command
var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Add new schedule",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		testId, _ := cmd.Flags().GetInt("tid")
		//	rawOutput, _ := cmd.Flags().GetBool("raw")
		addSchedule(testId)
	},
}

func init() {
	AddCmd.AddCommand(scheduleCmd)
	scheduleCmd.Flags().Int("tid", 0, "Provide a test id to create a schedule for")
}

func cronPrompt() string {
	validate := func(input string) error {
		if len(input) <= 8 {
			return errors.New("invalid crone")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:       "Cron Expression: ",
		HideEntered: true,
		Validate:    validate,
	}
	resultCronEx, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}
	return resultCronEx
}

type addShedulesResponse struct {
	Result addScheduleResult `json:"result"`
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
	json.Unmarshal(bodyText, &responseBodyAddShedules)
	fmt.Printf("\n%-28s %-10s %-50s %-15s", "SCHEDULE ID", "ENABLED", "CRON", "NEXT RUN")
	scheduleId := responseBodyAddShedules.Result.Id
	//sheduleOwn := responseBodyAddShedules.Result.CreatedById
	sheduleCron := responseBodyAddShedules.Result.Cron
	cd, _ := crondescriptor.NewCronDescriptor(sheduleCron)
	sheduleCronStr, _ := cd.GetDescription(crondescriptor.Full)
	scheduleEnabled := responseBodyAddShedules.Result.Enabled
	sheduleNextrunEP := int64(responseBodyAddShedules.Result.Created)
	sheduleNextRun := time.Unix(sheduleNextrunEP, 0)
	sheduleNextRunStr := fmt.Sprint(sheduleNextRun)
	fmt.Printf("\n%-28s %-10t %-50s %-15s\n\n", scheduleId, scheduleEnabled, *sheduleCronStr, sheduleNextRunStr[0:16])
}
