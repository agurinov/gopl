package stack

type Interface[T any] interface {
	Peek()
	Pop()
	Push()
}
