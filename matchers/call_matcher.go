package matchers

import (
	"errors"
	"fmt"
	"reflect"
)

// Specificity levels for matchers (higher = more specific)
const (
	SpecificityFallback  = 0  // MOCKfallback_ stubs - lowest priority
	SpecificityAnything  = 10 // Anything() or MOCKany_ - catch-all
	SpecificityPartial   = 50 // Matchers like ContainSubstring, HasPrefix, etc.
	SpecificityExact     = 100 // Equal() or exact value matchers - highest priority
)

// CallMatcher holds matchers for a complete method call's arguments.
// It handles matching, capturing, and copy-into operations for all arguments.
type CallMatcher struct {
	matchers   []Matcher
	isFallback bool // true for MOCKfallback_ stubs (lowest priority, excluded from VerifyAllMatchersCalled)
}

func (c *CallMatcher) Invoke(args []any) (err error) {
	for i, matcher := range c.matchers {
		if i >= len(args) {
			return fmt.Errorf("matcher index %d out of bounds for args (len=%d)", i, len(args))
		}
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

// CreateCallMatcher creates a new CallMatcher with the given matchers.
func CreateCallMatcher(matchers []Matcher) *CallMatcher {
	return &CallMatcher{
		matchers: matchers,
	}
}

// MarkAsFallback marks this matcher as a fallback stub (lowest priority).
// Fallback stubs are excluded from VerifyAllMatchersCalled checks.
func (c *CallMatcher) MarkAsFallback() {
	c.isFallback = true
}

// IsFallback returns true if this is a fallback stub.
func (c *CallMatcher) IsFallback() bool {
	return c.isFallback
}

// Specificity calculates the overall specificity score for this CallMatcher.
// Higher scores indicate more specific matchers that should take priority.
// The score is the sum of individual matcher specificities.
func (c *CallMatcher) Specificity() int {
	if c.isFallback {
		return SpecificityFallback
	}
	if len(c.matchers) == 0 {
		return SpecificityAnything
	}

	total := 0
	for _, m := range c.matchers {
		total += matcherSpecificity(m)
	}
	return total
}

// matcherSpecificity returns the specificity score for a single matcher.
func matcherSpecificity(m Matcher) int {
	if m == nil {
		return SpecificityAnything
	}

	// Check for birb's built-in catch-all matchers
	switch m.(type) {
	case *Anything, *MatchTheRestOfTheArguments:
		return SpecificityAnything
	}

	// Check if it's an Equal matcher (most specific)
	if getEqualMatcherValue(m) != nil {
		return SpecificityExact
	}

	// For other Gomega matchers (ContainSubstring, HavePrefix, etc.)
	// treat them as partial matches - more specific than Anything but less than Equal
	return SpecificityPartial
}

// Describe returns a human-readable description of this CallMatcher for error messages.
func (c *CallMatcher) Describe() string {
	if c == nil || len(c.matchers) == 0 {
		return "(no matchers)"
	}

	result := "("
	for i, m := range c.matchers {
		if i > 0 {
			result += ", "
		}
		result += describeMatcher(m)
	}
	result += ")"
	return result
}

// describeMatcher returns a human-readable description of a single matcher.
func describeMatcher(m Matcher) string {
	if m == nil {
		return "nil"
	}

	switch m.(type) {
	case *Anything:
		return "Anything()"
	case *MatchTheRestOfTheArguments:
		return "MatchTheRestOfTheArguments()"
	}

	// Check for Equal matcher
	if val := getEqualMatcherValue(m); val != nil {
		return fmt.Sprintf("Equal(%#v)", val)
	}

	// For other Gomega matchers, try to get the type name
	t := reflect.TypeOf(m)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name() + "(...)"
}

// Equivalent checks if two CallMatchers have equivalent matchers.
// This is used to find existing stubs when setting up multiple WhenCalling
// for the same method/argument combination.
// Note: A fallback stub is never equivalent to a non-fallback stub.
func (c *CallMatcher) Equivalent(other *CallMatcher) bool {
	if c == nil || other == nil {
		return c == other
	}
	// Fallback stubs are never equivalent to non-fallback stubs
	// This ensures MOCKfallback_ creates a separate stub from MOCK_ or MOCKany_
	if c.isFallback != other.isFallback {
		return false
	}
	if len(c.matchers) != len(other.matchers) {
		return false
	}
	for i := range c.matchers {
		if !matchersEquivalent(c.matchers[i], other.matchers[i]) {
			return false
		}
	}
	return true
}

// matchersEquivalent checks if two matchers are functionally equivalent.
func matchersEquivalent(a, b Matcher) bool {
	if a == nil || b == nil {
		return a == b
	}

	// Same instance
	if a == b {
		return true
	}

	// Check by type for birb's built-in matchers
	switch a.(type) {
	case *Anything:
		_, ok := b.(*Anything)
		return ok
	case *MatchTheRestOfTheArguments:
		_, ok := b.(*MatchTheRestOfTheArguments)
		return ok
	}

	// For Gomega Equal matchers, compare the underlying expected values
	// Both gomega.Equal() and our CreateMatcherFromValue() return *matchers.EqualMatcher
	aVal := getEqualMatcherValue(a)
	bVal := getEqualMatcherValue(b)
	if aVal != nil && bVal != nil {
		return reflect.DeepEqual(aVal, bVal)
	}

	// For other matchers (custom, Captor, CopyInto), use pointer equality
	return a == b
}

// getEqualMatcherValue extracts the expected value from a Gomega EqualMatcher.
// Returns nil if the matcher is not an EqualMatcher.
func getEqualMatcherValue(m Matcher) any {
	// Gomega's EqualMatcher has an "Expected" field we can access via reflection
	v := reflect.ValueOf(m)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}
	expected := v.FieldByName("Expected")
	if !expected.IsValid() || !expected.CanInterface() {
		return nil
	}
	return expected.Interface()
}
