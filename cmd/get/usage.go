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
	"time"

	"github.com/spf13/cobra"
)

// usageCmd represents the usage command
var usageCmd = &cobra.Command{
	Use:   "usage",
	Short: "Get Usage report for an account",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ac, _ := cmd.Flags().GetBool("ac")
		var accountId int
		if ac {
			accountId = defaultAccount()
		} else {
			accountId, _ = cmd.Flags().GetInt("accountid")
		}
		fromDate, _ := cmd.Flags().GetString("from")
		toDate, _ := cmd.Flags().GetString("to")
		switch {
		case fromDate != "" && toDate != "":
			getUsage(accountId, fromDate, toDate)
		default:
			fmt.Println("\nPlease provide a correct account Id, from data & to date to get the report")
			fmt.Println("[bmgo get -a <account_id> usage --from YYYY-MM-DD --to YYYY-MM-DD] OR [bmgo get --ac <account_id> usage --from YYYY-MM-DD --to YYYY-MM-DD")
		}
	},
}

func init() {
	GetCmd.AddCommand(usageCmd)
	usageCmd.Flags().String("from", "", "Provide the start/from date of report [YYYY-MM-DD]")
	usageCmd.Flags().String("to", "", "Provide the finish/to date of report [YYYY-MM-DD]")
	usageCmd.MarkFlagRequired("from")
	usageCmd.MarkFlagRequired("to")
}

func getUsage(accountId int, fromDate, toDate string) {
	apiId, apiSecret := Getapikeys()
	accountIdStr := strconv.Itoa(accountId)
	defaultTime := "T1:00:00+00:00"
	// from time converted to epoch
	fromHtime := fmt.Sprint(fromDate + defaultTime)
	fromTime, e := time.Parse(time.RFC3339, fromHtime)
	if e != nil {
		panic("Can't parse time format")
	}
	fromEpoch := fromTime.Unix()
	fromEpochStr := fmt.Sprint(fromEpoch)
	// To time converted to epoch
	toHtime := fmt.Sprint(toDate + defaultTime)
	toTime, e := time.Parse(time.RFC3339, toHtime)
	if e != nil {
		panic("Can't parse time format")
	}
	toEpoch := toTime.Unix()
	toEpochStr := fmt.Sprint(toEpoch)

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/accounts/"+accountIdStr+"/reports/usage/tests/credits?daysInterval=1&toDate="+toEpochStr+"&fromDate="+fromEpochStr+"&download=true", nil)
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
