package handlers

import (
	"fmt"
	"reflect"
	"slices"

	"github.com/onsi/gomega"
	"github.com/switch/birb/matchers"
)

type callData struct {
	MethodCalls int
}

type callCollection struct {
	MethodName string
	Calls      []*call
}

type call struct {
	MethodName string
	Args       []any
}

func (c *callCollection) matches(method reflect.Method, args []any) (calls int) {
	for _, invocation := range c.Calls {
		if invocation.matches(method, args) {
			calls++
		}
	}

	return
}

func (c *call) matches(method reflect.Method, args []any) bool {
	ignoreTheRestOfTheArguments := false
	// if method.Type.IsVariadic() {
	if len(args) > 0 {
		_, ignoreTheRestOfTheArguments = args[len(args)-1].(*matchers.MatchTheRestOfTheArguments)
	}

	if !ignoreTheRestOfTheArguments && len(c.Args) != len(args) {
		return false
	}

	for i, arg := range c.Args {
		m := i
		if m >= len(args) {
			m = len(args) - 1
		}
		if matcher, ok := args[m].(matchers.Matcher); ok {
			if success, _ := matcher.Match(arg); !success {
				return false
			}
		} else if success, _ := gomega.Equal(args[m]).Match(arg); !success {
			return false
		}
	}
	return true
}

func Times(n int) *CallVerifier {
	return MethodVerifierFromCallData(func(data *callData) error {
		if data.MethodCalls != n {
			return fmt.Errorf("expected num method calls: %d, got : %d", n, data.MethodCalls)
		}
		return nil
	})
}

func AtLeast(n int) *CallVerifier {
	return MethodVerifierFromCallData(func(data *callData) error {
		if n >= data.MethodCalls {
			return fmt.Errorf("expected at least num method calls: %d, got : %d", n, data.MethodCalls)
		}
		return nil
	})
}

func AtMost(n int) *CallVerifier {
	return MethodVerifierFromCallData(func(data *callData) error {
		if n <= data.MethodCalls {
			return fmt.Errorf("expected at most num method calls: %d, got : %d", n, data.MethodCalls)
		}
		return nil
	})
}

func MethodVerifierFromCallData(f func(data *callData) error) *CallVerifier {
	return &CallVerifier{
		f: f,
	}
}

type CallVerifier struct {
	f func(data *callData) error
}

func (fcv *CallVerifier) Verify(data *callData) error {
	return fcv.f(data)
}

func newCallCollection(method string) *callCollection {
	return &callCollection{
		MethodName: method,
		Calls:      make([]*call, 0),
	}
}

func newCall(method string, args []any) *call {
	return &call{
		MethodName: method,
		Args:       slices.Clone(args),
	}
}

func (c *callCollection) String() string {
	rv := c.MethodName + ":\n"
	for _, fc := range c.Calls {
		rv += fmt.Sprintf("  - %s\n", fc.String())
	}

	return rv
}

func (c *call) String() string {
	return fmt.Sprintf("%v", c.Args)
}
