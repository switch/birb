package handlers

import (
	"errors"
	"fmt"
	"github.com/switch/birb/types"
	"reflect"
	"sync"
)

type Verifier struct {
	mutex   sync.Mutex
	wip     *CallVerifier
	handler *Handler
}

func (f *Verifier) getMethodHandler(method string) *MethodHandler {
	return f.handler.methodHandlers.GetMethodHandler(method)
}

func (f *Verifier) getMethodHandlers() map[string]*MethodHandler {
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
			err = errors.Join(err, fmt.Errorf("Verifier: expected no unhandled invocations for method %s, but got %d unhandled invocations: %v",
				method, len(methodHandler.unhandledInvocations), methodHandler.unhandledInvocations))
		}
	}

	if err != nil {
		f.testingT().Fatalf("Expected no other interactions:\n\t%s", err)
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
			err = errors.Join(err, fmt.Errorf("Verifier: expected no interactions with method %s, but got %d invocations", method, calls))
		}
	}

	if err != nil {
		f.testingT().Fatalf("Verifier: %v", err)
	}
}

func (f *Verifier) PrepCallVerifier(verifier *CallVerifier) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if f.wip != nil {
		f.testingT().Fatalf("Verifier: call verifier already set")
	}
	f.wip = verifier
}

func (f *Verifier) VerifyCall(method reflect.Method, args []reflect.Value) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.testingT().Helper()

	if f.wip == nil {
		f.testingT().Fatalf("Verifier: call verifier not set")
	}

	verifier := f.wip
	f.wip = nil

	mh := f.getMethodHandler(method.Name)
	fcd := &callData{
		MethodCalls: 0,
	}

	// Verifying CALLED_ functions will NOT be variadic
	fcd.MethodCalls = mh.callInvocations.matches(method, mungReflectValuesToAny(method, args))

	err := verifier.Verify(fcd)
	if err != nil {
		f.testingT().Fatalf("Verifier: %v\n%s", err, mh.callInvocations.String())
	}
}

func NewVerifier(handler *Handler) *Verifier {
	return &Verifier{
		mutex:   sync.Mutex{},
		handler: handler,
	}
}
