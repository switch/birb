package handlers

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/switch/birb/matchers"
	"github.com/switch/birb/types"
)

// MethodHandler manages stubs and invocations for a single mocked method.
// It tracks call history and matches incoming calls against configured stubs.
// Thread safety is provided by Handler.mutex which wraps all Mock() and Handle() calls.
type MethodHandler struct {
	methodName           string
	handler              *Handler
	BirbMockers          []*BirbMocker
	callInvocations      *callCollection
	defaultReturn        []any
	unhandledInvocations []any
}

func (mh *MethodHandler) testingT() types.TestingT {
	return mh.handler.TestingT()
}

func (mh *MethodHandler) provisionBirbMocker(args []matchers.Matcher) (sm *BirbMocker) {
	mh.testingT().Helper()
	return mh.provisionBirbMockerWithOptions(args, false)
}

// provisionBirbMockerAsFallback creates a new fallback stub (always creates new, never reuses).
func (mh *MethodHandler) provisionBirbMockerAsFallback(args []matchers.Matcher) (sm *BirbMocker) {
	mh.testingT().Helper()
	return mh.provisionBirbMockerWithOptions(args, true)
}

func (mh *MethodHandler) provisionBirbMockerWithOptions(args []matchers.Matcher, isFallback bool) (sm *BirbMocker) {
	mh.testingT().Helper()

	// Check if there's an existing stub with equivalent matchers
	newMatcher := matchers.CreateCallMatcher(args)
	if isFallback {
		newMatcher.MarkAsFallback()
	}

	for _, existing := range mh.BirbMockers {
		if existing.matcher.Equivalent(newMatcher) {
			return existing
		}
	}

	// No equivalent stub found, create a new one
	sm = newBirbMocker(mh, args)
	if isFallback {
		sm.MarkAsFallback()
	}
	mh.BirbMockers = append(mh.BirbMockers, sm)

	return
}

func (mh *MethodHandler) Handle(method reflect.Method, args []any) []any {
	mh.testingT().Helper()
	mh.Invoked(args)

	// Sort stubs by specificity (highest first) so more specific matchers take priority
	sortedMockers := mh.stubsBySpecificity()

	for _, mocker := range sortedMockers {
		if mocker.matches(method, args) {
			return mocker.returnValues(method, args)
		}
	}

	// No matching stub found
	mh.unhandledInvocations = append(mh.unhandledInvocations, args)

	// If there are registered stubs but none matched, panic with helpful message
	// (This means all stubs are either exhausted or have non-matching matchers)
	if len(mh.BirbMockers) > 0 {
		mh.panicNoMatchingStub(method, args)
	}

	return mh.defaultReturn
}

// stubsBySpecificity returns stubs sorted by specificity (highest first).
// This ensures more specific matchers are checked before catch-all matchers.
func (mh *MethodHandler) stubsBySpecificity() []*BirbMocker {
	// Create a copy to avoid modifying the original slice
	sorted := make([]*BirbMocker, len(mh.BirbMockers))
	copy(sorted, mh.BirbMockers)

	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Specificity() > sorted[j].Specificity()
	})

	return sorted
}

// panicNoMatchingStub panics with a helpful error message when no stub matches.
func (mh *MethodHandler) panicNoMatchingStub(method reflect.Method, args []any) {
	msg := mh.buildNoMatchingStubMessage(method, args)
	panic(msg)
}

// buildNoMatchingStubMessage builds a detailed error message for debugging.
func (mh *MethodHandler) buildNoMatchingStubMessage(method reflect.Method, args []any) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("birb: no matching stub for %s.%s\n\n", mh.handler.mockName, method.Name))

	// Show what was called with
	sb.WriteString("  Called with:\n")
	if len(args) == 0 {
		sb.WriteString("    (no arguments)\n")
	} else {
		for i, arg := range args {
			sb.WriteString(fmt.Sprintf("    arg[%d]: %#v (%T)\n", i, arg, arg))
		}
	}
	sb.WriteString("\n")

	// Show registered stubs and their status
	sb.WriteString(fmt.Sprintf("  Registered stubs for %s:\n", method.Name))
	hasFallback := false
	for _, mocker := range mh.BirbMockers {
		prefix := "✗"
		if mocker.IsFallback() {
			hasFallback = true
			prefix = "⚡" // fallback indicator
		}
		sb.WriteString(fmt.Sprintf("    %s %s %s\n", prefix, mocker.Describe(), mocker.Status()))
	}
	sb.WriteString("\n")

	// Provide actionable advice
	if !hasFallback {
		sb.WriteString("  No fallback defined. Consider adding:\n")
		sb.WriteString(fmt.Sprintf("    WhenCalling(mock.MOCKfallback_%s()).ThenReturn(...)\n\n", method.Name))
	}

	sb.WriteString("  Tip: Use .Times(n) to limit calls, or omit for unlimited.\n")

	return sb.String()
}

func (mh *MethodHandler) Invoked(args []any) {
	mh.testingT().Helper()
	mh.callInvocations.Calls = append(mh.callInvocations.Calls, newCall(mh.methodName, args))
}

func (mh *MethodHandler) Calls() int {
	mh.testingT().Helper()
	return len(mh.callInvocations.Calls)
}

// Reset clears all stubs and call history for this method.
func (mh *MethodHandler) Reset() {
	mh.BirbMockers = make([]*BirbMocker, 0)
	mh.callInvocations.Calls = make([]*call, 0)
	mh.unhandledInvocations = make([]any, 0)
}

func newMethodHandler(h *Handler, name string, m reflect.Type) *MethodHandler {
	h.TestingT().Helper()
	return &MethodHandler{
		methodName:           m.Name(),
		handler:              h,
		BirbMockers:          make([]*BirbMocker, 0),
		callInvocations:      newCallCollection(name),
		defaultReturn:        defaultReturnValues(m),
		unhandledInvocations: make([]any, 0),
	}
}

func defaultReturnValues(m reflect.Type) []any {
	retVal := make([]any, m.NumOut())
	for i := 0; i < m.NumOut(); i++ {
		outType := m.Out(i)
		switch outType.Kind() {
		case reflect.Bool:
			retVal[i] = false
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			retVal[i] = 0
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			retVal[i] = 0
		case reflect.Float32, reflect.Float64:
			retVal[i] = 0.0
		case reflect.String:
			retVal[i] = ""
		case reflect.Slice, reflect.Array:
			retVal[i] = reflect.MakeSlice(outType, 0, 0).Interface()
		case reflect.Map:
			retVal[i] = reflect.MakeMap(outType).Interface()
		case reflect.Struct:
			retVal[i] = reflect.Zero(outType).Interface()
		case reflect.Interface, reflect.Ptr, reflect.Chan, reflect.Func:
			retVal[i] = nil
		}
	}

	return retVal
}
