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

func (c *callCollection) matches(method reflect.Method, args []any) (calls int, matchErr error) {
	for _, invocation := range c.Calls {
		matched, err := invocation.matches(method, args)
		if err != nil {
			matchErr = err
			return
		}
		if matched {
			calls++
		}
	}

	return
}

func (c *call) matches(method reflect.Method, args []any) (bool, error) {
	ignoreTheRestOfTheArguments := false
	if len(args) > 0 {
		_, ignoreTheRestOfTheArguments = args[len(args)-1].(*matchers.MatchTheRestOfTheArguments)
	}

	if !ignoreTheRestOfTheArguments && len(c.Args) != len(args) {
		return false, nil
	}

	for i, arg := range c.Args {
		m := i
		if m >= len(args) {
			m = len(args) - 1
		}
		if matcher, ok := args[m].(matchers.Matcher); ok {
			success, err := matcher.Match(arg)
			if err != nil {
				return false, fmt.Errorf("argument %d: matcher error: %w", i, err)
			}
			if !success {
				return false, nil
			}
		} else {
			success, err := gomega.Equal(args[m]).Match(arg)
			if err != nil {
				return false, fmt.Errorf("argument %d: equality check error: %w", i, err)
			}
			if !success {
				return false, nil
			}
		}
	}
	return true, nil
}

// Times creates a verifier that expects exactly n method calls.
func Times(n int) *CallVerifier {
	return MethodVerifierFromCallData(func(data *callData) error {
		if data.MethodCalls != n {
			return fmt.Errorf("expected %d call(s), got %d", n, data.MethodCalls)
		}
		return nil
	})
}

// AtLeast creates a verifier that expects at least n method calls.
func AtLeast(n int) *CallVerifier {
	return MethodVerifierFromCallData(func(data *callData) error {
		if data.MethodCalls < n {
			return fmt.Errorf("expected at least %d call(s), got %d", n, data.MethodCalls)
		}
		return nil
	})
}

// AtMost creates a verifier that expects at most n method calls.
func AtMost(n int) *CallVerifier {
	return MethodVerifierFromCallData(func(data *callData) error {
		if data.MethodCalls > n {
			return fmt.Errorf("expected at most %d call(s), got %d", n, data.MethodCalls)
		}
		return nil
	})
}

// Between creates a verifier that expects between min and max method calls (inclusive).
func Between(min, max int) *CallVerifier {
	return MethodVerifierFromCallData(func(data *callData) error {
		if data.MethodCalls < min || data.MethodCalls > max {
			return fmt.Errorf("expected between %d and %d call(s), got %d", min, max, data.MethodCalls)
		}
		return nil
	})
}

// MethodVerifierFromCallData creates a custom CallVerifier from a validation function.
func MethodVerifierFromCallData(f func(data *callData) error) *CallVerifier {
	return &CallVerifier{
		f: f,
	}
}

// CallVerifier validates method call counts during verification.
// Use Times, AtLeast, AtMost, or MethodVerifierFromCallData to create one.
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
