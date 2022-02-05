package cgame

import (
	"time"
)

type Clock struct {
	originTime           time.Time
	totalPausedDuration  time.Duration
	latestPauseStartTime time.Time
	paused               bool
}

func newClock() *Clock {
	return &Clock{originTime: time.Now()}
}

func (c *Clock) Now() time.Duration {
	if c.IsPaused() {
		return c.latestPauseStartTime.Sub(c.originTime) - c.totalPausedDuration
	}
	return time.Since(c.originTime) - c.totalPausedDuration
}

func (c *Clock) Pause() {
	if !c.IsPaused() {
		c.latestPauseStartTime = time.Now()
		c.paused = true
	}
}

func (c *Clock) Resume() {
	if c.IsPaused() {
		c.totalPausedDuration += time.Since(c.latestPauseStartTime)
		c.latestPauseStartTime = time.Time{}
		c.paused = false
	}
}

func (c *Clock) IsPaused() bool {
	return c.paused
}
