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

	"github.com/spf13/cobra"
)

// oplsCmd represents the opls command
var oplsCmd = &cobra.Command{
	Use:   "opls",
	Short: "Get a list of Private locations in the account",
	Long: `Use the command to list Private locations within a specified workspace or account. Private locations are the on-premise solution when you need to test applications or create Mock Services behind a firewall. The output includes Private location NAME, ID, FUNCTIONALITIES, Agents, etc.

	For example: [bmgo get -w <workspace id> opls] OR 
	             [bmgo get -a <account id> opls]
	For default: [bmgo get --ws opls] OR 
	             [bmgo get --ac opls]`,
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
		if (workspaceId != 0) && (accountId == 0) {
			getOplsWS(workspaceId, rawOutput)
		} else if (accountId != 0) && (workspaceId == 0) {
			getOpls(accountId, rawOutput)
		} else {
			cmd.Help()
		}
	},
}

func init() {
	GetCmd.AddCommand(oplsCmd)
}

type oplsResponse struct {
	Result []oplsResult `json:"result"`
	Error  errorResult  `json:"error"`
}

type oplsResult struct {
	Id               string   `json:"id"`
	Name             string   `json:"name"`
	ThreadsPerEngine int      `json:"threadsPerEngine"`
	Slots            int      `json:"slots"`
	FuncIds          []string `json:"funcIds"`
	ShipsId          []string `json:"shipsId"`
	Ships            []ships
	WorkspacesId     []int `json:"workspacesId"`
}
type ships struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	State         string `json:"state"`
	LastHeartBeat int    `json:"lastHeartBeat"`
}

func getOpls(accountId int, rawOutput bool) {
	apiId, apiSecret := Getapikeys()
	accountIdStr := strconv.Itoa(accountId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/private-locations?limit=0&accountId="+accountIdStr, nil)
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
		//fmt.Printf("\n%-25s %-7s %-7s %-7s %-20s %-10s\n", "HARBOUR ID", "TPE", "EPA", "AGENTS", "NAME", "WORKSPACES")
		var responseBodyAcOpls oplsResponse
		json.Unmarshal(bodyText, &responseBodyAcOpls)
		if responseBodyAcOpls.Error.Code == 0 {
			tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			// Print headers
			fmt.Fprintln(tabWriter, "HARBOUR_ID\tTPE\tEPA\tAGENTS\tNAME\tWORKSPACES")
			for i := 0; i < len(responseBodyAcOpls.Result); i++ {
				harbourID := responseBodyAcOpls.Result[i].Id
				oplName := responseBodyAcOpls.Result[i].Name
				threadsPerEngine := responseBodyAcOpls.Result[i].ThreadsPerEngine
				enginePerAgent := responseBodyAcOpls.Result[i].Slots
				totalWorkspaces := []int{}
				for wd := 0; wd < len(responseBodyAcOpls.Result[i].WorkspacesId); wd++ {
					workspaceList := responseBodyAcOpls.Result[i].WorkspacesId[wd]
					totalWorkspaces = append(totalWorkspaces, workspaceList)
				}
				totalWorkspacesDup := removeDuplicateValuesInt(totalWorkspaces)
				//fmt.Printf("\n%-25s %-7v %-7v %-7v %-25s %-5v", harbourID, threadsPerEngine, enginePerAgent, len(responseBodyAcOpls.Result[i].ShipsId), oplName, totalWorkspacesDup)
				fmt.Fprintf(tabWriter, "%s\t%d\t%d\t%d\t%s\t%d\n", harbourID, threadsPerEngine, enginePerAgent, len(responseBodyAcOpls.Result[i].ShipsId), oplName, totalWorkspacesDup)
			}
			tabWriter.Flush()
			fmt.Println("\n-")

			//fmt.Printf("%-20s %-20s\n", "NAME", "FUNCTIONALITIES SUPPORTED")
			tabWriterFunc := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			// Print headers
			fmt.Fprintln(tabWriterFunc, "NAME\tFUNCTIONALITIES_SUPPORTED")
			for i := 0; i < len(responseBodyAcOpls.Result); i++ {
				oplName := responseBodyAcOpls.Result[i].Name
				functAgent := responseBodyAcOpls.Result[i].FuncIds
				//	fmt.Printf("\n%-20s %-5s\n", oplName, functAgent)
				fmt.Fprintf(tabWriterFunc, "%s\t%s\n", oplName, functAgent)
			}
			tabWriterFunc.Flush()
			fmt.Println("-")
		} else {
			errorCode := responseBodyAcOpls.Error.Code
			errorMessage := responseBodyAcOpls.Error.Message
			fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
		}
	}

}

