package matchers

import (
	"slices"
)

type Captor interface {
	Matcher
	Capture(t any)
	GetValues() []any
}

type captor struct {
	values []any
}

func (c *captor) Match(any) (success bool, err error) {
	// captor always matches, as it is used to capture arguments for later use
	return true, nil
}

func (c *captor) Capture(t any) {
	c.values = append(c.values, t)
}

func (c *captor) GetValues() []any {
	return slices.Clone(c.values)
}

func NewCaptor() Captor {
	return &captor{values: make([]any, 0)}
}
