package shell

import (
	"fmt"
	"log"
	"strings"

	etcd_client "etcd-shell/etcd"

	"github.com/c-bata/go-prompt"
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

type Completer struct {
	etcd   *etcd_client.Etcd
	Prompt *prompt.Prompt
}

func NewCompleter() (*Completer, error) {
	etcd := &etcd_client.Etcd{}
	Pwd = ""
	if Endpointlist != "" {
		fmt.Printf("connecting to endpoints %v\n", Endpointlist)
		etcd.Connect(Endpointlist, User, Password, UseTls)
		res, err := etcd.GetObject(Pwd)
		if err == nil {
			setRoot(res)
		} else {
			log.Fatal(err)
		}
	}
	return &Completer{
		etcd: etcd,
	}, nil
}

func (c *Completer) Complete(d prompt.Document) []prompt.Suggest {
	if d.TextBeforeCursor() == "" {
		if d.LastKeyStroke().String() == "Tab" {
			return commands
		}
		//return commands
		return []prompt.Suggest{}
	}
	return c.argumentsCompleter(d)

}

func setRoot(res map[string][]byte) {
	for d := range res {
		if strings.HasPrefix(d, SEP) {
			RootPath = SEP
		} else {
			dirs := strings.Split(d, SEP)
			if len(dirs) > 0 {
				RootPath = dirs[0]
			}
		}
	}
}
