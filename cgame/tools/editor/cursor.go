package main

import (
	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cterm"
	"github.com/jf-tech/console/cwin"
	"github.com/jf-tech/go-corelib/maths"
)

func (e *editor) resetCursor() {
	x, y := 0, 0
	if e.cursor != nil {
		x, y = e.cursor.Rect().X, e.cursor.Rect().Y
		e.removeCursor()
	}
	e.cursor = cgame.NewSpriteBase(e.g, e.winMap, cursorName,
		e.createCursorFrame(cwin.Rect{X: x, Y: y, W: 1, H: 1}), x, y)
	e.moveCursor(0, 0)
}

func (e *editor) removeCursor() {
	// TODO update ruler
	e.g.SpriteMgr.DeleteSprite(e.cursor)
}

func (e *editor) createCursorFrame(r cwin.Rect) cgame.Frame {
	var f cgame.Frame
	adjust := func(c, def cterm.Attribute) cterm.Attribute {
		if c == cterm.ColorDefault {
			return def
		}
		return c
	}
	r.Range(func(x, y int) {
		chx := e.winMap.GetClient(x, y)
		chx.Attr = cwin.Attr{
			Fg: adjust(chx.Attr.Bg, cterm.ColorBlack),
			Bg: adjust(chx.Attr.Fg, cterm.ColorWhite),
		}
		f = append(f, cgame.Cell{X: x - r.X, Y: y - r.Y, Chx: chx})
	})
	return f
}

func (e *editor) moveCursor(dx, dy int) {
	r := e.cursor.Rect()
	newR := cwin.Rect{
		X: r.X + dx,
		Y: r.Y + dy,
		W: r.W,
		H: r.H,
	}
	if o, _ := e.winMap.ClientRect().ToOrigin().Overlap(newR); o != newR {
		return
	}
	e.cursor.Update(cgame.UpdateArg{
		DXY: &cwin.Point{
			X: dx,
			Y: dy,
		},
		F:   e.createCursorFrame(newR),
		IBC: cgame.InBoundsCheckNone,
		CD:  cgame.CollisionDetectionOff,
	})
	// now move e.winMap so that cursor is always visible
	xOffsetFromFrame := e.winMap.Rect().X + e.winMap.ClientRect().X
	yOffsetFromFrame := e.winMap.Rect().Y + e.winMap.ClientRect().Y

	dx, dy = 0, 0
	if xOffsetFromFrame+newR.X < 0 {
		dx = -(xOffsetFromFrame + newR.X)
	} else if xOffsetFromFrame+newR.X+newR.W-1 >= e.winMapFrame.ClientRect().W {
		dx = e.winMapFrame.ClientRect().W - (xOffsetFromFrame + newR.X + newR.W)
	}
	if yOffsetFromFrame+newR.Y < 0 {
		dy = -(yOffsetFromFrame + newR.Y)
	} else if yOffsetFromFrame+newR.Y+newR.H-1 >= e.winMapFrame.ClientRect().H {
		dy = e.winMapFrame.ClientRect().H - (yOffsetFromFrame + newR.Y + newR.H)
	}
	e.winMap.SetPosRel(dx, dy)
}

func (e *editor) putAtCursor(chx cwin.Chx) {
	e.cursor.Rect().Range(func(x, y int) {
		e.winMap.PutClient(x, y, chx)
	})
	e.moveCursor(0, 0) // update the cursor frame with the newly set rune(s)
}

func (e *editor) createSelection() {
	if e.selection != nil {
		return
	}
	pt := cwin.Point{X: e.cursor.Rect().X, Y: e.cursor.Rect().Y}
	e.selection = &selection{
		e:       e,
		origin:  pt,
		current: pt,
	}
}

func (e *editor) abortSelection() {
	if e.selection == nil {
		return
	}
	e.selection = nil
	e.resetCursor()
}

// TODO: rethink about selection. should always be in selection, a single char cursor is a
// single char selection. this way the programming model is consistent.

type selection struct {
	e               *editor
	origin, current cwin.Point
}

func (s *selection) move(delta cwin.Point) {
	newCurrent := cwin.Point{X: s.current.X + delta.X, Y: s.current.Y + delta.Y}
	if !s.e.winMap.ClientRect().ToOrigin().Contain(newCurrent.X, newCurrent.Y) {
		return
	}
	s.current = newCurrent
}

func (s *selection) currentRect() cwin.Rect {
	r := cwin.Rect{
		X: maths.MinInt(s.origin.X, s.current.X),
		Y: maths.MinInt(s.origin.Y, s.current.Y),
	}
	r.W = maths.MaxInt(s.origin.X, s.current.X) - r.X + 1
	r.H = maths.MaxInt(s.origin.Y, s.current.Y) - r.Y + 1
	return r
}
