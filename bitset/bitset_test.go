package bitset_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/agurinov/gopl/bitset"
	pl_testing "github.com/agurinov/gopl/testing"
)

const (
	bit1 = 1 << iota
	bit2
	bit3
	bit4
)

func TestBitSet(t *testing.T) {
	pl_testing.Init(t)

	type (
		args struct {
			initial []uint8
			set     []uint8
			toClear []uint8
			toggle  []uint8
		}
		results struct {
			expectedSet map[uint8]bool
		}
	)

	cases := map[string]struct {
		args    args
		results results
		pl_testing.TestCase
	}{
		"case00: set, clear, toggle, has": {
			args: args{
				initial: []uint8{bit1, bit3},
				set:     []uint8{bit2},
				toClear: []uint8{bit1},
				toggle:  []uint8{bit3, bit4},
			},
			results: results{
				expectedSet: map[uint8]bool{
					bit1: false,
					bit2: true,
					bit3: false,
					bit4: true,
				},
			},
		},
	}

	for name := range cases {
		name, tc := name, cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			var bs bitset.BitSet[uint8]

			for _, i := range tc.args.initial {
				bs.Set(i)
			}

			for _, i := range tc.args.set {
				bs.Set(i)
			}

			for _, i := range tc.args.toClear {
				bs.Clear(i)
			}

			for _, i := range tc.args.toggle {
				bs.Toggle(i)
			}

			for k, v := range tc.results.expectedSet {
				require.Equal(t, v, bs.Has(k))
			}
		})
	}
}
