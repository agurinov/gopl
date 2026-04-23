//nolint:lll
package kafka

const (
	KafkaConsumerHandlerDurationHistogramName      = "gopl_kafka_consumer_handler_duration_seconds"
	KafkaConsumerHandlerBatchDurationHistogramName = "gopl_kafka_consumer_handler_batch_duration_seconds"
	KafkaNotMyRecordCounterName                    = "gopl_kafka_consumer_not_my_record_cnt"
	KafkaConsumerIdleDurationHistogramName         = "gopl_kafka_consumer_idle_duration_seconds"
	KafkaConsumerDiscardedCounterName              = "gopl_kafka_consumer_discarded_cnt"
	KafkaConsumerCommittedCounterName              = "gopl_kafka_consumer_committed_cnt"
	KafkaConsumerPollingDurationHistogramName      = "gopl_kafka_consumer_polling_duration_seconds"
)
