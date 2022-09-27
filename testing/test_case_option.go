package pl_testing

// TODO(a.gurinov): Autogen fmt.Stringer interface for alias type
type TestCaseOption = byte

const (
	TESTING_NO_DOTENV_FILE TestCaseOption = 1 << iota
	TESTING_NO_PARALLEL
)
