//go:build failing

package failing

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/switch/birb"
	"github.com/switch/birb/internal/mocks"
)

// These tests are DESIGNED TO FAIL to demonstrate the improved error messages.
// Run with: make failing
//
// Each test shows a different error scenario with the new improved formatting:
// - Method names included in all errors
// - Source location (file:line) for stub-related errors
// - Consistent prefixes (Mock: vs Verifier:)
// - No double-prefixing
// - Proper error propagation (no swallowed errors)

var _ = Describe("Error Message Examples", func() {
	var (
		mockGreeter *mocks.MockGreeter
	)

	BeforeEach(func() {
		mockGreeter = mocks.NewMockGreeter(GinkgoT())
	})

	Describe("Verification Failures", func() {
		It("call count mismatch - expected more calls", func() {
			// Call once, but verify twice was expected
			mockGreeter.Greet()

			// ERROR: Verifier: Greet: expected 2 call(s), got 1
			Verify(mockGreeter, Twice()).CALLED_Greet()
		})

		It("call count mismatch - expected fewer calls", func() {
			// Call twice, but verify once was expected
			mockGreeter.Greet()
			mockGreeter.Greet()

			// ERROR: Verifier: Greet: expected 1 call(s), got 2
			Verify(mockGreeter, Once()).CALLED_Greet()
		})

		It("call count mismatch - expected none", func() {
			// Call once, but verify never was expected
			mockGreeter.Greet()

			// ERROR: Verifier: Greet: expected 0 call(s), got 1
			Verify(mockGreeter, Never()).CALLED_Greet()
		})

		It("at least constraint not met", func() {
			// Call once, but verify at least 3 was expected
			mockGreeter.Greet()

			// ERROR: Verifier: Greet: expected at least 3 call(s), got 1
			Verify(mockGreeter, AtLeast(3)).CALLED_Greet()
		})

		It("at most constraint exceeded", func() {
			// Call 3 times, but verify at most 1 was expected
			mockGreeter.Greet()
			mockGreeter.Greet()
			mockGreeter.Greet()

			// ERROR: Verifier: Greet: expected at most 1 call(s), got 3
			Verify(mockGreeter, AtMost(1)).CALLED_Greet()
		})
	})

	Describe("Return Value Count Mismatch", func() {
		// NOTE: Type validation has been relaxed - we only validate return value COUNT.
		// Wrong types will result in zero values at runtime (via type assertion), not errors.
		// This allows common Go patterns like returning mocks for interfaces.

		It("wrong number of return values", func() {
			// Stub returns 3 values, but method expects 1
			// ERROR: Mock: Greet: return type mismatch: expected 1 return values, got 3
			//        (stub defined at error_examples_test.go:XX)
			WhenCalling(mockGreeter.MOCK_Greet()).ThenReturn("a", "b", "c")

			mockGreeter.Greet()
		})
	})

	Describe("Argument Matching Errors", func() {
		It("missing matchers for variadic method", func() {
			// No matchers provided for method with arguments
			// ERROR: Mock: AllTheGreets: argument matching error: expected at least one matcher
			//        (stub defined at error_examples_test.go:XX)
			WhenCalling(mockGreeter.MOCK_AllTheGreets()).ThenReturn("result", nil)

			_, _ = mockGreeter.AllTheGreets("arg1", "arg2")
		})

		It("argument value mismatch - panics with helpful message", func() {
			// Stub expects specific argument, but called with different value
			// NEW BEHAVIOR: Panics with detailed message showing registered stubs
			WhenCalling(mockGreeter.MOCK_PersonalGreet("Alice")).ThenReturn("Hello Alice!", nil)

			// Call with different name - stub won't match
			// PANIC: birb: no matching stub for Greeter.PersonalGreet
			//   Called with: arg[0]: "Bob" (string)
			//   Registered stubs: (Equal("Alice")) [UNLIMITED, used 0 times]
			//   No fallback defined. Consider adding MOCKfallback_PersonalGreet()
			_, _ = mockGreeter.PersonalGreet("Bob")
		})
	})

	Describe("Stub Exhaustion Errors", func() {
		It("exhausted stub with no fallback panics", func() {
			// Stub limited to 1 call, no fallback defined
			WhenCalling(mockGreeter.MOCK_Greet()).ThenReturn("hello").Once()

			// First call succeeds
			mockGreeter.Greet()

			// Second call panics because stub is exhausted and no fallback
			// PANIC: birb: no matching stub for Greeter.Greet
			//   Called with: (no arguments)
			//   Registered stubs: (no matchers) [EXHAUSTED after 1 call(s)]
			//   No fallback defined. Consider adding MOCKfallback_Greet()
			mockGreeter.Greet()
		})
	})

	Describe("VerifyAllMatchersCalled Errors", func() {
		It("unused stub fails verification", func() {
			// Set up stub but never call it
			WhenCalling(mockGreeter.MOCK_Greet()).ThenReturn("hello")

			// Don't call mockGreeter.Greet()

			// ERROR: VerifyAllMatchersCalled: the following stubs were never invoked:
			//   - Greet(no matchers) [UNLIMITED, used 0 times]
			//   Tip: Remove unused stubs or use MOCKfallback_ for optional stubs.
			VerifyAllMatchersCalled(mockGreeter)
		})
	})

	Describe("Unexpected Interactions", func() {
		// NOTE: These tests demonstrate correct error messages, but the stack trace
		// points to api.go instead of the test file. This is a known limitation
		// because the api.go wrapper functions can't call t.Helper() without a testing.T.

		It("unhandled method calls", func() {
			// Call methods without setting up stubs, then verify no other interactions
			mockGreeter.Greet()
			mockGreeter.Greet()

			// ERROR: Verifier: unexpected interactions:
			//          Greet: 2 unhandled invocation(s): [[] []]
			// NOTE: Stack trace points to api.go (known limitation)
			VerifyNoOtherInteractions(mockGreeter)
		})

		It("verify never called but was called", func() {
			// Call a method, then verify nothing was called
			mockGreeter.Greet()

			// ERROR: Verifier: expected no interactions:
			//          Greet: expected 0 call(s), got 1
			// NOTE: Stack trace points to api.go (known limitation)
			VerifyNeverCalled(mockGreeter)
		})
	})

	Describe("Frozen Mock Errors", func() {
		It("calling frozen mock", func() {
			// Freeze the mock, then try to call it
			Freeze(mockGreeter)

			// ERROR: Mock: Greet: handler is frozen, cannot invoke method
			mockGreeter.Greet()
		})
	})
})
