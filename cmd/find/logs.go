/*
Copyright © 2024 Manan Patel - github.com/immnan
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
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Find and download logs",
	Long: ` Use this command to find the logs for a specified session or list of sessions part of a master. Use --download or -d flag can be used to download the artifacts.zip automatically. Similarly, --raw -r flag is available too, however, it is only applicable for listing logs for a perticular session.
	For example: [bmgo find logs --sid <Session ID>] OR 
	             [bmgo find logs --sid <Session ID> --download]
	             [bmgo find logs --mid <Master ID>] OR  
				 [bmgo find logs --mid <Master ID> --download]

	             [bmgo find logs --mockid <Mock service ID>]`,
	Run: func(cmd *cobra.Command, args []string) {
		sessionId, _ := cmd.Flags().GetString("sid")
		masterId, _ := cmd.Flags().GetInt("mid")
		rawOutput, _ := cmd.Flags().GetBool("raw")
		mockId, _ := cmd.Flags().GetInt("mockid")
		download, _ := cmd.Flags().GetBool("download")

		switch {
		case sessionId != "" && mockId == 0 && masterId == 0:
			findlogSession(sessionId, download, rawOutput)
		case sessionId == "" && mockId != 0 && masterId == 0:
			findlogMockservice(mockId)
		case sessionId == "" && mockId == 0 && masterId != 0:
			var wg sync.WaitGroup
			masterChan := make(chan string, 100)
			wg.Add(1)
			go func() {
				defer wg.Done()
				findMasterlogsConc(masterId, masterChan)
				close(masterChan)
			}()
			for sessionId := range masterChan {
				wg.Add(1)
				go findlogSessionConcOutput(sessionId, download, &wg)
			}
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

type mastersResponseLogs struct {
	Result mastersResultLogs `json:"result"`
	Error  errorResult       `json:"error"`
}

type mastersResultLogs struct {
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
	TestId      int    `json:"testId"`
}

func findlogSession(sessionId string, download, rawOutput bool) {
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
	if rawOutput {
		fmt.Printf("%s\n", bodyText)
	} else {
		var responseObjectLogs findLogsResponse
		json.Unmarshal(bodyText, &responseObjectLogs)
		if responseObjectLogs.Error.Code == 0 {
			var logsFileName, logsDataUrl string
			var wgDownload sync.WaitGroup
			for i := 0; i < len(responseObjectLogs.Result.Data); i++ {
				if download {
					fileName := responseObjectLogs.Result.Data[i].Filename
					fileURL := responseObjectLogs.Result.Data[i].DataUrl
					if fileName == "jmeter.log" || fileName == "bzt.log" {
						continue
					} else {
						downloadFileSessions(fileURL, fileName, sessionId)
					}
				} else {
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
			}
			wgDownload.Wait()
		} else {
			errorCode := responseObjectLogs.Error.Code
			errorMessage := responseObjectLogs.Error.Message
			fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
		}
	}
}

// This function is used to download the logs for a perticular mock service
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

// This function is used to download or view the logs for all sessions within a master
func findMasterlogsConc(masterId int, ch chan<- string) {
	apiId, apiSecret := Getapikeys()
	masterIdStr := strconv.Itoa(masterId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/masters/"+masterIdStr, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}
	req.SetBasicAuth(apiId, apiSecret)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return
	}
	var responseObjectMaster mastersResponseLogs
	err = json.Unmarshal(bodyText, &responseObjectMaster)
	if err != nil {
		log.Printf("Error unmarshalling response: %v", err)
		return
	}
	if responseObjectMaster.Error.Code == 0 {
		for _, sessionId := range responseObjectMaster.Result.SessionId {
			ch <- sessionId
		}
	} else {
		log.Printf("Error code: %v, Error message: %v", responseObjectMaster.Error.Code, responseObjectMaster.Error.Message)
	}
}

// This function will download the files from the sessionId and print the logs depending on the download flag
func findlogSessionConcOutput(sessionId string, download bool, wg *sync.WaitGroup) {
	defer wg.Done()

	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/sessions/"+sessionId+"/reports/logs", nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}
	req.SetBasicAuth(apiId, apiSecret)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return
	}
	var responseObjectLogs findLogsResponse
	err = json.Unmarshal(bodyText, &responseObjectLogs)
	if err != nil {
		log.Printf("Error unmarshalling response: %v", err)
		return
	}
	if responseObjectLogs.Error.Code == 0 {
		if download {
			var wgDownload sync.WaitGroup
			for _, data := range responseObjectLogs.Result.Data {
				if data.Filename == "jmeter.log" || data.Filename == "bzt.log" {
					continue
				}
				wgDownload.Add(1)
				go func(url, filename, sid string) {
					defer wgDownload.Done()
					downloadFileSessions(url, filename, sid)
				}(data.DataUrl, data.Filename, sessionId)
			}
			wgDownload.Wait()
		} else {
			for _, data := range responseObjectLogs.Result.Data {
				if !termlink.SupportsHyperlinks() {
					fmt.Printf("\nFile name: %s\n%v", data.Filename, data.DataUrl)
				} else {
					fmt.Printf("Session: %v: %v\n", sessionId, termlink.Link(data.Filename, data.DataUrl))
				}
			}
		}
	} else {
		log.Printf("Error code: %v, Error message: %v", responseObjectLogs.Error.Code, responseObjectLogs.Error.Message)
	}
}

// This function is to download the files from the sessions of the provided master Id. It will download all the files except jmeter.log and bzt.log.
func downloadFileSessions(fileURL, fileName, sessionId string) {
	file, err := os.Create(sessionId + fileName)
	if err != nil {
		log.Printf("Error creating file: %v", err)
		return
	}
	defer file.Close()

	client := &http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	resp, err := client.Get(fileURL)
	if err != nil {
		log.Printf("Error downloading file: %v", err)
		return
	}
	defer resp.Body.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"downloading",
	)
	_, err = io.Copy(io.MultiWriter(file, bar), resp.Body)
	if err != nil {
		log.Printf("Error saving file: %v", err)
		return
	}
	fmt.Printf("Downloaded %s for session %s\n", fileName, sessionId)
}
