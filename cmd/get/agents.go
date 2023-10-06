/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package get

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// agentsCmd represents the agents command
var agentsCmd = &cobra.Command{
	Use:   "agents",
	Short: "Get agents within a private location",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		harbourId, _ := cmd.Flags().GetString("hid")
		getAgentsOpl(harbourId)
	},
}

func init() {
	GetCmd.AddCommand(agentsCmd)
	agentsCmd.Flags().String("hid", "", "Provide the harbour id")
}
func getAgentsOpl(harbourId string) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/private-locations/"+harbourId+"/servers", nil)
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
