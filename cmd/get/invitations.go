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

// invitationsCmd represents the invitations command
var invitationsCmd = &cobra.Command{
	Use:   "invitations",
	Short: "Get a list of pending invitations within your account",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("invitations called")
		accountId, _ := cmd.Flags().GetInt("accountid")
		getInvitations(accountId)
	},
}

func init() {
	GetCmd.AddCommand(invitationsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// invitationsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// invitationsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getInvitations(accountId int) {
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
