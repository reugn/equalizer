package equalizer

import (
	"math/rand"
	"sync/atomic"
)

// Offset is the Equalizer's tape offset manager interface.
type Offset interface {
	NextIndex() int
}

// RandomOffset is the random based offset manager.
type RandomOffset struct {
	Len int
}

// NewRandomOffset allocates and returns a new RandomOffset.
// len is the bitmap length.
func NewRandomOffset(len int) *RandomOffset {
	return &RandomOffset{len}
}

// NextIndex returns the next random index.
func (ro *RandomOffset) NextIndex() int {
	return rand.Intn(ro.Len)
}

// StepOffset is the step based offset manager.
type StepOffset struct {
	Len           int
	Step          int64
	previousIndex int64
}

// NewStepOffset allocates and returns a new StepOffset.
// len is the bitmap length.
// step is the offset from the previous index.
func NewStepOffset(len int, step int64) *StepOffset {
	return &StepOffset{len, step, 0}
}

// NextIndex returns the next index in the Round-robin way.
func (so *StepOffset) NextIndex() int {
	return int(atomic.AddInt64(&so.previousIndex, so.Step)) % so.Len
}
