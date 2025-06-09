package convertago

import (
	"reflect"
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
