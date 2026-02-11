package kafka

type RecordMapper[R any, V any] interface {
	FromVendor(V) R
	ToVendor(R) V
}
