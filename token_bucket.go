package equalizer

import (
	"runtime"
	"sync/atomic"
	"time"

	"github.com/reugn/equalizer/internal"
)

// TokenBucket rate limiter
type TokenBucket struct {
	t *tokenBucket
}

type tokenBucket struct {
	capacity       int32
	refillInterval time.Duration
	issued         int32
	sweeper        *internal.Sweeper
}

// NewTokenBucket returns new TokenBucket instance
func NewTokenBucket(capacity int32, refillInterval time.Duration) *TokenBucket {
	tokenBucket := &tokenBucket{
		capacity:       capacity,
		refillInterval: refillInterval,
		issued:         0,
	}
	tb := &TokenBucket{tokenBucket}
	runSweeper(tokenBucket, refillInterval)
	runtime.SetFinalizer(tb, stopSweeper)
	return tb
}

func (tb *tokenBucket) refill() {
	atomic.StoreInt32(&tb.issued, 0)
}

// Claim quota
func (tb *TokenBucket) Claim() bool {
	if tb.t.issued < tb.t.capacity {
		atomic.AddInt32(&tb.t.issued, 1)
		return true
	}
	return false
}

func stopSweeper(t *TokenBucket) {
	t.t.sweeper.Stop <- struct{}{}
}

func runSweeper(t *tokenBucket, d time.Duration) {
	s := &internal.Sweeper{
		Interval: d,
		Stop:     make(chan interface{}),
	}
	t.sweeper = s
	go s.Run(t.refill)
}
