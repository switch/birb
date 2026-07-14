package handlers

import (
	"fmt"
	"reflect"

	"github.com/switch/birb/matchers"
	"github.com/switch/birb/types"
)

// Answer is a function type that produces return values from method arguments.
// Used by ThenAnswer to provide dynamic return values based on call arguments.
type Answer = func([]any) []any

// BirbMocker configures stubbed behavior for a mocked method call.
// Created by WhenCalling and configured with ThenReturn or ThenAnswer.
type BirbMocker struct {
	methodHandler *MethodHandler
	matcher       *matchers.CallMatcher
	unanswered    []Answer
	lastAnswer    Answer
	sourceFile    string // file where stub was defined (for debugging)
	sourceLine    int    // line where stub was defined (for debugging)

	// Exhaustion tracking
	timesLimit    int  // 0 = unlimited (sticky), >0 = limited calls before exhaustion
	timesUsed     int  // number of times this stub has been used
	wasInvoked    bool // true if this stub was ever matched (for VerifyAllMatchersCalled)
}

// SetSourceLocation records where this stub was defined for debugging purposes.
func (m *BirbMocker) SetSourceLocation(file string, line int) {
	m.sourceFile = file
	m.sourceLine = line
}

// sourceLocation returns a formatted string of where this stub was defined.
func (m *BirbMocker) sourceLocation() string {
	if m.sourceFile == "" {
		return ""
	}
	return fmt.Sprintf(" (stub defined at %s:%d)", m.sourceFile, m.sourceLine)
}

func (m *BirbMocker) testingT() types.TestingT {
	return m.methodHandler.testingT()
}

func (m *BirbMocker) returnValues(method reflect.Method, args []any) []any {
	m.testingT().Helper()

	// Track invocation
	m.wasInvoked = true
	m.timesUsed++

	err := m.matcher.Invoke(args)
	if err != nil {
		m.testingT().Fatalf("Mock: %s: matcher invocation error: %v%s", method.Name, err, m.sourceLocation())
	}
	if method.Type.NumOut() == 0 {
		return nil
	}

	answer := m.getNextAnswer()
	result := answer(args)

	if err := validateReturnTypes(method, result); err != nil {
		m.testingT().Fatalf("Mock: %s: return type mismatch: %v%s", method.Name, err, m.sourceLocation())
	}

	return result
}

// validateReturnTypes checks that the return value count matches the method signature.
// We intentionally only validate the count, not the types, because:
//   - Go's common pattern is "return struct, accept interface" - mocks often return
//     concrete types or mocks that implement interfaces
//   - The generated mock code does type assertions anyway, so type mismatches will
//     surface as clear panics at the call site
//   - Being too strict prevents valid patterns like returning mocks for interfaces
func validateReturnTypes(method reflect.Method, values []any) error {
	numOut := method.Type.NumOut()

	if len(values) != numOut {
		return fmt.Errorf("expected %d return values, got %d", numOut, len(values))
	}

	return nil
}

func (m *BirbMocker) getNextAnswer() Answer {
	m.testingT().Helper()
	if len(m.unanswered) == 0 {
		return m.lastAnswer
	}
	last := m.unanswered[0]
	m.unanswered = m.unanswered[1:]
	m.lastAnswer = last
	return last
}

func (m *BirbMocker) matches(method reflect.Method, args []any) bool {
	m.testingT().Helper()

	// Check if this stub is exhausted (times limit reached)
	if m.IsExhausted() {
		return false
	}

	ok, err := m.matcher.Match(method, args)
	if err != nil {
		m.testingT().Fatalf("Mock: %s: argument matching error: %v%s", method.Name, err, m.sourceLocation())
	}

	return ok
}

// IsExhausted returns true if this stub has reached its times limit.
func (m *BirbMocker) IsExhausted() bool {
	if m.timesLimit == 0 {
		return false // unlimited
	}
	return m.timesUsed >= m.timesLimit
}

// WasInvoked returns true if this stub was ever matched and used.
func (m *BirbMocker) WasInvoked() bool {
	return m.wasInvoked
}

// Specificity returns the specificity score of this stub's matcher.
// Higher scores indicate more specific matchers that should take priority.
func (m *BirbMocker) Specificity() int {
	return m.matcher.Specificity()
}

// IsFallback returns true if this is a fallback stub (lowest priority).
func (m *BirbMocker) IsFallback() bool {
	return m.matcher.IsFallback()
}

// MarkAsFallback marks this stub as a fallback (lowest priority).
func (m *BirbMocker) MarkAsFallback() *BirbMocker {
	m.matcher.MarkAsFallback()
	return m
}

// Describe returns a human-readable description of this stub for error messages.
func (m *BirbMocker) Describe() string {
	return m.matcher.Describe()
}

// Status returns a human-readable status string for error messages.
func (m *BirbMocker) Status() string {
	if m.timesLimit == 0 {
		return fmt.Sprintf("[UNLIMITED, used %d times]", m.timesUsed)
	}
	if m.IsExhausted() {
		return fmt.Sprintf("[EXHAUSTED after %d call(s)]", m.timesUsed)
	}
	return fmt.Sprintf("[%d/%d calls used]", m.timesUsed, m.timesLimit)
}

// TimesUsed returns the number of times this stub has been invoked.
func (m *BirbMocker) TimesUsed() int {
	return m.timesUsed
}

// TimesLimit returns the configured times limit (0 = unlimited).
func (m *BirbMocker) TimesLimit() int {
	return m.timesLimit
}

// ThenReturn configures the mock to return the given values when called.
// Can be chained to provide different return values on successive calls.
func (m *BirbMocker) ThenReturn(values ...any) *BirbMocker {
	m.testingT().Helper()
	m.ThenAnswer(func(args []any) []any {
		return values
	})
	return m
}

// ThenAnswer configures the mock to compute return values dynamically.
// The answer function receives the call arguments and returns the values.
func (m *BirbMocker) ThenAnswer(answer Answer) *BirbMocker {
	m.testingT().Helper()
	m.unanswered = append(m.unanswered, answer)
	return m
}

// Times limits this stub to the given number of invocations before exhaustion.
// After exhaustion, the stub will no longer match and fall through to less specific stubs.
// Times(0) will panic as it's almost certainly a bug.
func (m *BirbMocker) Times(n int) *BirbMocker {
	m.testingT().Helper()
	if n <= 0 {
		panic("birb: Times(n) requires n > 0; use Times(1) for a single use or don't use Times() for unlimited")
	}
	m.timesLimit = n
	return m
}

// Once is syntactic sugar for Times(1).
// The stub will match exactly once, then fall through to less specific stubs.
func (m *BirbMocker) Once() *BirbMocker {
	return m.Times(1)
}

func newBirbMocker(m *MethodHandler, args []matchers.Matcher) *BirbMocker {
	m.testingT().Helper()
	return &BirbMocker{
		methodHandler: m,
		matcher:       matchers.CreateCallMatcher(args),
	}
}
