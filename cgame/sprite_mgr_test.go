package cgame

import (
	"testing"

	"github.com/jf-tech/console/cwin"
	"github.com/stretchr/testify/assert"
)

func TestDetectCollision(t *testing.T) {
	sm := &SpriteManager{
		collisionDetectionBuf: make([]bool, 0),
	}
	// case 1: rects not overlapping.
	w1 := cwin.NewWin(nil, cwin.WinCfg{R: cwin.Rect{X: 0, Y: 0, W: 5, H: 5}, NoBorder: true})
	w2 := cwin.NewWin(nil, cwin.WinCfg{R: cwin.Rect{X: 5, Y: 0, W: 5, H: 5}, NoBorder: true})
	assert.False(t, sm.detectCollision(w1, w2))

	// case 2: rects overlapping, but only transparent cells intersect
	w1 = cwin.NewWin(nil, cwin.WinCfg{R: cwin.Rect{X: 0, Y: 0, W: 5, H: 5}, NoBorder: true})
	FrameToWin(FrameFromString(`
12345
12
1234
12
12345`, cwin.ChAttr{}), w1)
	w2 = cwin.NewWin(nil, cwin.WinCfg{R: cwin.Rect{X: 2, Y: 1, W: 3, H: 3}, NoBorder: true})
	FrameToWin(FrameFromString(`
abc
  c
abc`, cwin.ChAttr{}), w2)
	assert.False(t, sm.detectCollision(w1, w2))

	// case 3: rects overlapping, and one non-transparent cell intersects
	w1 = cwin.NewWin(nil, cwin.WinCfg{R: cwin.Rect{X: 0, Y: 0, W: 3, H: 3}, NoBorder: true})
	FrameToWin(FrameFromString(`
123
1 3
123
`, cwin.ChAttr{}), w1)
	w2 = cwin.NewWin(nil, cwin.WinCfg{R: cwin.Rect{X: 2, Y: 2, W: 2, H: 2}, NoBorder: true})
	FrameToWin(FrameFromString(`
ab
ab
`, cwin.ChAttr{}), w2)
	assert.True(t, sm.detectCollision(w1, w2))
}
