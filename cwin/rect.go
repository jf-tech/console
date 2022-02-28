package cwin

import (
	"fmt"

	"github.com/jf-tech/go-corelib/maths"
)

type Rect struct {
	X, Y, W, H int
}

func (r Rect) Contain(x, y int) bool {
	return x >= r.X && x < r.X+r.W && y >= r.Y && y < r.Y+r.H
}

func (r Rect) Overlap(other Rect) (Rect, bool) {
	var overlapped Rect
	overlapped.X = maths.MaxInt(r.X, other.X)
	overlapped.W = maths.MinInt(r.X+r.W, other.X+other.W) - overlapped.X
	if overlapped.W <= 0 {
		return overlapped, false
	}
	overlapped.Y = maths.MaxInt(r.Y, other.Y)
	overlapped.H = maths.MinInt(r.Y+r.H, other.Y+other.H) - overlapped.Y
	return overlapped, overlapped.H > 0
}

func (r Rect) MoveDelta(dx, dy int) Rect {
	return Rect{r.X + dx, r.Y + dy, r.W, r.H}
}

func (r Rect) ToOrigin() Rect {
	return r.MoveDelta(-r.X, -r.Y)
}

func (r Rect) Range(f func(x, y int)) {
	for y := 0; y < r.H; y++ {
		for x := 0; x < r.W; x++ {
			f(r.X+x, r.Y+y)
		}
	}
}

func (r Rect) String() string {
	return fmt.Sprintf("(%d,%d,%d,%d)", r.X, r.Y, r.W, r.H)
}
