/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package add

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// oplCmd represents the opl command
var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Add agent into an OPL",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		harbourId, _ := cmd.Flags().GetString("hid")
		rawOutput, _ := cmd.Flags().GetBool("raw")
		if harbourId != "" && rawOutput {
			addAgentraw(harbourId)
		} else if harbourId != "" {
			addAgent(harbourId)
		} else {
			fmt.Println("\nPlease provide a correct Harbour ID to add agent")
			fmt.Println("[bmgo add agent --hid <harbour id>]")
		}
	},
}

func init() {
	AddCmd.AddCommand(agentCmd)
	agentCmd.Flags().String("hid", "", "Provide Harbour ID")
	agentCmd.MarkFlagRequired("hid")
}

func agentNamePrompt() string {
	validate := func(input string) error {
		if len(input) <= 2 {
			return errors.New("invalid name")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:       "Name your agent:",
		HideEntered: true,
		Validate:    validate,
	}
	resultWsId, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}
	return resultWsId
}

type addAgentResponse struct {
	Result addAgentResult `json:"result"`
}
type addAgentResult struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	State string `json:"state"`
}
type getAgentCmdResponse struct {
	Result getAgentcmdResult `json:"result"`
}
type getAgentcmdResult struct {
	DockerCommand string `json:"dockerCommand"`
}

func addAgent(harbourId string) {
	agentName := agentNamePrompt()
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	data := fmt.Sprintf(`{"name":"%s","address":"127.0.0.1"}`, agentName)
	reqBodydata := strings.NewReader(data)
	req, err := http.NewRequest("POST", "https://a.blazemeter.com/api/v4/private-locations/"+harbourId+"/servers", reqBodydata)
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
	var responseBodyaddAgent addAgentResponse
	json.Unmarshal(bodyText, &responseBodyaddAgent)

	fmt.Printf("\n%-30s %-20s %-10s\n", "SHIP-ID", "NAME", "STATE")
	shipId := responseBodyaddAgent.Result.Id
	shipName := responseBodyaddAgent.Result.Name
	shipstate := responseBodyaddAgent.Result.State
	fmt.Printf("\n%-30s %-20s %-10s\n", shipId, shipName, shipstate)

	req1, err := http.NewRequest("POST", "https://a.blazemeter.com/api/v4/private-locations/"+harbourId+"/ships/"+shipId+"/docker-command", nil)
	if err != nil {
		log.Fatal(err)
	}
	req1.Header.Set("Content-Type", "application/json")
	req1.SetBasicAuth(apiId, apiSecret)
	resp1, err := client.Do(req1)
	if err != nil {
		log.Fatal(err)
	}
	defer resp1.Body.Close()
	bodyText1, err := io.ReadAll(resp1.Body)
	if err != nil {
		log.Fatal(err)
	}
	//	fmt.Printf("%s\n", bodyText1)
	var responseBodyaddAgentCmd getAgentCmdResponse
	json.Unmarshal(bodyText1, &responseBodyaddAgentCmd)
	dockerRun := responseBodyaddAgentCmd.Result.DockerCommand
	fmt.Println("\n---------------------------------------------------------------------------------------------")
	fmt.Printf("Docker RUN COMMAND:\n %s\n", dockerRun)
	fmt.Println("\n-")
}
func addAgentraw(harbourId string) {
	agentName := agentNamePrompt()
	apiId, apiSecret := Getapikeys()
	client := &http.Client{}
	data := fmt.Sprintf(`{"name":"%s","address":"127.0.0.1"}`, agentName)
	reqBodydata := strings.NewReader(data)
	req, err := http.NewRequest("POST", "https://a.blazemeter.com/api/v4/private-locations/"+harbourId+"/servers", reqBodydata)
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
	var responseBodyaddAgent addAgentResponse
	json.Unmarshal(bodyText, &responseBodyaddAgent)

	fmt.Printf("\n%-30s %-20s %-10s\n", "SHIP-ID", "NAME", "STATE")
	shipId := responseBodyaddAgent.Result.Id
	fmt.Printf("%s\n", bodyText)

	req1, err := http.NewRequest("POST", "https://a.blazemeter.com/api/v4/private-locations/"+harbourId+"/ships/"+shipId+"/docker-command", nil)
	if err != nil {
		log.Fatal(err)
	}
	req1.Header.Set("Content-Type", "application/json")
	req1.SetBasicAuth(apiId, apiSecret)
	resp1, err := client.Do(req1)
	if err != nil {
		log.Fatal(err)
	}
	defer resp1.Body.Close()
	bodyText1, err := io.ReadAll(resp1.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n%s\n", bodyText1)
}
