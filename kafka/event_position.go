package kafka

import "fmt"

const (
	UknownPartition int32 = -1
)

var EmptyEventPosition = EventPosition{}

type EventPosition struct {
	Topic     string
	Partition int32
	Offset    int64
}

func (ep EventPosition) ValidateWith(configured EventPosition) error {
	if ep.Topic != configured.Topic {
		return fmt.Errorf(
			"%w: unexpected topic; %q instead of %q",
			ErrInvalidEventPosition,
			ep.Topic,
			configured.Topic,
		)
	}

	if configured.Partition == UknownPartition {
		return nil
	}

	if ep.Partition != configured.Partition {
		return fmt.Errorf(
			"%w: unexpected partition; %d instead of %d",
			ErrInvalidEventPosition,
			ep.Partition,
			configured.Partition,
		)
	}

	// TODO(a.gurinov): Why lte?
	if ep.Offset <= configured.Offset {
		return fmt.Errorf(
			"%w: unexpected offset=%d (configured gt %d)",
			ErrInvalidEventPosition,
			ep.Offset,
			configured.Offset,
		)
	}

	return nil
}
