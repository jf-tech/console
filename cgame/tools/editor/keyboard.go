package main

import (
	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cterm"
	"github.com/jf-tech/console/cwin"
)

func (e *editor) handleKeyArrow(key cterm.Key) cwin.EventResponse {
	delta := cwin.DirOffSetXY[cwin.Key2Dir[key]]
	if e.selection == nil {
		e.moveCursor(delta.X, delta.Y)
		return cwin.EventHandled
	}
	e.selection.move(delta)
	r := e.selection.currentRect()
	e.removeCursor()
	e.cursor = cgame.NewSpriteBase(e.g, e.winMap, cursorName, e.createCursorFrame(r), r.X, r.Y)
	e.moveCursor(0, 0) // needed for winMap shifts to keep cursor always visible. // TODO can you always?
	return cwin.EventHandled
}

func (e *editor) handleKeyHomeEnd(key cterm.Key) cwin.EventResponse {
	e.abortSelection()
	dx := -e.cursor.Rect().X
	if key == cterm.KeyCtrlE {
		dx += e.winMap.ClientRect().W - 1
	}
	e.moveCursor(dx, 0)
	return cwin.EventHandled
}

func (e *editor) handleKeyCopyPaste(key cterm.Key) cwin.EventResponse {
	if key == cterm.KeyCtrlC {
		e.copy()
		e.abortSelection() // first copy selected stuff, then remove selection.
	} else {
		e.abortSelection() // for paste, remove selection if any, then paste.
		e.paste()
	}
	return cwin.EventHandled
}

func (e *editor) handleKeyDelete(key cterm.Key) cwin.EventResponse {
	del := func(r cwin.Rect) {
		for y := 0; y < r.H; y++ {
			for x := 0; x < r.W; x++ {
				e.winMap.PutClient(r.X+x, r.Y+y, cwin.TransparentChx())
			}
		}
		e.moveCursor(0, 0)
	}
	if key == cterm.KeyCtrlK {
		e.abortSelection() // ctrl+k doesn't support selection.
		del(cwin.Rect{
			X: e.cursor.Rect().X,
			Y: e.cursor.Rect().Y,
			W: e.winMap.ClientRect().W - e.cursor.Rect().X,
			H: 1,
		})
		return cwin.EventHandled
	}
	// KeyBackspace2
	if e.selection == nil {
		if e.cursor.Rect().X == e.winMap.ClientRect().W-1 {
			if e.winMap.GetClient(e.cursor.Rect().X, e.cursor.Rect().Y) != cwin.TransparentChx() {
				del(cwin.Rect{X: e.cursor.Rect().X, Y: e.cursor.Rect().Y, W: 1, H: 1})
				return cwin.EventHandled
			}
		}
		if e.cursor.Rect().X > 0 {
			del(cwin.Rect{X: e.cursor.Rect().X - 1, Y: e.cursor.Rect().Y, W: 1, H: 1})
			e.moveCursor(-1, 0)
		}
	} else {
		del(e.selection.currentRect())
		e.abortSelection()
		e.resetCursor()
	}
	return cwin.EventHandled
}

func (e *editor) handleKeyEnter() cwin.EventResponse {
	e.abortSelection()
	e.moveCursor(-e.cursor.Rect().X, 1)
	return cwin.EventHandled
}

func (e *editor) handleSelection() cwin.EventResponse {
	if e.selection != nil {
		return cwin.EventHandled
	}
	e.createSelection()
	return cwin.EventHandled
}

func (e *editor) handleKeyFocusChange() cwin.EventResponse {
	if e.g.WinSys.GetFocused().Same(e.winMain) {
		e.g.WinSys.SetFocus(e.winToolbox)
	} else {
		e.g.WinSys.SetFocus(e.winMain)
	}
	return cwin.EventHandled
}

func (e *editor) handleRunePress(ch rune) cwin.EventResponse {
	if ch == cwin.RuneSpace {
		e.putAtCursor(cwin.TransparentChx())
	} else {
		e.putAtCursor(cwin.Chx{Ch: ch})
	}
	e.abortSelection()
	e.moveCursor(1, 0)
	return cwin.EventHandled
}

func (e *editor) handleKeyDebug() cwin.EventResponse {
	e.debug()
	return cwin.EventHandled
}
