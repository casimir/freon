package serialize

import (
	"fmt"
	"reflect"
	"strings"
)

const TagKey = "desc"

type Field struct {
	Name      string `json:"name"`
	Value     any    `json:"value"`
	Readonly  bool   `json:"readonly"`
	Obfuscate bool   `json:"obfuscate"`
}

func Describe(obj any) ([]Field, error) {
	ret := []Field{}
	v := reflect.Indirect(reflect.ValueOf(obj))
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() == reflect.Struct && v.Type().Field(i).Anonymous {
			inner, err := Describe(v.Field(i).Interface())
			if err != nil {
				return nil, err
			}
			ret = append(ret, inner...)
			continue
		}

		tag := v.Type().Field(i).Tag.Get(TagKey)
		params, err := parseTag(tag)
		if err != nil {
			return nil, err
		}

		if params.hidden {
			continue
		}

		ret = append(ret, Field{
			Name:      v.Type().Field(i).Name,
			Value:     v.Field(i).Interface(),
			Readonly:  params.readonly,
			Obfuscate: params.obfuscate,
		})
	}
	return ret, nil
}

type tagParams struct {
	hidden    bool
	readonly  bool
	obfuscate bool
}

func parseTag(tag string) (*tagParams, error) {
	if tag == "" {
		return &tagParams{}, nil
	}

	params := &tagParams{}
	for _, param := range strings.Split(tag, ",") {
		switch param {
		case "hidden":
			params.hidden = true
		case "readonly":
			params.readonly = true
		case "obfuscate":
			params.obfuscate = true
		default:
			return nil, fmt.Errorf("%s: unknown parameter: %s", TagKey, param)
		}
	}
	return params, nil
}
