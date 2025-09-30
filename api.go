// Package birb provides a simple interface for creating and verifying mock interactions.
package birb

import (
	"github.com/switch/birb/handlers"
	"github.com/switch/birb/matchers"
)

// Birb is an interface that defines the methods required for a Birb mock.
type Birb interface {
	BirbHandler() *handlers.Handler
}

// WhenCalling is the API to use for setting up method stubs on a Birb mock.
func WhenCalling(ans *handlers.BirbMocker) *handlers.BirbMocker {
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

// Freeze verifies that no methods will be called on the Birb instance.
func Freeze[Mock Birb](mock Mock) {
	mock.BirbHandler().Freeze()
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

// DeepCopyInto
func DeepCopyInto[T matchers.DeepCopyInto[T]](toCopy T) matchers.CopyInto {
	return matchers.NewDeepCopyInto(toCopy)
}

// CopyIntoFunc
func CopyIntoFunc[T any](toCopy T, copyFunc matchers.CopyIntoFunc[T]) matchers.CopyInto {
	return matchers.NewCopyIntoFunc(toCopy, copyFunc)
}

// ShallowCopyInto is a convenience function to create a CopyInto matcher that copies the value from a pointer.
func ShallowCopyInto[T any](toCopy *T) matchers.CopyInto {
	return matchers.NewCopyIntoPointer(toCopy)
}
