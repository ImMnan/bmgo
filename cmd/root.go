/*
Copyright © 2024 Manan Patel - github.com/immnan
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/immnan/bmgo/cmd/add"
	"github.com/immnan/bmgo/cmd/find"
	"github.com/immnan/bmgo/cmd/get"
	"github.com/immnan/bmgo/cmd/update"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bmgo",
	Short: "A Blazemeter CLI tool for Blazemeter Admins",
	Long:  "A Cli tool used for performing actions within Blazemeter platform",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		version, _ := cmd.Flags().GetBool("version")
		license, _ := cmd.Flags().GetBool("license")
		config, _ := cmd.Flags().GetBool("config")
		if version {
			fmt.Println("Version: bmgo-0.1.1_beta")
			fmt.Println("Author:  https://github.com/ImMnan/")
			fmt.Println(`
bmgo  Copyright (C) 2024 Manan Patel
This program comes with ABSOLUTELY NO WARRANTY; 
This is free software, and you are welcome to redistribute it
under certain conditions; type "bmgo --license" for details.
	`)
		} else if license {
			fmt.Println(`bmgo - An OpenSource CLI tool for Blazemeter Admins
Copyright (C) 2024  Manan Patel

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details. 

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.`)
		} else if config {
			writeConfig()
		} else {
			cmd.Help()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func addSubCommand() {
	rootCmd.AddCommand(get.GetCmd)
	rootCmd.AddCommand(find.FindCmd)
	rootCmd.AddCommand(add.AddCmd)
	rootCmd.AddCommand(update.UpdateCmd)
}
func init() {
	//	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bmgo.yaml)")
	// rootCmd.PersistentFlags().StringP("author", "a", "Manan Patel", "author name for copyright attribution")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().BoolP("version", "v", false, "Version of the installed bmgo!")
	rootCmd.Flags().BoolP("license", "l", false, "Show license details")
	rootCmd.Flags().BoolP("config", "c", false, "Add config details")
	addSubCommand()
}
