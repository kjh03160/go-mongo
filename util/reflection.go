package util

import "reflect"

func GetInterfaceType(v interface{}) reflect.Type {
	var t reflect.Type
	if xt, ok := v.(reflect.Type); ok {
		t = xt
	} else {
		t = reflect.TypeOf(v)
	}
	return t
}
