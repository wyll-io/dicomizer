package functions

import (
	"reflect"
	"strings"
)

func CreateMap(args ...interface{}) map[string]interface{} {
	out := make(map[string]interface{})

	for i := 0; i < len(args); i++ {
		k := args[i].(string)
		i++
		v := args[i]

		if strings.Contains(k, ".") {
			// TODO: support more than one level of nesting
			keys := strings.SplitN(k, ".", 1)
			k := keys[0]
			subK := keys[1]

			if subV, ok := out[k].(map[string]interface{})[subK]; ok {
				switch reflect.TypeOf(subV).Kind() {
				case reflect.Slice:
					out[k].(map[string]interface{})[subK] = append(subV.([]interface{}), v)
				case reflect.Map:
					// * overwrite existing map using new value
					out[k].(map[string]interface{})[subK] = v
				default:
					// * convert already existing value to slice and append new value
					out[k].(map[string]interface{})[subK] = append([]interface{}{subV}, v)
				}

				continue
			}

			out[k].(map[string]interface{})[subK] = v
			continue
		}

		if subV, ok := out[k]; ok {
			switch reflect.TypeOf(subV).Kind() {
			case reflect.Slice:
				out[k] = append(subV.([]interface{}), v)
			case reflect.Map:
				// * overwrite existing map using new value
				out[k] = v
			default:
				// * convert already existing value to slice and append new value
				out[k] = append([]interface{}{subV}, v)
			}

			continue
		}

		out[k] = v
	}

	return out
}

func Isset(v interface{}) bool {
	return v == nil || v == "" || v == 0
}
