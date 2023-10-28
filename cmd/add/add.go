/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package add

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

// AddCmd represents the add command
var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Use Add command for Adding resources to Blazemeter",
	Long: `Use Add command to add/create resources in your Blazemeter account. Use --help throughout subcommands to get an idea of how these commands are structured, though they are all straight forward. Add command has a lot of prompts from the bmgo, so make sure you have the information handy to respond to these prompts
	
	For defaults: Make sure you have a file named bmConfig.yaml specifying the defaults.
	The file can be in working DIR or in $HOME`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	AddCmd.PersistentFlags().IntP("accountid", "a", 0, " Provide Account ID to add a resource to")
	AddCmd.PersistentFlags().IntP("workspaceid", "w", 0, " Provide Workspace ID to add a resource to")
	AddCmd.PersistentFlags().BoolP("raw", "r", false, "[OPTIONAL] If set, the output will be raw json")
	AddCmd.PersistentFlags().Bool("ac", false, "Use default account Id (bmConfig.yaml)")
	AddCmd.PersistentFlags().Bool("ws", false, "Use default workspace Id (bmConfig.yaml)")
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type errorResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
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

// Getting default account Id in case user uses --ac
func defaultAccount() int {
	vp := viper.New()
	vp.SetConfigName("bmConfig")
	vp.SetConfigType("yaml")
	vp.AddConfigPath(".")
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
	vp.AddConfigPath(".")
	vp.AddConfigPath("$HOME")
	err := vp.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	workspaceId := vp.GetInt("workspaceId")
	return workspaceId
}

// Helper functions added here
// Prompt to user-input for agent name
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
	resultAgentName, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}
	return resultAgentName
}

// Prompt to user to configure OPL
func oplconfigPrompt() (int, int) {
	validate := func(input string) error {
		_, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return errors.New("invalid number or nan")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:       "Number of Engines per Agent",
		HideEntered: true,
		Validate:    validate,
	}
	prompt1 := promptui.Prompt{
		Label:       "Number of Threads per Engine",
		HideEntered: true,
		Validate:    validate,
	}
	resultEPAstr, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}
	resultTPEstr, err := prompt1.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}
	resultEPA, _ := strconv.Atoi(resultEPAstr)
	resultTPE, _ := strconv.Atoi(resultTPEstr)
	return resultEPA, resultTPE
}

// Prompt to user-input for cron expression
func cronPrompt() string {
	validate := func(input string) error {
		if len(input) <= 8 {
			return errors.New("invalid crone")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:       "Cron [example: 0 0 * * 1-5 ]: ",
		HideEntered: true,
		Validate:    validate,
	}
	resultCronEx, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}
	return resultCronEx
}

// Prompt to user-input for selection of account level user roles
func userRoleSelectorA() (string, bool) {
	prompt := promptui.Select{
		Label:        "Select Account Role",
		Items:        []string{"admin", "standard", "user_manager", "billing"},
		HideSelected: true,
	}
	prompt1 := promptui.Select{
		Label:        "attachAutomatically",
		Items:        []bool{true, false},
		HideSelected: true,
	}
	_, roleSelected, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}
	_, attachAuto, err := prompt1.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}
	boolVal, _ := strconv.ParseBool(attachAuto)
	return roleSelected, boolVal
}

// Prompt to user-input for selection of workspace level user roles
func userRoleSelectorWs() string {
	prompt := promptui.Select{
		Label:        "Select Workspace Role",
		Items:        []string{"tester", "manager", "viewer"},
		HideSelected: true,
	}
	_, roleSelected, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}
	return roleSelected
}

// Prompt to user for workspace ids
func workspaceIdPrompt() string {
	prompt := promptui.Prompt{
		Label:       "Provide Workspace/s[separated by commas]",
		HideEntered: true,
	}
	resultWsId, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}
	return resultWsId
}
