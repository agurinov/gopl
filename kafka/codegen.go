//go:build neverbuild

//go:generate -command stringer go run ../vendor/golang.org/x/tools/cmd/stringer
//go:generate -command mockgen go run ../vendor/github.com/golang/mock/mockgen
package kafka

//go:generate mockgen -source=lib_iface.go     -destination=mock/lib_iface.go           -package=mock -mock_names=ConsumerLibrary=ConsumerLibrary,ProducerLibrary=ProducerLibrary
//go:generate mockgen -source=event_handler.go -destination=mock/event_handler_iface.go -package=mock -mock_names=EventHandler=EventHandler,EventBatchHandler=EventBatchHandler
//go:generate stringer -type=EventHandleStrategy .
