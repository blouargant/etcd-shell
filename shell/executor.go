package shell

import (
	"context"
	"encoding/json"
	"etcd-shell/tools"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
)

const SEP string = "/"

func (c *Completer) Executor(s string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return
	} else if s == "quit" || s == "exit" {
		fmt.Println("Bye!")
		os.Exit(0)
		return
	}

	args := strings.Split(s, " ")
	if Endpointlist == "" {
		if args[0] == "connect" {
			c.connect(args)
		}
	} else {
		switch args[0] {
		case "help":
			help()
		case "cp":
			c.cp(args)
		case "rm":
			c.rm(args)
		case "set":
			c.set(args)
		case "ls":
			c.ls(args)
		case "cat":
			c.cat(args)
		case "watch":
			c.watch(args)
		case "cd":
			c.cd(args)
		case "dump":
			c.dump(args)
		case "disconnect":
			c.disconnect()
		case "pwd":
			fmt.Println(Pwd)
		}
	}
}

func help() {
	fmt.Println("Available commands:")
	for _, cmd := range commands {
		tab := "\t"
		if len(cmd.Text) < 5 {
			tab += "\t"
		}
		fmt.Printf("  %s:%s%s\n", cmd.Text, tab, cmd.Description)
	}
}

func (c *Completer) connect(args []string) {
	var i int = 0
	max := len(args)
	for i < max {
		if args[i] == "-e" || args[i] == "--endpoints" {
			if i+1 < max {
				Endpointlist = args[i+1]
			} else {
				log.Fatal("missing enpoints settings")
			}
		} else if args[i] == "-u" || args[i] == "--user" {
			if i+1 < max {
				User = args[i+1]
			} else {
				log.Fatal("missing user setting")
			}
		} else if args[i] == "-p" || args[i] == "--Password" {
			if i+1 < max {
				Password = args[i+1]
			} else {
				log.Fatal("missing password setting")
			}
		} else if args[i] == "-t" || args[i] == "--tls" {
			UseTls = true
		}
		i += 1
	}
	connectEtcd(c.etcd)
}

func (c *Completer) disconnect() {
	c.etcd = nil
	fmt.Printf("Closed connection to %s\n", Endpointlist)
	Endpointlist = ""
}

func (c *Completer) cp(args []string) {
	var recursiv bool
	var final_args []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			if arg == "-r" {
				recursiv = true
			}
		} else {
			final_args = append(final_args, arg)
		}
	}
	var in_path, out_path string
	if len(final_args) > 2 {
		in_path = getInputPath(final_args[1])
		out_path = getInputPath(final_args[2])
	} else {
		return
	}
	if !recursiv {
		res, err := c.etcd.GetValue(in_path)
		if err == nil {
			err = c.etcd.Put(out_path, string(res))
			if err != nil {
				fmt.Println(err)
			} else {
				tools.PrintKeyValue(out_path, string(res))
			}
		} else {
			fmt.Println(err)
		}
	} else {
		res, err := c.etcd.GetObject(in_path)
		if err != nil {
			return
		}
		for key := range res {
			shorted := strings.TrimPrefix(key, in_path)
			cp_path := filepath.Join(out_path, shorted)
			err = c.etcd.Put(cp_path, string(res[key]))
			if err != nil {
				fmt.Println(err)
				break
			} else {
				tools.PrintKeyValue(cp_path, string(res[key]))
			}
		}
	}
}

func (c *Completer) Delete(args []string) {
	c.rm(args)
}

func (c *Completer) rm(args []string) {
	var confirmed, forced, recursiv bool
	var final_args []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			if arg == "-f" {
				confirmed = true
				forced = true
			} else if arg == "-r" {
				recursiv = true
			}
		} else {
			final_args = append(final_args, arg)
		}
	}
	var path string
	if len(final_args) > 1 {
		path = getInputPath(final_args[1])
	} else {
		return
	}
	res, err := c.etcd.GetObject(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(res) > 1 && !recursiv {
		fmt.Println("Please use \"-r\" switch to delete a directory")
		return
	}

	if !forced {
		confirmed = askForConfirmation(fmt.Sprintf("Deleting %v ?", path), "n")
	}
	if confirmed {
		err := c.etcd.Delete(path)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("%s deleted\n", path)
		}
	}
}

