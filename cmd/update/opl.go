/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package update

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// oplCmd represents the opl command
var oplCmd = &cobra.Command{
	Use:   "opl",
	Short: "Update OPL- Add or Remove workspace from private location",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ws, _ := cmd.Flags().GetBool("ws")
		var workspaceId int
		if ws {
			workspaceId = defaultWorkspace()
		} else {
			workspaceId, _ = cmd.Flags().GetInt("workspaceid")
		}
		harbourId, _ := cmd.Flags().GetString("hid")
		removeWs, _ := cmd.Flags().GetBool("remove")
		addWs, _ := cmd.Flags().GetBool("add")
		switch {
		case workspaceId != 0 && harbourId != "" && addWs:
			updateOplAddWs(workspaceId, harbourId)
		case workspaceId != 0 && harbourId != "" && removeWs:
			updateOplRemWs(workspaceId, harbourId)
		}
	},
}

func init() {
	UpdateCmd.AddCommand(oplCmd)
	oplCmd.Flags().String("hid", "", "Provide the Harbour ID to update")
	oplCmd.Flags().Bool("add", false, "To add the declared workspace to Private location")
	oplCmd.Flags().Bool("remove", false, "To remove the declared workspace from Private location")
}

func updateOplAddWs(workspaceId int, harbourId string) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	data := fmt.Sprintf(`{"workspaceId": %v}`, workspaceId)
	requestdata := strings.NewReader(data)
	req, err := http.NewRequest("POST", "https://a.blazemeter.com/api/v4/private-locations/"+harbourId+"/add-workspace", requestdata)
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

func updateOplRemWs(workspaceId int, harbourId string) {
	apiId, apiSecret := Getapikeys()
	workspaceIdStr := strconv.Itoa(workspaceId)
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", "https://a.blazemeter.com/api/v4/private-locations/"+harbourId+"/workspaces/"+workspaceIdStr, nil)
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
