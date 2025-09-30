package matchers

import (
	"errors"
	"fmt"
	"reflect"
)

type CallMatcher struct {
	matchers []Matcher
}

func (c *CallMatcher) Invoke(args []any) (err error) {
	for i, matcher := range c.matchers {
		if captor, ok := matcher.(Captor); ok {
			captor.Capture(args[i])
		}
		if copyInto, ok := matcher.(CopyInto); ok {
			err = errors.Join(err, copyInto.CopyInto(args[i]))
		}
	}

	return
}

func (c *CallMatcher) Match(method reflect.Method, args []any) (ok bool, err error) {
	if len(c.matchers) > len(args) {
		return false, fmt.Errorf("expected %d matchers, got %d", len(c.matchers), len(args))
	}

	if len(c.matchers) == 0 && len(args) != 0 {
		return false, fmt.Errorf("expected at least one matcher")
	}

	ok = true

	for i, arg := range args {
		m := i
		if i >= len(c.matchers) {
			m = len(c.matchers) - 1
		}
		ok, err = c.matchers[m].Match(arg)
		if !ok || err != nil {
			return
		}
	}

	return
}

func CreateCallMatcher(matchers []Matcher) *CallMatcher {
	return &CallMatcher{
		matchers: matchers,
	}
}
