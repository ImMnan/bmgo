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

// accountCmd represents the account command
var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Get details about the account, use with other sub-commands to get specific info",
	Long:  `.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("account called")
		accountId, _ := cmd.Flags().GetInt("id")
		getAccountId(accountId)
	},
}

func getAccountId(accountId int) {
	apiId, apiSecret := Getapikeys()

	accountIdStr := strconv.Itoa(accountId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/accounts/"+accountIdStr, nil)
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

func init() {
	GetCmd.AddCommand(accountCmd)
	accountCmd.PersistentFlags().Int("id", 0, " [*Required] Confirm the account id")
	accountCmd.MarkPersistentFlagRequired("id")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// accountCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// accountCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
