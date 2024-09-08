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
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

// projectsCmd represents the projects command
var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Get the list of Projects under workspace",
	Long: `Use the command to list Projects within a specified workspace or account. Projects are designed to organize tests and reports and track usage within a Workspace. The output includes Project NAME, ID, Test count, etc.

	For example: [bmgo get -w <workspace id> projects] OR 
	             [bmgo get -a <account id> projects]
	For default: [bmgo get --ac projects] OR 
	             [bmgo get --ws projects]`,
	Run: func(cmd *cobra.Command, args []string) {
		ac, _ := cmd.Flags().GetBool("ac")
		ws, _ := cmd.Flags().GetBool("ws")
		var accountId, workspaceId int
		if ac {
			accountId = defaultAccount()
		} else {
			accountId, _ = cmd.Flags().GetInt("accountid")
		}
		if ws {
			workspaceId = defaultWorkspace()
		} else {
			workspaceId, _ = cmd.Flags().GetInt("workspaceid")
		}
		rawOutput, _ := cmd.Flags().GetBool("raw")
		switch {
		case workspaceId == 0 && accountId != 0:
			getProjectsA(accountId, rawOutput)
		case workspaceId != 0 && accountId == 0:
			getProjectsWs(workspaceId, rawOutput)
		default:
			cmd.Help()
		}
	},
}

func init() {
	GetCmd.AddCommand(projectsCmd)
}

type results struct {
	Result []projectsResult `json:"result"`
	Error  errorResult      `json:"error"`
}
type projectsResult struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	TestsCount int    `json:"testsCount"`
	Created    int    `json:"created"`
}

func getProjectsWs(workspaceId int, rawOutput bool) {
	apiId, apiSecret := Getapikeys()
	workspaceIdstr := strconv.Itoa(workspaceId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/projects?workspaceId="+workspaceIdstr+"&limit=0", nil)
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
		var response results
		json.Unmarshal(bodyText, &response)
		if response.Error.Code == 0 {
			tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(tabWriter, "ID\tNAME\tTESTS\tCREATED")
			for i := 0; i < len(response.Result); i++ {
				projectId := response.Result[i].Id
				projectName := response.Result[i].Name
				projectTests := response.Result[i].TestsCount
				pCreatedEpoch := int64(response.Result[i].Created)
				projectCreated := fmt.Sprint(time.Unix(pCreatedEpoch, 0))
				fmt.Fprintf(tabWriter, "%d\t%s\t%d\t%s\n", projectId, projectName, projectTests, projectCreated[0:10])
			}
			tabWriter.Flush()
			fmt.Println("-")
		} else {
			errorCode := response.Error.Code
			errorMessage := response.Error.Message
			fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
		}
	}
}

func getProjectsA(accountId int, rawOutput bool) {
	apiId, apiSecret := Getapikeys()
	accountIdstr := strconv.Itoa(accountId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/projects?accountId="+accountIdstr+"&limit=0", nil)
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
		var response results
		json.Unmarshal(bodyText, &response)
		if response.Error.Code == 0 {
			tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(tabWriter, "ID\tNAME\tTESTS\tCREATED")
			for i := 0; i < len(response.Result); i++ {
				projectId := response.Result[i].Id
				projectName := response.Result[i].Name
				projectTests := response.Result[i].TestsCount
				pCreatedEpoch := int64(response.Result[i].Created)
				projectCreated := fmt.Sprint(time.Unix(pCreatedEpoch, 0))
				fmt.Fprintf(tabWriter, "%d\t%s\t%d\t%s\n", projectId, projectName, projectTests, projectCreated[0:10])
			}
			tabWriter.Flush()
			fmt.Println("-")
		} else {
			errorCode := response.Error.Code
			errorMessage := response.Error.Message
			fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
		}
	}
}
