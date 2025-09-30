package mocks

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/switch/birb"
)

var _ = Describe("Mocks Event Recorder", func() {
	var SUT *MockEventRecorder

	BeforeEach(func() {
		SUT = NewMockEventRecorder(GinkgoT())
	})

	Context("Events", func() {
		It("should allow mocking of AnnotatedEventf", func() {
			WhenCalling(SUT.MOCKany_AnnotatedEventf())

			SUT.AnnotatedEventf(nil,
				nil,
				"Normal",
				"TestReason",
				"This is a test event",
				//"arg1", "arg2",
			)

			Verify(SUT, Once()).CALLED_AnnotatedEventf(WithAnyArgs())
		})
	})
})