func getOplsWS(workspaceId int, rawOutput bool) {
	apiId, apiSecret := Getapikeys()
	workspaceIdStr := strconv.Itoa(workspaceId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/private-locations?workspaceId="+workspaceIdStr+"&limit=0", nil)
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
		//fmt.Printf("\n%-25s %-7s %-7s %-7s %-10s %-20s \n", "HARBOUR ID", "TPE", "EPA", "AGENTS", "WORKSPACES", "NAME")
		var responseBodyWsOpls oplsResponse
		json.Unmarshal(bodyText, &responseBodyWsOpls)
		if responseBodyWsOpls.Error.Code == 0 {
			tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			// Print headers
			fmt.Fprintln(tabWriter, "HARBOUR_ID\tTPE\tEPA\tAGENTS\tNAME\tWORKSPACES")
			for i := 0; i < len(responseBodyWsOpls.Result); i++ {
				harbourID := responseBodyWsOpls.Result[i].Id
				oplName := responseBodyWsOpls.Result[i].Name
				threadsPerEngine := responseBodyWsOpls.Result[i].ThreadsPerEngine
				enginePerAgent := responseBodyWsOpls.Result[i].Slots
				totalWorkspaces := []int{}
				for wd := 0; wd < len(responseBodyWsOpls.Result[i].WorkspacesId); wd++ {
					workspaceList := responseBodyWsOpls.Result[i].WorkspacesId[wd]
					totalWorkspaces = append(totalWorkspaces, workspaceList)
				}
				totalWorkspacesDup := removeDuplicateValuesInt(totalWorkspaces)
				//	fmt.Printf("\n%-25s %-7v %-7v %-7v %-5v %-20s", harbourID, threadsPerEngine, enginePerAgent,len(responseBodyWsOpls.Result[i].ShipsId), totalWorkspacesDup, oplName)
				fmt.Fprintf(tabWriter, "%s\t%d\t%d\t%d\t%s\t%d\n", harbourID, threadsPerEngine, enginePerAgent, len(responseBodyWsOpls.Result[i].ShipsId), oplName, totalWorkspacesDup)
			}
			tabWriter.Flush()
			fmt.Println("\n-")

			fmt.Printf("%-20s %-20s\n", "NAME", "FUNCTIONALITIES SUPPORTED")
			for i := 0; i < len(responseBodyWsOpls.Result); i++ {
				oplName := responseBodyWsOpls.Result[i].Name
				functAgent := responseBodyWsOpls.Result[i].FuncIds
				fmt.Printf("\n%-20s %-5s", oplName, functAgent)
			}
			fmt.Println("\n-")

			tabWriterFunc := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			// Print headers
			fmt.Fprintln(tabWriterFunc, "OPL_NAME\tSHIP_ID\tSTATE\tNAME")
			//	fmt.Printf("%-30s %-27s %-8s %-10s\n", "NAME", "SHIP ID", "STATE", "SHIP NAME")
			for i := 0; i < len(responseBodyWsOpls.Result); i++ {
				oplName := responseBodyWsOpls.Result[i].Name
				for f := 0; f < len(responseBodyWsOpls.Result[i].Ships); f++ {
					shipId := responseBodyWsOpls.Result[i].Ships[f].Id
					shipName := responseBodyWsOpls.Result[i].Ships[f].Name
					shipStatus := responseBodyWsOpls.Result[i].Ships[f].State
					//	fmt.Printf("\n%-30s %-27s %-8s %-10s", oplName, shipId, shipStatus, shipName)
					fmt.Fprintf(tabWriterFunc, "%s\t%s\t%s\t%s\n", oplName, shipId, shipStatus, shipName)
				}
			}
			tabWriterFunc.Flush()
			fmt.Println("-")
		} else {
			errorCode := responseBodyWsOpls.Error.Code
			errorMessage := responseBodyWsOpls.Error.Message
			fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
		}
	}
}
