package equalizer

import (
	"testing"
	"time"
)

func TestTokenBucket(t *testing.T) {
	tokenBucket := NewTokenBucket(32, time.Millisecond*100)
	var quota bool
	for i := 0; i < 32; i++ {
		quota = tokenBucket.Ask()
	}
	assertEqual(t, quota, true)

	quota = tokenBucket.Ask()
	assertEqual(t, quota, false)

	time.Sleep(time.Millisecond * 110)

	quota = tokenBucket.Ask()
	assertEqual(t, quota, true)

	tokenBucket.Take()
}
