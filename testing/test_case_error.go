package pl_testing

type TestCaseViolationError struct {
	Condition string
}

var (
	ErrViolationIs = TestCaseViolationError{Condition: "Is"}
	ErrViolationAs = TestCaseViolationError{Condition: "As"}
)

func (e TestCaseViolationError) Error() string {
	return "test case error violates " + e.Condition + " condition"
}
