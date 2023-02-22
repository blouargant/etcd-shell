package shell

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/c-bata/go-prompt"
	"github.com/c-bata/go-prompt/completer"
)

var (
	Endpointlist string
	User         string
	Password     string
	UseTls       bool
	Pwd          string
	AllPaths     []string
	RootPath     string
)

func handleExit() {
	rawModeOff := exec.Command("/bin/stty", "-raw", "echo")
	rawModeOff.Stdin = os.Stdin
	_ = rawModeOff.Run()
	rawModeOff.Wait()
}

func RunShell() {
	defer handleExit()

	c, err := NewCompleter()
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
		prompt.OptionPrefixTextColor(prompt.Brown),
		prompt.OptionInputTextColor(prompt.Yellow),
		prompt.OptionCompletionWordSeparator(completer.FilePathCompletionSeparator),
	)
	p.Run()
}

func livePrefix() (prefix string, useLivePrefix bool) {
	dot := "❖"
	if Endpointlist == "" {
		dot = "?"
	}
	if len(Pwd) > 0 {
		prefix = fmt.Sprintf("%v %s ", Pwd, dot)
	} else {
		prefix = fmt.Sprintf("~ %s ", dot)
	}
	useLivePrefix = true
	return
}
