/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package get

import (
	"encoding/json"
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
	Short: "Get details about the account, use with other sub-commands to get specific/detailed info",
	Long:  `.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("account called")
		accountId, _ := cmd.Flags().GetInt("id")
		getAccountId(accountId)
	},
}

type responseBody struct {
	Result result
}
type result struct {
	Name            string `json:"name"`
	Owner           owner
	MembersCount    int `json:"membersCount"`
	WorkspacesCount int `json:"workspacesCount"`
	Plan            plan
	CloudProviders  []string
}

type owner struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type plan struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	ReportRetention  int    `json:"reportRetention"`
	ThreadsPerEngine int    `json:"threadsPerEngine"`
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
	//fmt.Printf("%s\n", bodyText)
	var responseObject responseBody
	json.Unmarshal(bodyText, &responseObject)

	accountName := responseObject.Result.Name
	ownerEmail := responseObject.Result.Owner.Email
	workspaceCount := responseObject.Result.WorkspacesCount
	memberCount := responseObject.Result.MembersCount

	accountPlanId := responseObject.Result.Plan.Id
	accountPlanName := responseObject.Result.Plan.Name
	accountReportRet := responseObject.Result.Plan.ReportRetention
	accountThreadsPE := responseObject.Result.Plan.ThreadsPerEngine

	fmt.Printf("\n%-20s %-30s %-10s %-10s\n", "NAME", "OWNER", "WORKSPACES", "USERS")
	fmt.Printf("%-20s %-30s %-10d %-10d\n\n", accountName, ownerEmail, workspaceCount, memberCount)

	fmt.Printf("PLan details for account %s (%v)\n", accountName, accountId)

	fmt.Printf("%-20s %-30s %-15s %-10s\n", "PLAN ID", "PLAN NAME", "REPORT RETENT.", "THREADS/ENGINE")
	fmt.Printf("%-20s %-30s %-15d %-10d\n", accountPlanId, accountPlanName, accountReportRet, accountThreadsPE)

	cloudProviders := []string{}
	for i := 0; i < len(responseObject.Result.CloudProviders); i++ {
		cloudProlist := responseObject.Result.CloudProviders[i]
		cloudProviders = append(cloudProviders, cloudProlist)
	}
	fmt.Printf("\n Available cloud provider for %s (%v): \n", accountName, accountId)
	fmt.Println(cloudProviders)
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
