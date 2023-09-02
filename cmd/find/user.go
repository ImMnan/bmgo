/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package find

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "User details",
	Short: "Get details about the user",
	Long:  `.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("user find called")
		userEmail, _ := cmd.Flags().GetString("email")
		getUserByEmail(userEmail)
	},
}

func getUserByEmail(userEmail int) {
	apiId, apiSecret := Getapikeys()

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/admin/users?email="+userEmail+, nil)
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
	FindCmd.AddCommand(userCmd)
	userCmd.PersistentFlags().string("email", "-e", "", " [*Required] Confirm the user email address")
	userCmd.MarkPersistentFlagRequired("email")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// userCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
