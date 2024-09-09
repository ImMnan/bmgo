/*
Copyright © 2024 Manan Patel - github.com/immnan
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
	Long: `Use the command to get CSV based output of the account usage report by specifying the account Id. BlazeMeter Usage Reports provide metrics about the utilization of BlazeMeter in your organization. These metrics include high-level utilization data that is available either in graphs or JSON payloads, aggregated daily. You can also download detailed (per-test) utilization data in CSV format using this command. You will need to specify the From data and To data of this report & bmgo will output the report for specified period in CSV. The output can be copied to a csv file by using append function as the output is CSV, it is dificult to read through the terminal.  

	For example: [bmgo get -a <account id> usage --from YYYY-MM-DD --to YYYY-MM-DD] OR
                 [bmgo get -a <account id> usage --from YYYY-MM-DD --to YYYY-MM-DD > usage.csv]
	For default: [bmgo get --ac usage --from YYYY-MM-DD --to YYYY-MM-DD] OR
	             [bmgo get --ac usage --from YYYY-MM-DD --to YYYY-MM-DD > usage.csv]`,
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
			cmd.Help()
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
	//fmt.Println(fromEpoch, toEpochStr)
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
