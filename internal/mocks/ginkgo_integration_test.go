package mocks

import (
	"context"
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/switch/birb"
)

var _ = Describe("ginkgo integration - Greeter", func() {
	var mockGreeter *MockGreeter

	BeforeEach(func() {
		mockGreeter = NewMockGreeter(GinkgoT())
	})

	Describe("Default returns", func() {
		It("should return empty string", func() {
			actual := mockGreeter.Greet()

			Expect(actual).To(BeEmpty())
			Expect(actual).To(Equal(""))
			Verify(mockGreeter, Once()).CALLED_Greet()
		})

		It("should return empty array", func() {
			actual, err := mockGreeter.AllTheGreets()

			Expect(err).NotTo(HaveOccurred())
			Expect(actual).To(BeEmpty())
			Expect(actual).To(HaveLen(0))
			Expect(actual).To(Equal([]string{}))
			Verify(mockGreeter, Once()).CALLED_AllTheGreets()
		})
	})

	Describe("All Functions verified", func() {
		AfterEach(func() {
			VerifyNoOtherInteractions(mockGreeter)
		})

		Describe("Greet", func() {
			It("should verify a mock hasn't been touched", func() {
				VerifyNeverCalled(mockGreeter)
				Freeze(mockGreeter)
			})

			It("should run a simple test", func() {
				WhenCalling(mockGreeter.MOCK_Greet()).ThenReturn("Hello, John")

				actual := mockGreeter.Greet()

				Expect(actual).To(Equal("Hello, John"))

				Verify(mockGreeter, Once()).CALLED_Greet()
			})
		})

		Describe("PersonalGreet", func() {
			It("should run a simple test", func() {
				WhenCalling(mockGreeter.MOCK_PersonalGreet("John")).ThenReturn("Hello, John", nil)

				actual, err := mockGreeter.PersonalGreet("John")

				Expect(actual).To(Equal("Hello, John"))
				Expect(err).NotTo(HaveOccurred())

				Verify(mockGreeter, Once()).CALLED_PersonalGreet("John")
			})

			It("should work with multiple returns", func() {
				WhenCalling(mockGreeter.MOCK_PersonalGreet(Anything())).
					ThenReturn("Hello, Foo", nil).
					ThenReturn("Hello, Bar", nil).
					ThenReturn("Hello, Baz", nil).
					ThenReturn("Hello, Qux", nil)

				actual, err := mockGreeter.PersonalGreet("John 1")
				Expect(actual).To(Equal("Hello, Foo"))
				Expect(err).NotTo(HaveOccurred())

				actual, err = mockGreeter.PersonalGreet("John 2")
				Expect(actual).To(Equal("Hello, Bar"))
				Expect(err).NotTo(HaveOccurred())

				actual, err = mockGreeter.PersonalGreet("John 3")
				Expect(actual).To(Equal("Hello, Baz"))
				Expect(err).NotTo(HaveOccurred())

				actual, err = mockGreeter.PersonalGreet("qux")
				Expect(actual).To(Equal("Hello, Qux"))
				Expect(err).NotTo(HaveOccurred())

				actual, err = mockGreeter.PersonalGreet("qux")
				Expect(actual).To(Equal("Hello, Qux"))
				Expect(err).NotTo(HaveOccurred())

				Verify(mockGreeter, Once()).CALLED_PersonalGreet("John 1")
				Verify(mockGreeter, Once()).CALLED_PersonalGreet("John 2")
				Verify(mockGreeter, Once()).CALLED_PersonalGreet("John 3")
				Verify(mockGreeter, Twice()).CALLED_PersonalGreet("qux")
				Verify(mockGreeter, AtLeast(1)).CALLED_PersonalGreet("qux")
				Verify(mockGreeter, AtMost(3)).CALLED_PersonalGreet("qux")

				Verify(mockGreeter, Never()).CALLED_PersonalGreet("awesome-o-dne")
			})

			It("should work with multiple single returns", func() {
				WhenCalling(mockGreeter.MOCK_PersonalGreet("John 1")).
					ThenReturn("Hello, Foo", nil)
				WhenCalling(mockGreeter.MOCK_PersonalGreet("John 2")).
					ThenReturn("Hello, Bar", nil)
				WhenCalling(mockGreeter.MOCK_PersonalGreet("John 3")).
					ThenReturn("Hello, Baz", nil)
				WhenCalling(mockGreeter.MOCK_PersonalGreet("qux")).
					ThenReturn("Hello, Qux", nil)

				actual, err := mockGreeter.PersonalGreet("John 1")
				Expect(actual).To(Equal("Hello, Foo"))
				Expect(err).NotTo(HaveOccurred())

				actual, err = mockGreeter.PersonalGreet("John 2")
				Expect(actual).To(Equal("Hello, Bar"))
				Expect(err).NotTo(HaveOccurred())

				actual, err = mockGreeter.PersonalGreet("John 3")
				Expect(actual).To(Equal("Hello, Baz"))
				Expect(err).NotTo(HaveOccurred())

				actual, err = mockGreeter.PersonalGreet("qux")
				Expect(actual).To(Equal("Hello, Qux"))
				Expect(err).NotTo(HaveOccurred())

				actual, err = mockGreeter.PersonalGreet("qux")
				Expect(actual).To(Equal("Hello, Qux"))
				Expect(err).NotTo(HaveOccurred())

				Verify(mockGreeter, Once()).CALLED_PersonalGreet("John 1")
				Verify(mockGreeter, Once()).CALLED_PersonalGreet("John 2")
				Verify(mockGreeter, Once()).CALLED_PersonalGreet("John 3")
				Verify(mockGreeter, Twice()).CALLED_PersonalGreet("qux")

				Verify(mockGreeter, Never()).CALLED_PersonalGreet("awesome-o-dne")
			})

			It("should work with function returns", func() {
				WhenCalling(mockGreeter.MOCK_PersonalGreet(Anything())).
					ThenAnswer(func(args []any) []any {
						return []any{"Hello, " + args[0].(string), nil}
					})

				actual, err := mockGreeter.PersonalGreet("Foo")
				Expect(actual).To(Equal("Hello, Foo"))
				Expect(err).NotTo(HaveOccurred())

				actual, err = mockGreeter.PersonalGreet("Bar")
				Expect(actual).To(Equal("Hello, Bar"))
				Expect(err).NotTo(HaveOccurred())

				actual, err = mockGreeter.PersonalGreet("Baz")
				Expect(actual).To(Equal("Hello, Baz"))
				Expect(err).NotTo(HaveOccurred())

				actual, err = mockGreeter.PersonalGreet("qux")
				Expect(actual).To(Equal("Hello, qux"))
				Expect(err).NotTo(HaveOccurred())

				actual, err = mockGreeter.PersonalGreet("qux")
				Expect(actual).To(Equal("Hello, qux"))
				Expect(err).NotTo(HaveOccurred())

				Verify(mockGreeter, Once()).CALLED_PersonalGreet("Foo")
				Verify(mockGreeter, Once()).CALLED_PersonalGreet("Bar")
				Verify(mockGreeter, Once()).CALLED_PersonalGreet("Baz")
				Verify(mockGreeter, Twice()).CALLED_PersonalGreet("qux")

				Verify(mockGreeter, Never()).CALLED_PersonalGreet("awesome-o-dne")
			})

			It("should work with function returns", func() {
				captor := Captor()
				WhenCalling(mockGreeter.MOCK_PersonalGreet(captor)).
					ThenAnswer(func(args []any) []any {
						return []any{"Hello, " + args[0].(string), nil}
					})

				actual, err := mockGreeter.PersonalGreet("Foo")
				Expect(actual).To(Equal("Hello, Foo"))
				Expect(err).NotTo(HaveOccurred())

				actual, err = mockGreeter.PersonalGreet("Bar")
				Expect(actual).To(Equal("Hello, Bar"))
				Expect(err).NotTo(HaveOccurred())

				actual, err = mockGreeter.PersonalGreet("Baz")
				Expect(actual).To(Equal("Hello, Baz"))
				Expect(err).NotTo(HaveOccurred())

				actual, err = mockGreeter.PersonalGreet("qux")
				Expect(actual).To(Equal("Hello, qux"))
				Expect(err).NotTo(HaveOccurred())

				actual, err = mockGreeter.PersonalGreet("qux")
				Expect(actual).To(Equal("Hello, qux"))
				Expect(err).NotTo(HaveOccurred())

				Verify(mockGreeter, Once()).CALLED_PersonalGreet("Foo")
				Verify(mockGreeter, Once()).CALLED_PersonalGreet("Bar")
				Verify(mockGreeter, Once()).CALLED_PersonalGreet("Baz")
				Verify(mockGreeter, Twice()).CALLED_PersonalGreet("qux")

				Expect(captor.GetValues()).To(ConsistOf("Foo", "Bar", "Baz", "qux", "qux"))

				Verify(mockGreeter, Never()).CALLED_PersonalGreet("awesome-o-dne")
			})

			It("should work with context matching anything", func() {
				WhenCalling(mockGreeter.MOCK_GreetWithContext(Anything())).ThenReturn("howdy")

				actual := mockGreeter.GreetWithContext(context.Background())

				Expect(actual).To(Equal("howdy"))
				Verify(mockGreeter, Once()).CALLED_GreetWithContext(Anything())
			})

			It("should run a simple error test", func() {
				WhenCalling(mockGreeter.MOCK_PersonalGreet("John")).ThenReturn("", errors.New("oh no!"))

				_, err := mockGreeter.PersonalGreet("John")

				Expect(err).To(HaveOccurred())

				Verify(mockGreeter, Once()).CALLED_PersonalGreet("John")
				Verify(mockGreeter, Never()).CALLED_PersonalGreet("awesome-o-dne")
			})
		})
	})
})

