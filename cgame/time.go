package cgame

import "time"

type GameClock struct {
	originTime           time.Time
	totalPausedDuration  time.Duration
	latestPauseStartTime time.Time
}

func newGameClock() *GameClock {
	return &GameClock{originTime: time.Now()}
}

func (gmc *GameClock) SinceOrigin() time.Duration {
	if gmc.isPaused() {
		return gmc.latestPauseStartTime.Sub(gmc.originTime) - gmc.totalPausedDuration
	}
	return time.Since(gmc.originTime) - gmc.totalPausedDuration
}

func (gmc *GameClock) pause() {
	if !gmc.isPaused() {
		gmc.latestPauseStartTime = time.Now()
	}
}

func (gmc *GameClock) resume() {
	if gmc.isPaused() {
		gmc.totalPausedDuration += time.Since(gmc.latestPauseStartTime)
		gmc.latestPauseStartTime = time.Time{}
	}
}

func (gmc *GameClock) isPaused() bool {
	return gmc.latestPauseStartTime != time.Time{}
}