func (c *Completer) Set(args []string) {
	c.set(args)
}

func (c *Completer) set(args []string) {
	var path string
	if len(args) > 2 {
		new_val := strings.Join(args[2:], " ")
		new_val = strings.TrimSpace(new_val)
		path = getInputPath(args[1])

		err := c.etcd.Put(path, new_val)
		if err != nil {
			fmt.Println(err)
			return
		}
		tools.PrintKeyValue(path, new_val)
	} else {
		return
	}
}

func (c *Completer) cd(args []string) {
	var path string
	if len(args) > 1 {
		path = getInputPath(args[1])
	} else {
		return
	}
	_, err := c.etcd.GetObject(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(path)
	Pwd = path
}

func (c *Completer) List(args []string) {
	c.ls(args)
}

func (c *Completer) ls(args []string) {
	var path string
	if len(args) > 1 {
		path = getInputPath(args[1])
	} else {
		path = getInputPath("")
	}
	res, err := c.etcd.GetObject(path)
	if err != nil {
		//tools.Error("Error: %v\n ", err)
		return
	}
	if len(res) == 1 {
		//tools.Error("Error: %v is a key\n ", path)
		return
	}
	var sorted, all, kinds []string
	for d := range res {
		d_args := strings.Split(d, SEP)
		if len(d_args) > 0 {
			if !tools.Contains(all, d_args[0]) {
				all = append(all, d_args[0])
				sorted = append(sorted, d_args[0])
				if len(d_args) > 1 {
					kinds = append(kinds, "dir")
				} else {
					kinds = append(kinds, "file")
				}
			}
		}
	}
	sort.Strings(sorted)
	for _, s := range sorted {
		idx := tools.IndexOf(all, s)
		switch kinds[idx] {
		case "dir":
			tools.PrintDir(s)
		case "file":
			tools.PrintKey(s)
		}
	}
}

func (c *Completer) KeyCompletion(toComplete string) []string {
	result := []string{}
	//res, err := c.etcd.GetObject(toComplete)
	res, err := c.etcd.ListKeys(toComplete)
	if err != nil {
		return result
	}
	keys := strings.Split(toComplete, "/")
	idx := len(keys) - 1
	if idx < 0 {
		idx = 0
	}
	var sorted, all, kinds []string
	for _, d := range res {
		d_args := strings.Split(d, SEP)
		if len(d_args) > 0 {
			var prefix string
			if len(d_args) < idx+1 {
				prefix = d
			} else {
				prefix = strings.Join(d_args[:idx+1], SEP)
			}
			if !tools.Contains(all, prefix) {
				all = append(all, prefix)
				sorted = append(sorted, prefix)
				if len(d_args)-idx > 1 {
					kinds = append(kinds, "dir")
				} else {
					kinds = append(kinds, "file")
				}
			}
		}
	}
	sort.Strings(sorted)
	for _, s := range sorted {
		idx := tools.IndexOf(all, s)
		switch kinds[idx] {
		case "dir":
			result = append(result, s+"/")
		case "file":
			result = append(result, s)
		}
	}
	return result
}

func (c *Completer) Show(args []string) {
	c.cat(args)
}

func (c *Completer) cat(args []string) {
	var path string
	if len(args) > 1 {
		path = getInputPath(args[1])
	} else {
		path = getInputPath("")
	}
	res, err := c.etcd.GetValue(path)
	if err == nil {
		if len(res) == 1 {
			fmt.Printf("%v is not a key\n ", path)
			return
		}
		tools.PrintValue(string(res))
	} else {
		res, err := c.etcd.GetObject(path)
		if err != nil {
			//tools.Error("Error: %v\n ", err)
			return
		}
		var sorted []string
		for d := range res {
			sorted = append(sorted, d)
		}
		sort.Strings(sorted)
		to_trim := filepath.Dir(path)
		for _, key := range sorted {
			key = strings.TrimPrefix(key, to_trim)
			tmp_val := string(res[key])
			tools.PrintKeyValue(key, tmp_val)
		}
	}
}

func (c *Completer) Watch(args []string) {
	c.watch(args)
}

func (c *Completer) watch(args []string) {
	var path string
	if len(args) > 1 {
		path = getInputPath(args[1])
	} else {
		path = getInputPath("")
	}
	handleOutput()
	fmt.Printf("watching %s ...\n", path)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go watch(ctx, path, c)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	<-ctx.Done()
	stop()
	fmt.Println("watcher stopped")
}

func watch(ctx context.Context, path string, c *Completer) {
	rch := c.etcd.Watch(ctx, path)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			//fmt.Printf("%s %s : %s\n", ev.Type.String(), string(ev.Kv.Key), string(ev.Kv.Value))
			tools.PrintActionKeyValue(ev.Type.String(), string(ev.Kv.Key), string(ev.Kv.Value))
		}
	}
}

