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
	GetCmd.PersistentFlags().BoolP("raw", "r", false, "[Optional] If set, the output will be raw json")
	GetCmd.PersistentFlags().Bool("ac", false, "Use default account Id (defaults.yaml)")
	GetCmd.PersistentFlags().Bool("ws", false, "Use default workspace Id (defaults.yaml)")
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

func defaultAccount() int {
	vp := viper.New()
	vp.SetConfigName("defaults")
	vp.SetConfigType("yaml")
	vp.AddConfigPath(".")
	err := vp.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	accountId := vp.GetInt("accountId")
	return accountId
}
func defaultWorkspace() int {
	vp := viper.New()
	vp.SetConfigName("defaults")
	vp.SetConfigType("yaml")
	vp.AddConfigPath(".")
	err := vp.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	workspaceId := vp.GetInt("workspaceId")
	return workspaceId
}

// Helper functions added here
func removeDuplicateValuesInt(slice []int) []int {
	keys := make(map[int]bool)
	list := []int{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
