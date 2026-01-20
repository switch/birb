package matchers

import (
	"reflect"

	"github.com/onsi/gomega"
)

// Matcher is the interface for argument matching in mock method calls.
// All matchers must implement this interface to be used with birb mocks.
type Matcher interface {
	Match(actual any) (bool, error)
}

// CreateMatcherFromValue wraps a value in a Gomega Equal matcher.
// This is used internally to convert literal values to matchers.
func CreateMatcherFromValue(val any) Matcher {
	return gomega.Equal(val)
}

// MungToMatchers converts reflect.Value arguments to Matcher instances.
// Values that are already Matchers are used directly; others are wrapped with Equal.
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
