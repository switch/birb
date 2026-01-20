package matchers

import (
	"slices"
	"sync"
)

// Captor is a matcher that captures arguments for later inspection.
// It is thread-safe and can be used in concurrent tests.
type Captor interface {
	Matcher
	Capture(t any)
	GetValues() []any
}

type captor struct {
	mutex  sync.Mutex
	values []any
}

func (c *captor) Match(any) (success bool, err error) {
	// captor always matches, as it is used to capture arguments for later use
	return true, nil
}

func (c *captor) Capture(t any) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.values = append(c.values, t)
}

func (c *captor) GetValues() []any {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return slices.Clone(c.values)
}

func NewCaptor() Captor {
	return &captor{values: make([]any, 0)}
}
