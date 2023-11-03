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

	"github.com/savioxavier/termlink"
	"github.com/spf13/cobra"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Find and download logs",
	Long: `
	For example: [bmgo find logs --sid <Session ID>] OR
	             [bmgo find logs --mockid <Mock service ID>]`,
	Run: func(cmd *cobra.Command, args []string) {
		sessionId, _ := cmd.Flags().GetString("sid")
		masterId, _ := cmd.Flags().GetInt("mid")
		rawOutput, _ := cmd.Flags().GetBool("raw")
		mockId, _ := cmd.Flags().GetInt("mockid")
		switch {
		case sessionId != "" && mockId == 0 && rawOutput:
			findlogSessionraw(sessionId)
		case sessionId != "" && mockId == 0:
			findlogSession(sessionId)
		case sessionId == "" && mockId != 0 && rawOutput:
			findlogMockservice(mockId)
		case sessionId == "" && mockId != 0:
			findlogMockservice(mockId)
		case sessionId == "" && mockId == 0 && masterId != 0:
			findMasterlogs(masterId)
		default:
			cmd.Help()
		}
	},
}

func init() {
	FindCmd.AddCommand(logsCmd)
	logsCmd.Flags().String("sid", "", "Provide session Id to pull logs for")
	logsCmd.Flags().Int("mockid", 0, "Provide the mock service id")
	logsCmd.Flags().Int("mid", 0, "Provide the master Id to get logs for all sessions")
}

type findLogsResponse struct {
	Result logsResult  `json:"result"`
	Error  errorResult `json:"error"`
}
type logsResult struct {
	DataUrl    string     `json:"dataUrl"`
	Data       []dataList `json:"data"`
	LocationId string     `json:"locationId"`
}
type dataList struct {
	Filename string `json:"filename"`
	DataUrl  string `json:"dataUrl"`
}

func findlogSession(sessionId string) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/sessions/"+sessionId+"/reports/logs", nil)
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
	//	fmt.Printf("%s\n", bodyText)
	var responseObjectLogs findLogsResponse
	json.Unmarshal(bodyText, &responseObjectLogs)
	if responseObjectLogs.Error.Code == 0 {
		dataUrl := responseObjectLogs.Result.DataUrl
		fmt.Printf("DataUrl: %s\n\n", dataUrl)
		var logsFileName, logsDataUrl string
		for i := 0; i < len(responseObjectLogs.Result.Data); i++ {

			if !termlink.SupportsHyperlinks() {
				//	fmt.Printf("\n%-20s %-20s", logsFileName, logsDataUrl)
				logsFileName = responseObjectLogs.Result.Data[i].Filename
				logsDataUrl = responseObjectLogs.Result.Data[i].DataUrl
				fmt.Printf("\nFile name: %s\n", logsFileName)
				fmt.Println(logsDataUrl)
			} else {
				logsFileName = responseObjectLogs.Result.Data[i].Filename
				logsDataUrl = responseObjectLogs.Result.Data[i].DataUrl
				fmt.Println("DOWNLOAD: ", termlink.Link(logsFileName, logsDataUrl))
			}
		}
		fmt.Println("\n-")
	} else {
		errorCode := responseObjectLogs.Error.Code
		errorMessage := responseObjectLogs.Error.Message
		fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
	}
}
func findlogSessionraw(sessionId string) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/sessions/"+sessionId+"/reports/logs", nil)
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

func findlogMockservice(mockId int) {
	workspaceIdStr := workspaceIdPrompt()
	mockIdStr := strconv.Itoa(mockId)
	fmt.Println(workspaceIdStr, mockIdStr)
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://mock.blazemeter.com/api/v1/workspaces/1294894/service-mocks/119057/log?limit=50", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("accept", "application/json;charset=UTF-8")
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

func findMasterlogs(masterId int) {
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
		totalLocations := []string{}
		for l := 0; l < len(responseObjectMaster.Result.Locations); l++ {
			locations := responseObjectMaster.Result.Locations[l]
			totalLocations = append(totalLocations, locations)
		}
		fmt.Printf("MASTER:    %d\nLOCATIONS: %s\n", masterId, totalLocations)
		//totalSessions := []string{}
		for rv := 0; rv < len(responseObjectMaster.Result.SessionId); rv++ {
			sessionsId := responseObjectMaster.Result.SessionId[rv]
			fmt.Printf("SESSION ID [%d]: %s\n", rv, sessionsId)
			findlogSession(sessionsId)
		}
	} else {
		errorCode := responseObjectMaster.Error.Code
		errorMessage := responseObjectMaster.Error.Message
		fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
	}
}
