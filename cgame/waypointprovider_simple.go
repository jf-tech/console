package cgame

type waypointProviderSimple struct {
	wps  []Waypoint
	idx  int
	loop bool
}

func (s *waypointProviderSimple) Next() (Waypoint, bool) {
	if s.idx >= len(s.wps) {
		return Waypoint{}, false
	}
	wp := s.wps[s.idx]
	s.idx++
	if s.loop {
		s.idx = s.idx % len(s.wps)
	}
	return wp, true
}

func NewWaypointProviderSimple(wps []Waypoint) *waypointProviderSimple {
	return &waypointProviderSimple{wps: wps}
}

func NewWaypointProviderSimpleLoop(wps []Waypoint) *waypointProviderSimple {
	return &waypointProviderSimple{wps: wps, loop: true}
}
