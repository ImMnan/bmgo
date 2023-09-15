/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package add

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmd represents the get command
var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Use get command Adding resources to Blazemeter",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	AddCmd.PersistentFlags().IntP("accountid", "a", 0, "Provide Account ID to add a resource to")
	AddCmd.PersistentFlags().IntP("workspaceid", "w", 0, "Provide Workspace ID to add a resource to")
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
