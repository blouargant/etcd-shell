package tools

import (
	"reflect"
)

func CastToMap(data interface{}) map[string]interface{} {
	new_map := make(map[string]interface{})
	iter := reflect.ValueOf(data).MapRange()
	for iter.Next() {
		key := iter.Key().Interface()
		value := iter.Value().Interface()
		if str, ok := key.(string); ok {
			new_map[string(str)] = value
		}
	}
	return new_map
}

func MakeNestedMap(keys []string, value any, dic map[string]interface{}) map[string]interface{} {
	if len(keys) == 1 {
		dic[keys[0]] = value
	} else {
		new_dic := make(map[string]interface{})
		_, exist := dic[keys[0]]
		if exist {
			new_dic = CastToMap(dic[keys[0]])
		}
		var new_keys []string
		for i := 1; i < len(keys); i++ {
			new_keys = append(new_keys, keys[i])
		}
		tmp := MakeNestedMap(new_keys, value, new_dic)
		dic[keys[0]] = tmp
	}
	return dic
}

type ListType interface {
	string | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

func Contains[L ListType](lst []L, el L) bool {
	for _, v := range lst {
		if v == el {
			return true
		}
	}
	return false
}

func IndexOf[L ListType](lst []L, el L) int {
	for i, v := range lst {
		if v == el {
			return i
		}
	}
	return -1
}

func RemoveIndex[L any](s []L, index int) []L {
	return append(s[:index], s[index+1:]...)
}

func Equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

type Number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

func SumArray[N Number](s []N) N {
	var res N
	for _, v := range s {
		res += v
	}
	return res
}
