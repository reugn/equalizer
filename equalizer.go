package equalizer

import (
	"math/big"
	"sync"
)

// An Equalizer represents a bitmap based adaptive rate limiter.
// The quota management algorithm is based on a Round-robin bitmap tape with a moving head.
//
// An Equalizer is safe for use by multiple goroutines simultaneously.
//
// Use Ask function to request quota.
// Use Notify to update the bitmap tape with previous requests' statuses.
//
type Equalizer struct {
	sync.RWMutex

	// tape is the underlying bitmap tape
	tape *big.Int

	// mask is the positive bits mask
	mask *big.Int

	// offset is the next index offset manager
	offset Offset

	// size is the bitmap tape size
	size int

	// reserved is the number of reserved positive bits
	reserved int
}

// NewEqualizer allocates and returns a new Equalizer rate limiter.
//
// len is the size of the bitmap.
// reserved	is the number of reserved positive bits.
// offset is the equalizer.Offset strategy instance.
func NewEqualizer(size int, reserved int, offset Offset) *Equalizer {
	// init the bitmap tape
	var tape big.Int
	fill(&tape, 0, size, 1)

	// init the positive bits mask
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

// Ask moves the tape head to the next index and returns the value.
func (eq *Equalizer) Ask() bool {
	eq.RLock()
	defer eq.RUnlock()

	head := eq.next()
	return eq.tape.Bit(head) > 0
}

// Notify shifts the tape left by n bits and appends the specified value n times.
// Use n > 1 to update the same value in a bulk to gain performance.
// value is the boolean representation of a bit value (false -> 0, true -> 1).
// n is the number of bits to set.
func (eq *Equalizer) Notify(value bool, n uint) {
	eq.Lock()
	defer eq.Unlock()

	bit := boolToUint(value)
	var shift big.Int
	fill(&shift, 0, int(n), bit)
	eq.tape.Lsh(eq.tape, n)
	fill(eq.tape, eq.size, eq.size+int(n), 0)
	eq.tape.Or(eq.tape, &shift).Or(eq.tape, eq.mask)
}

// ResetPositive resets the tape with positive bits.
func (eq *Equalizer) ResetPositive() {
	eq.Lock()
	defer eq.Unlock()

	fill(eq.tape, 0, eq.size, 1)
}

// Reset resets the tape to initial state.
func (eq *Equalizer) Reset() {
	eq.Lock()
	defer eq.Unlock()

	fill(eq.tape, 0, eq.reserved, 1)
	fill(eq.tape, eq.reserved, eq.size, 0)
}

// next returns the next index of the tape head.
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
