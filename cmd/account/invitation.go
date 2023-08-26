/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package account

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
)

// invitationCmd represents the invitation command
var invitationCmd = &cobra.Command{
	Use:   "invitation",
	Short: "A brief description of your command",
	Long:  `A longer description that spans multiple lines and likely`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("invitation called")
		accountId, _ := cmd.Flags().GetInt("id")
		invitations(accountId)
	},
}

func init() {
	AccountCmd.AddCommand(invitationCmd)
}

func invitations(accountId int) {
	apiId, apiSecret := Getapikeys()

	client := &http.Client{}
	accountIdStr := strconv.Itoa(accountId)
	fmt.Println("Account scanned is: ", accountIdStr)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/accounts/"+accountIdStr+"/invitations", nil)
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
