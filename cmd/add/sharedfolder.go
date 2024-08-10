/*
Copyright © 2024 Manan Patel - github.com/immnan
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
	Short: "Add shared folder into workspace",
	Long: `If you have files that are used across multiple tests, then using shared folders will reduce the need to re-upload files every time. For more information on what shared folders are and when to use them go to help.blazemeter.com Add shared folder into your existing workspace by specifying the name of the shared folder. Add a shared folder to your existing workspace using this command. 
	
	For example: [bmgo add -w <workspace id> sharedfolder --name <folder name>]
	For default: [bmgo add --ws sharedfolder --name <folder name>]`,
	Run: func(cmd *cobra.Command, args []string) {
		ws, _ := cmd.Flags().GetBool("ws")
		var workspaceId int
		if ws {
			workspaceId = defaultWorkspace()
		} else {
			workspaceId, _ = cmd.Flags().GetInt("workspaceid")
		}
		folderName, _ := cmd.Flags().GetString("name")
		rawOutput, _ := cmd.Flags().GetBool("raw")
		switch {
		case (workspaceId != 0) && (folderName != ""):
			addSharedfolder(folderName, workspaceId, rawOutput)
		default:
			cmd.Help()
		}
	},
}

func init() {
	AddCmd.AddCommand(sharedfolderCmd)
	sharedfolderCmd.Flags().String("name", "", "Name your Shared folder")
	sharedfolderCmd.MarkFlagRequired("name")
	sharedfolderCmd.Flags().IntP("workspaceid", "w", 0, " Provide Workspace ID to add a resource to")
	sharedfolderCmd.Flags().Bool("ws", false, "Use default workspace Id (bmConfig.yaml)")
}

type addFolderResponse struct {
	Result addfolderResult `json:"result"`
	Error  errorResult     `json:"error"`
}
type addfolderResult struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func addSharedfolder(folderName string, workspaceId int, rawOutput bool) {
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
	if rawOutput {
		fmt.Printf("%s\n", bodyText)
	} else {
		var responseBodyAddFolder addFolderResponse
		json.Unmarshal(bodyText, &responseBodyAddFolder)
		if responseBodyAddFolder.Error.Code == 0 {
			folderIdres := responseBodyAddFolder.Result.Id
			folderNameres := responseBodyAddFolder.Result.Name
			fmt.Printf("\n%-30s %-15s", "Folder ID", "NAME")
			fmt.Printf("\n%-30s %-15s", folderIdres, folderNameres)
			fmt.Println("\n-")
		} else {
			errorCode := responseBodyAddFolder.Error.Code
			errorMessage := responseBodyAddFolder.Error.Message
			fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
		}
	}
}
