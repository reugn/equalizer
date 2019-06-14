package equalizer

import (
	"math/rand"
	"sync/atomic"
)

// Offset is an equalizer tape offset manager interface
type Offset interface {
	NextIndex() int
}

// RandomOffset - random based offset manager
type RandomOffset struct {
	Len int
}

// NewRandomOffset allocates and returns a new RandomOffset
// len - max length
func NewRandomOffset(len int) *RandomOffset {
	return &RandomOffset{len}
}

// NextIndex returns next random index
func (ro *RandomOffset) NextIndex() int {
	return rand.Intn(ro.Len)
}

// StepOffset - step based offset manager
type StepOffset struct {
	Len           int
	Step          int64
	previousIndex int64
}

// NewStepOffset allocates and returns a new StepOffset
// len  - max length
// step - offset from previous index
func NewStepOffset(len int, step int64) *StepOffset {
	return &StepOffset{len, step, 0}
}

// NextIndex returns next index in circular way
func (so *StepOffset) NextIndex() int {
	return int(atomic.AddInt64(&so.previousIndex, so.Step)) % so.Len
}
