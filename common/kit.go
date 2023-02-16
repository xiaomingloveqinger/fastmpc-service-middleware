package common

import (
	"reflect"
	"strconv"
	"strings"
)

func SplitAndTrim(input string) []string {
	result := strings.Split(input, ",")
	for i, r := range result {
		result[i] = strings.TrimSpace(r)
	}
	return result
}

//Map2Struct convert map into struct
//Field name must match
func Map2Struct(src map[string]interface{}, destStrct interface{}) {
	value := reflect.ValueOf(destStrct)
	e := value.Elem()
	for k, v := range src {
		f := e.FieldByName(strings.ToUpper(k[:1]) + k[1:])
		if !f.IsValid() {
			continue
		}
		if !f.CanSet() {
			continue
		}
		mv := reflect.ValueOf(v)
		// map value type
		mvt := mv.Kind().String()
		// struct field type
		sft := f.Kind().String()
		if sft != mvt {
			if mvt == "string" && (strings.Index(sft, "int") != -1) {
				if sft == "int64" {
					i, err := strconv.ParseInt(v.(string), 10, 64)
					if err == nil {
						f.Set(reflect.ValueOf(i))
					}
				} else if sft == "int32" {
					i, err := strconv.ParseInt(v.(string), 10, 32)
					r := int32(i)
					if err == nil {
						f.Set(reflect.ValueOf(r))
					}
				} else if sft == "int" {
					i, err := strconv.Atoi(v.(string))
					if err == nil {
						f.Set(reflect.ValueOf(i))
					}
				} else if sft == "uint64" {
					i, err := strconv.ParseUint(v.(string), 10, 64)
					if err == nil {
						f.Set(reflect.ValueOf(i))
					}
				} else if sft == "uint32" {
					i, err := strconv.ParseUint(v.(string), 10, 32)
					r := uint32(i)
					if err == nil {
						f.Set(reflect.ValueOf(r))
					}
				} else if sft == "uint" {
					i, err := strconv.ParseUint(v.(string), 10, 0)
					r := uint(i)
					if err == nil {
						f.Set(reflect.ValueOf(r))
					}
				}
			}

			// make string and string[] more friendly
			if mvt == "string" && sft == "slice" {
				_, ok := f.Interface().([]string)
				if ok {
					f.Set(reflect.ValueOf(strings.Split(v.(string), ",")))
				}
			}

			// make string and float more friendly
			if mvt == "string" && (strings.Index(sft, "float") != -1) {
				i, err := strconv.ParseFloat(v.(string), 64)
				if err == nil {
					f.Set(reflect.ValueOf(i))
				}
			}

			// make int to bool more friendly
			if mvt == "string" && sft == "bool" {
				i, err := strconv.Atoi(v.(string))
				if err == nil {
					if i == 1 {
						f.Set(reflect.ValueOf(true))
					} else if i == 0 {
						f.Set(reflect.ValueOf(false))
					}
				}
			}
			continue
		}
		f.Set(mv)
	}
}
