package matchers

type Anything struct{}

func (a *Anything) Match(any) (bool, error) {
	return true, nil
}

type MatchTheRestOfTheArguments struct{}

func (a *MatchTheRestOfTheArguments) Match(any) (bool, error) {
	return true, nil
}
