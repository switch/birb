package mocks

import (
	"context"
	"errors"
	"reflect"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/switch/birb"
)

var _ = Describe("copy into", func() {
	var mockSomeBirb *MockSomeBirb
	var mockContext context.Context
	var owl *Owl

	BeforeEach(func() {
		mockSomeBirb = NewMockSomeBirb(GinkgoT())
		mockContext = context.Background()
		owl = &Owl{
			SomeBirbImpl: SomeBirbImpl{
				DaMap: map[string]SomeBirb{
					"falcon": &Falcon{
						SomeBirbImpl: SomeBirbImpl{
							DaMap: map[string]SomeBirb{},
						},
					},
				},
			},
		}
	})

	Describe("duck type?", func() {
		It("should quack", func() {
			WhenCalling(mockSomeBirb.MOCK_Talk()).ThenReturn("quack")
			actual := mockSomeBirb.Talk()

			Expect(actual).To(Equal("quack"))
		})
	})

	Describe("Copy into", func() {
		It("should copy the wanted object into the argument - reference", func() {
			newOwl := &Owl{}
			WhenCalling(mockSomeBirb.MOCKany_Dive()).ThenAnswer(copyIntoAnswer(owl))
			err := mockSomeBirb.Dive(mockContext, newOwl)

			Expect(err).ToNot(HaveOccurred())
			Expect(newOwl).To(Equal(owl))
			Expect(newOwl.DaMap).To(HaveKey("falcon"))
		})

		It("should copy the wanted object into the argument - pointer only", func() {
			newOwl := &Owl{}
			WhenCalling(mockSomeBirb.MOCK_Dive(Anything(), ShallowCopyInto(owl))).ThenReturn(nil)
			err := mockSomeBirb.Dive(mockContext, newOwl)

			Expect(err).ToNot(HaveOccurred())
			Expect(newOwl).To(Equal(owl))
			Expect(newOwl.DaMap).To(HaveKey("falcon"))
		})

		It("should copy the wanted object into the argument - CopyIntoFunc", func() {
			newOwl := &Owl{}
			WhenCalling(mockSomeBirb.MOCK_Dive(Anything(), CopyIntoFunc(owl, copyIntoHelper))).ThenReturn(nil)
			err := mockSomeBirb.Dive(mockContext, newOwl)

			Expect(err).ToNot(HaveOccurred())
			Expect(newOwl).To(Equal(owl))
			Expect(newOwl.DaMap).To(HaveKey("falcon"))
		})

		It("should copy the wanted object into the argument - DeepCopyInto", func() {
			newOwl := &Owl{}
			WhenCalling(mockSomeBirb.MOCK_Dive(Anything(), DeepCopyInto(owl))).ThenReturn(nil)
			err := mockSomeBirb.Dive(mockContext, newOwl)

			Expect(err).ToNot(HaveOccurred())
			Expect(newOwl).To(Equal(owl))
			Expect(newOwl.DaMap).To(HaveKey("falcon"))
		})
	})
})

// this is nasty... use DeepCopyInto way if at all possible :)
func copyIntoHelper[T SomeBirb](destination, source T) (err error) {
	dstReflect := reflect.ValueOf(destination)
	srcReflect := reflect.ValueOf(source)
	switch dst := dstReflect.Interface().(type) {
	case *Owl:
		src := srcReflect.Interface().(*Owl)
		someBirb := SomeBirbImpl{
			DaMap: make(map[string]SomeBirb),
		}
		for key, val := range src.DaMap {
			switch val.(type) {
			case *Owl:
				someBirb.DaMap[key] = &Owl{}
			case *Falcon:
				someBirb.DaMap[key] = &Falcon{}
			}
			err = copyIntoHelper(someBirb.DaMap[key], val)
		}
		dst.SomeBirbImpl = someBirb
	case *Falcon:
		src := srcReflect.Interface().(*Falcon)
		someBirb := SomeBirbImpl{
			DaMap: make(map[string]SomeBirb),
		}
		for key, val := range src.DaMap {
			switch val.(type) {
			case *Owl:
				someBirb.DaMap[key] = &Owl{}
			case *Falcon:
				someBirb.DaMap[key] = &Falcon{}
			}
			err = copyIntoHelper(someBirb.DaMap[key], val)
		}
		dst.SomeBirbImpl = someBirb
	default:
		err = errors.New("unknown type")
	}
	return
}

func copyIntoAnswer(source SomeBirb) func(args []any) []any {
	return func(args []any) []any {
		return []any{copyIntoHelper(args[1].(SomeBirb), source)}
	}
}
