/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package connect

import (
	"etcd-shell/cmd"
	"etcd-shell/shell"

	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "connect to etcd",
	Run: func(cmd *cobra.Command, args []string) {
		shell.RunShell()
	},
}

func init() {
	cmd.RootCmd.AddCommand(connectCmd)
	connectCmd.PersistentFlags().StringVarP(&shell.Endpointlist, "endpoints", "e", "", "Comma separated list of endpoints")
	connectCmd.MarkPersistentFlagRequired("endpoints")

	connectCmd.PersistentFlags().StringVarP(&shell.User, "user", "u", "", "ETCD user")
	connectCmd.PersistentFlags().StringVarP(&shell.Password, "password", "p", "", "ETCD password")
	connectCmd.MarkFlagsRequiredTogether("user", "password")

	connectCmd.PersistentFlags().BoolVarP(&shell.UseTls, "tls", "t", false, "Use TLS for ETCD connection")
}
