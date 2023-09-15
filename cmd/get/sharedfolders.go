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
	"time"

	"github.com/spf13/cobra"
)

// sharedfoldersCmd represents the sharedfolders command
var sharedfoldersCmd = &cobra.Command{
	Use:   "sharedfolders",
	Short: "Get details of shared folders within workspace",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		workspaceId, _ := cmd.Flags().GetInt("id")
		rawrawOutput, _ := cmd.Flags().GetBool("raw")
		if rawrawOutput {
			getSharedFolderWsRaw(workspaceId)
		} else {
			fmt.Printf("Getting sharedfolders within workspace %v ...\n", workspaceId)
			getSharedFolderWs(workspaceId)
		}
	},
}

func init() {
	workspaceCmd.AddCommand(sharedfoldersCmd)
}

type sfolderResponseWS struct {
	Result []sfoldersResult `json:"result"`
}

type sfoldersResult struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Created int    `json:"created"`
	Hidden  bool   `json:"hidden"`
}

func getSharedFolderWs(workspaceId int) {
	workspaceIdStr := strconv.Itoa(workspaceId)
	client := &http.Client{}
	apiId, apiSecret := Getapikeys()
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/folders?workspaceId="+workspaceIdStr+"&limit=200", nil)
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
	var responseBodyWsSfolders sfolderResponseWS
	json.Unmarshal(bodyText, &responseBodyWsSfolders)
	fmt.Printf("\n%-25s %-15s %-32s %-10s", "ID", "NAME", "CREATED ON", "HIDDEN")
	for i := 0; i < len(responseBodyWsSfolders.Result); i++ {
		sFolderId := responseBodyWsSfolders.Result[i].Id
		sFolderName := responseBodyWsSfolders.Result[i].Name
		sFolderCreated := int64(responseBodyWsSfolders.Result[i].Created)
		sFHidden := responseBodyWsSfolders.Result[i].Hidden
		epochCreated := time.Unix(sFolderCreated, 0)
		fmt.Printf("\n%-25s %-15s %-32v %-10t", sFolderId, sFolderName, epochCreated, sFHidden)
	}
	fmt.Println("\n-")
}

func getSharedFolderWsRaw(workspaceId int) {
	workspaceIdStr := strconv.Itoa(workspaceId)
	client := &http.Client{}
	apiId, apiSecret := Getapikeys()
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/folders?workspaceId="+workspaceIdStr+"&limit=200", nil)
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
