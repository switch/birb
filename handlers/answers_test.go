package handlers

import (
	"errors"
	"reflect"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Test interfaces for validateReturnTypes testing
type (
	StringMethod         interface{ Method() string }
	IntMethod            interface{ Method() int }
	ErrorMethod          interface{ Method() error }
	StringErrorMethod    interface{ Method() (string, error) }
	IntStringErrorMethod interface{ Method() (int, string, error) }
	SliceMethod          interface{ Method() []string }
	MapMethod            interface{ Method() map[string]int }
	PointerMethod        interface{ Method() *string }
	ChanMethod           interface{ Method() chan int }
	FuncMethod           interface{ Method() func() }
	VoidMethod           interface{ Method() }
)

func getMethod[T any]() reflect.Method {
	var t T
	return reflect.TypeOf(&t).Elem().Method(0)
}

var _ = Describe("validateReturnTypes", func() {
	Describe("correct return types", func() {
		It("should accept correct string return", func() {
			method := getMethod[StringMethod]()
			err := validateReturnTypes(method, []any{"hello"})
			Expect(err).ToNot(HaveOccurred())
		})

		It("should accept correct int return", func() {
			method := getMethod[IntMethod]()
			err := validateReturnTypes(method, []any{42})
			Expect(err).ToNot(HaveOccurred())
		})

		It("should accept nil for error interface", func() {
			method := getMethod[ErrorMethod]()
			err := validateReturnTypes(method, []any{nil})
			Expect(err).ToNot(HaveOccurred())
		})

		It("should accept error value for error interface", func() {
			method := getMethod[ErrorMethod]()
			err := validateReturnTypes(method, []any{errors.New("test error")})
			Expect(err).ToNot(HaveOccurred())
		})

		It("should accept correct multiple return values", func() {
			method := getMethod[StringErrorMethod]()
			err := validateReturnTypes(method, []any{"result", nil})
			Expect(err).ToNot(HaveOccurred())
		})

		It("should accept correct multiple return values with error", func() {
			method := getMethod[StringErrorMethod]()
			err := validateReturnTypes(method, []any{"result", errors.New("oops")})
			Expect(err).ToNot(HaveOccurred())
		})

		It("should accept nil for slice types", func() {
			method := getMethod[SliceMethod]()
			err := validateReturnTypes(method, []any{nil})
			Expect(err).ToNot(HaveOccurred())
		})

		It("should accept actual slice for slice types", func() {
			method := getMethod[SliceMethod]()
			err := validateReturnTypes(method, []any{[]string{"a", "b"}})
			Expect(err).ToNot(HaveOccurred())
		})

		It("should accept nil for map types", func() {
			method := getMethod[MapMethod]()
			err := validateReturnTypes(method, []any{nil})
			Expect(err).ToNot(HaveOccurred())
		})

		It("should accept nil for pointer types", func() {
			method := getMethod[PointerMethod]()
			err := validateReturnTypes(method, []any{nil})
			Expect(err).ToNot(HaveOccurred())
		})

		It("should accept actual pointer for pointer types", func() {
			method := getMethod[PointerMethod]()
			s := "hello"
			err := validateReturnTypes(method, []any{&s})
			Expect(err).ToNot(HaveOccurred())
		})

		It("should accept nil for channel types", func() {
			method := getMethod[ChanMethod]()
			err := validateReturnTypes(method, []any{nil})
			Expect(err).ToNot(HaveOccurred())
		})

		It("should accept nil for function types", func() {
			method := getMethod[FuncMethod]()
			err := validateReturnTypes(method, []any{nil})
			Expect(err).ToNot(HaveOccurred())
		})

		It("should accept empty slice for void methods", func() {
			method := getMethod[VoidMethod]()
			err := validateReturnTypes(method, []any{})
			Expect(err).ToNot(HaveOccurred())
		})
	})

	// NOTE: Type validation has been intentionally relaxed.
	// We only validate return value COUNT, not types, because:
	// - Go's pattern of "return struct, accept interface" means mocks often
	//   return concrete types or mocks that implement interfaces
	// - The generated mock code does type assertions, so mismatches surface
	//   as zero values (not panics) which tests can still catch
	// - Being too strict prevents valid patterns like returning mocks for interfaces
	Describe("relaxed type validation", func() {
		It("should accept wrong type for string return (validation relaxed)", func() {
			method := getMethod[StringMethod]()
			err := validateReturnTypes(method, []any{42})
			Expect(err).ToNot(HaveOccurred()) // Type mismatch allowed
		})

		It("should accept nil for non-nullable types (validation relaxed)", func() {
			method := getMethod[IntMethod]()
			err := validateReturnTypes(method, []any{nil})
			Expect(err).ToNot(HaveOccurred()) // Nil allowed, will become zero value at runtime
		})
	})

	Describe("wrong number of return values", func() {
		It("should reject too few return values", func() {
			method := getMethod[StringErrorMethod]()
			err := validateReturnTypes(method, []any{"only one"})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("expected 2 return values, got 1"))
		})

		It("should reject too many return values", func() {
			method := getMethod[StringMethod]()
			err := validateReturnTypes(method, []any{"first", "second", "third"})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("expected 1 return values, got 3"))
		})

		It("should reject any return values for void methods", func() {
			method := getMethod[VoidMethod]()
			err := validateReturnTypes(method, []any{"unexpected"})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("expected 0 return values, got 1"))
		})
	})

	// Multiple return value type validation removed - only count is validated
})
