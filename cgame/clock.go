package cgame

import "time"

type Clock struct {
	originTime           time.Time
	totalPausedDuration  time.Duration
	latestPauseStartTime time.Time
}

func newClock() *Clock {
	return &Clock{originTime: time.Now()}
}

func (c *Clock) Now() time.Duration {
	if c.isPaused() {
		return c.latestPauseStartTime.Sub(c.originTime) - c.totalPausedDuration
	}
	return time.Since(c.originTime) - c.totalPausedDuration
}

func (c *Clock) pause() {
	if !c.isPaused() {
		c.latestPauseStartTime = time.Now()
	}
}

func (c *Clock) resume() {
	if c.isPaused() {
		c.totalPausedDuration += time.Since(c.latestPauseStartTime)
		c.latestPauseStartTime = time.Time{}
	}
}

func (c *Clock) isPaused() bool {
	return c.latestPauseStartTime != time.Time{}
}

type ClockManager struct {
	clocks map[*Clock]bool
}

func (cm *ClockManager) createClock() *Clock {
	c := newClock()
	cm.clocks[c] = true
	return c
}

func (cm *ClockManager) deleteClock(c *Clock) {
	delete(cm.clocks, c)
}

func (cm *ClockManager) PauseAll() {
	for c, _ := range cm.clocks {
		c.pause()
	}
}

func (cm *ClockManager) ResumeAll() {
	for c, _ := range cm.clocks {
		c.resume()
	}
}

func newClockManager() *ClockManager {
	return &ClockManager{
		clocks: make(map[*Clock]bool),
	}
}
