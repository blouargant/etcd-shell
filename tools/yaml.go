package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"gopkg.in/yaml.v3"
)

func ParseYaml(source string, output interface{}) (err error) {
	err = ValidatePath(source)
	if err != nil {
		return err
	}
	file, err := os.Open(source)
	if err != nil {
		return err
	}
	defer file.Close()

	dec := yaml.NewDecoder(file)
	data_list := reflect.Indirect(reflect.ValueOf(output))
	if data_list.Kind() != reflect.Slice {
		return fmt.Errorf("output must be a list")
	}
	target := data_list.Type()
	t_elem := target.Elem()
	new_slice := reflect.MakeSlice(target, 0, 0)
	empty := reflect.New(t_elem).Interface()
	for {
		data := reflect.New(t_elem).Interface()
		if dec.Decode(data) != nil {
			break
		}
		if !reflect.DeepEqual(data, empty) {
			new_slice = reflect.Append(new_slice, reflect.Indirect(reflect.ValueOf(data)))
		}
	}
	data_list.Set(new_slice)
	return
}

func ValidatePath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	fileExtension := filepath.Ext(path)
	if fileExtension == ".yaml" || fileExtension == ".yml" {
		return nil
	}
	return fmt.Errorf("file must have 'yaml' or 'yml' extention")
}