func (c *Completer) dump(args []string) {
	var path string
	var value_json bool
	var final_args []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			if arg == "-j" {
				value_json = true
			}
		} else {
			final_args = append(final_args, arg)
		}
	}

	if len(final_args) > 1 {
		path = getInputPath(final_args[1])
	} else {
		path = getInputPath("")
	}
	res, err := c.etcd.GetObject(path)
	if err != nil {
		//tools.Error("Error: %v\n ", err)
		return
	}
	var sorted []string
	for d := range res {
		sorted = append(sorted, d)
	}
	sort.Strings(sorted)
	to_trim := filepath.Dir(path)
	var display = make(map[string]interface{})
	for _, key := range sorted {
		key = strings.TrimPrefix(key, to_trim)
		key = strings.TrimPrefix(key, SEP)
		sub := strings.Split(key, SEP)

		tmp_val := string(res[key])
		if strings.HasPrefix(tmp_val, "{") {
			if value_json {
				value := make(map[string]interface{})
				json.Unmarshal(res[key], &value)
				display = tools.MakeNestedMap(sub, value, display)
			} else {
				value := strings.Replace(tmp_val, "\\", "", -1)
				display = tools.MakeNestedMap(sub, value, display)
			}
			//else {
			//	value = strings.ReplaceAll(tmp_val, "\", "")
			//}
		} else {
			value := tmp_val
			display = tools.MakeNestedMap(sub, value, display)
		}
	}
	//fmt.Printf("%#v \n", display)
	tools.PrettyPrint2(display)
}

func askForConfirmation(ask, defaultVal string) bool {
	var response, question string

	switch strings.ToLower(defaultVal) {
	case "y":
		question = ask + " [Y/n]:"
	case "yes":
		question = ask + " [Yes/no]:"
	case "n":
		question = ask + " [y/N]:"
	case "no":
		question = ask + " [yes/No]:"
	default:
		question = ask + " :"
	}

	handleOutput()
	tools.Ask(question)
	_, err := fmt.Scanln(&response)
	if err != nil && err.Error() != "unexpected newline" {
		log.Fatal(err)
	}

	switch strings.ToLower(response) {
	case "y", "yes":
		return true
	case "n", "no":
		return false
	case "":
		switch strings.ToLower(defaultVal) {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		default:
			fmt.Println("I'm sorry but I didn't get what you meant, please type (y)es or (n)o and then press enter:")
			return askForConfirmation(ask, defaultVal)
		}
	default:
		fmt.Println("I'm sorry but I didn't get what you meant, please type (y)es or (n)o and then press enter:")
		return askForConfirmation(ask, defaultVal)
	}
}

func handleOutput() {
	rawModeOff := exec.Command("/bin/stty", "-raw", "echo")
	rawModeOff.Stdin = os.Stdin
	_ = rawModeOff.Run()
	rawModeOff.Wait()
}
