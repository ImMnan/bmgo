/*
Copyright © 2024 Manan Patel - github.com/immnan
*/
package find

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// findCmd represents the find command
var FindCmd = &cobra.Command{
	Use:   "find",
	Short: "Use find command to free-search resources in Blazemeters",
	Long: `Use get command to Find details about the resources within Blazemeter account. Use --help throughout subcommands to get an idea of how these commands are structured. Navigate to help.blazemeter.com for detailed info about these resources. 
	
	For defaults: Make sure you have a file named bmConfig.yaml specifying the defaults.
	The file can be in the current working DIR or in $HOME
	
	Also, use this command to get information about a workspace, team or an account. The command lists number of users, plan details, owners, status, etc. 
	
	For example: [bmgo find -a <account_id>] OR 
	             [bmgo find -w <workspace_id>] OR
				 [bmgo find -t <team_uuid>]
	For default: [bmgo find --ac] OR 
	             [bmgo find --ws] OR
				 [bmgo find --tm]`,
	Run: func(cmd *cobra.Command, args []string) {
		ac, _ := cmd.Flags().GetBool("ac")
		ws, _ := cmd.Flags().GetBool("ws")
		tm, _ := cmd.Flags().GetBool("tm")

		testId, _ := cmd.Flags().GetInt("tid")
		var accountId, workspaceId int
		var teamId string
		if ac {
			accountId = defaultAccount()
		} else {
			accountId, _ = cmd.Flags().GetInt("accountid")
		}
		if ws {
			workspaceId = defaultWorkspace()
		} else {
			workspaceId, _ = cmd.Flags().GetInt("workspaceid")
		}
		if tm {
			teamId = defaultTeam()
		} else {
			teamId, _ = cmd.Flags().GetString("teamid")
		}
		rawOutput, _ := cmd.Flags().GetBool("raw")
		switch {
		case (workspaceId != 0) && (accountId == 0) && teamId == "":
			getWorkspace(workspaceId, rawOutput)
		case (accountId != 0) && (workspaceId == 0) && teamId == "":
			getAccountId(accountId, rawOutput)
		case (accountId == 0) && (workspaceId == 0) && teamId != "":
			getTeamInfo(teamId, rawOutput)
		case testId != 0:
			findTest(testId, rawOutput)
			findTestFiles(testId, rawOutput)
		default:
			cmd.Help()
		}
	},
}

func init() {
	FindCmd.PersistentFlags().BoolP("raw", "r", false, "[Optional] If set, the output will be raw json")
	FindCmd.Flags().IntP("accountid", "a", 0, "Provide Account ID to list a resources within an account")
	FindCmd.Flags().IntP("workspaceid", "w", 0, "Provide Workspace ID to list a resource within a workspace")
	FindCmd.Flags().Bool("ac", false, "Use default account Id (bmConfig.yaml)")
	FindCmd.Flags().Bool("ws", false, "Use default workspace Id (bmConfig.yaml)")
	FindCmd.Flags().StringP("teamid", "t", "", "[>]Provide Team UID to list resources within a team")
	FindCmd.Flags().Bool("tm", false, "[>]Use default team UId (bmConfig.yaml)")
	FindCmd.Flags().Int("tid", 0, "Provide the Test ID to view the Test details")
}

func Getapikeys() (string, string) {
	vp := viper.New()
	vp.SetConfigName("bmConfig")
	vp.SetConfigType("yaml")
	vp.AddConfigPath("$HOME")
	//	vp.AddConfigPath(".")
	err := vp.ReadInConfig()
	if err != nil {
		log.Fatal(err, "\nPlease add your Blazemeter configurations to bmConfig.yaml file in your home directory")
	}
	apiId := vp.GetString("id")
	apiSecret := vp.GetString("secret")
	return apiId, apiSecret
}

func defaultAccount() int {
	vp := viper.New()
	vp.SetConfigName("bmConfig")
	vp.SetConfigType("yaml")
	//	vp.AddConfigPath(".")
	vp.AddConfigPath("$HOME")
	err := vp.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	accountId := vp.GetInt("accountId")
	return accountId
}

// Getting default workspace Id in case user uses --ws
func defaultWorkspace() int {
	vp := viper.New()
	vp.SetConfigName("bmConfig")
	vp.SetConfigType("yaml")
	//	vp.AddConfigPath(".")
	vp.AddConfigPath("$HOME")
	err := vp.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	workspaceId := vp.GetInt("workspaceId")
	return workspaceId
}
func GetPersonalAccessToken() string {
	vp := viper.New()
	vp.SetConfigName("bmConfig")
	vp.SetConfigType("yaml")
	vp.AddConfigPath("$HOME")
	//	vp.AddConfigPath(".")
	err := vp.ReadInConfig()
	if err != nil {
		log.Fatal(err, "\nPlease add your Blazemeter API credentials to bmConfig.yaml file in your home directory")
	}
	pat := vp.GetString("pat")
	return pat
}
func defaultTeam() string {
	vp := viper.New()
	vp.SetConfigName("bmConfig")
	vp.SetConfigType("yaml")
	//vp.AddConfigPath(".")
	vp.AddConfigPath("$HOME")
	err := vp.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	teamId := vp.GetString("teamId")
	return teamId
}

type errorResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func promtHid() string {
	validate := func(input string) error {
		if len(input) <= 20 {
			return errors.New("invalid crone")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:       "Provide Harbour Id: ",
		HideEntered: true,
		Validate:    validate,
	}
	resultHid, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}
	return resultHid
}
func workspaceIdPrompt() string {
	validate := func(input string) error {
		_, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return errors.New("invalid workspace")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:       "Provide Workspace Id: ",
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
