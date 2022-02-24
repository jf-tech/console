package cwin

import (
	"fmt"
	"strings"
	"unsafe"

	"github.com/jf-tech/console/cterm"
	"github.com/jf-tech/go-corelib/maths"
)

type WinBase struct {
	cfg WinCfg
	sys *Sys

	buf []Chx

	// Note this clientR.X/Y are the actual X/Y of the client region relative to
	// this window's Rect.
	clientR Rect
	title   string
	focused bool

	parent, prev, next, childFirst, childLast *WinBase
}

func (wb *WinBase) Cfg() WinCfg {
	return wb.cfg
}

func (wb *WinBase) Sys() *Sys {
	return wb.sys
}

func (wb *WinBase) Base() *WinBase {
	return wb
}

func (wb *WinBase) Same(other Win) bool {
	return wb == other.Base()
}

func (wb *WinBase) This() Win {
	return wb.Sys().Find(wb)
}

func (wb *WinBase) Parent() Win {
	return wb.winBaseToWin(wb.parent)
}

func (wb *WinBase) Prev() Win {
	return wb.winBaseToWin(wb.prev)
}
func (wb *WinBase) Next() Win {
	return wb.winBaseToWin(wb.next)
}
func (wb *WinBase) ChildFirst() Win {
	return wb.winBaseToWin(wb.childFirst)
}

func (wb *WinBase) ChildLast() Win {
	return wb.winBaseToWin(wb.childLast)
}

func (wb *WinBase) Rect() Rect {
	return wb.cfg.R
}

func (wb *WinBase) ClientRect() Rect {
	return wb.clientR
}

func (wb *WinBase) SetPosAbs(x, y int) {
	wb.cfg.R.X = x
	wb.cfg.R.Y = y
}

func (wb *WinBase) SetPosRel(dx, dy int) {
	wb.cfg.R.X += dx
	wb.cfg.R.Y += dy
}

func (wb *WinBase) SetEventHandler(evHandler EventHandler) {
	wb.cfg.EventHandler = evHandler
}

func (wb *WinBase) SetTitle(format string, a ...interface{}) {
	if wb.cfg.NoBorder {
		return
	}
	if wb.cfg.NoTitle {
		return
	}
	wb.title = fmt.Sprintf(format, a...)
	// this is needed in case we're setting a new title that is shorter than (or even
	// to none) than the previously one
	wb.fill(Rect{1, 0, wb.cfg.R.W - 2, 1}, Chx{wb.borderRunes()[BorderRuneH], wb.borderAttr()})
	padding := 1
	if wb.cfg.NoPaddingTitle {
		padding = 0
	}
	if wb.clientR.W-2*padding <= 0 {
		return
	}
	t := []rune(wb.title)
	tlen := len(t)
	if tlen <= 0 {
		return
	}
	tlenActual := maths.MinInt(tlen+2*padding, wb.clientR.W)
	startX := wb.clientR.X
	switch wb.cfg.TitleAlign {
	case AlignCenter:
		startX += (wb.clientR.W - tlenActual) / 2
	case AlignRight:
		startX += wb.clientR.W - tlenActual
	}
	if padding > 0 {
		wb.put(startX, 0, Chx{RuneSpace, wb.borderAttr()})
		startX++
	}
	for i := 0; i < tlenActual-2*padding; i++ {
		wb.put(startX, 0, Chx{t[i], wb.borderAttr()})
		startX++
	}
	if padding > 0 {
		wb.put(startX, 0, Chx{RuneSpace, wb.borderAttr()})
		startX++
	}
}

func (wb *WinBase) SetText(format string, a ...interface{}) {
	wb.FillClient(wb.ClientRect().ToOrigin(), Chx{Ch: RuneSpace, Attr: wb.cfg.ClientAttr})
	lines := strings.Split(fmt.Sprintf(format, a...), "\n")
	for cy := 0; cy < maths.MinInt(wb.clientR.H, len(lines)); cy++ {
		wb.SetLine(cy, wb.cfg.ClientAttr, lines[cy])
	}
}

