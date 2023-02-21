/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "etcd-shell",
	Short: "ETCD shell Interface",
	Long:  `etcd-shell is a shell to interact with your ETCD cluster.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {
	// 	if len(args) > 0 {
	// 		if args[0] != "completion" {
	// 			runShell()
	// 		}
	// 	} else {
	// 		runShell()
	// 	}
	// },
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.CompletionOptions.HiddenDefaultCmd = true
}
