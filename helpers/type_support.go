package helpers

import "reflect"

func isFunction(f any) bool {
	actual := reflect.TypeOf(f)
	return actual.Kind() == reflect.Func && actual.NumIn() == 0 && actual.NumOut() > 0
}

func isChan(a any) bool {
	if isNil(a) {
		return false
	}
	return reflect.TypeOf(a).Kind() == reflect.Chan
}

func isNil(a any) bool {
	if a == nil {
		return true
	}

	switch reflect.TypeOf(a).Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return reflect.ValueOf(a).IsNil()
	}

	return false
}
