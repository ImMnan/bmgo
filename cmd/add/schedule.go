/*
Copyright © 2024 Manan Patel - github.com/immnan
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
	Long: `The test scheduler command will allow you to assign a cron schedule to your test to have it run on a schedule. You can read more about this feature at help.blazemeter.com. Add a schedule to your test using this command with the help of a cron expression. You will be prompted to provide the cron expression when the command is run.
	
	For example: [bmgo add schedule --tid <test id>]`,
	Run: func(cmd *cobra.Command, args []string) {
		testId, _ := cmd.Flags().GetInt("tid")
		rawOutput, _ := cmd.Flags().GetBool("raw")
		if testId != 0 {
			addSchedule(testId, rawOutput)
		} else {
			cmd.Help()
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

func addSchedule(testId int, rawOutput bool) {
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
	if rawOutput {
		fmt.Printf("%s\n", bodyText)
	} else {
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
}
