package types

// TestingT is the minimal testing interface required by birb mocks.
// Both *testing.T and GinkgoT() satisfy this interface.
type TestingT interface {
	Cleanup(func())
	Fatalf(format string, args ...any)
	Helper()
}
