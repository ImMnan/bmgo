/*
Copyright © 2024 Manan Patel - github.com/immnan
*/
package find

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

// userCmd represents the user command
var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Find details about the agent",
	Long: ` Use the command to find agents within a specified Private location (Harbour Id). Provide the agent id (Ship Id) to find the agent. Global Flag --raw can be used for raw Json output. 
	For example: [bmgo find agent --aid <agent/ship ID>]
	             [bmgo find agent --aid <agent/ship ID> --hid <harbour ID>]`,
	Run: func(cmd *cobra.Command, args []string) {
		agentId, _ := cmd.Flags().GetString("aid")
		rawOutput, _ := cmd.Flags().GetBool("raw")
		harbourId, _ := cmd.Flags().GetString("hid")
		if agentId != "" {
			findAgent(agentId, harbourId, rawOutput)
		} else {
			cmd.Help()
		}
	},
}

func init() {
	FindCmd.AddCommand(agentCmd)
	agentCmd.Flags().String("aid", "", "Provide the Agent ID or called Ship ID")
	agentCmd.MarkFlagRequired("aid")
	agentCmd.Flags().String("hid", "", "Provide the harbour ID, if not provided through flag the cli will ask in prompt")
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

func findAgent(agentId, harbourId string, rawOutput bool) {
	apiId, apiSecret := Getapikeys()
	var resp *http.Response
	client := &http.Client{}
	if harbourId != "" {
		req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/private-locations/"+harbourId+"/servers/"+agentId, nil)
		if err != nil {
			log.Fatal(err)
		}
		req.SetBasicAuth(apiId, apiSecret)
		resp, err = client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		harbourId = promtHid()
		req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/private-locations/"+harbourId+"/servers/"+agentId, nil)
		if err != nil {
			log.Fatal(err)
		}
		req.SetBasicAuth(apiId, apiSecret)
		resp, err = client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if rawOutput {
		fmt.Printf("%s\n", bodyText)
	} else {
		var responseBodyAgent agentResponse
		json.Unmarshal(bodyText, &responseBodyAgent)
		if responseBodyAgent.Error.Code == 0 {
			//	fmt.Printf("\n%-20s %-8s %-20s %-10s\n", "NAME", "STATE", "LAST HEART_BEAT", "CRANE")
			tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			// Print headers
			fmt.Fprintln(tabWriter, "NAME\tSTATE\tLAST_BEAT\tCRANE")

			agentName := responseBodyAgent.Result.Name
			agentState := responseBodyAgent.Result.State
			agentHbEp := int64(responseBodyAgent.Result.LastHeartBeat)
			agentHbEpStr := fmt.Sprint(time.Unix(agentHbEp, 0))
			craneVersion := responseBodyAgent.Result.Crane
			//fmt.Printf("%-20s %-8v %-20v %-10v\n", agentName, agentState, agentHbEpStr[0:16], craneVersion)
			fmt.Fprintf(tabWriter, "%s\t%s\t%s\t%s\n", agentName, agentState, agentHbEpStr[0:16], craneVersion)
			tabWriter.Flush()
			fmt.Println("\n-")
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
			fmt.Println("-")
		} else {
			errorCode := responseBodyAgent.Error.Code
			errorMessage := responseBodyAgent.Error.Message
			fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
		}
	}
}
