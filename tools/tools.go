package tools

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/fatih/color"
)

func Error(format string, a ...interface{}) {
	err_print := color.New(color.FgHiRed, color.Bold)
	if a != nil {
		err_print.Printf(format, a...)
	} else {
		err_print.Println(format)
	}
}

func PrintDir(format string, a ...interface{}) {
	d_print := color.New(color.FgBlue, color.Bold)
	if a != nil {
		d_print.Printf(format, a...)
	} else {
		d_print.Println(format)
	}
}

func PrintKey(format string, a ...interface{}) {
	if a != nil {
		color.HiCyan(format, a...)
	} else {
		color.HiCyan(format)
	}
}

func PrintValue(format string, a ...interface{}) {
	if a != nil {
		color.HiWhite(format, a...)
	} else {
		color.HiWhite(format)
	}
}

func PrintKeyValue(key, value string) {
	k := color.New(color.FgHiCyan)
	v := color.New(color.FgHiWhite)
	k.Printf("%s: ", key)
	v.Println(value)
}

func PrintActionKeyValue(action, key, value string) {
	a := color.New(color.BgBlue, color.FgHiYellow)
	k := color.New(color.FgHiCyan)
	v := color.New(color.FgHiWhite)
	a.Print(action)
	k.Printf(" %s: ", key)
	v.Println(value)
}

func Ask(question string) {
	k := color.New(color.FgHiCyan)
	k.Printf(question)
}

func PrettyPrint(in interface{}) {
	_json, err := json.MarshalIndent(in, "", "  ")
	if err == nil {
		color.HiWhite("%s\n", string(_json))
	}
}

func PrettyPrint2(in map[string]interface{}) {
	color.HiWhite("{")
	pretty_print(in, "  ")
	color.HiWhite("}")
}

func pretty_print(in map[string]interface{}, ident string) {
	keys := make([]string, 0, len(in))
	for k := range in {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	nb := len(keys)
	for i, s := range keys {
		suffix := ","
		if i+1 == nb {
			suffix = ""
		}
		val_kind := reflect.TypeOf(in[s]).Kind()
		if val_kind == reflect.Map {
			color.HiWhite("%s\"%s\": {\n", ident, s)
			pretty_print(CastToMap(in[s]), ident+"  ")
			color.HiWhite("%s}%s\n", ident, suffix)
		} else if val_kind == reflect.Slice {
			tmp := fmt.Sprintf("%q\n", in[s])
			tmp = strings.TrimPrefix(tmp, "[")
			tmp = strings.TrimSuffix(tmp, "]\n")
			if len(tmp) == 0 {
				color.HiWhite("%s\"%s\": []%s\n", ident, s, suffix)
			} else {
				tokens := strings.Split(tmp, " ")
				// value := strings.Join(tokens, ", ")
				// color.HiWhite("%s\"%s\": %v%s\n", ident, s, value, suffix)
				color.HiWhite("%s\"%s\": [\n", ident, s)
				tok_len := len(tokens) - 1
				tmp_ident := ident + "  "
				for j := range tokens {
					if j == tok_len {
						color.HiWhite("%s%s\n", tmp_ident, tokens[j])
					} else {
						color.HiWhite("%s%s,\n", tmp_ident, tokens[j])
					}
				}
				color.HiWhite("%s]%s\n", ident, suffix)
			}
		} else {
			val := fmt.Sprintf("%v", in[s])
			val = strings.TrimSuffix(val, "\n")
			val = strings.ReplaceAll(val, "\n", "\n"+ident+"  ")
			color.HiWhite("%s\"%s\": %v%s\n", ident, s, val, suffix)
		}
	}
}
