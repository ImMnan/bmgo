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

// dedicatedipCmd represents the dedicatedip command
var dedicatedipCmd = &cobra.Command{
	Use:   "dedicatedip",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("dedicatedip called")
		accountId, _ := cmd.Flags().GetInt("id")
		getDedicatedIp(accountId)
	},
}

func init() {
	accountCmd.AddCommand(dedicatedipCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dedicatedipCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dedicatedipCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type workspaceList struct {
	Result []workspaceResult `json:"result"`
	Name   string            `json:"name"`
	Total  int               `JSON:"total"`
}

type workspaceResult struct {
	Id     int `json:"id"`
	Userid int `json:"userid"`
}

func getDedicatedIps(accountId int) []int {
	apiId, apiSecret := Getapikeys()

	client := &http.Client{}
	accountIdStr := strconv.Itoa(accountId)
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/workspaces?accountId="+accountIdStr+"&limit=200", nil)
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
	slice := []int{}
	var responseObject workspaceList
	json.Unmarshal(bodyText, &responseObject)
	fmt.Println("Total workspace in this account: ", len(responseObject.Result))
	for i := 0; i < len(responseObject.Result); i++ {
		arr := responseObject.Result[i].Id
		//	workspaceIdStr := strconv.Itoa(arr)
		slice = append(slice, arr)
	}
	// Append element 4 to slice
	fmt.Println(slice) // [1 2 3 4]
	return slice
}

type workspaceList1 struct {
	Result []ipResult `json:"result"`
	Total  int        `JSON:"total"`
}

type ipResult struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

func getDedicatedIp(accountId int) {
	workspaceIds := getDedicatedIps(accountId)
	fmt.Println(accountId)
	apiId, apiSecret := Getapikeys()
	slice := []int{}

	for i := 0; i < len(workspaceIds); i++ {
		workspaceIdStr := strconv.Itoa(workspaceIds[i])
		client := &http.Client{}
		req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/workspaces/"+workspaceIdStr+"/users", nil)
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
		//	fmt.Printf("%s\n", bodyText)
		var responseObject workspaceList1
		json.Unmarshal(bodyText, &responseObject)

		for i := 0; i < len(responseObject.Result); i++ {
			userArr := responseObject.Result[i].Id
			slice = append(slice, userArr)
		}
		//fmt.Println("Total users in ", workspaceIds[i], responseObject.Total,)

	}
	fmt.Println(slice)
}
