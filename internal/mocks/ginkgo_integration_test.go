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
