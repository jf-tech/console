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
