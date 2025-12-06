package failing

import (
	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/ginkgo/v2/types"
	. "github.com/onsi/gomega"
	. "github.com/switch/birb"
	"github.com/switch/birb/internal/mocks"
)

var _ = Describe("Test Files", func() {
	var (
		mockGreeter *mocks.MockGreeter
		SutTestingT *mocks.MockTestingT
	)

	BeforeEach(func() {
		SutTestingT = mocks.NewMockTestingT(GinkgoT())
		mockGreeter = mocks.NewMockGreeter(SutTestingT)
	})

	Describe("failed validation file name should be the test file", func() {

		// The problem that this is trying to validate is that when a verification fails, the file name
		// in the stack trace is the mockery generated file, or the handler functions, and not the test file.

		// This is a problem because it makes it very difficult to track down the test that is failing
		// inside of th test suite you're working in.

		It("should not fail please", func() {
			defer func() {
				// this does not work.. the test has already been marked as failing by this time
				r := recover()
				Expect(r).To(HaveOccurred())
				ginkgoErr := r.(types.GinkgoError)
				Expect(ginkgoErr.CodeLocation.FileName).To(ContainSubstring("verify_any__awesomeness__fils_calls_test.go"))
			}()

			// this was never called, so it will fail.
			// it's registered with the mocked testingT, so it does NOT panic
			Verify(mockGreeter, Once()).CALLED_Greet()

			// this should fail, but huh...
			Verify(SutTestingT, Never()).CALLED_Fatalf(WithAnyArgs())
		})
	})
})
