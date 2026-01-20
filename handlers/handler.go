package handlers

import (
	"reflect"
	"sync"

	"github.com/switch/birb/matchers"
	"github.com/switch/birb/types"
)

// Handler is the core mock handler that manages method stubs and verifications.
// Each generated mock contains a Handler that coordinates all mock behavior.
type Handler struct {
	mutex          sync.Mutex
	methodHandlers *MethodHandlerCollection
	testingT       types.TestingT
	verifier       *Verifier
	frozen         bool
	mockName       string // Name of the mock type for error messages
}

func (h *Handler) Mock(method reflect.Method, args []reflect.Value) *BirbMocker {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.testingT.Helper()

	mh := h.getMethodHandler(method)
	return mh.provisionBirbMocker(matchers.MungToMatchers(method, args...))
}

// MockFallback creates a fallback stub for the given method.
// Fallback stubs have the lowest priority and are excluded from VerifyAllMatchersCalled.
func (h *Handler) MockFallback(method reflect.Method, args []reflect.Value) *BirbMocker {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.testingT.Helper()

	mh := h.getMethodHandler(method)
	return mh.provisionBirbMockerAsFallback(matchers.MungToMatchers(method, args...))
}

func (h *Handler) Handle(method reflect.Method, args []any) []any {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.testingT.Helper()

	if h.frozen {
		h.TestingT().Fatalf("Mock: %s: handler is frozen, cannot invoke method", method.Name)
	}

	mh := h.getMethodHandler(method)
	// the Handle function could possibly be a Variadic function
	return mh.Handle(method, mungAnyVariadic(method, args))
}

func (h *Handler) Verify(method reflect.Method, args []reflect.Value) {
	h.testingT.Helper()
	h.Verifier().VerifyCall(method, args)
}

func (h *Handler) Verifier() *Verifier {
	h.testingT.Helper()
	if h.verifier == nil {
		h.verifier = NewVerifier(h)
	}
	return h.verifier
}

func (h *Handler) TestingT() types.TestingT {
	return h.testingT
}

func (h *Handler) getMethodHandler(method reflect.Method) *MethodHandler {
	return h.methodHandlers.GetMethodHandler(method.Name)
}

func (h *Handler) Freeze() {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.testingT.Helper()

	h.frozen = true
}

func (h *Handler) Unfreeze() {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.testingT.Helper()

	h.frozen = false
}

// Reset clears all stubs, call history, and verification state.
// This allows reusing a mock across multiple test cases without creating a new instance.
// After Reset(), the mock behaves as if it was freshly created.
func (h *Handler) Reset() {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.testingT.Helper()

	// Reset all method handlers
	for _, mh := range h.methodHandlers.GetMethodHandlers() {
		mh.Reset()
	}

	// Reset verifier state
	h.verifier = nil

	// Reset frozen state
	h.frozen = false
}

// NewHandler creates a new Handler instance
func NewHandler(t types.TestingT, tp reflect.Type) *Handler {
	t.Helper()
	h := Handler{
		mutex:    sync.Mutex{},
		testingT: t,
		frozen:   false,
		mockName: tp.Name(),
	}

	h.methodHandlers = NewMethodHandlerCollection(&h, tp)

	return &h
}

// MockName returns the name of the mock type for error messages.
func (h *Handler) MockName() string {
	return h.mockName
}
