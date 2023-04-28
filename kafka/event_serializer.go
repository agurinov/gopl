package kafka

import (
	"encoding/json"
)

type EventSerializer[E Event] func([]byte) (E, error)

func JsonEventSerializer[E Event](data []byte) (E, error) {
	var e E

	if err := json.Unmarshal(data, &e); err != nil {
		return e, err
	}

	return e, nil
}

func ProtoEventSerializer[E Event](data []byte) (E, error) {
	var e E

	// TODO(a.gurinov): fix it
	// if err := proto.Unmarshal(data, &e); err != nil {
	// 	return e, err
	// }

	return e, nil
}
