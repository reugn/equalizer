package equalizer

import (
	"runtime"
	"sync/atomic"
	"time"

	"github.com/reugn/equalizer/internal"
)

// A Slider represents a rate limiter which is based on a sliding window
// with a specified quota capacity. Implements the equalizer.Limiter interface.
//
// A Slider is safe for use by multiple goroutines simultaneously.
//
// The underlying slider instance ensures that the background goroutine does
// not keep the main Slider object from being garbage collected. When it is
// garbage collected, the finalizer stops the background goroutine, after
// which the underlying slider can be collected.
//
type Slider struct {
	slider *slider
}

type slider struct {
	window          time.Duration
	slidingInterval time.Duration
	capacity        int
	permits         chan token
	issued          chan int64
	windowStart     int64
	sweeper         *internal.Sweeper
}

// NewSlider allocates and returns a new Slider rate limiter.
//
// window is the fixed duration of the sliding window.
// slidingInterval controls how frequently a new sliding window is started.
// capacity is the quota limit for the window.
func NewSlider(window time.Duration, slidingInterval time.Duration, capacity int) *Slider {
	underlying := &slider{
		window:          window,
		slidingInterval: slidingInterval,
		capacity:        capacity,
		permits:         make(chan token, capacity),
		issued:          make(chan int64, capacity),
	}
	underlying.initSlider()

	slider := &Slider{
		slider: underlying,
	}

	// start the sliding goroutine
	goSlide(underlying, slidingInterval)
	// the finalizer may run as soon as the slider becomes unreachable.
	runtime.SetFinalizer(slider, stopSliderSweeper)

	return slider
}

// initSlider prefills the permits channel.
func (s *slider) initSlider() {
	for i := 0; i < s.capacity; i++ {
		s.permits <- token{}
	}
}

// doSlide starts a new sliding window.
func (s *slider) doSlide() {
	atomic.StoreInt64(&s.windowStart, nowNano())
	for ts := range s.issued {
		s.permits <- token{}
		if ts > s.windowStart {
			break
		}
	}
}

// Ask requires a permit.
// It is a non blocking call, returns true or false.
func (s *Slider) Ask() bool {
	select {
	case <-s.slider.permits:
		s.slider.issued <- atomic.LoadInt64(&s.slider.windowStart)
		return true
	default:
		return false
	}
}

// Take blocks to get a permit.
func (s *Slider) Take() {
	<-s.slider.permits
	s.slider.issued <- atomic.LoadInt64(&s.slider.windowStart)
}

func nowNano() int64 {
	return time.Now().UTC().UnixNano()
}

// stopSliderSweeper is the callback to stop the sliding goroutine.
func stopSliderSweeper(s *Slider) {
	s.slider.sweeper.Stop <- struct{}{}
}

// goSlide starts the sliding goroutine.
func goSlide(slider *slider, slidingInterval time.Duration) {
	sweeper := &internal.Sweeper{
		Interval: slidingInterval,
		Stop:     make(chan interface{}),
	}
	slider.sweeper = sweeper
	go sweeper.Run(slider.doSlide)
}

// Verify Slider satisfies the equalizer.Limiter interface.
var _ Limiter = (*Slider)(nil)
