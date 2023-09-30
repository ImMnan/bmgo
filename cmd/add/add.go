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
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	AddCmd.PersistentFlags().IntP("accountid", "a", 0, " Provide Account ID to add a resource to")
	AddCmd.PersistentFlags().IntP("workspaceid", "w", 0, " Provide Workspace ID to add a resource to")
	AddCmd.PersistentFlags().BoolP("raw", "r", false, "[OPTIONAL] If set, the output will be raw json")
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func Getapikeys() (string, string) {
	vp := viper.New()
	vp.SetConfigName("api-key")
	vp.SetConfigType("json")
	vp.AddConfigPath(".")
	err := vp.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	apiId := vp.GetString("id")
	apiSecret := vp.GetString("secret")
	return apiId, apiSecret
}

// Functions to support the subcommands
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

func workspaceIdPrompt() string {
	validate := func(input string) error {
		_, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return errors.New("invalid workspace")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:       "Provide Workspace/s-[Array supported]",
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
