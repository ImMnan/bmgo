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
		accountId, _ := cmd.Flags().GetInt("id")
		rawOutput, _ := cmd.Flags().GetBool("raw")
		if rawOutput {
			getOplsRaw(accountId)
		} else {
			getOpls(accountId)
		}
		getOpls(accountId)
	},
}

func init() {
	accountCmd.AddCommand(oplsCmd)
}

type oplsResponseAC struct {
	Result []oplsResultAC `json:"result"`
}

type oplsResultAC struct {
	Id               string   `json:"id"`
	Name             string   `json:"name"`
	ThreadsPerEngine int      `json:"threadsPerEngine"`
	Slots            int      `json:"slots"`
	FuncIds          []string `json:"funcIds"`
	ShipsId          []string `json:"shipsId"`
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
	var responseBodyAcOpls oplsResponseAC
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
	fmt.Println("\n")
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
