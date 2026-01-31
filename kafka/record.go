package kafka

import (
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

type (
	Record       = *kgo.Record
	RecordMapper interface {
		FromVendor(*kgo.Record) Record
		ToVendor(Record) *kgo.Record
	}
	kgoRecordMapper struct{}
)

func (kgoRecordMapper) FromVendor(v *kgo.Record) Record {
	return v
}

func (kgoRecordMapper) ToVendor(r Record) *kgo.Record {
	return r
}

func RecordLogFields(record Record) []zap.Field {
	return []zap.Field{
		zap.String("record.topic", record.Topic),
		zap.Int32("record.partition", record.Partition),
		zap.Int64("record.offset", record.Offset),
		zap.ByteString("record.key", record.Key),
	}
}
