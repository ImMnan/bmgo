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
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Find Project using Project ID",
	Long: ` Use this command to find details about the specified project (--pid). Global Flag --raw can be used for raw Json output. 
	For example: [bmgo find project --pid <Project id>]`,
	Run: func(cmd *cobra.Command, args []string) {
		projectId, _ := cmd.Flags().GetInt("pid")
		rawOutput, _ := cmd.Flags().GetBool("raw")

		if projectId != 0 {
			var wg sync.WaitGroup
			wg.Add(1)
			go listTestsProject(projectId, &wg)
			wg.Wait()
			findProject(projectId, rawOutput)
		} else {
			cmd.Help()
		}
	},
}

func init() {
	FindCmd.AddCommand(projectCmd)
	projectCmd.Flags().Int("pid", 0, "Provide the project id to find")
	projectCmd.MarkFlagRequired("pid")
}

type ProjectResponse struct {
	Result projectResult `json:"result"`
	Error  errorResult   `json:"error"`
}
type projectResult struct {
	Name        string `json:"name"`
	WorkspaceId int    `json:"workspaceId"`
	Created     int    `json:"created"`
}
type ListTestsResponse struct {
	Result []listTestsResult `json:"result"`
	Error  errorResult       `json:"error"`
}
type listTestsResult struct {
	Name        string `json:"name"`
	Id          int    `json:"id"`
	LastRunTime int    `json:"lastRunTime"`
}

func findProject(projectId int, rawOutput bool) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	projectIdStr := strconv.Itoa(projectId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/projects/"+projectIdStr, nil)
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
		var responseObjectProject ProjectResponse
		json.Unmarshal(bodyText, &responseObjectProject)
		if responseObjectProject.Error.Code == 0 {
			//	fmt.Printf("%-25s %-10s %-15s\n", "PROJECT NAME", "WORKSPACE", "CREATED")
			tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			// Print headers
			fmt.Fprintln(tabWriter, "PROJECT_NAME\tWORKSPACE\tCREATED")
			projectName := responseObjectProject.Result.Name
			projectWorkspace := responseObjectProject.Result.WorkspaceId
			projectCreatedEp := int64(responseObjectProject.Result.Created)
			projectCreatedStr := fmt.Sprint(time.Unix(projectCreatedEp, 0))
			//	fmt.Printf("%-25s %-10v %-15v", projectName, projectWorkspace, projectCreatedStr[0:16])
			fmt.Fprintf(tabWriter, "%s\t%d\t%s\n", projectName, projectWorkspace, projectCreatedStr[0:16])
			tabWriter.Flush()
			fmt.Println("-")
		} else {
			errorCode := responseObjectProject.Error.Code
			errorMessage := responseObjectProject.Error.Message
			fmt.Printf("Error code: %v\nError Message: %v\n\n", errorCode, errorMessage)
		}
	}
}

func listTestsProject(projectId int, wg *sync.WaitGroup) {
	defer wg.Done()
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	projectIdStr := strconv.Itoa(projectId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/tests?projectId="+projectIdStr+"&limit=0", nil)
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
	var responseObjectListTests ListTestsResponse
	json.Unmarshal(bodyText, &responseObjectListTests)
	if responseObjectListTests.Error.Code == 0 {
		//	fmt.Printf("\n%-10s %-20s %-15s\n", "TEST ID", "LAST RUN", "TEST NAME")
		tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		// Print headers
		fmt.Fprintln(tabWriter, "TEST_ID\tLAST_RUN\tTEST_NAME")
		for i := 0; i < len(responseObjectListTests.Result); i++ {
			testName := responseObjectListTests.Result[i].Name
			testId := responseObjectListTests.Result[i].Id
			testLastRunEp1 := responseObjectListTests.Result[i].LastRunTime
			testLastRunEp := int64(responseObjectListTests.Result[i].LastRunTime)
			if testLastRunEp1 != 0 {
				testLastRun := time.Unix(testLastRunEp, 0)
				testLastRunSp := fmt.Sprint(testLastRun)
				//fmt.Printf("\n%-10v %-20s %-15s", testId, testLastRunSp[0:16], testName)
				fmt.Fprintf(tabWriter, "%d\t%s\t%s\n", testId, testLastRunSp[0:16], testName)
			} else {
				testLastRun := testLastRunEp1
				//	fmt.Printf("\n%-10v %-20v %-15s", testId, testLastRun, testName)
				fmt.Fprintf(tabWriter, "%d\t%d\t%s\n", testId, testLastRun, testName)
			}
		}
		tabWriter.Flush()
		fmt.Println("-")
	} else {
		errorCode := responseObjectListTests.Error.Code
		errorMessage := responseObjectListTests.Error.Message
		fmt.Printf("Error code: %v\nError Message: %v\n\n", errorCode, errorMessage)
	}
}
