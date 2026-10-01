package validation

import (
	"net/url"
	"reflect"
	"strings"
)

// Value unwraps pointers and interfaces without calling model methods.
func Value(input any) any {
	value := reflect.ValueOf(input)
	for value.IsValid() && (value.Kind() == reflect.Pointer || value.Kind() == reflect.Interface) {
		if value.IsNil() {
			return nil
		}
		value = value.Elem()
	}
	if !value.IsValid() {
		return nil
	}
	return value.Interface()
}

// AbsoluteURI accepts an absolute URI with a destination and optional scheme restrictions.
func AbsoluteURI(text string, schemes ...string) bool {
	u, err := url.Parse(text)
	if err != nil || u.Scheme == "" || (u.Host == "" && u.Opaque == "" && u.Path == "") {
		return false
	}
	scheme := strings.ToLower(u.Scheme)
	if len(schemes) != 0 {
		found := false
		for _, allowed := range schemes {
			found = found || scheme == allowed
		}
		if !found {
			return false
		}
	}
	if scheme == "http" || scheme == "https" {
		_, err := url.ParseRequestURI(text)
		return err == nil && u.Host != ""
	}
	return true
}
