/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package find

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// findCmd represents the find command
var FindCmd = &cobra.Command{
	Use:   "find",
	Short: "Use find command to free-search resources in Blazemeters",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	FindCmd.PersistentFlags().BoolP("raw", "r", false, "[Optional] If set, the output will be raw json")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
<<<<<<< Updated upstream
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
=======

type errorResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
>>>>>>> Stashed changes
}
