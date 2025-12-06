package handlers

import (
	"reflect"

	"github.com/switch/birb/matchers"
	"github.com/switch/birb/types"
)

type Answer = func([]any) []any

type BirbMocker struct {
	methodHandler *MethodHandler
	matcher       *matchers.CallMatcher
	unanswered    []Answer
	lastAnswer    Answer
}

func (m *BirbMocker) testingT() types.TestingT {
	return m.methodHandler.testingT()
}

func (m *BirbMocker) returnValues(method reflect.Method, args []any) []any {
	m.testingT().Helper()
	err := m.matcher.Invoke(args)
	if err != nil {
		m.testingT().Fatalf("Mock: matcher error: %v", err)
	}
	if method.Type.NumOut() == 0 {
		return nil
	}

	answer := m.getNextAnswer()

	return answer(args)
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
	ok, err := m.matcher.Match(method, args)
	if err != nil {
		m.testingT().Fatalf("Mock: matcher error: %v", err)
	}

	return ok
}

func (m *BirbMocker) ThenReturn(values ...any) *BirbMocker {
	m.testingT().Helper()
	m.ThenAnswer(func(args []any) []any {
		return values
	})
	return m
}

func (m *BirbMocker) ThenAnswer(answer Answer) *BirbMocker {
	m.testingT().Helper()
	m.unanswered = append(m.unanswered, answer)
	return m
}

func newBirbMocker(m *MethodHandler, args []matchers.Matcher) *BirbMocker {
	m.testingT().Helper()
	return &BirbMocker{
		methodHandler: m,
		matcher:       matchers.CreateCallMatcher(args),
	}
}
