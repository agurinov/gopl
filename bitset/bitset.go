package bitset

import (
	"golang.org/x/exp/constraints"
)

type BitSet[B constraints.Unsigned] struct {
	set B
}

func New[B constraints.Unsigned](set B) BitSet[B] {
	return BitSet[B]{
		set: set,
	}
}

func (b *BitSet[B]) Set(bit B)     { b.set |= bit }
func (b *BitSet[B]) Clear(bit B)   { b.set &^= bit }
func (b BitSet[B]) Has(bit B) bool { return b.set&bit != 0 }
func (b *BitSet[B]) Toggle(bit B)  { b.set ^= bit }
