/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"etcd-shell/cmd"
	//_ "etcd-shell/cmd/completion"
	_ "etcd-shell/cmd/connect"
)

func main() {
	cmd.Execute()
}
