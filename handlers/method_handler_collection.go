package handlers

import (
	"maps"
	"reflect"

	"github.com/switch/birb/types"
)

// MethodHandlerCollection manages all method handlers for a mock instance.
// It creates handlers lazily based on the interface type's methods.
type MethodHandlerCollection struct {
	handler        *Handler
	methodHandlers map[string]*MethodHandler
}

func (mhc *MethodHandlerCollection) testingT() types.TestingT {
	return mhc.handler.TestingT()
}

func (mhc *MethodHandlerCollection) GetMethodHandler(method string) *MethodHandler {
	mhc.testingT().Helper()
	mh, exists := mhc.methodHandlers[method]
	if !exists {
		mhc.testingT().Fatalf("Mock: method %s not registered in handler", method)
	}

	return mh
}

func (mhc *MethodHandlerCollection) GetMethodHandlers() map[string]*MethodHandler {
	mhc.testingT().Helper()
	return maps.Clone(mhc.methodHandlers)
}

// NewMethodHandlerCollection creates a collection for the given interface type.
func NewMethodHandlerCollection(h *Handler, tp reflect.Type) *MethodHandlerCollection {
	mhc := &MethodHandlerCollection{
		handler:        h,
		methodHandlers: make(map[string]*MethodHandler),
	}

	for i := 0; i < tp.NumMethod(); i++ {
		method := tp.Method(i)
		mhc.methodHandlers[method.Name] = newMethodHandler(h, method.Name, method.Type)
	}

	return mhc
}
