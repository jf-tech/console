package cgame

import "time"

type WaypointType int

const (
	WaypointAbs WaypointType = iota
	WaypointRelative
)

// WaypointAbs: from current position to (X, Y) using time T
// WaypointRel: from current position to (curX + X, curY + Y) using time T. Note a trick here
// is to use X=Y=0, which means keep the sprite at the current location for time T.
type Waypoint struct {
	Type WaypointType
	X, Y int
	T    time.Duration
}

type WaypointProvider interface {
	Next() (Waypoint, bool)
}

type simpleWaypoints struct {
	wps []Waypoint
	idx int
}

func (sw *simpleWaypoints) Next() (Waypoint, bool) {
	if sw.idx >= len(sw.wps) {
		return Waypoint{}, false
	}
	wp := sw.wps[sw.idx]
	sw.idx++
	return wp, true
}

func NewSimpleWaypoints(wps []Waypoint) *simpleWaypoints {
	return &simpleWaypoints{wps: wps}
}
