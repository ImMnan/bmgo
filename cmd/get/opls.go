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

// oplsCmd represents the opls command
var oplsCmd = &cobra.Command{
	Use:   "opls",
	Short: "Get a list of Private locations in the account",
	Long:  `.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("opls called")
		accountId, _ := cmd.Flags().GetInt("id")
		getOpls(accountId)
	},
}

func init() {
	accountCmd.AddCommand(oplsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// oplsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// oplsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getOpls(accountId int) {
	apiId, apiSecret := Getapikeys()

	accountIdStr := strconv.Itoa(accountId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/private-locations?limit=100&accountId="+accountIdStr, nil)
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
