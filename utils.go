package client

import (
	"reflect"
	"strconv"
)

func struct2map(i interface{}) RequestParams {
	values := RequestParams{}
	iVal := reflect.ValueOf(i).Elem()
	typ := iVal.Type()
	for i := 0; i < iVal.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := iVal.Field(i)
		fieldName, hasTag := field.Tag.Lookup("json")
		if !hasTag {
			fieldName = field.Name
		}
		// You ca use tags here...
		// tag := typ.Field(i).Tag.Get("tagname")
		// Convert each type into a string for the url.Values string map
		var v string
		switch fieldValue.Interface().(type) {
		case int, int8, int16, int32, int64:
			v = strconv.FormatInt(fieldValue.Int(), 10)
		case uint, uint8, uint16, uint32, uint64:
			v = strconv.FormatUint(fieldValue.Uint(), 10)
		case float32:
			v = strconv.FormatFloat(fieldValue.Float(), 'f', 4, 32)
		case float64:
			v = strconv.FormatFloat(fieldValue.Float(), 'f', 4, 64)
		case []byte:
			v = string(fieldValue.Bytes())
		case string:
			v = fieldValue.String()
		case map[string]string:
			for _, prm := range fieldValue.MapKeys() {
				values[sprintf("%s[%s]", fieldName, prm)] = fieldValue.MapIndex(prm).String()
			}
			continue
		}
		values[fieldName] = v
	}
	return values
}
