package convertago

import (
	"fmt"
	"reflect"
	"strings"
)

func unwrapValue(v any) reflect.Value {
	value := reflect.ValueOf(v)
	if !value.IsValid() {
		return value
	}

	return dereference(value)
}

func dereference(v reflect.Value) reflect.Value {
	for v.IsValid() &&
		(v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface) {
		if v.IsNil() {
			break
		}
		v = v.Elem()
	}
	return v
}

func sprintf(format string, v reflect.Value) string {
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			placeholder := fmt.Sprintf("{%d}", i)
			format = strings.Replace(format, placeholder, fmt.Sprintf("%v", v.Index(i)), -1)
		}
		return format
	case reflect.Map:
		iter := v.MapRange()
		for iter.Next() {
			placeholder := fmt.Sprintf("{%v}", iter.Key())
			format = strings.Replace(format, placeholder, fmt.Sprintf("%v", iter.Value()), -1)
		}
		return format
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fieldT, fieldV := v.Type().Field(i), dereference(v.Field(i))
			placeholder := fmt.Sprintf("{%v}", fieldT.Name)
			format = strings.Replace(format, placeholder, fmt.Sprintf("%v", fieldV), -1)
		}
		return format
	default:
		return fmt.Sprintf(format, v)
	}
}
