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

	Context("CallMatcher.Invoke", func() {
		It("should return error when matchers exceed args count", func() {
			captor := NewCaptor()
			callMatcher := CreateCallMatcher([]Matcher{captor, captor, captor})

			// Only 2 args but 3 matchers
			err := callMatcher.Invoke([]any{"arg1", "arg2"})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("out of bounds"))
		})

		It("should capture all args when counts match", func() {
			captor := NewCaptor()
			callMatcher := CreateCallMatcher([]Matcher{captor})

			err := callMatcher.Invoke([]any{"captured"})
			Expect(err).ToNot(HaveOccurred())
			Expect(captor.GetValues()).To(Equal([]any{"captured"}))
		})

		It("should handle empty matchers and args", func() {
			callMatcher := CreateCallMatcher([]Matcher{})

			err := callMatcher.Invoke([]any{})
			Expect(err).ToNot(HaveOccurred())
		})

		It("should work when args exceed matchers count", func() {
			captor := NewCaptor()
			callMatcher := CreateCallMatcher([]Matcher{captor})

			// 3 args but only 1 matcher - should only capture first
			err := callMatcher.Invoke([]any{"first", "second", "third"})
			Expect(err).ToNot(HaveOccurred())
			Expect(captor.GetValues()).To(Equal([]any{"first"}))
		})
	})

	Context("CallMatcher.Equivalent", func() {
		It("should return true for empty matchers", func() {
			a := CreateCallMatcher([]Matcher{})
			b := CreateCallMatcher([]Matcher{})
			Expect(a.Equivalent(b)).To(BeTrue())
		})

		It("should return false for different lengths", func() {
			a := CreateCallMatcher([]Matcher{&Anything{}})
			b := CreateCallMatcher([]Matcher{&Anything{}, &Anything{}})
			Expect(a.Equivalent(b)).To(BeFalse())
		})

		It("should return true for matching Anything matchers", func() {
			a := CreateCallMatcher([]Matcher{&Anything{}, &Anything{}})
			b := CreateCallMatcher([]Matcher{&Anything{}, &Anything{}})
			Expect(a.Equivalent(b)).To(BeTrue())
		})

		It("should return true for matching MatchTheRestOfTheArguments matchers", func() {
			a := CreateCallMatcher([]Matcher{&MatchTheRestOfTheArguments{}})
			b := CreateCallMatcher([]Matcher{&MatchTheRestOfTheArguments{}})
			Expect(a.Equivalent(b)).To(BeTrue())
		})

		It("should return false for different matcher types", func() {
			a := CreateCallMatcher([]Matcher{&Anything{}})
			b := CreateCallMatcher([]Matcher{&MatchTheRestOfTheArguments{}})
			Expect(a.Equivalent(b)).To(BeFalse())
		})

		It("should return true for Equal matchers with same value", func() {
			a := CreateCallMatcher([]Matcher{CreateMatcherFromValue("hello")})
			b := CreateCallMatcher([]Matcher{CreateMatcherFromValue("hello")})
			Expect(a.Equivalent(b)).To(BeTrue())
		})

		It("should return false for Equal matchers with different values", func() {
			a := CreateCallMatcher([]Matcher{CreateMatcherFromValue("hello")})
			b := CreateCallMatcher([]Matcher{CreateMatcherFromValue("world")})
			Expect(a.Equivalent(b)).To(BeFalse())
		})

		It("should handle nil CallMatcher", func() {
			a := CreateCallMatcher([]Matcher{&Anything{}})
			Expect(a.Equivalent(nil)).To(BeFalse())

			var nilMatcher *CallMatcher
			Expect(nilMatcher.Equivalent(nil)).To(BeTrue())
			Expect(nilMatcher.Equivalent(a)).To(BeFalse())
		})

		It("should return true for same instance matchers", func() {
			anything := &Anything{}
			a := CreateCallMatcher([]Matcher{anything})
			b := CreateCallMatcher([]Matcher{anything})
			Expect(a.Equivalent(b)).To(BeTrue())
		})

		It("should return false for different custom matchers", func() {
			captor1 := NewCaptor()
			captor2 := NewCaptor()
			a := CreateCallMatcher([]Matcher{captor1})
			b := CreateCallMatcher([]Matcher{captor2})
			// Different captor instances should not be equivalent
			Expect(a.Equivalent(b)).To(BeFalse())
		})
	})
})
