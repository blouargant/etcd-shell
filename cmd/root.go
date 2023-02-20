/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"etcd-shell/shell"
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
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			if args[0] != "completion" {
				runShell()
			}
		} else {
			runShell()
		}
	},
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	//cobra.OnInitialize(initConfig)
	RootCmd.Flags().StringVarP(&shell.Endpointlist, "endpoints", "e", "", "Comma separated list of endpoints")
	RootCmd.Flags().StringVarP(&shell.User, "user", "u", "", "ETCD user")
	RootCmd.Flags().StringVarP(&shell.Password, "password", "p", "", "ETCD password")
	RootCmd.Flags().BoolVarP(&shell.UseTls, "tls", "t", false, "Use TLS for ETCD connection")
}
