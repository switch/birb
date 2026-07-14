package handlers

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Call Verifiers", func() {
	Describe("Times", func() {
		It("should pass when call count matches exactly", func() {
			verifier := Times(3)
			err := verifier.Verify(&callData{MethodCalls: 3})
			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail when call count is less than expected", func() {
			verifier := Times(3)
			err := verifier.Verify(&callData{MethodCalls: 2})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("expected 3 call(s), got 2"))
		})

		It("should fail when call count is more than expected", func() {
			verifier := Times(3)
			err := verifier.Verify(&callData{MethodCalls: 5})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("expected 3 call(s), got 5"))
		})

		It("should work with zero expected calls", func() {
			verifier := Times(0)
			err := verifier.Verify(&callData{MethodCalls: 0})
			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail when expecting zero but got calls", func() {
			verifier := Times(0)
			err := verifier.Verify(&callData{MethodCalls: 1})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("expected 0 call(s), got 1"))
		})
	})

	Describe("AtLeast", func() {
		It("should pass when call count equals minimum", func() {
			verifier := AtLeast(3)
			err := verifier.Verify(&callData{MethodCalls: 3})
			Expect(err).ToNot(HaveOccurred())
		})

		It("should pass when call count exceeds minimum", func() {
			verifier := AtLeast(3)
			err := verifier.Verify(&callData{MethodCalls: 5})
			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail when call count is below minimum", func() {
			verifier := AtLeast(3)
			err := verifier.Verify(&callData{MethodCalls: 2})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("expected at least 3 call(s), got 2"))
		})

		It("should pass with AtLeast(0) regardless of calls", func() {
			verifier := AtLeast(0)
			Expect(verifier.Verify(&callData{MethodCalls: 0})).ToNot(HaveOccurred())
			Expect(verifier.Verify(&callData{MethodCalls: 100})).ToNot(HaveOccurred())
		})

		It("should fail with AtLeast(1) when zero calls", func() {
			verifier := AtLeast(1)
			err := verifier.Verify(&callData{MethodCalls: 0})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("expected at least 1 call(s), got 0"))
		})
	})

	Describe("AtMost", func() {
		It("should pass when call count equals maximum", func() {
			verifier := AtMost(3)
			err := verifier.Verify(&callData{MethodCalls: 3})
			Expect(err).ToNot(HaveOccurred())
		})

		It("should pass when call count is below maximum", func() {
			verifier := AtMost(3)
			err := verifier.Verify(&callData{MethodCalls: 1})
			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail when call count exceeds maximum", func() {
			verifier := AtMost(3)
			err := verifier.Verify(&callData{MethodCalls: 4})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("expected at most 3 call(s), got 4"))
		})

		It("should pass with AtMost(0) when zero calls", func() {
			verifier := AtMost(0)
			err := verifier.Verify(&callData{MethodCalls: 0})
			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail with AtMost(0) when any calls", func() {
			verifier := AtMost(0)
			err := verifier.Verify(&callData{MethodCalls: 1})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("expected at most 0 call(s), got 1"))
		})
	})

	Describe("Between", func() {
		It("should pass when call count equals minimum", func() {
			verifier := Between(2, 5)
			err := verifier.Verify(&callData{MethodCalls: 2})
			Expect(err).ToNot(HaveOccurred())
		})

		It("should pass when call count equals maximum", func() {
			verifier := Between(2, 5)
			err := verifier.Verify(&callData{MethodCalls: 5})
			Expect(err).ToNot(HaveOccurred())
		})

		It("should pass when call count is between min and max", func() {
			verifier := Between(2, 5)
			err := verifier.Verify(&callData{MethodCalls: 3})
			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail when call count is below minimum", func() {
			verifier := Between(2, 5)
			err := verifier.Verify(&callData{MethodCalls: 1})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("expected between 2 and 5 call(s), got 1"))
		})

		It("should fail when call count exceeds maximum", func() {
			verifier := Between(2, 5)
			err := verifier.Verify(&callData{MethodCalls: 6})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("expected between 2 and 5 call(s), got 6"))
		})

		It("should work with same min and max (equivalent to Times)", func() {
			verifier := Between(3, 3)
			Expect(verifier.Verify(&callData{MethodCalls: 3})).ToNot(HaveOccurred())
			Expect(verifier.Verify(&callData{MethodCalls: 2})).To(HaveOccurred())
			Expect(verifier.Verify(&callData{MethodCalls: 4})).To(HaveOccurred())
		})

		It("should work with Between(0, n) allowing zero calls", func() {
			verifier := Between(0, 3)
			Expect(verifier.Verify(&callData{MethodCalls: 0})).ToNot(HaveOccurred())
			Expect(verifier.Verify(&callData{MethodCalls: 3})).ToNot(HaveOccurred())
			Expect(verifier.Verify(&callData{MethodCalls: 4})).To(HaveOccurred())
		})
	})

	Describe("MethodVerifierFromCallData", func() {
		It("should allow custom verification logic", func() {
			// Custom verifier that requires an even number of calls
			verifier := MethodVerifierFromCallData(func(data *callData) error {
				if data.MethodCalls%2 != 0 {
					return Errorf("expected even number of calls, got %d", data.MethodCalls)
				}
				return nil
			})

			Expect(verifier.Verify(&callData{MethodCalls: 0})).ToNot(HaveOccurred())
			Expect(verifier.Verify(&callData{MethodCalls: 2})).ToNot(HaveOccurred())
			Expect(verifier.Verify(&callData{MethodCalls: 4})).ToNot(HaveOccurred())
			Expect(verifier.Verify(&callData{MethodCalls: 1})).To(HaveOccurred())
			Expect(verifier.Verify(&callData{MethodCalls: 3})).To(HaveOccurred())
		})
	})

	Describe("callCollection", func() {
		It("should create an empty call collection", func() {
			cc := newCallCollection("TestMethod")
			Expect(cc.MethodName).To(Equal("TestMethod"))
			Expect(cc.Calls).To(BeEmpty())
		})

		It("should format String() output correctly", func() {
			cc := newCallCollection("TestMethod")
			cc.Calls = append(cc.Calls, newCall("TestMethod", []any{"arg1", 42}))
			cc.Calls = append(cc.Calls, newCall("TestMethod", []any{"arg2", 100}))

			str := cc.String()
			Expect(str).To(ContainSubstring("TestMethod:"))
			Expect(str).To(ContainSubstring("arg1"))
			Expect(str).To(ContainSubstring("42"))
		})
	})

	Describe("call", func() {
		It("should clone arguments when creating a call", func() {
			originalArgs := []any{"original", 42}
			c := newCall("TestMethod", originalArgs)

			// Modify original args
			originalArgs[0] = "modified"

			// Call should have the original value
			Expect(c.Args[0]).To(Equal("original"))
		})

		It("should format String() output correctly", func() {
			c := newCall("TestMethod", []any{"hello", 123})
			Expect(c.String()).To(Equal("[hello 123]"))
		})
	})
})

// Errorf is a helper to create formatted errors for tests
func Errorf(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}
