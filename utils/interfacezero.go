package utils

import "reflect"

func IsInterfaceZero(i interface{}) bool {
	return isInterfaceZero(reflect.ValueOf(i))
}

func isInterfaceZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Func, reflect.Map, reflect.Slice:
		if v.IsNil() {
			return true
		}
		if v.Len() == 0 {
			return true
		}
		return false
	case reflect.Array:
		z := true
		for i := 0; i < v.Len(); i++ {
			z = z && isInterfaceZero(v.Index(i))
		}
		return z
	case reflect.Struct:
		z := true
		for i := 0; i < v.NumField(); i++ {
			z = z && isInterfaceZero(v.Field(i))
		}
		return z
	}
	// Compare other types directly:
	z := reflect.Zero(v.Type())
	return v.Interface() == z.Interface()
}
