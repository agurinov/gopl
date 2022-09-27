package pl_testing

type TestCaseErrorViolation struct {
	Condition string
}

var (
	ErrViolationIs = TestCaseErrorViolation{Condition: "Is"}
	ErrViolationAs = TestCaseErrorViolation{Condition: "As"}
)

func (e TestCaseErrorViolation) Error() string {
	return "test case error violates " + e.Condition + " condition"
}
