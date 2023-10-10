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
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// masterCmd represents the master command
var masterCmd = &cobra.Command{
	Use:   "master",
	Short: "Find Master details",
	Long:  `.`,
	Run: func(cmd *cobra.Command, args []string) {
		masterId, _ := cmd.Flags().GetInt("mid")
		rawOutput, _ := cmd.Flags().GetBool("raw")
		if masterId != 0 && rawOutput {
			findMasterraw(masterId)
		} else if masterId != 0 {
			findMaster(masterId)
		} else {
			fmt.Println("\nPlease provide a valid master ID to find the master")
			fmt.Println("[bmgo find master --mid <Master id>")
		}
	},
}

func init() {
	FindCmd.AddCommand(masterCmd)
	masterCmd.Flags().Int("mid", 0, "Provide the Master Id")
}

type mastersResponse struct {
	Result mastersResult `json:"result"`
	Error  errorResult   `json:"error"`
}

type mastersResult struct {
	Id           int                `json:"id"`
	Status       string             `json:"reportStatus"`
	Created      int                `json:"created"`
	Ended        int                `json:"ended"`
	Locations    []string           `json:"locations"`
	SessionId    []string           `json:"sessionsId"`
	ProjectId    int                `json:"projectId"`
	RunnerUserId int                `json:"runnerUserId"`
	Executions   []masterExecutions `json:"executions"`
	TestId       int                `json:"testId"`
}
type masterExecutions struct {
	Concurrency int    `json:"concurrency"`
	HoldFor     string `json:"holdFor"`
	Rampup      string `json:"rampUp"`
	Executor    string `json:"executor"`
}

func findMaster(masterId int) {
	apiId, apiSecret := Getapikeys()
	masterIdStr := strconv.Itoa(masterId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/masters/"+masterIdStr, nil)
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
	var responseObjectMaster mastersResponse
	json.Unmarshal(bodyText, &responseObjectMaster)
	if responseObjectMaster.Error.Code == 0 {
		fmt.Printf("\n%-10s %-8s %-10s %-10s %-15s %-15s ", "TEST", "STATUS", "RUN_BY", "PROJECT", "START", "END")
		masterId := responseObjectMaster.Result.Id
		masterTestId := responseObjectMaster.Result.TestId
		masterStatus := responseObjectMaster.Result.Status
		masterProject := responseObjectMaster.Result.ProjectId
		masterCreatedEp := int64(responseObjectMaster.Result.Created)
		masterEndEp := int64(responseObjectMaster.Result.Ended)
		masterRunner := responseObjectMaster.Result.RunnerUserId
		if masterCreatedEp != 0 && masterEndEp != 0 {
			masterCreatedStr := fmt.Sprint(time.Unix(masterCreatedEp, 0))
			masterEndStr := fmt.Sprint(time.Unix(masterEndEp, 0))
			fmt.Printf("\n%-10d %-8s %-10d %-10d %-15s %-15s", masterTestId, masterStatus, masterRunner, masterProject, masterCreatedStr[2:16], masterEndStr[5:16])
		} else {
			fmt.Printf("\n%-10d %-8s %-10d %-10d %-15d %-15d", masterTestId, masterStatus, masterRunner, masterProject, masterCreatedEp, masterEndEp)
		}
		fmt.Println("\n\n---------------------------------------------------------------------------------------------")
		fmt.Printf("%-15s %-10s %-10s %-10s", "EXECUTOR", "VUs", "RAMP_UP", "HOLD_FOR")

		for e := 0; e < len(responseObjectMaster.Result.Executions); e++ {
			masterConcurrency := responseObjectMaster.Result.Executions[e].Concurrency
			masterExecutor := responseObjectMaster.Result.Executions[e].Executor
			masterRampUp := responseObjectMaster.Result.Executions[e].Rampup
			masterHoldFor := responseObjectMaster.Result.Executions[e].HoldFor
			fmt.Printf("\n%-15s %-10d %-10s %-10s", masterExecutor, masterConcurrency, masterRampUp, masterHoldFor)
		}
		fmt.Println("\n\n---------------------------------------------------------------------------------------------")
		totalLocations := []string{}
		for l := 0; l < len(responseObjectMaster.Result.Locations); l++ {
			locations := responseObjectMaster.Result.Locations[l]
			totalLocations = append(totalLocations, locations)
		}
		totalSessions := []string{}
		for rv := 0; rv < len(responseObjectMaster.Result.SessionId); rv++ {
			sessions := responseObjectMaster.Result.SessionId[rv]
			totalSessions = append(totalSessions, sessions)
		}
		fmt.Printf("MASTER:    %d\nLOCATIONS: %s\nSESSIONS:  %s\n\n", masterId, totalLocations, totalSessions)

	} else {
		errorCode := responseObjectMaster.Error.Code
		errorMessage := responseObjectMaster.Error.Message
		fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
	}
}
func findMasterraw(masterId int) {
	apiId, apiSecret := Getapikeys()
	masterIdStr := strconv.Itoa(masterId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/masters/"+masterIdStr, nil)
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
