package shell

import (
	"etcd-shell/tools"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/c-bata/go-prompt"
)

var commands = []prompt.Suggest{
	{Text: "ls"},
	{Text: "cd"},
	{Text: "pwd"},
	{Text: "dump"},
	{Text: "cp"},
	{Text: "set"},
	{Text: "rm"},
	{Text: "watch"},
	{Text: "exit"},
}
var dumpOpt = []prompt.Suggest{
	{Text: "-j", Description: "Interpret values as JSON objects"},
}

var rmOpt = []prompt.Suggest{
	{Text: "-r", Description: "Delete directories"},
	{Text: "-f", Description: "Do not ask confirmation"},
}

var cpOpt = []prompt.Suggest{
	{Text: "-r", Description: "Copy directories"},
}

func (c *Completer) argumentsCompleter(d prompt.Document) []prompt.Suggest {
	args := strings.Split(d.TextBeforeCursor(), " ")
	if len(args) <= 1 {
		return prompt.FilterHasPrefix(commands, args[0], true)
	}
	sug := []prompt.Suggest{}
	switch args[0] {
	case "set":
		if len(args) < 3 {
			sorted, all, kinds := c.getPaths(args[len(args)-1])
			for _, s := range sorted {
				idx := tools.IndexOf(all, s)
				if kinds[idx] == "dir" {
					sug = append(sug, prompt.Suggest{Text: s, Description: "Subpath"})
				} else {
					sug = append(sug, prompt.Suggest{Text: s, Description: "Key"})
				}
			}
		} else {
			if len(args) == 3 && len(args[2]) == 0 {
				_, all, kinds := c.getPaths(args[1])
				if len(all) > 1 {
					return sug
				} else if len(all) == 1 && kinds[0] == "file" {
					value, err := c.getValue(args[1])
					if err == nil {
						sug = append(sug, prompt.Suggest{Text: fmt.Sprintf("%v", value), Description: "Previous value"})
					}
				}
			}
		}
	case "ls":
		if len(args) < 3 {
			sorted, all, kinds := c.getPaths(args[len(args)-1])
			for _, s := range sorted {
				idx := tools.IndexOf(all, s)
				if kinds[idx] == "dir" {
					sug = append(sug, prompt.Suggest{Text: s})
				}
			}
		}
	case "dump":
		if len(args) < 4 {
			if len(args) == 2 && args[1] == "" {
				sug = append(sug, dumpOpt...)
			}
			if len(args) > 1 {
				for i, arg := range args {
					if strings.HasPrefix(arg, "-") {
						if arg == "-" {
							sug = append(sug, dumpOpt...)
							return sug
						}
						tools.RemoveIndex(args, i)
					}
				}
				sorted, all, kinds := c.getPaths(args[len(args)-1])
				for _, s := range sorted {
					idx := tools.IndexOf(all, s)
					if kinds[idx] == "dir" {
						sug = append(sug, prompt.Suggest{Text: s})
					}
				}
			}
		}
	case "cp":
		if len(args) < 5 {
			if len(args) == 4 {
				var stop bool = true
				for _, arg := range args {
					if strings.HasPrefix(arg, "-") {
						stop = false
					}
				}
				if stop {
					return cpOpt
				}
			}
			if len(args) == 2 && args[1] == "" {
				sug = append(sug, cpOpt...)
			}
			if len(args) > 1 {
				for i, arg := range args {
					if strings.HasPrefix(arg, "-") {
						if arg == "-" {
							sug = append(sug, cpOpt...)
							return sug
						}
						tools.RemoveIndex(args, i)
					}
				}
				sorted, all, kinds := c.getPaths(args[len(args)-1])
				for _, s := range sorted {
					idx := tools.IndexOf(all, s)
					if kinds[idx] == "dir" {
						sug = append(sug, prompt.Suggest{Text: s, Description: "Subpath"})
					} else {
						sug = append(sug, prompt.Suggest{Text: s, Description: "Key"})
					}
				}
			}
		}
	case "rm":
		if len(args) < 5 {
			var final_args []string
			if len(args) < 3 && args[len(args)-1] == "" {
				sug = append(sug, rmOpt...)
				sug = removeDuplicates(sug, args)
			}
			if len(args) > 1 {
				for _, arg := range args {
					if strings.HasPrefix(arg, "-") {
						if arg == "-" {
							sug = append(sug, rmOpt...)
							sug = removeDuplicates(sug, args)
							return sug
						}
					} else {
						final_args = append(final_args, arg)
					}
				}
				sorted, all, kinds := c.getPaths(final_args[len(final_args)-1])
				for _, s := range sorted {
					idx := tools.IndexOf(all, s)
					if kinds[idx] == "dir" {
						sug = append(sug, prompt.Suggest{Text: s, Description: "Subpath"})
					} else {
						sug = append(sug, prompt.Suggest{Text: s, Description: "Key"})
					}
				}
			}
		}
	case "cat", "watch":
		if len(args) < 3 {
			sorted, all, kinds := c.getPaths(args[len(args)-1])
			for _, s := range sorted {
				idx := tools.IndexOf(all, s)
				if kinds[idx] == "dir" {
					sug = append(sug, prompt.Suggest{Text: s, Description: "Subpath"})
				} else {
					sug = append(sug, prompt.Suggest{Text: s, Description: "Key"})
				}
			}
		}
	case "cd":
		if len(args) < 3 {
			sorted, all, kinds := c.getPaths(args[len(args)-1])
			for _, s := range sorted {
				idx := tools.IndexOf(all, s)
				if kinds[idx] == "dir" {
					sug = append(sug, prompt.Suggest{Text: s})
				}
			}
		}
	}
	return sug
}

func (c *Completer) getValue(input string) (value string, err error) {
	path := getInputPath(input)
	res, err := c.etcd.GetValue(path)
	if err != nil {
		return
	}
	if len(res) == 1 {
		return
	}
	return string(res), nil
}

func (c *Completer) getPaths(input string) (sorted, all, kinds []string) {
	var path string
	if len(input) > 0 {
		path = getInputPath(input)
	} else {
		path = getInputPath("")
	}
	res, err := c.etcd.GetObject(path)
	if err != nil {
		return
	}
	to_trim := filepath.Dir(path)
	for d := range res {
		if !strings.HasPrefix(d, path) && !strings.HasSuffix(path, SEP) {
			return
		} else {
			d = strings.TrimPrefix(d, to_trim)
			d = strings.TrimPrefix(d, SEP)
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
	}
	sort.Strings(sorted)
	return
}

func getInputPath(input string) (path string) {
	if len(input) == 0 {
		if len(Pwd) > 0 {
			path = Pwd + "/"
		} else {
			path = ""
		}
		return
	}
	if strings.HasPrefix(input, SEP) {
		if RootPath == SEP {
			path = input
		} else {
			path = strings.TrimPrefix(input, SEP)
		}
	} else if strings.HasPrefix(input, RootPath+SEP) {
		path = input
	} else {
		path = filepath.Join(Pwd, input)
		if strings.HasPrefix(path, "../") {
			path = ""
		} else if path == "." {
			path = ""
		}
	}
	return
}

func removeDuplicates(sug []prompt.Suggest, args []string) []prompt.Suggest {
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			for i, s := range sug {
				fmt.Println(s.Text)
				if s.Text == arg {
					sug = tools.RemoveIndex(sug, i)
					break
				}
			}
		}
	}
	return sug
}
