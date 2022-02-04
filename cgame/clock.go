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

type clockManager struct {
	clocks map[*Clock]bool
}

func (cm *clockManager) createClock() *Clock {
	c := newClock()
	cm.clocks[c] = true
	return c
}

func (cm *clockManager) deleteClock(c *Clock) {
	delete(cm.clocks, c)
}

func (cm *clockManager) pauseAll() {
	for c := range cm.clocks {
		c.Pause()
	}
}

func (cm *clockManager) resumeAll() {
	for c := range cm.clocks {
		c.Resume()
	}
}

func newClockManager() *clockManager {
	return &clockManager{
		clocks: make(map[*Clock]bool),
	}
}
