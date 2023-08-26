/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package account

import (
	"fmt"

	"github.com/spf13/cobra"
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
