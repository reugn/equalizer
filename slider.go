package equalizer

import (
	"sync"
	"time"
)

// Slider is sliding window based rate limiter
type Slider struct {
	sync.RWMutex
	window   time.Duration
	slide    time.Duration
	capacity int
	timeline []int64
	issued   int
}

// NewSlider returns new Slider instance
func NewSlider(window time.Duration, slide time.Duration, capacity int) *Slider {
	timelineArray := make([]int64, capacity)
	slider := &Slider{
		window:   window,
		slide:    slide,
		capacity: capacity,
		timeline: timelineArray,
		issued:   0,
	}
	go slider.init()
	return slider
}

// init sliding window
func (s *Slider) init() {
	ticker := time.NewTicker(s.slide)
	for range ticker.C {
		s.Lock()
		//clean outdated quotas and set counter
		threshold := nowNano() - s.window.Nanoseconds()
		s.timeline, s.issued = filter(s.timeline, threshold)
		s.Unlock()
	}
}

// Claim quota
func (s *Slider) Claim() bool {
	s.Lock()
	defer s.Unlock()
	if s.issued < s.capacity {
		s.timeline[s.issued] = nowNano()
		s.issued++
		return true
	}
	return false
}

func nowNano() int64 {
	return time.Now().UTC().UnixNano()
}

func filter(arr []int64, threshold int64) ([]int64, int) {
	rt := make([]int64, cap(arr))
	i := 0
	for _, ts := range arr {
		if ts > threshold {
			rt[i] = ts
			i++
		}
	}
	return rt, i
}
