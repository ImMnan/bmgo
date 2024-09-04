/*
Copyright © 2024 Manan Patel - github.com/immnan
*/
package get

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var bucketsCmd = &cobra.Command{
	Use:   "buckets",
	Short: "[>]Use find command to free-search resources in Blazemeters",
	Long: `The command returns a list of buckets in an API Monitoring team, you will need to provide a team id to run the command.
	
	For example: [bmgo get -t <team id> buckets ] OR 
	             [bmgo get --team <team id> buckets]
	For default: [bmgo get --tm buckets]]`,
	Run: func(cmd *cobra.Command, args []string) {
		var teamId string
		tm, _ := cmd.Flags().GetBool("tm")
		if tm {
			teamId = defaultTeam()
		} else {
			teamId, _ = cmd.Flags().GetString("teamid")
		}
		rawOutput, _ := cmd.Flags().GetBool("raw")
		switch {
		case (teamId != ""):
			getTeamInfo(teamId, rawOutput)
		default:
			cmd.Help()
		}
	},
}

func init() {
	GetCmd.AddCommand(bucketsCmd)
	bucketsCmd.Flags().StringP("teamid", "t", "", "[>]Provide Team UID to list resources within a team")
	bucketsCmd.Flags().Bool("tm", false, "[>]Use default team UId (bmConfig.yaml)")
}

type teamInfo struct {
	Name         string       `json:"name"`
	CreatedAt    string       `json:"created_at"`
	UserCount    int          `json:"user_count"`
	BucketsCount int          `json:"bucket_count"`
	Owned_by     owned_by     `json:"owned_by"`
	Plan         teamPlan     `json:"plan"`
	Error        errorResult  `json:"error"`
	Buckets      []bucketList `json:"buckets"`
}

type owned_by struct {
	OwnerEmail string `json:"email"`
	OwnerUUID  string `json:"uuid"`
}
type teamPlan struct {
	PlanUUID          string `json:"uuid"`
	Name              string `json:"name"`
	Max_requests      int    `json:"max_requests"`
	Max_collaborators int    `json:"max_collaborators"`
	Max_buckets       int    `json:"max_buckets"`
}
type bucketList struct {
	Key          string `json:"key"`
	Name         string `json:"name"`
	TriggerToken string `json:"trigger_token"`
	IsPublic     bool   `json:"is_public"`
}

func getTeamInfo(teamId string, rawOutput bool) {
	Bearer := fmt.Sprintf("Bearer %v", GetPersonalAccessToken())
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.runscope.com/teams/"+teamId, nil)
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
	if rawOutput {
		fmt.Printf("%s\n", bodyText)
	} else {
		var responseBodyTeamInfo teamInfo
		json.Unmarshal(bodyText, &responseBodyTeamInfo)
		if responseBodyTeamInfo.Error.Status == 0 {
			tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			// Print headers
			fmt.Fprintln(tabWriter, "BUCKET_KEY\tNAME\tPUBLIC\tTRIGGER_TOKEN")
			for i := 0; i < len(responseBodyTeamInfo.Buckets); i++ {
				bucketName := responseBodyTeamInfo.Buckets[i].Name
				bucketKey := responseBodyTeamInfo.Buckets[i].Key
				bucketPublic := responseBodyTeamInfo.Buckets[i].IsPublic
				bucketTriggerToken := responseBodyTeamInfo.Buckets[i].TriggerToken
				//fmt.Printf("\n%-15s %-35s %-8t %-20s", bucketKey, bucketName, bucketPublic, bucketTriggerToken)
				fmt.Fprintf(tabWriter, "%s\t%s\t%t\t%s\n", bucketKey, bucketName, bucketPublic, bucketTriggerToken)
			}
			tabWriter.Flush()
			fmt.Println("-")
		} else {
			errorCode := responseBodyTeamInfo.Error.Status
			errorMessage := responseBodyTeamInfo.Error.Message
			fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
		}
	}
}
