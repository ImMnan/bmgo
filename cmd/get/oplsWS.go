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
var oplsWSCmd = &cobra.Command{
	Use:   "opls",
	Short: "Get a list of Private locations in the account",
	Long:  `.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("opls called")
		workspaceId, _ := cmd.Flags().GetInt("id")
		getOplsWS(workspaceId)
	},
}

func init() {
	workspaceCmd.AddCommand(oplsWSCmd)
}

type oplsResponseWS struct {
	Result []oplsResult `json:"result"`
}

type oplsResult struct {
	Id               string   `json:"id"`
	Name             string   `json:"name"`
	ThreadsPerEngine int      `json:"threadsPerEngine"`
	Slots            int      `json:"slots"`
	FuncIds          []string `json:"funcIds"`
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
	fmt.Printf("\n%-25s %-20s %-10s %-10s\n", "HARBOUR ID", "NAME", "TPE", "EPA")
	var responseBodyWsOpls oplsResponseWS
	json.Unmarshal(bodyText, &responseBodyWsOpls)
	for i := 0; i < len(responseBodyWsOpls.Result); i++ {
		harbourID := responseBodyWsOpls.Result[i].Id
		oplName := responseBodyWsOpls.Result[i].Name
		threadsPerEngine := responseBodyWsOpls.Result[i].ThreadsPerEngine
		enginePerAgent := responseBodyWsOpls.Result[i].Slots
		fmt.Printf("\n%-25s %-20s %-10v %-10v", harbourID, oplName, threadsPerEngine, enginePerAgent)
	}
}
