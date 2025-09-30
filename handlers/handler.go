package handlers

import (
	"reflect"
	"sync"

	"github.com/switch/birb/matchers"
	"github.com/switch/birb/types"
)

type Handler struct {
	mutex          sync.Mutex
	methodHandlers *MethodHandlerCollection
	testingT       types.TestingT
	verifier       *Verifier
	frozen         bool
}

func (h *Handler) Mock(method reflect.Method, args []reflect.Value) *BirbMocker {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	mh := h.getMethodHandler(method)
	return mh.provisionBirbMocker(matchers.MungToMatchers(method, args...))
}

func (h *Handler) Handle(method reflect.Method, args []any) []any {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if h.frozen {
		h.TestingT().Fatalf("Mock: handler is frozen, cannot handle method %s", method.Name)
	}

	mh := h.getMethodHandler(method)
	// the Handle function could possibly be a Variadic function
	return mh.Handle(method, mungAnyVariadic(method, args))
}

func (h *Handler) Verify(method reflect.Method, args []reflect.Value) {
	h.Verifier().VerifyCall(method, args)
}

func (h *Handler) Verifier() *Verifier {
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

	h.frozen = true
}

func (h *Handler) Unfreeze() {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.frozen = false
}

// NewHandler creates a new Handler instance
func NewHandler(t types.TestingT, tp reflect.Type) *Handler {
	h := Handler{
		mutex:    sync.Mutex{},
		testingT: t,
		frozen:   false,
	}

	h.methodHandlers = NewMethodHandlerCollection(&h, tp)

	return &h
}
