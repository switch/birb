package mocks

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/switch/birb"
)

var _ = Describe("Specificity and Exhaustion Features", func() {
	var SutGreeter *MockGreeter

	BeforeEach(func() {
		SutGreeter = NewMockGreeter(GinkgoT())
	})

	Describe("Specificity-based matching", func() {
		It("specific matchers should win over Anything() regardless of registration order", func() {
			// Register catch-all FIRST
			WhenCalling(SutGreeter.MOCKany_PersonalGreet()).ThenReturn("default", nil)
			// Register specific matcher SECOND
			WhenCalling(SutGreeter.MOCK_PersonalGreet(Equal("Bob"))).ThenReturn("hello Bob", nil)

			// Specific should win even though it was registered after
			result, err := SutGreeter.PersonalGreet("Bob")
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("hello Bob"))

			// Non-matching should fall through to catch-all
			result, err = SutGreeter.PersonalGreet("Alice")
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("default"))
		})

		It("Equal matchers should have higher specificity than partial matchers", func() {
			// ContainSubstring is a partial matcher
			WhenCalling(SutGreeter.MOCK_PersonalGreet(ContainSubstring("ob"))).ThenReturn("partial match", nil)
			// Equal is more specific
			WhenCalling(SutGreeter.MOCK_PersonalGreet(Equal("Bob"))).ThenReturn("exact match", nil)

			// Exact match should win
			result, _ := SutGreeter.PersonalGreet("Bob")
			Expect(result).To(Equal("exact match"))

			// Other strings containing "ob" should use partial
			result, _ = SutGreeter.PersonalGreet("Jacob")
			Expect(result).To(Equal("partial match"))
		})
	})

	Describe(".Times(n) and .Once()", func() {
		It("should limit stub to exactly n invocations", func() {
			WhenCalling(SutGreeter.MOCK_Greet()).ThenReturn("limited").Times(2)
			WhenCalling(SutGreeter.MOCKfallback_Greet()).ThenReturn("fallback")

			// First two calls use the limited stub
			Expect(SutGreeter.Greet()).To(Equal("limited"))
			Expect(SutGreeter.Greet()).To(Equal("limited"))

			// Third call falls through to fallback
			Expect(SutGreeter.Greet()).To(Equal("fallback"))
		})

		It(".Once() should be equivalent to .Times(1)", func() {
			WhenCalling(SutGreeter.MOCK_Greet()).ThenReturn("once only").Once()
			WhenCalling(SutGreeter.MOCKfallback_Greet()).ThenReturn("fallback")

			Expect(SutGreeter.Greet()).To(Equal("once only"))
			Expect(SutGreeter.Greet()).To(Equal("fallback"))
			Expect(SutGreeter.Greet()).To(Equal("fallback"))
		})

		It("should panic when Times(0) is called", func() {
			Expect(func() {
				WhenCalling(SutGreeter.MOCK_Greet()).ThenReturn("bad").Times(0)
			}).To(PanicWith(ContainSubstring("Times(n) requires n > 0")))
		})

		It("unlimited stubs should repeat last value forever (default behavior)", func() {
			WhenCalling(SutGreeter.MOCK_Greet()).ThenReturn("a").ThenReturn("b")

			Expect(SutGreeter.Greet()).To(Equal("a"))
			Expect(SutGreeter.Greet()).To(Equal("b"))
			Expect(SutGreeter.Greet()).To(Equal("b")) // repeats
			Expect(SutGreeter.Greet()).To(Equal("b")) // still repeats
		})
	})

	Describe("MOCKfallback_ stubs", func() {
		It("should have lowest priority", func() {
			WhenCalling(SutGreeter.MOCKfallback_PersonalGreet()).ThenReturn("fallback", nil)
			WhenCalling(SutGreeter.MOCKany_PersonalGreet()).ThenReturn("any", nil)
			WhenCalling(SutGreeter.MOCK_PersonalGreet(Equal("Bob"))).ThenReturn("specific", nil)

			// Specific wins
			result, _ := SutGreeter.PersonalGreet("Bob")
			Expect(result).To(Equal("specific"))

			// Any wins over fallback
			result, _ = SutGreeter.PersonalGreet("Alice")
			Expect(result).To(Equal("any"))
		})

		It("fallback should be used when all other stubs are exhausted", func() {
			WhenCalling(SutGreeter.MOCK_Greet()).ThenReturn("limited").Times(1)
			WhenCalling(SutGreeter.MOCKfallback_Greet()).ThenReturn("fallback")

			Expect(SutGreeter.Greet()).To(Equal("limited"))
			Expect(SutGreeter.Greet()).To(Equal("fallback"))
		})
	})

	Describe("Exhaustion panic", func() {
		It("should panic with helpful message when no stub matches and no fallback exists", func() {
			WhenCalling(SutGreeter.MOCK_Greet()).ThenReturn("limited").Times(1)
			// Note: no fallback defined

			Expect(SutGreeter.Greet()).To(Equal("limited"))

			Expect(func() {
				SutGreeter.Greet()
			}).To(PanicWith(And(
				ContainSubstring("no matching stub for Greeter.Greet"),
				ContainSubstring("EXHAUSTED"),
				ContainSubstring("MOCKfallback_"),
			)))
		})
	})

	Describe("VerifyAllMatchersCalled", func() {
		It("should pass when all non-fallback stubs are invoked", func() {
			WhenCalling(SutGreeter.MOCK_Greet()).ThenReturn("hello")
			WhenCalling(SutGreeter.MOCKfallback_PersonalGreet()).ThenReturn("fallback", nil)

			SutGreeter.Greet() // invoke the stub

			// Should not panic - MOCK_Greet was invoked, fallback is excluded
			VerifyAllMatchersCalled(SutGreeter)
		})

		It("should fail when a non-fallback stub was never invoked", func() {
			// Create a mock with a separate TestingT so we can capture the failure
			mockT := NewMockTestingT(GinkgoT())
			greeter := NewMockGreeter(mockT)

			// Set up the mock to capture Fatalf and Helper
			WhenCalling(mockT.MOCKany_Fatalf()).ThenReturn()
			WhenCalling(mockT.MOCK_Helper()).ThenReturn()

			WhenCalling(greeter.MOCK_Greet()).ThenReturn("never called")

			// Don't call Greet - stub should be considered "dead"
			VerifyAllMatchersCalled(greeter)

			// Verify that Fatalf was called with appropriate message
			// Note: Fatalf(format string, args ...any) - we pass a single string containing the full message
			Verify(mockT, Once()).CALLED_Fatalf(
				And(
					ContainSubstring("VerifyAllMatchersCalled"),
					ContainSubstring("never invoked"),
				),
			)
		})

		It("should exclude fallback stubs from verification", func() {
			WhenCalling(SutGreeter.MOCK_Greet()).ThenReturn("hello")
			WhenCalling(SutGreeter.MOCKfallback_PersonalGreet()).ThenReturn("never used", nil)

			SutGreeter.Greet()
			// Don't call PersonalGreet

			// Should pass - fallback is excluded from verification
			VerifyAllMatchersCalled(SutGreeter)
		})
	})

	Describe("Complex scenarios", func() {
		It("should handle BeforeEach/test pattern correctly", func() {
			// Simulate BeforeEach - set up defaults
			WhenCalling(SutGreeter.MOCKany_PersonalGreet()).ThenReturn("default", nil)

			// Simulate test-specific override
			WhenCalling(SutGreeter.MOCK_PersonalGreet(Equal("special"))).ThenReturn("special result", nil).Once()

			// First call to "special" returns override
			result, _ := SutGreeter.PersonalGreet("special")
			Expect(result).To(Equal("special result"))

			// Second call to "special" falls through to default (override exhausted)
			result, _ = SutGreeter.PersonalGreet("special")
			Expect(result).To(Equal("default"))

			// Regular calls always use default
			result, _ = SutGreeter.PersonalGreet("regular")
			Expect(result).To(Equal("default"))
		})

		It("should respect Reset() for exhaustion state", func() {
			WhenCalling(SutGreeter.MOCK_Greet()).ThenReturn("limited").Times(1)
			WhenCalling(SutGreeter.MOCKfallback_Greet()).ThenReturn("fallback")

			Expect(SutGreeter.Greet()).To(Equal("limited"))
			Expect(SutGreeter.Greet()).To(Equal("fallback"))

			// Reset should clear exhaustion state
			Reset(SutGreeter)

			// Re-setup stubs
			WhenCalling(SutGreeter.MOCK_Greet()).ThenReturn("limited again").Times(1)
			WhenCalling(SutGreeter.MOCKfallback_Greet()).ThenReturn("fallback again")

			Expect(SutGreeter.Greet()).To(Equal("limited again"))
			Expect(SutGreeter.Greet()).To(Equal("fallback again"))
		})
	})
})
