package mocks_test

import (
	"sync"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/switch/birb"
	. "github.com/switch/birb/internal/mocks"
)

var _ = Describe("Concurrency", func() {
	Describe("Handler thread safety", func() {
		It("should handle concurrent method calls safely", func() {
			mockGreeter := NewMockGreeter(GinkgoT())
			WhenCalling(mockGreeter.MOCK_PersonalGreet(Anything())).
				ThenAnswer(func(args []any) []any {
					return []any{"Hello, " + args[0].(string), nil}
				})

			const numGoroutines = 100
			var wg sync.WaitGroup
			wg.Add(numGoroutines)

			for i := 0; i < numGoroutines; i++ {
				go func(n int) {
					defer wg.Done()
					defer GinkgoRecover()

					result, err := mockGreeter.PersonalGreet("User")
					Expect(err).ToNot(HaveOccurred())
					Expect(result).To(Equal("Hello, User"))
				}(i)
			}

			wg.Wait()

			// Verify all calls were recorded
			Verify(mockGreeter, Times(numGoroutines)).CALLED_PersonalGreet(Anything())
		})

		It("should handle concurrent calls with different argument values", func() {
			mockGreeter := NewMockGreeter(GinkgoT())

			// Set up stubs BEFORE concurrent calling (normal usage pattern)
			WhenCalling(mockGreeter.MOCK_PersonalGreet(Anything())).
				ThenAnswer(func(args []any) []any {
					return []any{"Hello, " + args[0].(string), nil}
				})

			const numGoroutines = 50
			var wg sync.WaitGroup
			wg.Add(numGoroutines)

			results := make([]string, numGoroutines)

			for i := 0; i < numGoroutines; i++ {
				go func(n int) {
					defer wg.Done()
					defer GinkgoRecover()

					result, err := mockGreeter.PersonalGreet("User" + string(rune('A'+n%26)))
					Expect(err).ToNot(HaveOccurred())
					results[n] = result
				}(i)
			}

			wg.Wait()

			// All results should be valid greetings
			for _, result := range results {
				Expect(result).To(HavePrefix("Hello, User"))
			}
		})

		It("should correctly count calls after concurrent execution", func() {
			mockGreeter := NewMockGreeter(GinkgoT())
			WhenCalling(mockGreeter.MOCK_Greet()).ThenReturn("hello")

			const numGoroutines = 100
			var wg sync.WaitGroup
			wg.Add(numGoroutines)

			// Make concurrent calls
			for i := 0; i < numGoroutines; i++ {
				go func() {
					defer wg.Done()
					defer GinkgoRecover()

					_ = mockGreeter.Greet()
				}()
			}

			wg.Wait()

			// Verify counts are accurate after concurrent execution
			// Note: Verify calls should be done sequentially (not concurrently)
			Verify(mockGreeter, Times(numGoroutines)).CALLED_Greet()
			Verify(mockGreeter, AtLeast(numGoroutines)).CALLED_Greet()
			Verify(mockGreeter, AtMost(numGoroutines)).CALLED_Greet()
		})
	})

	Describe("Captor thread safety", func() {
		It("should capture arguments safely from concurrent calls", func() {
			mockGreeter := NewMockGreeter(GinkgoT())
			captor := Captor()

			WhenCalling(mockGreeter.MOCK_PersonalGreet(captor)).
				ThenReturn("captured", nil)

			const numGoroutines = 100
			var wg sync.WaitGroup
			wg.Add(numGoroutines)

			for i := 0; i < numGoroutines; i++ {
				go func(n int) {
					defer wg.Done()
					defer GinkgoRecover()

					_, _ = mockGreeter.PersonalGreet("User")
				}(i)
			}

			wg.Wait()

			// Verify all arguments were captured
			values := captor.GetValues()
			Expect(values).To(HaveLen(numGoroutines))

			// All captured values should be "User"
			for _, v := range values {
				Expect(v).To(Equal("User"))
			}
		})

		It("should allow concurrent GetValues calls safely", func() {
			mockGreeter := NewMockGreeter(GinkgoT())
			captor := Captor()

			WhenCalling(mockGreeter.MOCK_PersonalGreet(captor)).
				ThenReturn("captured", nil)

			// Make some calls first
			for i := 0; i < 10; i++ {
				_, _ = mockGreeter.PersonalGreet("User")
			}

			const numGoroutines = 50
			var wg sync.WaitGroup
			wg.Add(numGoroutines)

			results := make([][]any, numGoroutines)

			for i := 0; i < numGoroutines; i++ {
				go func(n int) {
					defer wg.Done()
					defer GinkgoRecover()

					results[n] = captor.GetValues()
				}(i)
			}

			wg.Wait()

			// All results should have 10 values
			for _, result := range results {
				Expect(result).To(HaveLen(10))
			}
		})

		It("should handle concurrent capture and retrieval safely", func() {
			mockGreeter := NewMockGreeter(GinkgoT())
			captor := Captor()

			WhenCalling(mockGreeter.MOCK_PersonalGreet(captor)).
				ThenReturn("captured", nil)

			const numGoroutines = 50
			var wg sync.WaitGroup
			wg.Add(numGoroutines * 2)

			// Half capture, half retrieve
			for i := 0; i < numGoroutines; i++ {
				go func() {
					defer wg.Done()
					defer GinkgoRecover()

					_, _ = mockGreeter.PersonalGreet("Concurrent")
				}()

				go func() {
					defer wg.Done()
					defer GinkgoRecover()

					_ = captor.GetValues()
				}()
			}

			wg.Wait()

			// Final check - should have captured all calls
			finalValues := captor.GetValues()
			Expect(finalValues).To(HaveLen(numGoroutines))
		})
	})

	Describe("Freeze/Unfreeze thread safety", func() {
		It("should handle concurrent freeze operations safely", func() {
			mockGreeter := NewMockGreeter(GinkgoT())

			const numGoroutines = 50
			var wg sync.WaitGroup
			wg.Add(numGoroutines)

			for i := 0; i < numGoroutines; i++ {
				go func(n int) {
					defer wg.Done()
					defer GinkgoRecover()

					if n%2 == 0 {
						Freeze(mockGreeter)
					} else {
						mockGreeter.BirbHandler().Unfreeze()
					}
				}(i)
			}

			wg.Wait()
			// Test passes if no race conditions occurred
		})
	})
})
