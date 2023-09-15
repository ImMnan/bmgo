/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package get

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmd represents the get command
var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Use get command for listing the information",
	Long:  `This is a get subcommad for listing info, related to account info, workspaces, users, invitation etc.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	GetCmd.PersistentFlags().IntP("accountid", "a", 0, " [REQUIRED] Provide Account ID to add a resource to")
	GetCmd.PersistentFlags().IntP("workspaceid", "w", 0, " [REQUIRED] Provide Workspace ID to add a resource to")
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
