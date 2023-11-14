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
	"os"
	"strconv"
	"sync"

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
		download, _ := cmd.Flags().GetBool("download")

		switch {
		case sessionId != "" && mockId == 0 && rawOutput:
			findlogSessionraw(sessionId)
		case sessionId != "" && mockId == 0 && download:
			findlogSessionraw(sessionId)
		case sessionId != "" && mockId == 0:
			findlogSession(sessionId, download)
		case sessionId == "" && mockId != 0 && rawOutput:
			findlogMockservice(mockId)
		case sessionId == "" && mockId != 0:
			findlogMockservice(mockId)
		case sessionId == "" && mockId == 0 && masterId != 0:
			var wg sync.WaitGroup
			masterChan := make(chan string, 10)
			wg.Add(2)
			go findlogSessionConc(masterChan, &wg, download)
			findMasterlogsConc(masterId, masterChan, &wg)
			close(masterChan)
			wg.Wait()
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
	logsCmd.Flags().BoolP("download", "d", false, "Confirm if you would like to download automatically")
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

func downloadFileSessions(fileURL, fileName, sessionId string, wgDownload *sync.WaitGroup) {
	//Create blank file
	file, err := os.Create(sessionId + fileName)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	// Put content on file
	resp, err := client.Get(fileURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	size, _ := io.Copy(file, resp.Body)
	defer file.Close()
	fmt.Printf("Downloaded %s with size %d\n", fileName, size)
	defer wgDownload.Done()
}
func findlogSession(sessionId string, download bool) {
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
		var logsFileName, logsDataUrl string
		var wgDownload sync.WaitGroup
		for i := 0; i < len(responseObjectLogs.Result.Data); i++ {
			if download {
				fileName := responseObjectLogs.Result.Data[i].Filename
				fileURL := responseObjectLogs.Result.Data[i].DataUrl
				if fileName == "jmeter.log" || fileName == "bzt.log" {
					continue
				} else {
					wgDownload.Add(1)
					go downloadFileSessions(fileURL, fileName, sessionId, &wgDownload)
				}
			} else {
				fmt.Printf("DataUrl: %s\n\n", dataUrl)
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
				fmt.Println("\n-")
			}
		}
		wgDownload.Wait()
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

func findMasterlogsConc(masterId int, ch chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
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
		for rv := 0; rv < len(responseObjectMaster.Result.SessionId); rv++ {
			sessionsId := responseObjectMaster.Result.SessionId[rv]
			ch <- sessionsId
		}
	} else {
		errorCode := responseObjectMaster.Error.Code
		errorMessage := responseObjectMaster.Error.Message
		fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
	}
}
func findlogSessionConc(ch <-chan string, wg *sync.WaitGroup, download bool) {
	defer wg.Done()
	//apiId, apiSecret := Getapikeys()
	//client := &http.Client{}
	var wgSession sync.WaitGroup
	for sessionId := range ch {
		wgSession.Add(1)
		go findlogSessionConcOutput(sessionId, download, &wgSession)
	}
	wgSession.Wait()
}
func findlogSessionConcOutput(sessionId string, download bool, wgSession *sync.WaitGroup) {
	defer wgSession.Done()
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
		var logsFileName, logsDataUrl string
		var wgDownload sync.WaitGroup
		for i := 0; i < len(responseObjectLogs.Result.Data); i++ {
			if download {
				fileName := responseObjectLogs.Result.Data[i].Filename
				fileURL := responseObjectLogs.Result.Data[i].DataUrl
				if fileName == "jmeter.log" || fileName == "bzt.log" {
					continue
				} else {
					wgDownload.Add(1)
					go downloadFileSessions(fileURL, fileName, sessionId, &wgDownload)
				}
			} else {
				fmt.Printf("DataUrl: %s\n\n", dataUrl)
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
				fmt.Println("\n-")
			}
		}
		wgDownload.Wait()
	} else {
		errorCode := responseObjectLogs.Error.Code
		errorMessage := responseObjectLogs.Error.Message
		fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
	}
}
