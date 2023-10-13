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
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		sessionId, _ := cmd.Flags().GetString("sid")
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
		default:
			fmt.Println("\nPlease provide a correct session ID or mock service Id  to find the logs")
			fmt.Println("[bmgo find schedule --scheduleid <schedule ID>")
		}
	},
}

func init() {
	FindCmd.AddCommand(logsCmd)
	logsCmd.Flags().String("sid", "", "Provide session Id to pull logs for")
	logsCmd.Flags().Int("mockid", 0, "Provide the mock service id")
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
		fmt.Printf("DataUrl: %s\n", dataUrl)
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
				fmt.Println("\nDOWNLOAD LINKS: ", termlink.Link(logsFileName, logsDataUrl))
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
