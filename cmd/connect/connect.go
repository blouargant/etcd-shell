/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package connect

import (
	"etcd-shell/cmd"
	"etcd-shell/shell"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "interactive shell",
	Run: func(cmd *cobra.Command, args []string) {
		shell.RunShell()
	},
}

func init() {
	//cmd.RootCmd.AddCommand(connectCmd)
	cmd.RootCmd.PersistentFlags().StringVarP(&shell.Endpointlist, "endpoints", "e", "", "Comma separated list of endpoints")
	cmd.RootCmd.MarkPersistentFlagRequired("endpoints")
	cmd.RootCmd.PersistentFlags().StringVarP(&shell.User, "user", "u", "", "ETCD user")
	cmd.RootCmd.PersistentFlags().StringVarP(&shell.Password, "password", "p", "", "ETCD password")
	cmd.RootCmd.MarkFlagsRequiredTogether("user", "password")
	cmd.RootCmd.PersistentFlags().BoolVarP(&shell.UseTls, "tls", "t", false, "Use TLS for ETCD connection")

	cmd.RootCmd.AddCommand(shellCmd)
	cmd.RootCmd.AddCommand(lsCmd)
	cmd.RootCmd.AddCommand(catCmd)
	cmd.RootCmd.AddCommand(watchCmd)
	cmd.RootCmd.AddCommand(setCmd)
	cmd.RootCmd.AddCommand(rmCmd)
}

func keyCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	c, err := shell.NewCompleter()
	if err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}
	candidates := c.KeyCompletion(toComplete)
	return candidates, cobra.ShellCompDirectiveNoFileComp
}

func dirCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var res []string
	lst, opt := keyCompletion(cmd, args, toComplete)
	for _, c := range lst {
		res = append(res, c)
		if strings.HasSuffix(c, "/") {
			res = append(res, c)
		}
	}
	return res, opt
}

var lsCmd = &cobra.Command{
	Use:               "ls",
	Short:             "list directory",
	Long:              "List a directory.",
	ValidArgsFunction: dirCompletion,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := shell.NewCompleter()
		if err != nil {
			fmt.Println("error", err)
			os.Exit(1)
		}
		keys := []string{""}
		val := ""
		if len(args) > 0 {
			val = strings.Join(args[:], "/")
			keys = append(keys, val)
		}
		c.List(keys)
	},
}

var catCmd = &cobra.Command{
	Use:               "cat",
	Short:             "cat keys' value",
	Long:              "Show a directory content.",
	ValidArgsFunction: keyCompletion,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := shell.NewCompleter()
		if err != nil {
			fmt.Println("error", err)
			os.Exit(1)
		}
		keys := []string{""}
		val := ""
		if len(args) > 0 {
			val = strings.Join(args[:], "/")
			keys = append(keys, val)
		}
		c.Show(keys)
	},
}

var watchCmd = &cobra.Command{
	Use:               "Watch",
	Short:             "watch keys",
	Long:              "watch a directory for modifications.",
	ValidArgsFunction: keyCompletion,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := shell.NewCompleter()
		if err != nil {
			fmt.Println("error", err)
			os.Exit(1)
		}
		keys := []string{""}
		val := ""
		if len(args) > 0 {
			val = strings.Join(args[:], "/")
			keys = append(keys, val)
		}
		c.Watch(keys)
	},
}

func setCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var res []string
	if len(args) == 0 {
		return keyCompletion(cmd, args, toComplete)
	} else {
		c, err := shell.NewCompleter()
		if err != nil {
			fmt.Println("error", err)
			os.Exit(1)
		}
		val := c.GetValue(args[0])
		if val != "" {
			res = append(res, val)
		}
	}
	return res, cobra.ShellCompDirectiveNoFileComp
}

var setCmd = &cobra.Command{
	Use:               "set",
	Short:             "set a key's value",
	Long:              "Set key's value.",
	ValidArgsFunction: setCompletion,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := shell.NewCompleter()
		if err != nil {
			fmt.Println("error", err)
			os.Exit(1)
		}
		if len(args) == 2 {
			c.Set(args[0], args[1])
			return
		} else {
			fmt.Println("Please provide a key and a value")
			os.Exit(1)
		}
	},
}

var rmCmd = &cobra.Command{
	Use:               "rm",
	Short:             "delete keys",
	Long:              "Delete a key or a directory.",
	ValidArgsFunction: keyCompletion,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := shell.NewCompleter()
		if err != nil {
			fmt.Println("error", err)
			os.Exit(1)
		}
		keys := []string{""}
		val := ""
		if len(args) > 0 {
			val = strings.Join(args[:], "/")
			keys = append(keys, val)
		} else {
			fmt.Println("Please provide a key to delete")
			os.Exit(1)
		}
		c.Delete(keys)
	},
}
