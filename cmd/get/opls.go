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
	Long:  `.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("opls called")
		accountId, _ := cmd.Flags().GetInt("accountid")
		workspaceId, _ := cmd.Flags().GetInt("workspaceid")
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
			fmt.Println("\nPlease provide a correct workspace Id or Account Id to get the info")
			fmt.Println("[bmgo get -a <account_id>...] OR [bmgo get -w <workspace_id>...]")
		}
	},
}

func init() {
	GetCmd.AddCommand(oplsCmd)
}

type oplsResponse struct {
	Result []oplsResult `json:"result"`
}

type oplsResult struct {
	Id               string   `json:"id"`
	Name             string   `json:"name"`
	ThreadsPerEngine int      `json:"threadsPerEngine"`
	Slots            int      `json:"slots"`
	FuncIds          []string `json:"funcIds"`
	ShipsId          []string `json:"shipsId"`
	Ships            []ships
}

type ships struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	State string `json:"state"`
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
	fmt.Printf("\n%-25s %-20s %-10s %-10s %-10s %-10s \n", "HARBOUR ID", "NAME", "TPE", "EPA", "AGENTS", "CAP")
	var responseBodyAcOpls oplsResponse
	json.Unmarshal(bodyText, &responseBodyAcOpls)
	for i := 0; i < len(responseBodyAcOpls.Result); i++ {
		harbourID := responseBodyAcOpls.Result[i].Id
		oplName := responseBodyAcOpls.Result[i].Name
		threadsPerEngine := responseBodyAcOpls.Result[i].ThreadsPerEngine
		enginePerAgent := responseBodyAcOpls.Result[i].Slots
		fmt.Printf("\n%-25s %-20s %-10v %-10v %-10v %-10v", harbourID, oplName, threadsPerEngine, enginePerAgent, len(responseBodyAcOpls.Result[i].ShipsId), (threadsPerEngine * enginePerAgent * len(responseBodyAcOpls.Result[i].ShipsId)))
	}
	fmt.Println("\n\n---------------------------------------------------------------------------------------------")
	fmt.Printf("%-20s %-20s\n", "NAME", "FUNCTIONALITIES SUPPORTED")
	for i := 0; i < len(responseBodyAcOpls.Result); i++ {
		oplName := responseBodyAcOpls.Result[i].Name
		functAgent := responseBodyAcOpls.Result[i].FuncIds
		fmt.Printf("\n%-20s %-5s", oplName, functAgent)
	}
	fmt.Println("\n-")
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
	fmt.Printf("\n%-25s %-20s %-10s %-10s %-10s %-10s \n", "HARBOUR ID", "NAME", "TPE", "EPA", "AGENTS", "CAP")
	var responseBodyWsOpls oplsResponse
	json.Unmarshal(bodyText, &responseBodyWsOpls)
	for i := 0; i < len(responseBodyWsOpls.Result); i++ {
		harbourID := responseBodyWsOpls.Result[i].Id
		oplName := responseBodyWsOpls.Result[i].Name
		threadsPerEngine := responseBodyWsOpls.Result[i].ThreadsPerEngine
		enginePerAgent := responseBodyWsOpls.Result[i].Slots
		fmt.Printf("\n%-25s %-20s %-10v %-10v %-10v %-10v", harbourID, oplName, threadsPerEngine, enginePerAgent, len(responseBodyWsOpls.Result[i].ShipsId), (threadsPerEngine * enginePerAgent * len(responseBodyWsOpls.Result[i].ShipsId)))
	}
	fmt.Println("\n\n---------------------------------------------------------------------------------------------")
	fmt.Printf("%-20s %-20s\n", "NAME", "FUNCTIONALITIES SUPPORTED")
	for i := 0; i < len(responseBodyWsOpls.Result); i++ {
		oplName := responseBodyWsOpls.Result[i].Name
		functAgent := responseBodyWsOpls.Result[i].FuncIds
		fmt.Printf("\n%-20s %-5s", oplName, functAgent)
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
