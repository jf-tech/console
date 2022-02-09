package cgame

import (
	"math"
	"strings"
	"time"

	"github.com/jf-tech/console/cwin"
	"github.com/jf-tech/go-corelib/maths"
)

type Cell struct {
	X, Y int
	Chx  cwin.Chx
}

type Frame []Cell

func FrameFromStringEx(s string, attr cwin.ChAttr, spaceEqualsTransparency bool) Frame {
	s = strings.Trim(s, "\n")
	rect := cwin.TextDimension(s)
	var f Frame
	rs := []rune(s)
	rsLen := len(rs)
	for i, x, y := 0, 0, 0; i < rsLen; i++ {
		if rs[i] == '\n' {
			if !spaceEqualsTransparency {
				for ; x < rect.W; x++ {
					f = append(f, Cell{X: x, Y: y, Chx: cwin.Chx{Ch: ' ', Attr: attr}})
				}
			}
			x = 0
			y++
			continue
		}
		if rs[i] != ' ' || !spaceEqualsTransparency {
			f = append(f, Cell{X: x, Y: y, Chx: cwin.Chx{Ch: rs[i], Attr: attr}})
		}
		x++
	}
	return f
}

func FrameFromString(s string, attr cwin.ChAttr) Frame {
	return FrameFromStringEx(s, attr, true)
}

func FrameRect(f Frame) cwin.Rect {
	maxX, maxY := math.MinInt32, math.MinInt32
	for i := 0; i < len(f); i++ {
		maxX = maths.MaxInt(maxX, f[i].X)
		maxY = maths.MaxInt(maxY, f[i].Y)
	}
	return cwin.Rect{X: 0, Y: 0, W: maxX + 1, H: maxY + 1}
}

func FrameFromWin(w *cwin.Win) Frame {
	var f Frame
	for y := 0; y < w.ClientRect().H; y++ {
		for x := 0; x < w.ClientRect().W; x++ {
			chx := w.GetClient(x, y)
			if chx != cwin.TransparentChx() {
				f = append(f, Cell{X: x, Y: y, Chx: chx})
			}
		}
	}
	return f
}

func FrameToWin(f Frame, w *cwin.Win) {
	w.FillClient(w.ClientRect().ToOrigin(), cwin.TransparentChx())
	for i := 0; i < len(f); i++ {
		w.PutClient(f[i].X, f[i].Y, f[i].Chx)
	}
}

type Frames []Frame

func FramesFromString(ss []string, attr cwin.ChAttr) Frames {
	var ret Frames
	for _, s := range ss {
		ret = append(ret, FrameFromString(s, attr))
	}
	return ret
}

type FrameProvider interface {
	Next() (Frame, time.Duration, bool)
}

type simpleFrameProvider struct {
	frames Frames
	t      time.Duration
	loop   bool
	idx    int
}

func (sfp *simpleFrameProvider) Next() (Frame, time.Duration, bool) {
	if sfp.idx >= len(sfp.frames) {
		if !sfp.loop {
			return nil, 0, false
		}
		sfp.idx = 0
	}
	f := sfp.frames[sfp.idx]
	sfp.idx++
	return f, sfp.t, true
}

func NewSimpleFrameProvider(frames Frames, t time.Duration, loop bool) *simpleFrameProvider {
	return &simpleFrameProvider{
		frames: frames,
		t:      t,
		loop:   loop,
	}
}
