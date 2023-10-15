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

	"github.com/spf13/cobra"
)

// bucketsCmd represents the buckets command
var bucketsCmd = &cobra.Command{
	Use:   "buckets",
	Short: "Get list of buckets in team",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		rawOutput, _ := cmd.Flags().GetBool("raw")
		if rawOutput {
			getBucketsraw()
		} else if !rawOutput {
			getBuckets()
		} else {
			cmd.Help()
		}
	},
}

func init() {
	GetCmd.AddCommand(bucketsCmd)
}

type bucketsResponse struct {
	Data  []bucketsData `json:"data"`
	Error errorResult   `json:"error"`
}

type bucketsData struct {
	Key  string      `json:"key"`
	Name string      `json:"name"`
	Team teamDetails `json:"team"`
}
type teamDetails struct {
	Id   string `json:"id"`
	Name string `json:"name"`
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
	//fmt.Printf("%s\n", bodyText)
	var responseObjectBuckets bucketsResponse
	json.Unmarshal(bodyText, &responseObjectBuckets)
	if responseObjectBuckets.Error.Code == 0 {
		fmt.Printf("\n%-25s %-15s %-10s\n", "NAME", "BUCKET KEY", "TEAM UUID")
		for i := 0; i < len(responseObjectBuckets.Data); i++ {
			bucketName := responseObjectBuckets.Data[i].Name
			bucketKey := responseObjectBuckets.Data[i].Key
			bucketTeamId := responseObjectBuckets.Data[i].Team.Id
			fmt.Printf("\n%-25s %-15s %-10s", bucketName, bucketKey, bucketTeamId)
		}
		fmt.Printf("\n-")
	} else {
		errorCode := responseObjectBuckets.Error.Status
		errorMessage := responseObjectBuckets.Error.Message
		fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
	}
}
func getBucketsraw() {
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
