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
	"time"

	"github.com/spf13/cobra"
)

// mastersCmd represents the masters command
var mastersCmd = &cobra.Command{
	Use:   "masters",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("masters called")
		testId, _ := cmd.Flags().GetInt("tid")
		rawOutput, _ := cmd.Flags().GetBool("raw")
		if testId != 0 && rawOutput {
			getMastersraw(testId)
		} else if testId != 0 {
			getMasters(testId)
		} else {
			fmt.Println("\nPlease provide a valid Test ID to get list of Masters")
			fmt.Println("[bmgo get masters --tid <test id>")
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

func getMasters(testId int) {
	apiId, apiSecret := Getapikeys()
	testIdStr := strconv.Itoa(testId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/masters?testId="+testIdStr+"&limit=25", nil)
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
	var responseObjectMasters mastersResponse
	json.Unmarshal(bodyText, &responseObjectMasters)
	if responseObjectMasters.Error.Code == 0 {
		fmt.Printf("\n%-10s %-8s %-10s %-10s %-15s %-15s\n", "MASTER", "STATUS", "RUN_BY", "PROJECT", "START", "END")
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
				fmt.Printf("\n%-10d %-8s %-10d %-10d %-15s %-15s", masterId, masterStatus, masterRunner, masterProject, masterCreatedStr[2:16], masterEndStr[5:16])
			} else {
				fmt.Printf("\n%-10d %-8s %-10d %-10d %-15d %-15d", masterId, masterStatus, masterRunner, masterProject, masterCreatedEp, masterEndEp)
			}
		}
		fmt.Println("\n\n---------------------------------------------------------------------------------------------")

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
func getMastersraw(testId int) {
	apiId, apiSecret := Getapikeys()
	testIdStr := strconv.Itoa(testId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/masters?testId="+testIdStr+"&limit=25", nil)
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
