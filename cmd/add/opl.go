/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package add

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

// oplCmd represents the opl command
var oplCmd = &cobra.Command{
	Use:   "opl",
	Short: "Add Private location in account",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ac, _ := cmd.Flags().GetBool("ac")
		var accountId int
		if ac {
			accountId = defaultAccount()
		} else {
			accountId, _ = cmd.Flags().GetInt("accountid")
		}
		oplName, _ := cmd.Flags().GetString("name")
		rawOutput, _ := cmd.Flags().GetBool("raw")
		switch {
		case accountId != 0 && rawOutput:
			addOplraw(accountId, oplName)
		case accountId != 0:
			addOpl(accountId, oplName)
		default:
			fmt.Println("\nPlease provide a correct Account Id & Private location Name")
			fmt.Println("[bmgo add -a <account id> opl --name <private location name>]\n[bmgo add --ac opl --name <private location name>]")
		}
	},
}

func init() {
	AddCmd.AddCommand(oplCmd)
	oplCmd.Flags().String("name", "", "Name your Private location")
	oplCmd.MarkFlagRequired("name")
}

type oplResponse struct {
	Result oplResult   `json:"result"`
	Error  errorResult `json:"error"`
}
type oplResult struct {
	Id               string   `json:"id"`
	Name             string   `json:"name"`
	ThreadsPerEngine int      `json:"threadsPerEngine"`
	Slots            int      `json:"slots"`
	FuncIds          []string `json:"funcIds"`
	ShipsId          []string `json:"shipsId"`
	AccountId        int      `json:"accountid"`
	WorkspacesId     []int    `json:"workspacesId"`
}

func addOpl(accountId int, oplName string) {
	apiId, apiSecret := Getapikeys()
	resultWsId := workspaceIdPrompt()
	resultEPA, resultTPE := oplconfigPrompt()
	client := &http.Client{}
	data := fmt.Sprintf(`{"consoleXms":1024,"consoleXmx":4096,"engineXms":1024,"engineXmx":4096,
	"name":"%s","slots":%v,"threadsPerEngine":%v,"type":"small","overrideCPU":2,"overrideMemory":4096,
	"accountId":%v,"workspacesId":[%v]}`, oplName, resultEPA, resultTPE, accountId, resultWsId)
	reqBodydata := strings.NewReader(data)
	req, err := http.NewRequest("POST", "https://a.blazemeter.com/api/v4/private-locations", reqBodydata)
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
	//fmt.Printf("%s\n", bodyText)
	fmt.Printf("\n%-27s %-15s %-7s %-7s %-10s %-5s\n", "ID", "NAME", "TPE", "EPA", "ACCOUNT", "WORKSPACES")
	var responseBodyAddOpl oplResponse
	json.Unmarshal(bodyText, &responseBodyAddOpl)
	if responseBodyAddOpl.Error.Code == 0 {
		threadsPerEngine := responseBodyAddOpl.Result.ThreadsPerEngine
		enginePerAgent := responseBodyAddOpl.Result.Slots
		oplAccountId := responseBodyAddOpl.Result.AccountId
		oplWorkspaceId := responseBodyAddOpl.Result.WorkspacesId
		oplHarbourId := responseBodyAddOpl.Result.Id
		fmt.Printf("%-27v %-15s %-7v %-7v %-10v %-5v\n", oplHarbourId, oplName, threadsPerEngine, enginePerAgent, oplAccountId, oplWorkspaceId)

		fmt.Println("\n---------------------------------------------------------------------------------------------")
		fmt.Printf("%-30s\n\n", "FUNCTIONALITIES SUPPORTED")
		for i := 0; i < len(responseBodyAddOpl.Result.FuncIds); i++ {
			oplfunctionalities := responseBodyAddOpl.Result.FuncIds[i]
			fmt.Printf("%-30v\n", oplfunctionalities)
		}
		fmt.Println("\n-")
	} else {
		errorCode := responseBodyAddOpl.Error.Code
		errorMessage := responseBodyAddOpl.Error.Message
		fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
	}
}

func addOplraw(accountId int, oplName string) {
	apiId, apiSecret := Getapikeys()
	resultWsId := workspaceIdPrompt()
	resultEPA, resultTPE := oplconfigPrompt()
	client := &http.Client{}
	data := fmt.Sprintf(`{"consoleXms":1024,"consoleXmx":4096,"engineXms":1024,"engineXmx":4096,
	"name":"%s","slots":%v,"threadsPerEngine":%v,"type":"small","overrideCPU":2,"overrideMemory":4096,
	"accountId":%v,"workspacesId":[%v]}`, oplName, resultEPA, resultTPE, accountId, resultWsId)
	reqBodydata := strings.NewReader(data)
	req, err := http.NewRequest("POST", "https://a.blazemeter.com/api/v4/private-locations", reqBodydata)
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
	fmt.Printf("%s\n", bodyText)
}
