/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package find

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

// userCmd represents the user command
var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Find details about the agent",
	Long:  `.`,
	Run: func(cmd *cobra.Command, args []string) {
		agentId, _ := cmd.Flags().GetString("aid")
		rawOutput, _ := cmd.Flags().GetBool("raw")
		if agentId != "" && rawOutput {
			findAgentraw(agentId)
		} else if agentId != "" {
			findAgent(agentId)
		} else {
			fmt.Println("\nPlease provide a correct Harbour ID to find the Private Location - OPL")
			fmt.Println("[bmgo find opl --hid <harbour ID>")
		}
	},
}

func init() {
	FindCmd.AddCommand(agentCmd)
	agentCmd.Flags().String("aid", "", "Provide the Agent ID")
	agentCmd.MarkFlagRequired("aid")
}

type agentResponse struct {
	Result agentResult `json:"result"`
	Error  errorResult `json:"error"`
}
type agentResult struct {
	Name          string   `json:"name"`
	State         string   `json:"state"`
	LastHeartBeat int      `json:"lastHeartBeat"`
	Crane         string   `json:"installedVersion"`
	HostInfo      hostInfo `json:"hostInfo"`
}
type hostInfo struct {
	DiskSpace        diskSpace        `json:"diskSpace"`
	Platform         []string         `json:"platform"`
	ContainerManager containerManager `json:"containerManager"`
}
type diskSpace struct {
	Root root `json:"/"`
}
type root struct {
	FreeSpace   int     `json:"freeSpace"`
	FreePercent float32 `json:"freePercent"`
}
type containerManager struct {
	Type string `json:"type"`
	Info info   `json:"info"`
}
type info struct {
	OperatingSystem string `json:"operatingSystem"`
	Memory          int    `json:"memory"`
	Cpus            int    `json:"cpus"`
}

func findAgent(agentId string) {
	apiId, apiSecret := Getapikeys()
	harbourId := promtHid()
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/private-locations/"+harbourId+"/servers/"+agentId, nil)
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
	var responseBodyAgent agentResponse
	json.Unmarshal(bodyText, &responseBodyAgent)
	if responseBodyAgent.Error.Code == 0 {
		fmt.Printf("\n%-20s %-8s %-20s %-10s\n", "NAME", "STATE", "LAST HEART_BEAT", "CRANE")
		agentName := responseBodyAgent.Result.Name
		agentState := responseBodyAgent.Result.State
		agentHbEp := int64(responseBodyAgent.Result.LastHeartBeat)
		agentHbEpStr := fmt.Sprint(time.Unix(agentHbEp, 0))
		craneVersion := responseBodyAgent.Result.Crane
		fmt.Printf("%-20s %-8v %-20v %-10v\n", agentName, agentState, agentHbEpStr[0:16], craneVersion)

		fmt.Println("\n---------------------------------------------------------------------------------------------")
		agentPlatform := []string{}
		for p := 0; p < len(responseBodyAgent.Result.HostInfo.Platform); p++ {
			agentPlatformArr := responseBodyAgent.Result.HostInfo.Platform[p]
			agentPlatform = append(agentPlatform, agentPlatformArr)
		}
		agentDiskSpace := responseBodyAgent.Result.HostInfo.DiskSpace.Root.FreeSpace
		agentDiskPercent := responseBodyAgent.Result.HostInfo.DiskSpace.Root.FreePercent
		agentType := responseBodyAgent.Result.HostInfo.ContainerManager.Type
		agentMemory := responseBodyAgent.Result.HostInfo.ContainerManager.Info.Memory
		agentCpu := responseBodyAgent.Result.HostInfo.ContainerManager.Info.Cpus
		agentOs := responseBodyAgent.Result.HostInfo.ContainerManager.Info.OperatingSystem
		if agentType == "DockerManager" {
			fmt.Printf("%-15s %-10s %-8v %-15s %-10s %-8s %-10s\n", "PLATFORM", "DISK(GiB)", "DISK(%)", "TYPE", "MEM(GiB)", "CPU", "OS")
			fmt.Printf("%-5s  %-10d %-8v %-15s %-10d %-8d %-10v", agentPlatform[1:], agentDiskSpace/1e+9, int(agentDiskPercent), agentType, agentMemory/1e+9, agentCpu, agentOs)
		} else {
			fmt.Printf("%-15s %-10s %-8v %-15s\n", "PLATFORM", "DISK(GiB)", "DISK(%)", "TYPE")
			fmt.Printf("%-5s  %-10d %-8v %-15s", agentPlatform[1:], agentDiskSpace/1e+9, int(agentDiskPercent), agentType)
		}
		fmt.Println("\n-")
	} else {
		errorCode := responseBodyAgent.Error.Code
		errorMessage := responseBodyAgent.Error.Message
		fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
	}
}
func findAgentraw(agentId string) {
	apiId, apiSecret := Getapikeys()
	harbourId := promtHid()
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/private-locations/"+harbourId+"/servers/"+agentId, nil)
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
