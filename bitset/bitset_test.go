//go:build test_unit

package pl_bitset_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	pl_bitset "github.com/agurinov/gopl.git/bitset"
	pl_testing "github.com/agurinov/gopl.git/testing"
)

const (
	bit1 = 1 << iota
	bit2
)

func TestBitSet_Set(t *testing.T) {
	pl_testing.Init(t)

	bitset := pl_bitset.BitSet[uint8]{}

	require.False(t, bitset.Has(bit1))
	require.False(t, bitset.Has(bit2))

	bitset.Set(bit2)

	require.False(t, bitset.Has(bit1))
	require.True(t, bitset.Has(bit2))
}

func TestBitSet_Clear(t *testing.T) {
	pl_testing.Init(t)

	bitset := pl_bitset.BitSet[uint32]{}

	require.False(t, bitset.Has(bit1))
	require.False(t, bitset.Has(bit2))

	bitset.Set(bit1)
	bitset.Set(bit2)
	bitset.Clear(bit2)

	require.True(t, bitset.Has(bit1))
	require.False(t, bitset.Has(bit2))
}

func TestBitSet_Toggle(t *testing.T) {
	pl_testing.Init(t)

	bitset := pl_bitset.BitSet[uint]{}

	require.False(t, bitset.Has(bit1))
	require.False(t, bitset.Has(bit2))

	bitset.Toggle(bit1)
	bitset.Toggle(bit2)

	require.True(t, bitset.Has(bit1))
	require.True(t, bitset.Has(bit2))

	bitset.Toggle(bit1)
	bitset.Toggle(bit2)

	require.False(t, bitset.Has(bit1))
	require.False(t, bitset.Has(bit2))
}
