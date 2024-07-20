package shell

import (
	"fmt"
	"log"
	"strings"

	etcd_client "etcd-shell/etcd"

	"github.com/c-bata/go-prompt"
)

type Completer struct {
	etcd *etcd_client.Etcd
}

func connectEtcd(etcd *etcd_client.Etcd) {
	Pwd = ""
	etcd.Connect(Endpointlist, User, Password, UseTls)
	res, err := etcd.GetObject(Pwd)
	if err == nil {
		setRoot(res)
	} else {
		log.Fatal(err)
	}
}

func NewCompleter() (*Completer, error) {
	etcd := &etcd_client.Etcd{}
	if Endpointlist != "" {
		connectEtcd(etcd)
	}
	return &Completer{
		etcd: etcd,
	}, nil
}

func (c *Completer) Complete(d prompt.Document) []prompt.Suggest {
	if d.TextBeforeCursor() == "" {
		if d.LastKeyStroke().String() == "Tab" {
			if Endpointlist == "" {
				return connCmd
			}
			//return commands
			fmt.Println("help")
			help()
			return []prompt.Suggest{}
		}
		return []prompt.Suggest{}
	} else if d.TextBeforeCursor() == " " {
		return commands
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
