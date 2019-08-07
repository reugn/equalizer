package equalizer

import (
	"sync/atomic"
	"time"
)

// TokenBucket rate limiter
type TokenBucket struct {
	capacity       int32
	refillInterval time.Duration
	issued         int32
}

// NewTokenBucket returns new TokenBucket instance
func NewTokenBucket(capacity int32, refillInterval time.Duration) *TokenBucket {
	tokenBucket := &TokenBucket{
		capacity:       capacity,
		refillInterval: refillInterval,
		issued:         0,
	}
	go tokenBucket.init()
	return tokenBucket
}

// init refill loop
func (tb *TokenBucket) init() {
	ticker := time.NewTicker(tb.refillInterval)
	for range ticker.C {
		atomic.StoreInt32(&tb.issued, 0)
	}
}

// Claim quota
func (tb *TokenBucket) Claim() bool {
	if tb.issued < tb.capacity {
		atomic.AddInt32(&tb.issued, 1)
		return true
	}
	return false
}
