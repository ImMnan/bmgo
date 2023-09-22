/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package update

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// updateCmd represents the update command
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Use update command to modify exisiting resources in Blazemeter",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	UpdateCmd.PersistentFlags().IntP("accountid", "a", 0, " [REQUIRED] Account ID of the resource expected to being updated")
	UpdateCmd.PersistentFlags().IntP("workspaceid", "w", 0, " [REQUIRED] Workspace ID of the resource expected to being updated")
	UpdateCmd.PersistentFlags().BoolP("raw", "r", false, "[Optional] If set, the output will be raw json")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
