package matchers

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("matcher", func() {
	Context("CreateMatcherFromValue", func() {
		It("should create a matcher that matches the given value", func() {
			matcher := CreateMatcherFromValue(27)

			actual, err := matcher.Match(27)
			Expect(actual).To(BeTrue())
			Expect(err).ToNot(HaveOccurred())

			actual, err = matcher.Match(43)
			Expect(actual).To(BeFalse())
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
