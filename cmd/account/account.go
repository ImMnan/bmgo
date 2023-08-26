/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package account

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// accountCmd represents the account command
var AccountCmd = &cobra.Command{
	Use:   "account",
	Short: "This will get info about the account",
	Long:  `NOt much to say at least`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("account called")
		accountId, _ := cmd.Flags().GetInt("id")
		getAccountId(accountId)

	},
}

func init() {
	AccountCmd.PersistentFlags().Int("id", 0, "confirm the account id")
	AccountCmd.MarkPersistentFlagRequired("id")
}

func getAccountId(accountId int) {
	fmt.Println("this is the account ID used:", accountId)
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
