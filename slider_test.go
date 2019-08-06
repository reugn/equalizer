package equalizer

import (
	"testing"
	"time"
)

func TestSlider(t *testing.T) {
	slider := NewSlider(time.Second, time.Millisecond*100, 32)
	var quota bool
	for i := 0; i < 32; i++ {
		quota = slider.Claim()
	}
	assertEqual(t, quota, true)
	quota = slider.Claim()
	assertEqual(t, quota, false)
	time.Sleep(time.Millisecond * 1010)
	quota = slider.Claim()
	assertEqual(t, quota, true)
}
