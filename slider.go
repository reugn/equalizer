package equalizer

import (
	"runtime"
	"sync"
	"time"

	"github.com/reugn/equalizer/internal"
)

// Slider is a sliding window based rate limiter
type Slider struct {
	slider *slider
}

type slider struct {
	sync.RWMutex
	window   time.Duration
	slide    time.Duration
	capacity int
	timeline []int64
	issued   int
	sweeper  *internal.Sweeper
}

// NewSlider returns new Slider instance
func NewSlider(window time.Duration, slide time.Duration, capacity int) *Slider {
	timelineArray := make([]int64, capacity)
	slider := &slider{
		window:   window,
		slide:    slide,
		capacity: capacity,
		timeline: timelineArray,
		issued:   0,
	}
	sl := &Slider{slider}
	runSliderSweeper(slider, slide)
	runtime.SetFinalizer(sl, stopSliderSweeper)
	return sl
}

// slide window
func (s *slider) doSlide() {
	s.Lock()
	//clean outdated quotas and set counter
	threshold := nowNano() - s.window.Nanoseconds()
	s.timeline, s.issued = filter(s.timeline, threshold)
	s.Unlock()
}

// Claim quota
func (s *Slider) Claim() bool {
	s.slider.Lock()
	defer s.slider.Unlock()
	if s.slider.issued < s.slider.capacity {
		s.slider.timeline[s.slider.issued] = nowNano()
		s.slider.issued++
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

func stopSliderSweeper(s *Slider) {
	s.slider.sweeper.Stop <- struct{}{}
}

func runSliderSweeper(slider *slider, d time.Duration) {
	s := &internal.Sweeper{
		Interval: d,
		Stop:     make(chan interface{}),
	}
	slider.sweeper = s
	go s.Run(slider.doSlide)
}
