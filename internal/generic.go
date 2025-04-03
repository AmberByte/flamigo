package internal

import (
	"reflect"
)

func ParseOptionalParam[T any](param []T, defaultValue T) T {
	r := reflect.ValueOf(param)
	if r.Kind() != reflect.Slice {
		panic("ParseOptionalParam requires slice as param")
	}
	if r.Len() == 0 {
		return defaultValue
	}

	d := r.Index(0)
	if d.Kind() == reflect.Interface && d.IsNil() {
		return defaultValue
	}
	return d.Interface().(T)

}
