//go:build test_unit

package segmentio_test

import (
	"testing"

	sk "github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"

	"github.com/agurinov/gopl/kafka"
	"github.com/agurinov/gopl/kafka/segmentio"
	pl_testing "github.com/agurinov/gopl/testing"
)

func TestSingleAssignmentGroupBalancer_AssignGroups(t *testing.T) {
	pl_testing.Init(t)

	t1p0 := sk.Partition{Topic: "topic_1", ID: 0}
	t1p1 := sk.Partition{Topic: "topic_1", ID: 1}
	t1p2 := sk.Partition{Topic: "topic_1", ID: 2}
	t2p0 := sk.Partition{Topic: "topic_2", ID: 0}
	t2p1 := sk.Partition{Topic: "topic_2", ID: 1}
	t3p0 := sk.Partition{Topic: "topic_3", ID: 0}

	// member1 always win member2 (same assigned partition)
	// member3 and member4 share topic (member4 undecided about partition)
	// member5 always win (member6 have no any available partition)
	// member7 want non existent topic
	// member8 want non existent partition

	member1userdata, err := segmentio.SingleAssignmentGroupBalancer{
		Topic:     t1p2.Topic,
		Partition: t1p2.ID,
	}.UserData()
	require.NoError(t, err)
	member1 := sk.GroupMember{ID: "member_1", UserData: member1userdata}

	member2userdata, err := segmentio.SingleAssignmentGroupBalancer{
		Topic:     t1p2.Topic,
		Partition: t1p2.ID,
	}.UserData()
	require.NoError(t, err)
	member2 := sk.GroupMember{ID: "member_2", UserData: member2userdata}

	member3userdata, err := segmentio.SingleAssignmentGroupBalancer{
		Topic:     t2p1.Topic,
		Partition: t2p1.ID,
	}.UserData()
	require.NoError(t, err)
	member3 := sk.GroupMember{ID: "member_3", UserData: member3userdata}

	member4userdata, err := segmentio.SingleAssignmentGroupBalancer{
		Topic:     t2p1.Topic,
		Partition: int(kafka.UknownPartition),
	}.UserData()
	require.NoError(t, err)
	member4 := sk.GroupMember{ID: "member_4", UserData: member4userdata}

	member5userdata, err := segmentio.SingleAssignmentGroupBalancer{
		Topic:     t3p0.Topic,
		Partition: t3p0.ID,
	}.UserData()
	require.NoError(t, err)
	member5 := sk.GroupMember{ID: "member_5", UserData: member5userdata}

	member6userdata, err := segmentio.SingleAssignmentGroupBalancer{
		Topic:     t3p0.Topic,
		Partition: int(kafka.UknownPartition),
	}.UserData()
	require.NoError(t, err)
	member6 := sk.GroupMember{ID: "member_6", UserData: member6userdata}

	member7userdata, err := segmentio.SingleAssignmentGroupBalancer{
		Topic:     "topic_100",
		Partition: int(kafka.UknownPartition),
	}.UserData()
	require.NoError(t, err)
	member7 := sk.GroupMember{ID: "member_7", UserData: member7userdata}

	member8userdata, err := segmentio.SingleAssignmentGroupBalancer{
		Topic:     "topic_1",
		Partition: 100,
	}.UserData()
	require.NoError(t, err)
	member8 := sk.GroupMember{ID: "member_8", UserData: member8userdata}

	cases := map[string]struct {
		inputMembers             []sk.GroupMember
		inputPartitions          []sk.Partition
		expectedGroupAssignments sk.GroupMemberAssignments
		pl_testing.TestCase
	}{
		"case00: all decided members": {
			inputMembers: []sk.GroupMember{
				member1, member2, member3, member4, member5, member6, member7, member8,
			},
			inputPartitions: []sk.Partition{t1p0, t1p1, t1p2, t2p0, t2p1, t3p0},
			expectedGroupAssignments: sk.GroupMemberAssignments{
				"member_1": {"topic_1": {2}},
				"member_3": {"topic_2": {1}},
				"member_4": {"topic_2": {0}},
				"member_5": {"topic_3": {0}},
			},
		},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			var (
				balancer         = segmentio.SingleAssignmentGroupBalancer{}
				groupAssignments = balancer.AssignGroups(tc.inputMembers, tc.inputPartitions)
			)

			require.Equal(t,
				tc.expectedGroupAssignments,
				groupAssignments,
			)
		})
	}
}
