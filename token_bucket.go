package equalizer

import (
	"runtime"
	"sync/atomic"
	"time"

	"github.com/reugn/equalizer/internal"
)

// A TokenBucket represents a rate limiter based on a custom implementation of the
// token bucket algorithm with a refill interval. Implements the equalizer.Limiter
// interface.
//
// A TokenBucket is safe for use by multiple goroutines simultaneously.
//
// The underlying tokenBucket instance ensures that the background goroutine does
// not keep the main TokenBucket object from being garbage collected. When it is
// garbage collected, the finalizer stops the background goroutine, after
// which the underlying tokenBucket can be collected.
//
type TokenBucket struct {
	t *tokenBucket
}

type tokenBucket struct {
	capacity       int32
	refillInterval time.Duration
	permits        chan token
	issued         int32
	sweeper        *internal.Sweeper
}

// NewTokenBucket allocates and returns a new TokenBucket rate limiter.
//
// capacity is the token bucket capacity.
// refillInterval is the bucket refill interval.
func NewTokenBucket(capacity int32, refillInterval time.Duration) *TokenBucket {
	underlying := &tokenBucket{
		capacity:       capacity,
		refillInterval: refillInterval,
		permits:        make(chan token, capacity),
		issued:         0,
	}
	underlying.fillTokenBucket(int(capacity))

	tokenBucket := &TokenBucket{
		t: underlying,
	}

	// start the refilling goroutine
	goRefill(underlying, refillInterval)
	// the finalizer may run as soon as the tokenBucket becomes unreachable.
	runtime.SetFinalizer(tokenBucket, stopSweeper)

	return tokenBucket
}

// fillTokenBucket fills the permits channel.
func (tb *tokenBucket) fillTokenBucket(capacity int) {
	for i := 0; i < capacity; i++ {
		tb.permits <- token{}
	}
}

// refill refills the token bucket.
func (tb *tokenBucket) refill() {
	issued := atomic.LoadInt32(&tb.issued)
	atomic.StoreInt32(&tb.issued, 0)
	tb.fillTokenBucket(int(issued))
}

// Ask requires a permit.
// It is a non blocking call, returns true or false.
func (tb *TokenBucket) Ask() bool {
	select {
	case <-tb.t.permits:
		atomic.AddInt32(&tb.t.issued, 1)
		return true
	default:
		return false
	}
}

// Take blocks to get a permit.
func (tb *TokenBucket) Take() {
	<-tb.t.permits
	atomic.AddInt32(&tb.t.issued, 1)
}

// stopSweeper is the callback to stop the refilling goroutine.
func stopSweeper(t *TokenBucket) {
	t.t.sweeper.Stop <- struct{}{}
}

// goRefill starts the refilling goroutine.
func goRefill(t *tokenBucket, refillInterval time.Duration) {
	sweeper := &internal.Sweeper{
		Interval: refillInterval,
		Stop:     make(chan interface{}),
	}
	t.sweeper = sweeper
	go sweeper.Run(t.refill)
}

// Verify TokenBucket satisfies the equalizer.Limiter interface.
var _ Limiter = (*TokenBucket)(nil)
