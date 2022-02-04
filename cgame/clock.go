package cgame

import (
	"fmt"
	"time"
)

type Clock struct {
	originTime           time.Time
	totalPausedDuration  time.Duration
	latestPauseStartTime time.Time
	pauseCounter         int
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
	if c.pauseCounter == 0 {
		c.latestPauseStartTime = time.Now()
	}
	c.pauseCounter++
}

func (c *Clock) Resume() {
	c.pauseCounter--
	if c.pauseCounter < 0 {
		panic(fmt.Sprintf("clock pause counter less than zero: %d", c.pauseCounter))
	}
	if c.pauseCounter == 0 {
		c.totalPausedDuration += time.Since(c.latestPauseStartTime)
		c.latestPauseStartTime = time.Time{}
	}
}

func (c *Clock) IsPaused() bool {
	return c.pauseCounter > 0
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
