package metrics

const (
	NopanicHandlerCounterName = "gopl_nopanic_recovered_count"
)

const (
	HTTPServerHandlerDurationHistogramName = "gopl_http_server_handler_duration_seconds"
)

const (
	KafkaConsumerHandlerDurationHistogramName      = "gopl_kafka_consumer_handler_duration_seconds"
	KafkaConsumerHandlerBatchDurationHistogramName = "gopl_kafka_consumer_handler_batch_duration_seconds"
)
