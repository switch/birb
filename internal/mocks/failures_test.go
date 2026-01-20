package mocks

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/switch/birb"
)

// errorMatcher is a custom matcher that always returns an error
type errorMatcher struct {
	err error
}

func (e *errorMatcher) Match(actual any) (bool, error) {
	return false, e.err
}

var _ = Describe("failures", func() {
	var SutTestingT *MockTestingT
	var SutGreeter *MockGreeter
	var SutConvoluted *MockConvoluted

	BeforeEach(func() {
		SutTestingT = NewMockTestingT(GinkgoT())
		SutGreeter = NewMockGreeter(SutTestingT)
		SutConvoluted = NewMockConvoluted(SutTestingT)
	})

	It("should balk on the number of calls", func() {
		actual := SutGreeter.Greet()

		Expect(actual).To(BeEmpty())

		Verify(SutGreeter, Twice()).CALLED_Greet()

		// Format: "Verifier: %s: %v\n%s" with args (methodName, err, callHistory)
		Verify(SutTestingT, Once()).CALLED_Fatalf(
			HavePrefix("Verifier: %s:"),
			Equal("Greet"),
			MatchError(ContainSubstring("expected 2 call(s), got 1")),
			ContainSubstring("Greet:"),
		)
	})

	Describe("number of returns mismatches", func() {
		It("should balk on the number of returns", func() {
			WhenCalling(SutGreeter.MOCK_Greet()).ThenReturn(1, 2, 3)

			actual := SutGreeter.Greet()
			Expect(actual).To(BeEmpty())

			Verify(SutTestingT, Once()).CALLED_Fatalf(
				HavePrefix("Expected %d outputs,"),
				1,
				3,
			)
		})

		It("should balk on the number of returns", func(ctx SpecContext) {
			WhenCalling(SutGreeter.MOCKany_GreetWithContext()).ThenReturn(1, 2, 3)

			actual := SutGreeter.GreetWithContext(ctx)
			Expect(actual).To(BeEmpty())

			Verify(SutTestingT, Once()).CALLED_Fatalf(
				HavePrefix("Expected %d outputs,"),
				1,
				3,
			)
		})

		It("should balk on the number of returns", func() {
			WhenCalling(SutGreeter.MOCK_AllTheGreets(WithAnyArgs())).ThenReturn(1, 2, 3)

			actual, err := SutGreeter.AllTheGreets("whee")
			Expect(actual).To(BeEmpty())
			Expect(err).ToNot(HaveOccurred())

			Verify(SutTestingT, Once()).CALLED_Fatalf(
				HavePrefix("Expected %d outputs,"),
				2,
				3,
			)
		})

		It("should panic when stub has no matchers but method has parameters", func() {
			// When a stub is registered without matchers but the method expects arguments,
			// the stub can never match. The new behavior panics with a helpful message.
			WhenCalling(SutGreeter.MOCK_AllTheGreets()).ThenReturn(1, 2, 3)

			Expect(func() {
				_, _ = SutGreeter.AllTheGreets("whee")
			}).To(PanicWith(ContainSubstring("no matching stub for Greeter.AllTheGreets")))
		})

		It("should balk on the number of returns", func() {
			WhenCalling(SutGreeter.MOCKany_PersonalGreet()).ThenReturn(1, 2, 3)

			actual, err := SutGreeter.PersonalGreet("whee")
			Expect(actual).To(BeEmpty())
			Expect(err).ToNot(HaveOccurred())

			Verify(SutTestingT, Once()).CALLED_Fatalf(
				HavePrefix("Expected %d outputs,"),
				2,
				3,
			)
		})

		It("should balk on the number of returns", func() {
			WhenCalling(SutConvoluted.MOCK_Convolute(Anything(), Anything(), Anything(), Anything())).ThenReturn(1, 2, 3)

			actual, err := SutConvoluted.Convolute(SutGreeter, "whee", 42, map[string]string{"foo": "bar"})
			Expect(actual).To(BeEmpty())
			Expect(err).ToNot(HaveOccurred())

			Verify(SutTestingT, Once()).CALLED_Fatalf(
				HavePrefix("Expected %d outputs,"),
				2,
				3,
			)
		})
	})

	Describe("frozen mock", func() {
		It("should fail when calling a frozen mock", func() {
			WhenCalling(SutGreeter.MOCK_Greet()).ThenReturn("hello")

			// Freeze the mock
			Freeze(SutGreeter)

			// Call triggers Fatalf but continues (mock TestingT doesn't panic like real testing.T)
			_ = SutGreeter.Greet()

			// Verify Fatalf was called with frozen error message
			// Fatalf is called with (format, args...) where format is "Mock: %s: handler is frozen, cannot invoke method"
			Verify(SutTestingT, Once()).CALLED_Fatalf(
				Equal("Mock: %s: handler is frozen, cannot invoke method"),
				Equal("Greet"),
			)
		})

		It("should work again after unfreezing", func() {
			WhenCalling(SutGreeter.MOCK_Greet()).ThenReturn("hello")

			Freeze(SutGreeter)
			SutGreeter.BirbHandler().Unfreeze()

			actual := SutGreeter.Greet()
			Expect(actual).To(Equal("hello"))

			// No Fatalf should have been called
			Verify(SutTestingT, Never()).CALLED_Fatalf(WithAnyArgs())
		})
	})

	Describe("CopyInto type mismatches", func() {
		var mockSomeBirb *MockSomeBirb
		var mockContext context.Context

		BeforeEach(func() {
			mockSomeBirb = NewMockSomeBirb(SutTestingT)
			mockContext = context.Background()
		})

		It("should fail when ShallowCopyInto source is nil", func() {
			var nilOwl *Owl = nil
			WhenCalling(mockSomeBirb.MOCK_Dive(Anything(), ShallowCopyInto(nilOwl))).ThenReturn(nil)

			newOwl := &Owl{}
			err := mockSomeBirb.Dive(mockContext, newOwl)

			// Should succeed (nil source means nothing to copy)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("matcher errors during verification", func() {
		It("should fail when matcher returns an error during verification", func() {
			WhenCalling(SutGreeter.MOCK_PersonalGreet(Anything())).ThenReturn("hello", nil)

			// Make a call
			_, _ = SutGreeter.PersonalGreet("test")

			// Verify with a matcher that returns an error
			errMatcher := &errorMatcher{err: errors.New("custom matcher error")}
			Verify(SutGreeter, Once()).CALLED_PersonalGreet(errMatcher)

			// Should have called Fatalf with the matcher error
			// Format: "Verifier: matcher error during verification of %s: %v"
			Verify(SutTestingT, AtLeast(1)).CALLED_Fatalf(
				Equal("Verifier: matcher error during verification of %s: %v"),
				Equal("PersonalGreet"),
				MatchError(ContainSubstring("custom matcher error")),
			)
		})
	})

	// NOTE: Return type validation has been intentionally relaxed.
	// We only validate return value COUNT, not types, because:
	// - Go's pattern of "return struct, accept interface" means mocks often
	//   return concrete types or mocks that implement interfaces
	// - The generated mock code does type assertions, so mismatches surface
	//   as zero values (not panics) which tests can still catch
	// - Being too strict prevents valid patterns like returning mocks for interfaces
})
