package handlers

import (
	"reflect"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Test interfaces for defaultReturnValues testing
type (
	ReturnsBool       interface{ Method() bool }
	ReturnsInt        interface{ Method() int }
	ReturnsInt8       interface{ Method() int8 }
	ReturnsInt16      interface{ Method() int16 }
	ReturnsInt32      interface{ Method() int32 }
	ReturnsInt64      interface{ Method() int64 }
	ReturnsUint       interface{ Method() uint }
	ReturnsUint8      interface{ Method() uint8 }
	ReturnsUint16     interface{ Method() uint16 }
	ReturnsUint32     interface{ Method() uint32 }
	ReturnsUint64     interface{ Method() uint64 }
	ReturnsFloat32    interface{ Method() float32 }
	ReturnsFloat64    interface{ Method() float64 }
	ReturnsString     interface{ Method() string }
	ReturnsSlice      interface{ Method() []string }
	ReturnsIntSlice   interface{ Method() []int }
	ReturnsArray      interface{ Method() [3]int }
	ReturnsInterface  interface{ Method() error }
	ReturnsMultiple   interface{ Method() (string, error) }
	ReturnsNothing    interface{ Method() }
	ReturnsMap        interface{ Method() map[string]int }
	ReturnsPointer    interface{ Method() *string }
	ReturnsChan       interface{ Method() chan int }
	ReturnsStruct     interface{ Method() struct{ Name string } }
	ReturnsFunc       interface{ Method() func() }
)

func getMethodType[T any]() reflect.Type {
	var t T
	return reflect.TypeOf(&t).Elem().Method(0).Type
}

var _ = Describe("defaultReturnValues", func() {
	Describe("primitive types", func() {
		It("should return false for bool", func() {
			methodType := getMethodType[ReturnsBool]()
			result := defaultReturnValues(methodType)
			Expect(result).To(HaveLen(1))
			Expect(result[0]).To(Equal(false))
		})

		It("should return 0 for int types", func() {
			Expect(defaultReturnValues(getMethodType[ReturnsInt]())[0]).To(Equal(0))
			Expect(defaultReturnValues(getMethodType[ReturnsInt8]())[0]).To(Equal(0))
			Expect(defaultReturnValues(getMethodType[ReturnsInt16]())[0]).To(Equal(0))
			Expect(defaultReturnValues(getMethodType[ReturnsInt32]())[0]).To(Equal(0))
			Expect(defaultReturnValues(getMethodType[ReturnsInt64]())[0]).To(Equal(0))
		})

		It("should return 0 for uint types", func() {
			Expect(defaultReturnValues(getMethodType[ReturnsUint]())[0]).To(Equal(0))
			Expect(defaultReturnValues(getMethodType[ReturnsUint8]())[0]).To(Equal(0))
			Expect(defaultReturnValues(getMethodType[ReturnsUint16]())[0]).To(Equal(0))
			Expect(defaultReturnValues(getMethodType[ReturnsUint32]())[0]).To(Equal(0))
			Expect(defaultReturnValues(getMethodType[ReturnsUint64]())[0]).To(Equal(0))
		})

		It("should return 0.0 for float types", func() {
			Expect(defaultReturnValues(getMethodType[ReturnsFloat32]())[0]).To(Equal(0.0))
			Expect(defaultReturnValues(getMethodType[ReturnsFloat64]())[0]).To(Equal(0.0))
		})

		It("should return empty string for string", func() {
			result := defaultReturnValues(getMethodType[ReturnsString]())
			Expect(result).To(HaveLen(1))
			Expect(result[0]).To(Equal(""))
		})
	})

	Describe("composite types", func() {
		It("should return empty slice for slice types", func() {
			result := defaultReturnValues(getMethodType[ReturnsSlice]())
			Expect(result).To(HaveLen(1))
			Expect(result[0]).To(BeAssignableToTypeOf([]string{}))
			Expect(result[0]).To(HaveLen(0))
		})

		It("should return empty int slice for []int", func() {
			result := defaultReturnValues(getMethodType[ReturnsIntSlice]())
			Expect(result).To(HaveLen(1))
			Expect(result[0]).To(BeAssignableToTypeOf([]int{}))
			Expect(result[0]).To(HaveLen(0))
		})

		It("should return nil for interface types", func() {
			result := defaultReturnValues(getMethodType[ReturnsInterface]())
			Expect(result).To(HaveLen(1))
			Expect(result[0]).To(BeNil())
		})

		It("should return empty map for map return types", func() {
			result := defaultReturnValues(getMethodType[ReturnsMap]())
			Expect(result).To(HaveLen(1))
			Expect(result[0]).To(BeAssignableToTypeOf(map[string]int{}))
			Expect(result[0]).To(HaveLen(0))
		})

		It("should return nil for pointer return types", func() {
			result := defaultReturnValues(getMethodType[ReturnsPointer]())
			Expect(result).To(HaveLen(1))
			Expect(result[0]).To(BeNil())
		})

		It("should return nil for channel return types", func() {
			result := defaultReturnValues(getMethodType[ReturnsChan]())
			Expect(result).To(HaveLen(1))
			Expect(result[0]).To(BeNil())
		})

		It("should return zero-valued struct for struct return types", func() {
			result := defaultReturnValues(getMethodType[ReturnsStruct]())
			Expect(result).To(HaveLen(1))
			// Should be a zero-valued struct, not nil
			Expect(result[0]).ToNot(BeNil())
			s, ok := result[0].(struct{ Name string })
			Expect(ok).To(BeTrue())
			Expect(s.Name).To(Equal(""))
		})

		It("should return nil for function return types", func() {
			result := defaultReturnValues(getMethodType[ReturnsFunc]())
			Expect(result).To(HaveLen(1))
			Expect(result[0]).To(BeNil())
		})
	})

	Describe("multiple return values", func() {
		It("should return defaults for all return values", func() {
			result := defaultReturnValues(getMethodType[ReturnsMultiple]())
			Expect(result).To(HaveLen(2))
			Expect(result[0]).To(Equal(""))    // string default
			Expect(result[1]).To(BeNil())      // error (interface) default
		})
	})

	Describe("no return values", func() {
		It("should return empty slice for void methods", func() {
			result := defaultReturnValues(getMethodType[ReturnsNothing]())
			Expect(result).To(HaveLen(0))
		})
	})
})
