package mocks

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/switch/birb"
)

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

		Verify(SutTestingT, Once()).CALLED_Fatalf(
			HavePrefix("Verifier: "),
			MatchError(ContainSubstring("expected num method calls:")),
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

		It("should balk on the number of returns without a 'WithAnyArgs()'", func() {
			// TODO: fix this so it won't panic.. it should give a better error
			WhenCalling(SutGreeter.MOCK_AllTheGreets()).ThenReturn(1, 2, 3)

			actual, err := SutGreeter.AllTheGreets("whee")
			Expect(actual).To(BeEmpty())
			Expect(err).ToNot(HaveOccurred())

			Verify(SutTestingT, Once()).CALLED_Fatalf(
				HavePrefix("Mock: matcher error: "),
				MatchError(ContainSubstring("expected at least one matcher")),
			)
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
})
