package shared

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Pretty prints any value as indented JSON, converting []uint8 to string for readability.
func Pretty(label string, v interface{}) {
	processed := convertBytes(v)
	b, _ := json.MarshalIndent(processed, "", "  ")
	fmt.Printf("%s:\n%s\n", label, string(b))
}

func convertBytes(in interface{}) interface{} {
	switch val := in.(type) {
	case map[string]interface{}:
		m := map[string]interface{}{}
		for k, v := range val {
			m[k] = convertBytes(v)
		}
		return m
	case []interface{}:
		out := make([]interface{}, len(val))
		for i, v := range val {
			out[i] = convertBytes(v)
		}
		return out
	case []uint8:
		return string(val)
	default:
		rv := reflect.ValueOf(in)
		if rv.Kind() == reflect.Slice {
			// generic slice handling
			out := make([]interface{}, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				out[i] = convertBytes(rv.Index(i).Interface())
			}
			return out
		}
		return in
	}
}
