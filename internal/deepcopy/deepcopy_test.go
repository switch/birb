package deepcopy_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/switch/birb/internal/deepcopy"
)

func TestDeepCopy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DeepCopy Suite")
}

var _ = Describe("DeepCopyInto", func() {
	Describe("TestConfig", func() {
		var original *deepcopy.TestConfig

		BeforeEach(func() {
			original = &deepcopy.TestConfig{
				Name:     "test-config",
				Replicas: 3,
				Labels: map[string]string{
					"app":  "myapp",
					"tier": "backend",
				},
				Ports: []int{8080, 8443, 9090},
				Nested: &deepcopy.NestedConfig{
					Host:    "localhost",
					Port:    5432,
					Options: []string{"sslmode=disable", "timezone=UTC"},
				},
			}
		})

		It("should copy all scalar fields", func() {
			out := &deepcopy.TestConfig{}
			original.DeepCopyInto(out)

			Expect(out.Name).To(Equal("test-config"))
			Expect(out.Replicas).To(Equal(3))
		})

		It("should deep copy maps so modifications don't affect original", func() {
			out := &deepcopy.TestConfig{}
			original.DeepCopyInto(out)

			// Modify the copy's map
			out.Labels["new-key"] = "new-value"
			delete(out.Labels, "app")

			// Original should be unchanged
			Expect(original.Labels).To(HaveLen(2))
			Expect(original.Labels["app"]).To(Equal("myapp"))
			Expect(original.Labels).NotTo(HaveKey("new-key"))
		})

		It("should deep copy slices so modifications don't affect original", func() {
			out := &deepcopy.TestConfig{}
			original.DeepCopyInto(out)

			// Modify the copy's slice
			out.Ports[0] = 9999
			out.Ports = append(out.Ports, 1234)

			// Original should be unchanged
			Expect(original.Ports).To(Equal([]int{8080, 8443, 9090}))
		})

		It("should deep copy nested structs", func() {
			out := &deepcopy.TestConfig{}
			original.DeepCopyInto(out)

			// Modify the copy's nested struct
			out.Nested.Host = "modified-host"
			out.Nested.Port = 1234
			out.Nested.Options[0] = "modified-option"

			// Original nested struct should be unchanged
			Expect(original.Nested.Host).To(Equal("localhost"))
			Expect(original.Nested.Port).To(Equal(5432))
			Expect(original.Nested.Options[0]).To(Equal("sslmode=disable"))
		})

		It("should handle nil nested struct", func() {
			original.Nested = nil
			out := &deepcopy.TestConfig{}
			original.DeepCopyInto(out)

			Expect(out.Nested).To(BeNil())
		})

		It("should handle nil map", func() {
			original.Labels = nil
			out := &deepcopy.TestConfig{}
			original.DeepCopyInto(out)

			Expect(out.Labels).To(BeNil())
		})

		It("should handle nil slice", func() {
			original.Ports = nil
			out := &deepcopy.TestConfig{}
			original.DeepCopyInto(out)

			Expect(out.Ports).To(BeNil())
		})
	})

	Describe("DeepCopy convenience method", func() {
		It("should return a new deep copied instance", func() {
			original := &deepcopy.TestConfig{
				Name:     "original",
				Replicas: 5,
				Labels:   map[string]string{"key": "value"},
			}

			copied := original.DeepCopy()

			Expect(copied).NotTo(BeIdenticalTo(original))
			Expect(copied.Name).To(Equal("original"))
			Expect(copied.Replicas).To(Equal(5))

			// Modify copy, original unchanged
			copied.Labels["key"] = "modified"
			Expect(original.Labels["key"]).To(Equal("value"))
		})

		It("should return nil for nil receiver", func() {
			var nilConfig *deepcopy.TestConfig
			copied := nilConfig.DeepCopy()
			Expect(copied).To(BeNil())
		})
	})

	Describe("NestedConfig", func() {
		It("should deep copy Options slice", func() {
			original := &deepcopy.NestedConfig{
				Host:    "example.com",
				Port:    443,
				Options: []string{"opt1", "opt2"},
			}

			out := &deepcopy.NestedConfig{}
			original.DeepCopyInto(out)

			// Modify copy
			out.Options[0] = "modified"

			// Original unchanged
			Expect(original.Options[0]).To(Equal("opt1"))
		})

		It("should handle nil Options", func() {
			original := &deepcopy.NestedConfig{
				Host:    "example.com",
				Port:    443,
				Options: nil,
			}

			out := &deepcopy.NestedConfig{}
			original.DeepCopyInto(out)

			Expect(out.Options).To(BeNil())
		})
	})
})
