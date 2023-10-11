/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package find

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Find and download logs",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		sessionId, _ := cmd.Flags().GetString("sid")
		findlogSession(sessionId)
	},
}

func init() {
	FindCmd.AddCommand(logsCmd)
	logsCmd.Flags().String("sid", "", "Provide session Id to pull logs for")
}

func findlogSession(sessionId string) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/sessions/"+sessionId+"/reports/logs", nil)
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
