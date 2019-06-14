package equalizer

import (
	"math/big"
	"sync"
)

// Equalizer is a bitmap quota rate limiter
type Equalizer struct {
	sync.RWMutex
	//bitmap tape
	tape *big.Int
	//positive bits mask
	mask *big.Int
	//next index offset manager
	offset Offset
	//bitmap tape size
	size int
	//number of reserved positive bits
	reserved int
}

// NewEqualizer allocates and returns a new Equalizer
//
// len 		- size of the bitmap
// reserved	- number of reserved positive bits
// offset	- Offset instance
func NewEqualizer(size int, reserved int, offset Offset) *Equalizer {
	var tape big.Int
	fill(&tape, 0, size, 1)
	var mask big.Int
	fill(&mask, 0, reserved, 1)
	mask.Lsh(&mask, uint(size-reserved))
	return &Equalizer{
		tape:     &tape,
		mask:     &mask,
		offset:   offset,
		size:     size,
		reserved: reserved,
	}
}

// Claim quota
func (eq *Equalizer) Claim() bool {
	eq.RLock()
	defer eq.RUnlock()
	head := eq.next()
	if eq.tape.Bit(head) > 0 {
		return true
	}
	return false
}

// Notify equalizer (might be in bulk for performance sake)
// shifts bitmap left n bits before settle new ones
// set 	- bit value (false -> 0, true -> 1)
// n	- number of bits to set
func (eq *Equalizer) Notify(set bool, n uint) {
	eq.Lock()
	defer eq.Unlock()
	bit := boolToUint(set)
	var shift big.Int
	fill(&shift, 0, int(n), bit)
	eq.tape.Lsh(eq.tape, n)
	fill(eq.tape, eq.size, eq.size+int(n), 0)
	eq.tape.Or(eq.tape, &shift).Or(eq.tape, eq.mask)
}

// Fill the tape with positives
func (eq *Equalizer) Fill() {
	eq.Lock()
	defer eq.Unlock()
	fill(eq.tape, 0, eq.size, 1)
}

// Reset the tape to reserved mask
func (eq *Equalizer) Reset() {
	eq.Lock()
	defer eq.Unlock()
	fill(eq.tape, 0, eq.reserved, 1)
	fill(eq.tape, eq.reserved, eq.size, 0)
}

// get tape head next index
func (eq *Equalizer) next() int {
	return eq.offset.NextIndex()
}

func fill(n *big.Int, from int, to int, bit uint) {
	for i := from; i < to; i++ {
		n.SetBit(n, i, bit)
	}
}

func boolToUint(b bool) uint {
	if b {
		return 1
	}
	return 0
}
