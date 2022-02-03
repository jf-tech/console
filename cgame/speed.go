package cgame

import "time"

type ActionPerSec float64

type ActionPerSecTicker struct {
	clock          *GameClock
	aps            ActionPerSec
	startTime      time.Duration
	actsSinceStart int64
}

func (a *ActionPerSecTicker) Start() {
	if !a.IsRunning() {
		a.startTime = a.clock.SinceOrigin()
		a.actsSinceStart = 0
	}
}

func (a *ActionPerSecTicker) Stop() {
	a.startTime = -1
	a.actsSinceStart = 0
}

func (a *ActionPerSecTicker) HowMany() int64 {
	if !a.IsRunning() {
		panic("Start() not called")
	}
	moves := int64(float64(a.clock.SinceOrigin()-a.startTime) / float64(time.Second) * float64(a.aps))
	delta := moves - a.actsSinceStart
	a.actsSinceStart = moves
	return delta
}

func (a *ActionPerSecTicker) IsRunning() bool {
	return a.startTime >= 0
}

func NewActionPerSecTicker(c *GameClock, aps ActionPerSec, autoStart bool) *ActionPerSecTicker {
	a := &ActionPerSecTicker{clock: c, aps: aps, startTime: -1}
	if autoStart {
		a.Start()
	}
	return a
}
