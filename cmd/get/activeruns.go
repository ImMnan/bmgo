/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package get

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
)

// activerunsCmd represents the activeruns command
var activerunsCmd = &cobra.Command{
	Use:   "activeruns",
	Short: "Get active tests for a workspace",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var workspaceId int
		ws, _ := cmd.Flags().GetBool("ws")
		rawOutput, _ := cmd.Flags().GetBool("raw")
		if ws {
			workspaceId = defaultWorkspace()
		} else {
			workspaceId, _ = cmd.Flags().GetInt("workspaceid")
		}
		if workspaceId != 0 && rawOutput {
			getActiveruns(workspaceId)
		} else if workspaceId != 0 {
			getActiveruns(workspaceId)
		} else {
			fmt.Println("\nPlease provide a valid workspace ID to get list of active tests in a workspace")
			fmt.Println("[bmgo get -w <workspace Id>")
		}
	},
}

func init() {
	GetCmd.AddCommand(activerunsCmd)
}
func getActiveruns(workspaceId int) {
	apiId, apiSecret := Getapikeys()
	workspaceIdStr := strconv.Itoa(workspaceId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/workspaces/"+workspaceIdStr+"/active", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("accept", "application/json")
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
