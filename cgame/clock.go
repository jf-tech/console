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

type Stopwatch struct {
	clock     *Clock
	total     time.Duration
	startedOn time.Duration
}

func (s *Stopwatch) Start() {
	if s.started() {
		panic("cannot Start() a started counter")
	}
	s.startedOn = s.clock.Now()
}

func (s *Stopwatch) Stop() {
	if !s.started() {
		panic("cannot Stop() a stopped counter")
	}
	s.total += s.clock.Now() - s.startedOn
	s.startedOn = -1
}

func (s *Stopwatch) Reset() {
	if s.started() {
		panic("cannot Reset() a started counter")
	}
	s.total = 0
	s.startedOn = -1
}

func (s *Stopwatch) Total() time.Duration {
	if s.started() {
		panic("cannot get Total() on a started counter")
	}
	return s.total
}

func (s *Stopwatch) started() bool {
	return s.startedOn >= 0
}

func NewStopwatch(clock *Clock) *Stopwatch {
	return &Stopwatch{clock: clock, startedOn: -1}
}
