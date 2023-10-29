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
	             [bmgo get --ac projects]`,
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
		if (workspaceId != 0) && (accountId == 0) && rawOutput {
			getOplsWSRaw(workspaceId)
		} else if (accountId != 0) && (workspaceId == 0) && rawOutput {
			getOplsRaw(accountId)
		} else if (workspaceId != 0) && (accountId == 0) {
			getOplsWS(workspaceId)
		} else if (accountId != 0) && (workspaceId == 0) {
			getOpls(accountId)
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

func getOpls(accountId int) {
	apiId, apiSecret := Getapikeys()
	accountIdStr := strconv.Itoa(accountId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/private-locations?limit=100&accountId="+accountIdStr, nil)
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
	fmt.Printf("\n%-25s %-20s %-7s %-7s %-7s %-10s \n", "HARBOUR ID", "NAME", "TPE", "EPA", "AGENTS", "WORKSPACES")
	var responseBodyAcOpls oplsResponse
	json.Unmarshal(bodyText, &responseBodyAcOpls)
	if responseBodyAcOpls.Error.Code == 0 {
		totalWorkspaces := []int{}
		for i := 0; i < len(responseBodyAcOpls.Result); i++ {
			harbourID := responseBodyAcOpls.Result[i].Id
			oplName := responseBodyAcOpls.Result[i].Name
			threadsPerEngine := responseBodyAcOpls.Result[i].ThreadsPerEngine
			enginePerAgent := responseBodyAcOpls.Result[i].Slots
			for wd := 0; wd < len(responseBodyAcOpls.Result[i].WorkspacesId); wd++ {
				workspaceList := responseBodyAcOpls.Result[i].WorkspacesId[wd]
				totalWorkspaces = append(totalWorkspaces, workspaceList)
			}
			totalWorkspacesDup := removeDuplicateValuesInt(totalWorkspaces)
			fmt.Printf("\n%-25s %-20s %-7v %-7v %-7v %-5v", harbourID, oplName, threadsPerEngine, enginePerAgent,
				len(responseBodyAcOpls.Result[i].ShipsId), totalWorkspacesDup)
		}
		fmt.Println("\n\n---------------------------------------------------------------------------------------------")
		fmt.Printf("%-20s %-20s\n", "NAME", "FUNCTIONALITIES SUPPORTED")
		for i := 0; i < len(responseBodyAcOpls.Result); i++ {
			oplName := responseBodyAcOpls.Result[i].Name
			functAgent := responseBodyAcOpls.Result[i].FuncIds
			fmt.Printf("\n%-20s %-5s\n", oplName, functAgent)
		}
		fmt.Println("\n-")
	} else {
		errorCode := responseBodyAcOpls.Error.Code
		errorMessage := responseBodyAcOpls.Error.Message
		fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
	}

}
func getOplsRaw(accountId int) {
	apiId, apiSecret := Getapikeys()

	accountIdStr := strconv.Itoa(accountId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/private-locations?limit=100&accountId="+accountIdStr, nil)
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

func getOplsWS(workspaceId int) {
	apiId, apiSecret := Getapikeys()
	workspaceIdStr := strconv.Itoa(workspaceId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/private-locations?workspaceId="+workspaceIdStr+"&limit=100", nil)
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
	fmt.Printf("\n%-25s %-25s %-7s %-7s %-7s %-5s \n", "HARBOUR ID", "NAME", "TPE", "EPA", "AGENTS", "WORKSPACES")
	var responseBodyWsOpls oplsResponse
	json.Unmarshal(bodyText, &responseBodyWsOpls)
	if responseBodyWsOpls.Error.Code == 0 {
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
			fmt.Printf("\n%-25s %-25s %-7v %-7v %-7v %-5v", harbourID, oplName, threadsPerEngine, enginePerAgent, len(responseBodyWsOpls.Result[i].ShipsId), totalWorkspacesDup)
		}
		fmt.Println("\n\n---------------------------------------------------------------------------------------------")
		fmt.Printf("%-20s %-20s\n", "NAME", "FUNCTIONALITIES SUPPORTED")
		for i := 0; i < len(responseBodyWsOpls.Result); i++ {
			oplName := responseBodyWsOpls.Result[i].Name
			functAgent := responseBodyWsOpls.Result[i].FuncIds
			fmt.Printf("\n%-20s %-5s\n", oplName, functAgent)
		}
		fmt.Println("\n\n---------------------------------------------------------------------------------------------")
		fmt.Printf("%-20s %-20s %-25s %-10s\n", "NAME", "SHIP NAME", "SHIP ID", "STATE")
		for i := 0; i < len(responseBodyWsOpls.Result); i++ {
			oplName := responseBodyWsOpls.Result[i].Name
			for f := 0; f < len(responseBodyWsOpls.Result[i].Ships); f++ {
				shipId := responseBodyWsOpls.Result[i].Ships[f].Id
				shipName := responseBodyWsOpls.Result[i].Ships[f].Name
				shipStatus := responseBodyWsOpls.Result[i].Ships[f].State
				fmt.Printf("\n%-20s %-20s %-25s %-10s", oplName, shipName, shipId, shipStatus)
			}
		}
		fmt.Println("\n-")
	} else {
		errorCode := responseBodyWsOpls.Error.Code
		errorMessage := responseBodyWsOpls.Error.Message
		fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
	}
}

func getOplsWSRaw(workspaceId int) {
	apiId, apiSecret := Getapikeys()
	workspaceIdStr := strconv.Itoa(workspaceId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/private-locations?workspaceId="+workspaceIdStr+"&limit=100", nil)
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
