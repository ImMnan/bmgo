/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package get

import (
	"fmt"

	"github.com/spf13/cobra"
)

// accountCmd represents the account command
var demoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Demo test",
	Long:  `.`,
	Run: func(cmd *cobra.Command, args []string) {
		//accountId, _ := cmd.Flags().GetInt("accountid")
		//		workspaceId, _ := cmd.Flags().GetInt("workspaceid")
		ac, _ := cmd.Flags().GetBool("ac")
		//		ws, _ := cmd.Flags().GetBool("ws")
		var accountId int
		if ac {
			accountId = defaultAccount()
		} else {
			accountId, _ = cmd.Flags().GetInt("accountid")
		}
		getDemoac(accountId)
	},
}

func init() {
	GetCmd.AddCommand(demoCmd)
	demoCmd.Flags().Bool("ac", false, "look into account")
	demoCmd.Flags().Bool("ws", false, "look into workspace")
}

func getDemoac(accountId int) {
	fmt.Println(accountId)
}
