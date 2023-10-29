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
	Long: `Use get command to list resources in your Blazemeter account. Use --help throughout subcommands to get an idea of how these commands are structured. Navigate to help.blazemeter.com for detailed info about these resources. 
	
	For defaults: Make sure you have a file named bmConfig.yaml specifying the defaults.
	The file can be in working DIR or in $HOME`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	GetCmd.PersistentFlags().IntP("accountid", "a", 0, "Provide Account ID to list a resources within an account")
	GetCmd.PersistentFlags().IntP("workspaceid", "w", 0, "Provide Workspace ID to list a resource within a workspace")
	GetCmd.PersistentFlags().BoolP("raw", "r", false, "[Optional] If set, the output will be raw json")
	GetCmd.PersistentFlags().Bool("ac", false, "Use default account Id (bmConfig.yaml)")
	GetCmd.PersistentFlags().Bool("ws", false, "Use default workspace Id (bmConfig.yaml)")
	GetCmd.PersistentFlags().StringP("teamid", "t", "", "[#] Provide Team UID to list resources within a team")
	GetCmd.PersistentFlags().Bool("tm", false, "Use default team UId (bmConfig.yaml)")
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// Error handling struct
type errorResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// Getting the API Keys
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
func GetPersonalAccessToken() string {
	vp := viper.New()
	vp.SetConfigName("bmConfig")
	vp.SetConfigType("yaml")
	vp.AddConfigPath("$HOME")
	//	vp.AddConfigPath(".")
	err := vp.ReadInConfig()
	if err != nil {
		log.Fatal(err, "\nPlease add your Blazemeter configurations to bmConfig.yaml file in your home directory")
	}
	pat := vp.GetString("pat")
	return pat
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
func defaultTeam() string {
	vp := viper.New()
	vp.SetConfigName("bmConfig")
	vp.SetConfigType("yaml")
	vp.AddConfigPath(".")
	vp.AddConfigPath("$HOME")
	err := vp.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	teamId := vp.GetString("teamId")
	return teamId
}

// Helper functions added here
// Removing duplicates
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
