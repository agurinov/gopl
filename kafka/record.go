package kafka

import "github.com/twmb/franz-go/pkg/kgo"

type RecordMapper[V any, R any] interface {
	FromVendor(V) R
	ToVendor(R) V
}

type (
	Record struct{}
)

type kgoRecordMapper struct{}

func (kgoRecordMapper) FromVendor(*kgo.Record) Record {
	return Record{}
}

func (kgoRecordMapper) ToVendor(Record) *kgo.Record {
	return new(kgo.Record)
}
