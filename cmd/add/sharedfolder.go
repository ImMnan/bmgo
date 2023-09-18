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

// projectCmd represents the project command
var sharedfolderCmd = &cobra.Command{
	Use:   "sharedfolder",
	Short: "A brief description of your command",
	Long:  `.`,
	Run: func(cmd *cobra.Command, args []string) {
		workspaceId, _ := cmd.Flags().GetInt("workspaceid")
		folderName, _ := cmd.Flags().GetString("name")
		rawOutput, _ := cmd.Flags().GetBool("raw")
		switch {
		case (workspaceId != 0) && (folderName != "") && rawOutput:
			addSharedfolderraw(folderName, workspaceId)
		case (workspaceId != 0) && (folderName != ""):
			addSharedfolder(folderName, workspaceId)
		default:
			fmt.Println("\nPlease provide a correct Workspace Id & Project Name")
			fmt.Println("[bmgo add -w <workspace id> project --name <project name>]\n-")
		}
	},
}

func init() {
	AddCmd.AddCommand(sharedfolderCmd)
	sharedfolderCmd.Flags().String("name", "", "Name your Shared folder")
	sharedfolderCmd.MarkFlagRequired("name")
}

type addFolderResponse struct {
	Result addfolderResult `json:"result"`
}
type addfolderResult struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func addSharedfolder(folderName string, workspaceId int) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	data := fmt.Sprintf(`{
		"name": "%s",  
		"workspaceId": %v}`, folderName, workspaceId)
	reqBodyData := strings.NewReader(data)
	req, err := http.NewRequest("POST", "https://a.blazemeter.com/api/v4/folders", reqBodyData)
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
	var responseBodyAddFolder addFolderResponse
	json.Unmarshal(bodyText, &responseBodyAddFolder)
	folderIdres := responseBodyAddFolder.Result.Id
	folderNameres := responseBodyAddFolder.Result.Name
	fmt.Printf("\n%-30s %-15s", "Folder ID", "NAME")
	fmt.Printf("\n%-30s %-15s", folderIdres, folderNameres)
	fmt.Println("\n-")
}
func addSharedfolderraw(folderName string, workspaceId int) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	data := fmt.Sprintf(`{
		"name": "%s",  
		"workspaceId": %v}`, folderName, workspaceId)
	reqBodyData := strings.NewReader(data)
	req, err := http.NewRequest("POST", "https://a.blazemeter.com/api/v4/folders", reqBodyData)
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
