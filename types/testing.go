package types

type TestingT interface {
	Cleanup(func())
	Fatalf(format string, args ...any)
	Helper()
}
