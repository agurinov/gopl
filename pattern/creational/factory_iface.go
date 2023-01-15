package creational

type Factory[O Object] interface {
	NewObject() (O, error)
	MustNewObject() O
}
