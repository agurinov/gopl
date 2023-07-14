package segmentio

import (
	"bytes"
	"encoding/gob"

	"github.com/agurinov/gopl/kafka"
	segmentio "github.com/segmentio/kafka-go"
)

type userdata struct {
	Topic     string
	Partition int
}

type member struct {
	segmentio.GroupMember
	userdata
}

type SingleAssignmentGroupBalancer struct {
	Topic     string
	Partition int
}

var (
	_ segmentio.GroupBalancer = (*SingleAssignmentGroupBalancer)(nil)
	_ segmentio.GroupBalancer = SingleAssignmentGroupBalancer{}
)

func (b SingleAssignmentGroupBalancer) ProtocolName() string {
	return "single-assignment"
}

func (b SingleAssignmentGroupBalancer) UserData() ([]byte, error) {
	var (
		buf   bytes.Buffer
		udata = userdata{
			Topic:     b.Topic,
			Partition: b.Partition,
		}
	)

	if err := gob.NewEncoder(&buf).Encode(udata); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (b SingleAssignmentGroupBalancer) AssignGroups(
	members []segmentio.GroupMember,
	partitions []segmentio.Partition,
) segmentio.GroupMemberAssignments {
	var (
		groupAssignments    = segmentio.GroupMemberAssignments{}
		undecidedMembers    = make([]member, 0, len(members))
		availablePartitions = make(map[userdata]struct{})
	)

	for i := range partitions {
		udata := userdata{
			Topic:     partitions[i].Topic,
			Partition: partitions[i].ID,
		}
		availablePartitions[udata] = struct{}{}
	}

	assign := func(u userdata, m member) {
		_, available := availablePartitions[u]

		if !available {
			return
		}

		delete(availablePartitions, u)
		groupAssignments[m.ID] = map[string][]int{
			u.Topic: {u.Partition},
		}
	}

	for i := range members {
		var (
			buf   = bytes.NewBuffer(members[i].UserData)
			udata userdata
		)

		if err := gob.NewDecoder(buf).Decode(&udata); err != nil {
			// TODO(a.gurinov): fix it to use first topic from segmentio interface
			continue
		}

		m := member{
			members[i],
			udata,
		}

		switch {
		case udata.Partition == int(kafka.UknownPartition):
			undecidedMembers = append(undecidedMembers, m)
		default:
			assign(udata, m)
		}
	}

	for i := range undecidedMembers {
		for udata := range availablePartitions {
			if undecidedMembers[i].userdata.Topic == udata.Topic {
				assign(udata, undecidedMembers[i])
			}
		}
	}

	return groupAssignments
}
