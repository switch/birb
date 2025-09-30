package mocks

import (
	"strconv"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/switch/birb"
)

var _ = Describe("Any Calls", func() {
	var mockGreeter *MockGreeter
	var convoluted *MockConvoluted

	BeforeEach(func() {
		mockGreeter = NewMockGreeter(GinkgoT())
		convoluted = NewMockConvoluted(GinkgoT())
	})

	AfterEach(func() {
		VerifyNoOtherInteractions(mockGreeter)
		VerifyNoOtherInteractions(convoluted)
	})

	Describe("Any calls shenanigans", func() {
		It("should allow mocking of any calls", func() {
			WhenCalling(mockGreeter.MOCK_Greet()).ThenReturn("awesome")
			WhenCalling(convoluted.MOCKany_Convolute()).ThenAnswer(func(args []any) []any {
				retMap := args[3].(map[string]string)
				retMap["integer"] = strconv.Itoa(args[2].(int))
				retMap[args[1].(string)] = args[0].(Greeter).Greet()
				return []any{retMap, nil}
			})

			actual, err := convoluted.Convolute(mockGreeter, "something", 27, map[string]string{
				"here": "there",
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal(map[string]string{
				"integer":   "27",
				"something": "awesome",
				"here":      "there",
			}))
			Verify(mockGreeter, Once()).CALLED_Greet()
			Verify(convoluted, Once()).CALLED_Convolute(WithAnyArgs())
		})
	})

	Describe("Variadic shenanigans", func() {
		It("should be easy to validate a function was called with all arguments", func() {
			WhenCalling(mockGreeter.MOCK_AllTheGreets(WithAnyArgs())).ThenAnswer(func(args []any) []any {
				allArgs := make([]string, len(args))
				for i, arg := range args {
					allArgs[i] = arg.(string)
				}
				return []any{[]string{strings.Join(allArgs, " ")}, nil}
			})
			actual, err := mockGreeter.AllTheGreets("something", "awesome", "happened")

			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal([]string{"something awesome happened"}))
			Verify(mockGreeter, Once()).CALLED_AllTheGreets("something", "awesome", "happened")
		})

		It("should be easy to validate a function was called with any arguments", func() {
			WhenCalling(mockGreeter.MOCK_AllTheGreets(WithAnyArgs())).ThenAnswer(func(args []any) []any {
				allArgs := make([]string, len(args))
				for i, arg := range args {
					allArgs[i] = arg.(string)
				}
				return []any{[]string{strings.Join(allArgs, " ")}, nil}
			})
			actual, err := mockGreeter.AllTheGreets("something", "awesome", "happened")

			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal([]string{"something awesome happened"}))
			Verify(mockGreeter, Once()).CALLED_AllTheGreets(WithAnyArgs())
		})

		It("should be easy to validate a function was called with partial arguments defined then any", func() {
			WhenCalling(mockGreeter.MOCK_AllTheGreets(WithAnyArgs())).ThenAnswer(func(args []any) []any {
				allArgs := make([]string, len(args))
				for i, arg := range args {
					allArgs[i] = arg.(string)
				}
				return []any{[]string{strings.Join(allArgs, " ")}, nil}
			})
			actual, err := mockGreeter.AllTheGreets("something", "awesome", "happened")

			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal([]string{"something awesome happened"}))
			Verify(mockGreeter, Once()).CALLED_AllTheGreets("something", WithAnyArgs())
		})

		It("should be easy to validate a function was called with partial arguments defined then any", func() {
			WhenCalling(mockGreeter.MOCK_AllTheGreets(WithAnyArgs())).ThenAnswer(func(args []any) []any {
				allArgs := make([]string, len(args))
				for i, arg := range args {
					allArgs[i] = arg.(string)
				}
				return []any{[]string{strings.Join(allArgs, " ")}, nil}
			})
			actual, err := mockGreeter.AllTheGreets("something", "awesome", "happened")

			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal([]string{"something awesome happened"}))
			Verify(mockGreeter, Once()).CALLED_AllTheGreets("something", "awesome", WithAnyArgs())
		})
	})
})