var _ = Describe("Reset", func() {
	var mockGreeter *MockGreeter

	BeforeEach(func() {
		mockGreeter = NewMockGreeter(GinkgoT())
	})

	It("should clear stubs after reset", func() {
		// Set up a stub
		WhenCalling(mockGreeter.MOCK_Greet()).ThenReturn("stubbed")
		Expect(mockGreeter.Greet()).To(Equal("stubbed"))

		// Reset the mock
		Reset(mockGreeter)

		// Stub should be gone - returns default (empty string)
		Expect(mockGreeter.Greet()).To(Equal(""))
	})

	It("should clear call history after reset", func() {
		// Make some calls
		WhenCalling(mockGreeter.MOCK_Greet()).ThenReturn("hello")
		mockGreeter.Greet()
		mockGreeter.Greet()
		mockGreeter.Greet()

		// Verify 3 calls
		Verify(mockGreeter, Times(3)).CALLED_Greet()

		// Reset the mock
		Reset(mockGreeter)

		// Call history should be cleared
		Verify(mockGreeter, Never()).CALLED_Greet()
	})

	It("should allow new stubs after reset", func() {
		// Set up initial stub
		WhenCalling(mockGreeter.MOCK_Greet()).ThenReturn("first")
		Expect(mockGreeter.Greet()).To(Equal("first"))

		// Reset and set up new stub
		Reset(mockGreeter)
		WhenCalling(mockGreeter.MOCK_Greet()).ThenReturn("second")
		Expect(mockGreeter.Greet()).To(Equal("second"))
	})

	It("should clear frozen state after reset", func() {
		WhenCalling(mockGreeter.MOCK_Greet()).ThenReturn("hello")

		// Freeze the mock
		Freeze(mockGreeter)

		// Reset should unfreeze
		Reset(mockGreeter)

		// Should be able to call without error (no Fatalf)
		WhenCalling(mockGreeter.MOCK_Greet()).ThenReturn("after reset")
		Expect(mockGreeter.Greet()).To(Equal("after reset"))
	})
})

