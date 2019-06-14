package equalizer

import (
	"testing"
)

//go test -bench=. -benchmem
func BenchmarkEqualizerShortClaimStep(b *testing.B) {
	offset := NewStepOffset(96, 1)
	eq := NewEqualizer(96, 16, offset)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eq.Claim()
	}
}

func BenchmarkEqualizerShortClaimRandom(b *testing.B) {
	offset := NewRandomOffset(96)
	eq := NewEqualizer(96, 16, offset)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eq.Claim()
	}
}

func BenchmarkEqualizerShortNotify(b *testing.B) {
	offset := NewStepOffset(96, 1)
	eq := NewEqualizer(96, 16, offset)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eq.Notify(false, 1)
	}
}

func BenchmarkEqualizerLongClaimStep(b *testing.B) {
	offset := NewStepOffset(1048576, 1)
	eq := NewEqualizer(1048576, 16, offset)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eq.Claim()
	}
}

func BenchmarkEqualizerLongClaimRandom(b *testing.B) {
	offset := NewRandomOffset(1048576)
	eq := NewEqualizer(1048576, 16, offset)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eq.Claim()
	}
}

func BenchmarkEqualizerLongNotify(b *testing.B) {
	offset := NewStepOffset(1048576, 1)
	eq := NewEqualizer(1048576, 16, offset)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eq.Notify(false, 1)
	}
}
