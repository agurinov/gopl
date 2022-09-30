package pl_bitset

type bitset interface {
	uint8 | uint16 | uint32 | uint64 | uint
}

type BitSet[B bitset] struct {
	set B
}

func (b *BitSet[B]) Set(bit B)     { b.set |= bit }
func (b *BitSet[B]) Clear(bit B)   { b.set &^= bit }
func (b BitSet[B]) Has(bit B) bool { return b.set&bit != 0 }
func (b *BitSet[B]) Toggle(bit B)  { b.set ^= bit }
