<task>
	Investigate why record being fetched several times from broker.
</task>

<idea>
	The idea of new vision is behind the kafka feature of pause/resume mechanics. Kafka oficially documented this behaviour as normal to use.
	In fetch operation `PollTopicPartition` in @kafka/consumer_client.go each partition resumes only themself and got ONLY their partition records.
</idea>

<context>
	- Kafka
	- Golang
	- True concurrent consuming
	- Our vendor is franz-go/kgo
	- @kafka/ directory is the source code.
</context>

<scope>
	@kafka/ directory
</scope>

<problem>
	<initial>
		The problem is with the pattern "shared fetch" and channel of work: channels which send records to their goroutines per partition will blocks until the slowest partition frees the channel.
		Do not offer me that and similar solutions.
		Our topics is highly volatile, so we need true concurrency.
	</initial>

	<current>
		When partition changing shared client's state to "pause all and resume only me"
		vendor library flushes all buffered records for another partitions.
		As a result each "switch" leads to another network roundtrip again.
	</current>
</problem>

<constraints>
	- Don't suggest reverting to the initial implementation (shared fetch)
	- highly volatile topics, concurrency is matter
</constraints>

<output>
	Provide:
	1. Root cause hypothesis
	2. Files to investigate (in priority order)
	3. Suggest fixes with trade-offs. Even if it's in vendor code.
</output>
