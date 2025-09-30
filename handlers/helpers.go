package handlers

import (
	"reflect"
)

func mungReflectValuesToAny(method reflect.Method, args []reflect.Value) []any {
	a := make([]any, len(args))
	for i, arg := range args {
		a[i] = arg.Interface()
	}

	return a
}

func mungAnyVariadic(method reflect.Method, args []any) []any {
	a := make([]any, len(args))
	copy(a, args)

	// if the function is variadic, the last argument is the variadic arg
	if method.Type.IsVariadic() && len(a) >= method.Type.NumIn() {
		a = a[:len(a)-1]

		last := reflect.ValueOf(args[len(args)-1])
		switch last.Kind() {
		case reflect.Slice, reflect.Array:
			for i := 0; i < last.Len(); i++ {
				a = append(a, last.Index(i).Interface())
			}
		default:
		}
	}

	return a
}
