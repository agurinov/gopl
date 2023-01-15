package pl_factory

type Factory[O object] interface {
	NewObject() (O, error)
	MustNewObject() O
}
