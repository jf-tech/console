package main

import (
	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cterm"
	"github.com/jf-tech/console/cwin"
	"github.com/jf-tech/go-corelib/maths"
)

type cursor struct {
	*cgame.SpriteBase
	e                          *editor
	selectionFrom, selectionTo *cwin.Point
}

func createCursorFrame(e *editor, r cwin.Rect) cgame.Frame {
	var f cgame.Frame
	defaultTo := func(c, def cterm.Attribute) cterm.Attribute {
		if c == cterm.ColorDefault {
			return def
		}
		return c
	}
	r.Range(func(x, y int) {
		chx := e.winMap.GetClient(x, y)
		chx.Attr = cwin.Attr{
			Fg: defaultTo(chx.Attr.Bg, cterm.ColorBlack),
			Bg: defaultTo(chx.Attr.Fg, cterm.ColorWhite),
		}
		f = append(f, cgame.Cell{X: x - r.X, Y: y - r.Y, Chx: chx})
	})
	return f
}

func newCursor(e *editor, x, y int) *cursor {
	c := &cursor{
		SpriteBase: cgame.NewSpriteBase(e.g, e.winMap, cursorName,
			createCursorFrame(e, cwin.Rect{X: x, Y: y, W: 1, H: 1}), x, y),
		e: e,
	}
	c.move(0, 0) // this is to make sure cursor is visible.
	return c
}

func (c *cursor) reset() *cursor {
	r := c.Rect()
	c.Mgr().DeleteSprite(c)
	return newCursor(c.e, r.X, r.Y)
}

func (c *cursor) move(dx, dy int) {
	// First, make sure the new cursor rect is fully inside the map.
	r := c.Rect()
	newR := cwin.Rect{
		X: r.X + dx,
		Y: r.Y + dy,
		W: r.W,
		H: r.H,
	}
	if o, _ := c.e.winMap.ClientRect().ToOrigin().Overlap(newR); o != newR {
		return
	}
	c.Update(cgame.UpdateArg{
		DXY: &cwin.Point{
			X: dx,
			Y: dy,
		},
		F:   createCursorFrame(c.e, newR),
		IBC: cgame.InBoundsCheckNone,
		CD:  cgame.CollisionDetectionOff,
	})
	// Then try best to move e.winMap in its container e.winMapFrame so that cursor is
	// visible. Note it's not always possible to make the cursor entirely visible: if
	// the cursor is larger than e.winMapFrame's W or H or both. If that's the case,
	// we bias toward the visibility of the top left corner of the cursor rect.

	// Note given e.winMap is borderless, adding c.e.winMap.ClientRect().X/Y (which is
	// always 0) isn't strictly necessary, but rather for being logically sound.
	xOffsetFromFrame := c.e.winMap.Rect().X + c.e.winMap.ClientRect().X
	yOffsetFromFrame := c.e.winMap.Rect().Y + c.e.winMap.ClientRect().Y

	dx, dy = 0, 0
	if xOffsetFromFrame+newR.X < 0 {
		dx = -(xOffsetFromFrame + newR.X)
	} else if xOffsetFromFrame+newR.X+newR.W-1 >= c.e.winMapFrame.ClientRect().W {
		dx = c.e.winMapFrame.ClientRect().W - (xOffsetFromFrame + newR.X + newR.W)
	}
	if yOffsetFromFrame+newR.Y < 0 {
		dy = -(yOffsetFromFrame + newR.Y)
	} else if yOffsetFromFrame+newR.Y+newR.H-1 >= e.winMapFrame.ClientRect().H {
		dy = c.e.winMapFrame.ClientRect().H - (yOffsetFromFrame + newR.Y + newR.H)
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
