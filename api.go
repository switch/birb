// Package birb provides a simple interface for creating and verifying mock interactions.
package birb

import (
	"runtime"

	"github.com/switch/birb/handlers"
	"github.com/switch/birb/matchers"
)

// Birb is an interface that defines the methods required for a Birb mock.
type Birb interface {
	BirbHandler() *handlers.Handler
}

// WhenCalling is the API to use for setting up method stubs on a Birb mock.
// It captures the source location for better debugging when stub errors occur.
func WhenCalling(ans *handlers.BirbMocker) *handlers.BirbMocker {
	// Capture where the stub was defined (skip 1 = WhenCalling, caller is user code)
	if _, file, line, ok := runtime.Caller(1); ok {
		ans.SetSourceLocation(file, line)
	}
	return ans
}

// Verify will verify that the Birb instance has been called as expected.
func Verify[Mock Birb](mock Mock, v *handlers.CallVerifier) Mock {
	mock.BirbHandler().Verifier().PrepCallVerifier(v)
	return mock
}

// VerifyNeverCalled verifies that no methods were called on the Birb instance that were not already verified.
func VerifyNeverCalled[Mock Birb](mock Mock) {
	mock.BirbHandler().Verifier().VerifyNeverCalled()
}

// VerifyNoOtherInteractions verifies that no methods were called on the Birb instance that were not already verified.
func VerifyNoOtherInteractions[Mock Birb](mock Mock) {
	mock.BirbHandler().Verifier().VerifyNoOtherInteractions()
}

// VerifyAllMatchersCalled checks that every non-fallback stub was invoked at least once.
// This helps catch "dead" stubs that might indicate test setup errors or refactoring
// that made stubs obsolete. Fallback stubs (created with MOCKfallback_) are excluded.
//
// Example:
//
//	WhenCalling(mock.MOCK_Foo(Equal("x"))).ThenReturn("result")
//	// ... test code that should call mock.Foo("x") ...
//	VerifyAllMatchersCalled(mock) // fails if MOCK_Foo was never called
func VerifyAllMatchersCalled[Mock Birb](mock Mock) {
	mock.BirbHandler().Verifier().VerifyAllMatchersCalled()
}

// Freeze verifies that no methods will be called on the Birb instance.
func Freeze[Mock Birb](mock Mock) {
	mock.BirbHandler().Freeze()
}

// Reset clears all stubs, call history, and verification state on the mock.
// This allows reusing a mock across multiple test cases without creating a new instance.
// After Reset(), the mock behaves as if it was freshly created.
func Reset[Mock Birb](mock Mock) {
	mock.BirbHandler().Reset()
}

// Times verify that a method is called exactly n times.
func Times(n int) *handlers.CallVerifier {
	return handlers.Times(n)
}

// Never helper function to create a verifier that expects the method to have never been called.
func Never() *handlers.CallVerifier {
	return Times(0)
}

// Once helper function to create a verifier that expects the method to be called exactly once.
func Once() *handlers.CallVerifier {
	return Times(1)
}

// Twice helper function to create a verifier that expects the method to be called exactly twice.
func Twice() *handlers.CallVerifier {
	return Times(2)
}

// AtLeast is a helper function to create verifiers that expect the method to be called at least n times.
func AtLeast(n int) *handlers.CallVerifier {
	return handlers.AtLeast(n)
}

// AtMost is a helper function to create verifiers that expect the method to be called at most n times.
func AtMost(n int) *handlers.CallVerifier {
	return handlers.AtMost(n)
}

// Between is a helper function to create verifiers that expect the method to be called
// between min and max times (inclusive).
func Between(min, max int) *handlers.CallVerifier {
	return handlers.Between(min, max)
}

// Anything is a matcher that matches any SINGULAR value. This is useful when you don't care about the specific value passed to a method.
func Anything() matchers.Matcher {
	return &matchers.Anything{}
}

// WithAnyArgs is a matcher that matches any values. This is useful when you don't care about the specific value passed to a method.
func WithAnyArgs() matchers.Matcher {
	return &matchers.MatchTheRestOfTheArguments{}
}

// Captor is a matcher that captures the value passed to a method. This is useful when you want to verify the value passed to a method after it has been called.
func Captor() matchers.Captor {
	return matchers.NewCaptor()
}

// DeepCopyInto creates a matcher that deep-copies toCopy into the argument.
// The type T must implement the DeepCopyInto[T] interface (common in K8s types).
func DeepCopyInto[T matchers.DeepCopyInto[T]](toCopy T) matchers.CopyInto {
	return matchers.NewDeepCopyInto(toCopy)
}

// CopyIntoFunc creates a matcher using a custom copy function.
// Use this when you need control over how values are copied into arguments.
func CopyIntoFunc[T any](toCopy T, copyFunc matchers.CopyIntoFunc[T]) matchers.CopyInto {
	return matchers.NewCopyIntoFunc(toCopy, copyFunc)
}

// ShallowCopyInto is a convenience function to create a CopyInto matcher that copies the value from a pointer.
func ShallowCopyInto[T any](toCopy *T) matchers.CopyInto {
	return matchers.NewCopyIntoPointer(toCopy)
}
