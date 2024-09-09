/*
Copyright © 2024 Manan Patel - github.com/immnan
*/
package get

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

// testsCmd represents the tests command
var testsCmd = &cobra.Command{
	Use:   "tests",
	Short: "Get list of tests",
	Long: `Use the command to list Tests within a specified workspace. The output includes Test ID, Name, Project Id, etc. The output can be further filtered by specifying a project id by using the --pid flag.

	For example: [bmgo get -w <workspace id> tests] OR
	             [bmgo get -w <workspace id> tests --pid <project id>]
	For default: [bmgo get --ws tests] OR 
	             [bmgo get --ws tests --pid <project id>]`,
	Run: func(cmd *cobra.Command, args []string) {
		ws, _ := cmd.Flags().GetBool("ws")
		var workspaceId int
		if ws {
			workspaceId = defaultWorkspace()
		} else {
			workspaceId, _ = cmd.Flags().GetInt("workspaceid")
		}
		tm, _ := cmd.Flags().GetBool("tm")
		var teamId string
		if tm {
			teamId = defaultTeam()
		} else {
			teamId, _ = cmd.Flags().GetString("teamid")
		}
		rawOutput, _ := cmd.Flags().GetBool("raw")
		projectId, _ := cmd.Flags().GetInt("pid")

		switch {
		case (workspaceId != 0 || projectId != 0):
			listTestsWS(workspaceId, projectId, rawOutput)
		case (workspaceId == 0 && projectId == 0) && teamId != "":
			getTeamBuckets(teamId, rawOutput)
		default:
			cmd.Help()
		}
	},
}

func init() {
	GetCmd.AddCommand(testsCmd)
	testsCmd.Flags().Int("pid", 0, "Provide the project ID to filter the results")
}

type ListTestsResponse struct {
	Result []listTestsResult `json:"result"`
	Error  errorResult       `json:"error"`
}
type listTestsResult struct {
	Name          string             `json:"name"`
	Id            int                `json:"id"`
	LastRunTime   int                `json:"lastRunTime"`
	ProjectId     int                `json:"projectId"`
	Configuration config             `json:"configuration"`
	Executions    []executiondetails `json:"executions"`
}
type config struct {
	DedicatedIpsEnabled     bool `json:"dedicatedIpsEnabled"`
	EnableTestData          bool `json:"enableTestData"`
	EnableLoadConfiguration bool `json:"enableLoadConfiguration"`
}
type executiondetails struct {
	Locations []string `json:"locations"`
}

func listTestsWS(workspaceId, projectId int, rawOutput bool) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	var req *http.Request
	var err error
	if projectId != 0 {
		projectIdStr := strconv.Itoa(projectId)
		req, err = http.NewRequest("GET", "https://a.blazemeter.com/api/v4/tests?projectId="+projectIdStr+"&limit=0", nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		workspaceIdStr := strconv.Itoa(workspaceId)
		req, err = http.NewRequest("GET", "https://a.blazemeter.com/api/v4/tests?workspaceId="+workspaceIdStr+"&limit=0", nil)
		if err != nil {
			log.Fatal(err)
		}
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
		var responseObjectListTests ListTestsResponse
		json.Unmarshal(bodyText, &responseObjectListTests)
		if responseObjectListTests.Error.Code == 0 {
			// Create tabwriter
			tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			// Print headers
			fmt.Fprintln(tabWriter, "TEST ID\tLAST RUN\tPROJECT\tTEST NAME")

			for i := 0; i < len(responseObjectListTests.Result); i++ {
				testName := responseObjectListTests.Result[i].Name
				testId := responseObjectListTests.Result[i].Id
				testLastRunEp1 := responseObjectListTests.Result[i].LastRunTime
				testProjectId := responseObjectListTests.Result[i].ProjectId
				testLastRunEp := int64(responseObjectListTests.Result[i].LastRunTime)
				// This is because there are epoch values as "0", it converts to a time on 1970s, so we want to condition that here:
				if testLastRunEp1 != 0 {
					testLastRun := time.Unix(testLastRunEp, 0)
					testLastRunSp := fmt.Sprint(testLastRun)
					fmt.Fprintf(tabWriter, "%d\t%s\t%d\t%s\n", testId, testLastRunSp[0:16], testProjectId, testName)
				} else {
					fmt.Fprintf(tabWriter, "%d\t%d\t%d\t%s\n", testId, testLastRunEp1, testProjectId, testName)
				}
			}
			tabWriter.Flush()
			fmt.Println("\n-")
		} else {
			errorCode := responseObjectListTests.Error.Code
			errorMessage := responseObjectListTests.Error.Message
			fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
		}
	}
}

func getTeamBuckets(teamId string, rawOutput bool) {
	Bearer := fmt.Sprintf("Bearer %v", GetPersonalAccessToken())
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.runscope.com/teams/"+teamId, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", Bearer)
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
		var responseBodyTeamInfo teamInfo
		json.Unmarshal(bodyText, &responseBodyTeamInfo)
		var wg sync.WaitGroup
		results := make(chan StepDetail)

		for _, bucket := range responseBodyTeamInfo.Buckets {
			wg.Add(1)
			go func(bucketKey string) {
				defer wg.Done()
				listTestsTm(bucketKey, results)
			}(bucket.Key)
		}

		// Close the channel when all goroutines are done
		go func() {
			wg.Wait()
			close(results)
		}()

		// Print results in a table
		printResultsTable(results)
	}
}

type result struct {
	Data []bucketTests `json:"data"`
}
type bucketTests struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type testResult struct {
	Data dataTests `json:"data"`
}
type dataTests struct {
	Steps []steps `json:"steps"`
}
type steps struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

// This is for channel
type StepDetail struct {
	TestID string
	URL    string
}

func listTestsTm(bucketKey string, results chan<- StepDetail) {
	Bearer := fmt.Sprintf("Bearer %v", GetPersonalAccessToken())
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.runscope.com/buckets/"+bucketKey+"/tests", nil)
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Set("Authorization", Bearer)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	var response result
	json.Unmarshal(bodyText, &response)
	var wg sync.WaitGroup

	for _, test := range response.Data {
		wg.Add(1)
		go func(testId string) {
			defer wg.Done()
			findTestTm(bucketKey, testId, results)
		}(test.Id)
	}

	wg.Wait()
}

func findTestTm(bucketKey, testId string, results chan<- StepDetail) {
	Bearer := fmt.Sprintf("Bearer %v", GetPersonalAccessToken())
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.runscope.com/buckets/"+bucketKey+"/tests/"+testId, nil)
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Set("Authorization", Bearer)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	var response testResult
	json.Unmarshal(bodyText, &response)
	for _, step := range response.Data.Steps {
		results <- StepDetail{
			TestID: testId,
			URL:    step.Url,
		}
	}
}

func printResultsTable(results <-chan StepDetail) {
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(writer, "Test ID\tStep URL")
	//fmt.Fprintln(writer, "-------\t--------")

	for result := range results {
		fmt.Fprintf(writer, "%s\t%s\n", result.TestID, result.URL)
	}
	writer.Flush()
}
