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
	"time"

	"github.com/spf13/cobra"
)

// sharedfoldersCmd represents the sharedfolders command
var sharedfoldersCmd = &cobra.Command{
	Use:   "sharedfolders",
	Short: "Get details of shared folders within workspace",
	Long: `Use the command to list Shared folders within a specified workspace. You can use the same files across multiple tests. Upload the files to Shared Folders and include the folders in as many tests as you like. The output includes service ID, Name, etc.

	For example: [bmgo get -w <workspace id> sharedfolders] OR
	For default: [bmgo get --ws sharedfolders]`,
	Run: func(cmd *cobra.Command, args []string) {
		ws, _ := cmd.Flags().GetBool("ws")
		var workspaceId int
		if ws {
			workspaceId = defaultWorkspace()
		} else {
			workspaceId, _ = cmd.Flags().GetInt("workspaceid")
		}
		rawOutput, _ := cmd.Flags().GetBool("raw")
		if workspaceId != 0 {
			getSharedFolderWs(workspaceId, rawOutput)
		} else {
			cmd.Help()
		}
	},
}

func init() {
	GetCmd.AddCommand(sharedfoldersCmd)
}

type sfolderResponseWS struct {
	Result []sfoldersResult `json:"result"`
	Error  errorResult      `json:"error"`
}

type sfoldersResult struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Created int    `json:"created"`
	Hidden  bool   `json:"hidden"`
}

func getSharedFolderWs(workspaceId int, rawOutput bool) {
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
	if rawOutput {
		fmt.Printf("%s\n", bodyText)
	} else {
		var responseBodyWsSfolders sfolderResponseWS
		json.Unmarshal(bodyText, &responseBodyWsSfolders)
		if responseBodyWsSfolders.Error.Code == 0 {
			//fmt.Printf("\n%-27s %-30s %-18s %-10s", "ID", "NAME", "CREATED ON", "HIDDEN")
			tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			// Print headers
			fmt.Fprintln(tabWriter, "ID\tNAME\tCREATED_ON\tHIDDEN")
			for i := 0; i < len(responseBodyWsSfolders.Result); i++ {
				sFolderId := responseBodyWsSfolders.Result[i].Id
				sFolderName := responseBodyWsSfolders.Result[i].Name
				sFolderCreated := int64(responseBodyWsSfolders.Result[i].Created)
				sFHidden := responseBodyWsSfolders.Result[i].Hidden
				epochCreated := fmt.Sprint(time.Unix(sFolderCreated, 0))
				//	fmt.Printf("\n%-27s %-30s %-18v %-10t", sFolderId, sFolderName, epochCreated[0:16], sFHidden)
				fmt.Fprintf(tabWriter, "%s\t%s\t%s\t%t\n", sFolderId, sFolderName, epochCreated[0:16], sFHidden)
			}
			tabWriter.Flush()
			fmt.Println("-")
		} else {
			errorCode := responseBodyWsSfolders.Error.Code
			errorMessage := responseBodyWsSfolders.Error.Message
			fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
		}

	}
}
