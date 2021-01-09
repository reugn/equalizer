package internal

import "time"

// Callback represents a background task to run.
type Callback func()

// Sweeper verifies shutting down background processing.
type Sweeper struct {
	Interval time.Duration
	Stop     chan interface{}
}

// Run Sweeper with the specified callback.
func (s *Sweeper) Run(callback Callback) {
	ticker := time.NewTicker(s.Interval)
	for {
		select {
		case <-ticker.C:
			callback()
		case <-s.Stop:
			ticker.Stop()
			return
		}
	}
}
