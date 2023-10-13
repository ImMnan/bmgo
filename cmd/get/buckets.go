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

// bucketsCmd represents the buckets command
var bucketsCmd = &cobra.Command{
	Use:   "buckets",
	Short: "Get list of buckets in team",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("buckets called")
		getBuckets()
	},
}

func init() {
	GetCmd.AddCommand(bucketsCmd)
}

func getBuckets() {
	//	pat := GetPersonalAccessToken()
	Bearer := fmt.Sprintf("Bearer %v", GetPersonalAccessToken())
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.runscope.com/buckets", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", Bearer)
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
