package cmd

import (
	"etcd-shell/shell"
	"fmt"
	"os"
	"os/exec"

	"github.com/c-bata/go-prompt"
	"github.com/c-bata/go-prompt/completer"
)

func handleExit() {
	rawModeOff := exec.Command("/bin/stty", "-raw", "echo")
	rawModeOff.Stdin = os.Stdin
	_ = rawModeOff.Run()
	rawModeOff.Wait()
}

func runShell() {
	defer handleExit()

	c, err := shell.NewCompleter()
	if err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}

	fmt.Printf("etcd-shell\n")
	fmt.Println("Please use `exit` or `Ctrl-D` to exit this program.")
	defer fmt.Println("Bye!")
	p := prompt.New(
		c.Executor,
		c.Complete,
		prompt.OptionTitle("etcd-shell: interactive etcd client"),
		//prompt.OptionPrefix(">>> "),
		prompt.OptionLivePrefix(livePrefix),
		prompt.OptionPrefixTextColor(prompt.DarkBlue),
		prompt.OptionInputTextColor(prompt.Yellow),
		prompt.OptionCompletionWordSeparator(completer.FilePathCompletionSeparator),
	)
	p.Run()
}

func livePrefix() (prefix string, useLivePrefix bool) {
	if len(shell.Pwd) > 0 {
		prefix = fmt.Sprintf("%v ❖ ", shell.Pwd)
	} else {
		prefix = "~ ❖ "
	}
	useLivePrefix = true
	return
}