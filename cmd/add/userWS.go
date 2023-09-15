/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package add

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// userCmd represents the user command
var userWsCmd = &cobra.Command{
	Use:   "user",
	Short: "Get details about the user",
	Long:  `.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("user find called")
		userId, _ := cmd.Flags().GetInt("uid")
		workspaceId, _ := cmd.Flags().GetInt("workspaceid")
		//rawOutput, _ := cmd.Flags().GetBool("raw")

		//		if rawOutput {
		//	addUserByUidRaw(userId, workspaceId)
		//	} else {
		addUserByUid(userId, workspaceId)
		//	}
	},
}

func init() {
	AddCmd.AddCommand(userWsCmd)
	userWsCmd.Flags().Int("uid", 0, "User ID for the user")
	userWsCmd.MarkFlagRequired("uid")
	userWsCmd.Flags().BoolP("raw", "r", false, "[Optional] If set, the output will be raw json")
}

func addUserByUid(userId, workspaceId int) {
	apiId, apiSecret := Getapikeys()
	workspaceIdStr := strconv.Itoa(workspaceId)
	client := &http.Client{}
	//	var data = strings.NewReader(`{"usersIds":[%v],"roles": ["manager"]}`)
	data := fmt.Sprintf(`{"usersIds":[%v],"roles": ["manager"]}`, userId)
	var bodyData = strings.NewReader(data)
	req, err := http.NewRequest("POST", "https://a.blazemeter.com/api/v4/workspaces/"+workspaceIdStr+"/users", bodyData)
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
