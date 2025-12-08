package handlers

import (
	"reflect"
	"sync"

	"github.com/switch/birb/matchers"
	"github.com/switch/birb/types"
)

type MethodHandler struct {
	methodName           string
	handler              *Handler
	mutex                sync.Mutex
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
	sm = newBirbMocker(mh, args)
	mh.BirbMockers = append(mh.BirbMockers, sm)

	return
}

func (mh *MethodHandler) Handle(method reflect.Method, args []any) []any {
	mh.testingT().Helper()
	mh.Invoked(args)

	for _, matcher := range mh.BirbMockers {
		if matcher.matches(method, args) {
			return matcher.returnValues(method, args)
		}
	}

	mh.unhandledInvocations = append(mh.unhandledInvocations, args)
	return mh.defaultReturn
}

func (mh *MethodHandler) Invoked(args []any) {
	mh.testingT().Helper()
	mh.callInvocations.Calls = append(mh.callInvocations.Calls, newCall(mh.methodName, args))
}

func (mh *MethodHandler) Calls() int {
	mh.testingT().Helper()
	return len(mh.callInvocations.Calls)
}

func newMethodHandler(h *Handler, name string, m reflect.Type) *MethodHandler {
	h.TestingT().Helper()
	return &MethodHandler{
		methodName:           m.Name(),
		handler:              h,
		mutex:                sync.Mutex{},
		BirbMockers:          make([]*BirbMocker, 0),
		callInvocations:      newCallCollection(name),
		defaultReturn:        defaultReturnValues(m),
		unhandledInvocations: make([]any, 0),
	}
}

func defaultReturnValues(m reflect.Type) []any {
	retVal := make([]any, m.NumOut())
	for i := 0; i < m.NumOut(); i++ {
		switch m.Out(i).Kind() {
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
			retVal[i] = reflect.MakeSlice(m.Out(i), 0, 0).Interface()
		case reflect.Interface:
			retVal[i] = nil
			// TODO: needed ,,, ???
			// default:
			// 	retVal[i] = nil
		}
	}

	return retVal
}
