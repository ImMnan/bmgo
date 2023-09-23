/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package add

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// oplCmd represents the opl command
var oplCmd = &cobra.Command{
	Use:   "opl",
	Short: "Add Private location in account",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		accountId, _ := cmd.Flags().GetInt("accountid")
		oplName, _ := cmd.Flags().GetString("name")
		rawOutput, _ := cmd.Flags().GetBool("raw")
		switch {
		case accountId != 0 && rawOutput:
			addOplraw(accountId, oplName)
		case accountId != 0:
			addOpl(accountId, oplName)
		default:
			fmt.Println("\nPlease provide a correct Account Id & Private location Name")
			fmt.Println("[bmgo add -a <account id> opl --name <private location name>]\n-")
		}
	},
}

func init() {
	AddCmd.AddCommand(oplCmd)
	oplCmd.Flags().String("name", "", "Name your Private location")
	oplCmd.MarkFlagRequired("name")
}

type oplResponse struct {
	Result oplResult `json:"result"`
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

func oplconfigPrompt() (int, int) {
	validate := func(input string) error {
		_, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return errors.New("invalid number or nan")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:       "Number of Engines per Agent",
		HideEntered: true,
		Validate:    validate,
	}
	prompt1 := promptui.Prompt{
		Label:       "Number of Threads per Engine",
		HideEntered: true,
		Validate:    validate,
	}
	resultEPAstr, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}
	resultTPEstr, err := prompt1.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}
	resultEPA, _ := strconv.Atoi(resultEPAstr)
	resultTPE, _ := strconv.Atoi(resultTPEstr)
	return resultEPA, resultTPE
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
