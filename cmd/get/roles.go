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

// rolesCmd represents the roles command
var rolesCmd = &cobra.Command{
	Use:   "roles",
	Short: "Get a list of Roles in a team",
	Long:  `.`,
	Run: func(cmd *cobra.Command, args []string) {
		tm, _ := cmd.Flags().GetBool("tm")
		var teamId string
		if tm {
			teamId = defaultTeam()
		} else {
			teamId, _ = cmd.Flags().GetString("teamid")
		}
		//rawOutput, _ := cmd.Flags().GetBool("raw")
		getRolesTm(teamId)
	},
}

func init() {
	GetCmd.AddCommand(rolesCmd)
}

func getRolesTm(teamId string) {
	bearer := fmt.Sprintf("Bearer %v", GetPersonalAccessToken())
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.runscope.com/teams/"+teamId+"/roles", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", bearer)
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
