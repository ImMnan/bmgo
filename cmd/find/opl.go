/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package find

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// userCmd represents the user command
var oplCmd = &cobra.Command{
	Use:   "opl",
	Short: "Find details about the Private location",
	Long:  `.`,
	Run: func(cmd *cobra.Command, args []string) {
		harbourId, _ := cmd.Flags().GetString("hid")
		rawOutput, _ := cmd.Flags().GetBool("raw")
		if harbourId != "" && rawOutput {
			findOplraw(harbourId)
		} else if harbourId != "" {
			findOpl(harbourId)
		} else {
			fmt.Println("\nPlease provide a correct Harbour ID to find the Private Location - OPL")
			fmt.Println("[bmgo find opl --hid <harbour ID>")
		}
	},
}

func init() {
	FindCmd.AddCommand(oplCmd)
	oplCmd.Flags().String("hid", "", "Provide the Harbour ID")
	oplCmd.MarkFlagRequired("hid")
}

type oplResponse struct {
	Result oplResult   `json:"result"`
	Error  errorResult `json:"error"`
}
type oplResult struct {
	Name             string   `json:"name"`
	ThreadsPerEngine int      `json:"threadsPerEngine"`
	Slots            int      `json:"slots"`
	FuncIds          []string `json:"funcIds"`
	ShipsId          []string `json:"shipsId"`
	AccountId        int      `json:"accountid"`
	WorkspacesId     []int    `json:"workspacesId"`
	Ships            []ships
}
type ships struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	State string `json:"state"`
}

func findOpl(harbourId string) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/private-locations/"+harbourId, nil)
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
	var responseBodyOpl oplResponse
	json.Unmarshal(bodyText, &responseBodyOpl)
	if responseBodyOpl.Error.Code == 0 {
		fmt.Printf("\n%-20s %-7s %-7s %-7s %-10s %-5s\n", "NAME", "TPE", "EPA", "AGENTS", "ACCOUNT", "WORKSPACES")
		oplName := responseBodyOpl.Result.Name
		threadsPerEngine := responseBodyOpl.Result.ThreadsPerEngine
		enginePerAgent := responseBodyOpl.Result.Slots
		oplAccountId := responseBodyOpl.Result.AccountId
		oplWorkspaceId := responseBodyOpl.Result.WorkspacesId
		fmt.Printf("%-20s %-7v %-7v %-7v %-10v %-5v\n", oplName, threadsPerEngine, enginePerAgent, len(responseBodyOpl.Result.ShipsId), oplAccountId, oplWorkspaceId)

		fmt.Println("\n---------------------------------------------------------------------------------------------")
		fmt.Printf("%-30s\n\n", "FUNCTIONALITIES SUPPORTED")
		for i := 0; i < len(responseBodyOpl.Result.FuncIds); i++ {
			oplfunctionalities := responseBodyOpl.Result.FuncIds[i]
			fmt.Printf("%-30v\n", oplfunctionalities)
		}
		fmt.Println("\n---------------------------------------------------------------------------------------------")
		fmt.Printf("%-20s %-20s %-25s %-10s\n", "HARBOUR NAME", "SHIP NAME", "SHIP ID", "STATE")
		for f := 0; f < len(responseBodyOpl.Result.Ships); f++ {
			shipId := responseBodyOpl.Result.Ships[f].Id
			shipName := responseBodyOpl.Result.Ships[f].Name
			shipStatus := responseBodyOpl.Result.Ships[f].State
			fmt.Printf("\n%-20s %-20s %-25s %-10s", oplName, shipName, shipId, shipStatus)
		}
		fmt.Println("\n-")
	} else {
		errorCode := responseBodyOpl.Error.Code
		errorMessage := responseBodyOpl.Error.Message
		fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
	}
}
func findOplraw(harbourId string) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/private-locations/"+harbourId, nil)
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
