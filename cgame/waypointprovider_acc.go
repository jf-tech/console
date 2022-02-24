package cgame

import (
	"time"

	"github.com/jf-tech/console/cutil"
)

const (
	defaultDeltaT = time.Millisecond
)

type WaypointProviderAccelerationCfg struct {
	Clock      *cutil.Clock
	InitXSpeed CharPerSec
	InitYSpeed CharPerSec
	AccX       CharPerSecSec
	AccY       CharPerSecSec
	DeltaT     time.Duration // if 0, defaults to defaultDeltaT
}

type WaypointProviderAcceleration struct {
	cfg            WaypointProviderAccelerationCfg
	startTime      time.Duration
	dxDone, dyDone int
	curXSpeed      CharPerSec
	curYSpeed      CharPerSec
}

func (s *WaypointProviderAcceleration) Cfg() WaypointProviderAccelerationCfg {
	return s.cfg
}

func (s *WaypointProviderAcceleration) Next() (Waypoint, bool) {
	now := s.cfg.Clock.Now()
	dt := s.cfg.DeltaT - now%s.cfg.DeltaT
	if dt == 0 {
		dt = s.cfg.DeltaT
	}
	totalT := now + dt - s.startTime
	var totalDX, totalDY int
	totalDX, s.curXSpeed = s.calc(s.cfg.InitXSpeed, s.cfg.AccX, totalT)
	totalDY, s.curYSpeed = s.calc(s.cfg.InitYSpeed, s.cfg.AccY, totalT)
	wp := Waypoint{
		DX: totalDX - s.dxDone,
		DY: totalDY - s.dyDone,
		T:  dt,
	}
	s.dxDone = totalDX
	s.dyDone = totalDY
	return wp, true
}

func (s *WaypointProviderAcceleration) CurSpeed() (xSpeed, ySpeed CharPerSec) {
	return s.curXSpeed, s.curYSpeed
}

func (s *WaypointProviderAcceleration) calc(
	v CharPerSec, a CharPerSecSec, t time.Duration) (dx int, speed CharPerSec) {
	sec := float64(t) / float64(time.Second)
	// D = vt + (at^2)/2
	dx = int(float64(v)*sec + (float64(a)*sec*sec)/2)
	// S = v + at
	speed = CharPerSec(float64(v) + float64(a)*sec)
	return dx, speed
}

func NewWaypointProviderAcceleration(
	cfg WaypointProviderAccelerationCfg) *WaypointProviderAcceleration {
	if cfg.DeltaT == 0 {
		cfg.DeltaT = defaultDeltaT
	}
	return &WaypointProviderAcceleration{
		cfg:       cfg,
		startTime: cfg.Clock.Now(),
	}
}
