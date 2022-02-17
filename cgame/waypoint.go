package cgame

import "time"

// From current position to (curX + DX, curY + DY) using time T. Note a trick here
// is to use DX=DY=0 to keep the sprite at the current location for time T.
type Waypoint struct {
	DX, DY int
	T      time.Duration
}

type WaypointProvider interface {
	Next() (Waypoint, bool)
}

type simpleWaypoints struct {
	wps  []Waypoint
	idx  int
	loop bool
}

func (sw *simpleWaypoints) Next() (Waypoint, bool) {
	if sw.idx >= len(sw.wps) {
		return Waypoint{}, false
	}
	wp := sw.wps[sw.idx]
	sw.idx++
	if sw.loop {
		sw.idx = sw.idx % len(sw.wps)
	}
	return wp, true
}

func NewSimpleWaypoints(wps []Waypoint) *simpleWaypoints {
	return &simpleWaypoints{wps: wps}
}

func NewSimpleLoopWaypoints(wps []Waypoint) *simpleWaypoints {
	return &simpleWaypoints{wps: wps, loop: true}
}
