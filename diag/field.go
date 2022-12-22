package pl_diag

const (
	Message    = Field("msg")
	MethodName = Field("method_name")
)

type Field string

func (f Field) String() string { return string(f) }
