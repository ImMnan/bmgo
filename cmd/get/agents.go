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
	"time"

	"github.com/spf13/cobra"
)

// agentsCmd represents the agents command
var agentsCmd = &cobra.Command{
	Use:   "agents",
	Short: "Get agents within a private location",
	Long: `The command returns a list of created agents, you will need to provide a workspace id or a harborId to run the command. any server on which you install our agent is an agent within a Private location. These are your load generators. Formerly known as a 'ship'. The command returns a list of agents within a workspace or within a harborId id specified.  Outputs "SHIP ID", "STATE", etc.
	
	For example: [bmgo get -w <workspace id> agents --hid <harbour_id>] OR 
	             [bmgo get agents <harbour_id> --raw]
	             [bmgo get -w <workspace_id> agents]
	For default: [bmgo get --ws agents]
	             [bmgo get --ws agents --hid <harbour id>]`,
	Run: func(cmd *cobra.Command, args []string) {
		ws, _ := cmd.Flags().GetBool("ws")
		var workspaceId int
		if ws {
			workspaceId = defaultWorkspace()
		} else {
			workspaceId, _ = cmd.Flags().GetInt("workspaceid")
		}
		rawOutput, _ := cmd.Flags().GetBool("raw")
		harbourId, _ := cmd.Flags().GetString("hid")
		switch {
		case workspaceId == 0 && harbourId != "" && rawOutput:
			getAgentsOplraw(harbourId)
		case workspaceId != 0 && harbourId != "" && rawOutput:
			getAgentsOplraw(harbourId)
		case workspaceId != 0 && harbourId == "" && rawOutput:
			getAgentsWsraw(workspaceId)
		case workspaceId != 0 && harbourId == "":
			getAgentsWs(workspaceId)
		case workspaceId != 0 && harbourId != "":
			getAgentsOpl(workspaceId, harbourId)
		default:
			cmd.Help()
		}
	},
}

func init() {
	GetCmd.AddCommand(agentsCmd)
	agentsCmd.Flags().String("hid", "", "Provide the harbour id")
}

// This function is because the API response for listing agents using harbour id is a struct & not a list/array to iterate over.
func getAgentsOpl(workspaceId int, harbourId string) {
	apiId, apiSecret := Getapikeys()
	workspaceIdStr := strconv.Itoa(workspaceId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/private-locations?workspaceId="+workspaceIdStr+"&limit=0", nil)
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
	var responseBodyWsAgents oplsResponse
	json.Unmarshal(bodyText, &responseBodyWsAgents)
	if responseBodyWsAgents.Error.Code == 0 {
		for i := 0; i < len(responseBodyWsAgents.Result); i++ {
			oplId := responseBodyWsAgents.Result[i].Id
			if oplId == harbourId {
				fmt.Printf("For OPL/HARBOUR %v & NAMED %v:\n", oplId, responseBodyWsAgents.Result[i].Name)
				fmt.Printf("\n%-28s %-8s %-18s %-10s\n", "SHIP ID", "STATE", "LAST BEAT", "NAME")
				for f := 0; f < len(responseBodyWsAgents.Result[i].Ships); f++ {
					shipId := responseBodyWsAgents.Result[i].Ships[f].Id
					shipName := responseBodyWsAgents.Result[i].Ships[f].Name
					shipStatus := responseBodyWsAgents.Result[i].Ships[f].State
					shipLastHeartBeatEp := int64(responseBodyWsAgents.Result[i].Ships[f].LastHeartBeat)
					//	shipLastHeartBeat := time.Unix(shipLastHeartBeatEp, 0)
					if shipLastHeartBeatEp != 0 {
						shipLastHeartBeatStr := fmt.Sprint(time.Unix(shipLastHeartBeatEp, 0))
						fmt.Printf("\n%-28s %-8s %-18s %-10s", shipId, shipStatus, shipLastHeartBeatStr[0:16], shipName)
					} else {
						shipLastHeartBeat := shipLastHeartBeatEp
						fmt.Printf("\n%-28s %-8s %-18d %-10s", shipId, shipStatus, shipLastHeartBeat, shipName)
					}
				}
				fmt.Println("\n\n---------------------------------------------------------------------------------------------")
			} else {
				continue
			}
		}
	} else {
		errorCode := responseBodyWsAgents.Error.Code
		errorMessage := responseBodyWsAgents.Error.Message
		fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
	}
}
func getAgentsOplraw(harbourId string) {
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/private-locations/"+harbourId+"/servers", nil)
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

func getAgentsWs(workspaceId int) {
	apiId, apiSecret := Getapikeys()
	workspaceIdStr := strconv.Itoa(workspaceId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/private-locations?workspaceId="+workspaceIdStr+"&limit=0", nil)
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
	var responseBodyWsAgents oplsResponse
	json.Unmarshal(bodyText, &responseBodyWsAgents)
	if responseBodyWsAgents.Error.Code == 0 {
		for i := 0; i < len(responseBodyWsAgents.Result); i++ {
			fmt.Printf("For OPL/HARBOUR %v & NAMED %v:\n", responseBodyWsAgents.Result[i].Id, responseBodyWsAgents.Result[i].Name)
			fmt.Printf("\n%-28s %-8s %-18s %-10s\n", "SHIP ID", "STATE", "LAST BEAT", "NAME")
			for f := 0; f < len(responseBodyWsAgents.Result[i].Ships); f++ {
				shipId := responseBodyWsAgents.Result[i].Ships[f].Id
				shipName := responseBodyWsAgents.Result[i].Ships[f].Name
				shipStatus := responseBodyWsAgents.Result[i].Ships[f].State
				shipLastHeartBeatEp := int64(responseBodyWsAgents.Result[i].Ships[f].LastHeartBeat)
				//	shipLastHeartBeat := time.Unix(shipLastHeartBeatEp, 0)
				if shipLastHeartBeatEp != 0 {
					shipLastHeartBeatStr := fmt.Sprint(time.Unix(shipLastHeartBeatEp, 0))
					fmt.Printf("\n%-28s %-8s %-18s %-10s", shipId, shipStatus, shipLastHeartBeatStr[0:16], shipName)
				} else {
					shipLastHeartBeat := shipLastHeartBeatEp
					fmt.Printf("\n%-28s %-8s %-18d %-10s", shipId, shipStatus, shipLastHeartBeat, shipName)
				}
			}
			fmt.Println("\n\n---------------------------------------------------------------------------------------------")
		}
		fmt.Println("\n-")
	} else {
		errorCode := responseBodyWsAgents.Error.Code
		errorMessage := responseBodyWsAgents.Error.Message
		fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
	}
}
func getAgentsWsraw(workspaceId int) {
	apiId, apiSecret := Getapikeys()
	workspaceIdStr := strconv.Itoa(workspaceId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/private-locations?workspaceId="+workspaceIdStr+"&limit=0", nil)
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