func (wb *WinBase) SetLine(cy int, attr Attr, format string, a ...interface{}) {
	wb.FillClient(Rect{X: 0, Y: cy, W: wb.clientR.W, H: 1}, Chx{Ch: RuneSpace, Attr: attr})
	padding := 1
	if wb.cfg.NoPaddingText {
		padding = 0
	}
	if wb.clientR.W-2*padding <= 0 {
		return
	}
	l := []rune(fmt.Sprintf(format, a...))
	llen := len(l)
	if llen <= 0 {
		return
	}
	llenActual := maths.MinInt(llen+2*padding, wb.clientR.W)
	startCX := 0
	switch wb.cfg.TextAlign {
	case AlignCenter:
		startCX += (wb.clientR.W - llenActual) / 2
	case AlignRight:
		startCX += wb.clientR.W - llenActual
	}
	if padding > 0 {
		wb.PutClient(startCX, cy, Chx{RuneSpace, attr})
		startCX++
	}
	for i := 0; i < llenActual-2*padding; i++ {
		wb.PutClient(startCX, cy, Chx{l[i], attr})
		startCX++
	}
	if padding > 0 {
		wb.PutClient(startCX, cy, Chx{RuneSpace, attr})
		startCX++
	}
}

func (wb *WinBase) SetFocus(focused bool) {
	wb.focused = focused
	wb.renderBorderAndTitle()
}

func (wb *WinBase) SendToBottom(recursive bool) {
	parent := wb.parent
	if parent == nil {
		return
	}
	parent.removeChild(wb)
	wb.parent = parent
	wb.next = parent.childFirst
	if parent.childFirst != nil {
		parent.childFirst.prev = wb
	}
	parent.childFirst = wb
	if parent.childLast == nil {
		parent.childLast = wb
	}
	if recursive {
		parent.SendToBottom(recursive)
	}
}

func (wb *WinBase) SendToTop(recursive bool) {
	parent := wb.parent
	if parent == nil {
		return
	}
	parent.removeChild(wb)
	wb.parent = parent
	wb.prev = parent.childLast
	if parent.childLast != nil {
		parent.childLast.next = wb
	}
	parent.childLast = wb
	if parent.childFirst == nil {
		parent.childFirst = wb
	}
	if recursive {
		parent.SendToTop(recursive)
	}
}

func (wb *WinBase) Get(x, y int) Chx {
	return wb.buf[wb.bufIdx(x, y)]
}

func (wb *WinBase) GetClient(cx, cy int) Chx {
	return wb.Get(wb.clientR.X+cx, wb.clientR.Y+cy)
}

func (wb *WinBase) PutClient(cx, cy int, chx Chx) {
	wb.put(wb.clientR.X+cx, wb.clientR.Y+cy, chx)
}

func (wb *WinBase) PutClientCh(cx, cy int, ch rune) {
	wb.PutClient(cx, cy, Chx{ch, wb.cfg.ClientAttr})
}

func (wb *WinBase) FillClient(cr Rect, chx Chx) {
	for y := 0; y < cr.H; y++ {
		for x := 0; x < cr.W; x++ {
			wb.PutClient(cr.X+x, cr.Y+y, chx)
		}
	}
}

func (wb *WinBase) String() string {
	return fmt.Sprintf("WinBase['%s'|0x%X|%s]", wb.cfg.Name, uintptr(unsafe.Pointer(wb)), wb.Rect())
}

func (*WinBase) winBaseToWin(b *WinBase) Win {
	if b == nil {
		return nil
	}
	return b.This()
}

func (wb *WinBase) addChild(child *WinBase) {
	child.prev = wb.childLast
	child.next = nil
	if wb.childLast == nil {
		wb.childFirst = child
		wb.childLast = child
	} else {
		wb.childLast.next = child
		wb.childLast = child
	}
}

func (wb *WinBase) removeChild(child *WinBase) {
	childPrev := child.prev
	childNext := child.next
	if childPrev != nil {
		childPrev.next = childNext
	}
	if childNext != nil {
		childNext.prev = childPrev
	}
	if wb.childFirst == child {
		wb.childFirst = childNext
	}
	if wb.childLast == child {
		wb.childLast = childPrev
	}
	child.parent, child.prev, child.next = nil, nil, nil
}

