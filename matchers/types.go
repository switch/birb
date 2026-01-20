package matchers

// Anything is a matcher that matches any single argument value.
// Use this when you don't care about a specific argument in a method call.
type Anything struct{}

func (a *Anything) Match(any) (bool, error) {
	return true, nil
}

// MatchTheRestOfTheArguments is a matcher that matches all remaining arguments.
// Use this as the last matcher when you want to ignore trailing variadic arguments.
type MatchTheRestOfTheArguments struct{}

func (a *MatchTheRestOfTheArguments) Match(any) (bool, error) {
	return true, nil
}
