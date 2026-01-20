package handlers

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/switch/birb/types"
)

// Verifier tracks and validates method invocations on a mock.
// Used internally by Verify, VerifyNeverCalled, and VerifyNoOtherInteractions.
type Verifier struct {
	mutex   sync.Mutex
	wip     *CallVerifier
	handler *Handler
}

func (f *Verifier) getMethodHandler(method string) *MethodHandler {
	f.testingT().Helper()
	return f.handler.methodHandlers.GetMethodHandler(method)
}

func (f *Verifier) getMethodHandlers() map[string]*MethodHandler {
	f.testingT().Helper()
	return f.handler.methodHandlers.GetMethodHandlers()
}

func (f *Verifier) testingT() types.TestingT {
	return f.handler.TestingT()
}

func (f *Verifier) VerifyNoOtherInteractions() {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.testingT().Helper()

	var err error

	for method, methodHandler := range f.getMethodHandlers() {
		if len(methodHandler.unhandledInvocations) > 0 {
			err = errors.Join(err, fmt.Errorf("%s: %d unhandled invocation(s): %v",
				method, len(methodHandler.unhandledInvocations), methodHandler.unhandledInvocations))
		}
	}

	if err != nil {
		f.testingT().Fatalf("Verifier: unexpected interactions:\n  %s", err)
	}
}

func (f *Verifier) VerifyNeverCalled() {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.testingT().Helper()

	var err error
	for method, methodHandler := range f.getMethodHandlers() {
		calls := methodHandler.Calls()
		if calls > 0 {
			err = errors.Join(err, fmt.Errorf("%s: expected 0 call(s), got %d", method, calls))
		}
	}

	if err != nil {
		f.testingT().Fatalf("Verifier: expected no interactions:\n  %v", err)
	}
}

func (f *Verifier) PrepCallVerifier(verifier *CallVerifier) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.testingT().Helper()

	if f.wip != nil {
		f.testingT().Fatalf("Verifier: internal error: call verifier already set (did you call Verify twice without a CALLED_ method?)")
	}
	f.wip = verifier
}

func (f *Verifier) VerifyCall(method reflect.Method, args []reflect.Value) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.testingT().Helper()

	if f.wip == nil {
		f.testingT().Fatalf("Verifier: internal error: call verifier not set (use Verify(mock, Times(n)).CALLED_Method() pattern)")
	}

	verifier := f.wip
	f.wip = nil

	mh := f.getMethodHandler(method.Name)
	fcd := &callData{
		MethodCalls: 0,
	}

	// Verifying CALLED_ functions will NOT be variadic
	var matchErr error
	fcd.MethodCalls, matchErr = mh.callInvocations.matches(method, mungReflectValuesToAny(method, args))
	if matchErr != nil {
		f.testingT().Fatalf("Verifier: matcher error during verification of %s: %v", method.Name, matchErr)
	}

	err := verifier.Verify(fcd)
	if err != nil {
		f.testingT().Fatalf("Verifier: %s: %v\n%s", method.Name, err, mh.callInvocations.String())
	}
}

// NewVerifier creates a new Verifier for the given Handler.
func NewVerifier(handler *Handler) *Verifier {
	return &Verifier{
		mutex:   sync.Mutex{},
		handler: handler,
	}
}

// VerifyAllMatchersCalled checks that every non-fallback stub was invoked at least once.
// This helps catch "dead" stubs that might indicate test setup errors or refactoring
// that made stubs obsolete.
// Fallback stubs (created with MOCKfallback_) are excluded from this check.
func (f *Verifier) VerifyAllMatchersCalled() {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.testingT().Helper()

	var uninvoked []stubInfo

	for methodName, methodHandler := range f.getMethodHandlers() {
		for _, mocker := range methodHandler.BirbMockers {
			// Skip fallback stubs - they're optional by design
			if mocker.IsFallback() {
				continue
			}

			if !mocker.WasInvoked() {
				uninvoked = append(uninvoked, stubInfo{
					methodName: methodName,
					describe:   mocker.Describe(),
					status:     mocker.Status(),
				})
			}
		}
	}

	if len(uninvoked) > 0 {
		msg := "VerifyAllMatchersCalled: the following stubs were never invoked:\n"
		for _, stub := range uninvoked {
			msg += fmt.Sprintf("  - %s%s %s\n", stub.methodName, stub.describe, stub.status)
		}
		msg += "\nTip: Remove unused stubs or use MOCKfallback_ for optional stubs."
		f.testingT().Fatalf(msg)
	}
}

// stubInfo holds information about a stub for error reporting.
type stubInfo struct {
	methodName string
	describe   string
	status     string
}