func (wb *WinBase) borderRunes() *BorderRunes {
	if wb.focused {
		if wb.cfg.InFocusBorderRunes != nil {
			return wb.cfg.InFocusBorderRunes
		}
		if wb.cfg.BorderRunes != nil {
			return wb.cfg.BorderRunes
		}
		return &DoubleLineBorderRunes
	} else {
		if wb.cfg.BorderRunes != nil {
			return wb.cfg.BorderRunes
		}
		return &SingleLineBorderRunes
	}
}

func (wb *WinBase) borderAttr() Attr {
	if wb.focused {
		if wb.cfg.InFocusBorderAttr != nil {
			return *wb.cfg.InFocusBorderAttr
		}
		return Attr{Fg: cterm.ColorLightYellow}
	}
	return wb.cfg.BorderAttr
}

func (wb *WinBase) bufIdx(x, y int) int {
	return y*wb.cfg.R.W + x
}

func (wb *WinBase) put(x, y int, chx Chx) {
	wb.buf[wb.bufIdx(x, y)] = chx
}

// Note we don't have a fillCh because it's possible the rect cross both
// the border and the client region, thus cannot default to either
// wb.cfg.BorderAttr or wb.cfg.ClientAttr.
func (wb *WinBase) fill(r Rect, chx Chx) {
	for y := 0; y < r.H; y++ {
		for x := 0; x < r.W; x++ {
			wb.put(r.X+x, r.Y+y, chx)
		}
	}
}

func (wb *WinBase) renderBorderAndTitle() {
	if wb.cfg.NoBorder {
		return
	}
	if wb.cfg.R.W < 2 || wb.cfg.R.H < 2 {
		return
	}
	borderRunes := wb.borderRunes()
	borderAttr := wb.borderAttr()
	// UL
	wb.put(0, 0, Chx{borderRunes[BorderRuneUL], borderAttr})
	// UR
	wb.put(wb.cfg.R.W-1, 0, Chx{borderRunes[BorderRuneUR], borderAttr})
	// LR
	wb.put(wb.cfg.R.W-1, wb.cfg.R.H-1, Chx{borderRunes[BorderRuneLR], borderAttr})
	// LL
	wb.put(0, wb.cfg.R.H-1, Chx{borderRunes[BorderRuneLL], borderAttr})
	// top/bottom horizontal lines
	wb.fill(Rect{1, 0, wb.cfg.R.W - 2, 1}, Chx{borderRunes[BorderRuneH], borderAttr})
	wb.fill(Rect{1, wb.cfg.R.H - 1, wb.cfg.R.W - 2, 1}, Chx{borderRunes[BorderRuneH], borderAttr})
	// left/right vertical lines
	wb.fill(Rect{0, 1, 1, wb.cfg.R.H - 2}, Chx{borderRunes[BorderRuneV], borderAttr})
	wb.fill(Rect{wb.cfg.R.W - 1, 1, 1, wb.cfg.R.H - 2}, Chx{borderRunes[BorderRuneV], borderAttr})
	wb.SetTitle(wb.title)
}

func NewWinBase(sys *Sys, parent Win, c WinCfg) *WinBase {
	cw := &WinBase{sys: sys, cfg: c, title: c.Name, focused: false}
	if parent != nil {
		cw.parent = parent.Base()
	}
	cw.clientR = Rect{0, 0, cw.cfg.R.W, cw.cfg.R.H}
	if !cw.cfg.NoBorder {
		cw.clientR.X++
		cw.clientR.Y++
		cw.clientR.W -= 2
		cw.clientR.H -= 2
	}
	cw.buf = make([]Chx, cw.cfg.R.W*cw.cfg.R.H)
	cw.renderBorderAndTitle()
	cw.fill(cw.clientR, Chx{RuneSpace, cw.cfg.ClientAttr})
	return cw
}
