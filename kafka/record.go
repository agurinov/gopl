package kafka

import "github.com/twmb/franz-go/pkg/kgo"

type Record struct{}

func recordFromKgo(*kgo.Record) Record {
	return Record{}
}
