package matchers

import (
	"github.com/onsi/gomega"
	"reflect"
)

type Matcher interface {
	Match(actual any) (bool, error)
}

func CreateMatcherFromValue(val any) Matcher {
	return gomega.Equal(val)
}

func MungToMatchers(method reflect.Method, args ...reflect.Value) []Matcher {
	fixedArgs := make([]Matcher, len(args))
	for i, arg := range args {
		if matcher, ok := arg.Interface().(Matcher); ok {
			fixedArgs[i] = matcher
		} else {
			fixedArgs[i] = CreateMatcherFromValue(arg.Interface())
		}
	}

	// if the function is variadic, the last argument is the variadic arg
	if method.Type.IsVariadic() && len(fixedArgs) >= method.Type.NumIn() {
		fixedArgs = fixedArgs[:len(fixedArgs)-1]

		last := args[len(args)-1]
		switch last.Kind() {
		case reflect.Slice, reflect.Array:
			for i := 0; i < last.Len(); i++ {
				possibleMatcher := last.Index(i).Interface()
				matcher, ok := possibleMatcher.(Matcher)
				if ok {
					fixedArgs = append(fixedArgs, matcher)
				} else {
					fixedArgs = append(fixedArgs, CreateMatcherFromValue(possibleMatcher))
				}
			}
		default:
		}

	}

	return fixedArgs
}
