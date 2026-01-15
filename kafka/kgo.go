package kafka

import "github.com/twmb/franz-go/pkg/kgo"

type (
	kgoConsumer[R any]     consumer[R, *kgo.Record]
	kgoRecordMapper[R any] struct{}
)

func (kgoRecordMapper[R]) FromVendor(*kgo.Record) R {
	var r R

	return r
}

func (kgoRecordMapper[R]) ToVendor(R) *kgo.Record {
	return new(kgo.Record)
}