var _ = Describe("Stub Merging", func() {
	var mockGreeter *MockGreeter

	BeforeEach(func() {
		mockGreeter = NewMockGreeter(GinkgoT())
	})

	Describe("multiple WhenCalling with same matchers", func() {
		It("should merge stubs with no arguments", func() {
			// Two separate WhenCalling calls should act like chained ThenReturns
			WhenCalling(mockGreeter.MOCK_Greet()).ThenReturn("first")
			WhenCalling(mockGreeter.MOCK_Greet()).ThenReturn("second")

			// Should get sequential values, then repeat last
			Expect(mockGreeter.Greet()).To(Equal("first"))
			Expect(mockGreeter.Greet()).To(Equal("second"))
			Expect(mockGreeter.Greet()).To(Equal("second"))
		})

		It("should merge stubs with Anything() matcher", func() {
			WhenCalling(mockGreeter.MOCK_PersonalGreet(Anything())).ThenReturn("first", nil)
			WhenCalling(mockGreeter.MOCK_PersonalGreet(Anything())).ThenReturn("second", nil)
			WhenCalling(mockGreeter.MOCK_PersonalGreet(Anything())).ThenReturn("third", nil)

			result, _ := mockGreeter.PersonalGreet("anyone")
			Expect(result).To(Equal("first"))

			result, _ = mockGreeter.PersonalGreet("someone")
			Expect(result).To(Equal("second"))

			result, _ = mockGreeter.PersonalGreet("else")
			Expect(result).To(Equal("third"))

			// Should repeat last value
			result, _ = mockGreeter.PersonalGreet("more")
			Expect(result).To(Equal("third"))
		})

		It("should merge stubs with same literal argument", func() {
			WhenCalling(mockGreeter.MOCK_PersonalGreet("John")).ThenReturn("Hello John 1", nil)
			WhenCalling(mockGreeter.MOCK_PersonalGreet("John")).ThenReturn("Hello John 2", nil)

			result, _ := mockGreeter.PersonalGreet("John")
			Expect(result).To(Equal("Hello John 1"))

			result, _ = mockGreeter.PersonalGreet("John")
			Expect(result).To(Equal("Hello John 2"))

			result, _ = mockGreeter.PersonalGreet("John")
			Expect(result).To(Equal("Hello John 2"))
		})

		It("should NOT merge stubs with different literal arguments", func() {
			WhenCalling(mockGreeter.MOCK_PersonalGreet("John")).ThenReturn("Hello John", nil)
			WhenCalling(mockGreeter.MOCK_PersonalGreet("Jane")).ThenReturn("Hello Jane", nil)

			// Different arguments should use different stubs
			result, _ := mockGreeter.PersonalGreet("John")
			Expect(result).To(Equal("Hello John"))

			result, _ = mockGreeter.PersonalGreet("Jane")
			Expect(result).To(Equal("Hello Jane"))

			// Each stub keeps returning its last value
			result, _ = mockGreeter.PersonalGreet("John")
			Expect(result).To(Equal("Hello John"))

			result, _ = mockGreeter.PersonalGreet("Jane")
			Expect(result).To(Equal("Hello Jane"))
		})
	})

	Describe("combining chained and separate WhenCalling", func() {
		It("should work with mix of chained and separate calls", func() {
			// First two via chaining
			WhenCalling(mockGreeter.MOCK_Greet()).
				ThenReturn("first").
				ThenReturn("second")
			// Third via separate WhenCalling (should merge)
			WhenCalling(mockGreeter.MOCK_Greet()).ThenReturn("third")

			Expect(mockGreeter.Greet()).To(Equal("first"))
			Expect(mockGreeter.Greet()).To(Equal("second"))
			Expect(mockGreeter.Greet()).To(Equal("third"))
			Expect(mockGreeter.Greet()).To(Equal("third"))
		})
	})
})
