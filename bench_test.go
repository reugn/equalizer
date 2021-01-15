package equalizer

import (
	"testing"
	"time"
)

// go test -bench=. -benchmem
func BenchmarkEqualizerShortAskStep(b *testing.B) {
	offset := NewStepOffset(96, 1)
	eq := NewEqualizer(96, 16, offset)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eq.Ask()
	}
}

func BenchmarkEqualizerShortAskRandom(b *testing.B) {
	offset := NewRandomOffset(96)
	eq := NewEqualizer(96, 16, offset)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eq.Ask()
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

func BenchmarkEqualizerLongAskStep(b *testing.B) {
	offset := NewStepOffset(1048576, 1)
	eq := NewEqualizer(1048576, 16, offset)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eq.Ask()
	}
}

func BenchmarkEqualizerLongAskRandom(b *testing.B) {
	offset := NewRandomOffset(1048576)
	eq := NewEqualizer(1048576, 16, offset)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eq.Ask()
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

func BenchmarkSliderShortWindow(b *testing.B) {
	slider := NewSlider(time.Millisecond*100, time.Millisecond*10, 32)
	for i := 0; i < b.N; i++ {
		slider.Ask()
	}
}

func BenchmarkSliderLongerWindow(b *testing.B) {
	slider := NewSlider(time.Second, time.Millisecond*100, 32)
	for i := 0; i < b.N; i++ {
		slider.Ask()
	}
}

func BenchmarkTokenBucketDenseRefill(b *testing.B) {
	tokenBucket := NewTokenBucket(32, time.Millisecond*10)
	for i := 0; i < b.N; i++ {
		tokenBucket.Ask()
	}
}

func BenchmarkTokenBucketSparseRefill(b *testing.B) {
	tokenBucket := NewTokenBucket(32, time.Second)
	for i := 0; i < b.N; i++ {
		tokenBucket.Ask()
	}
}
