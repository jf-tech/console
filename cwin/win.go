package cwin

import (
	"fmt"
	"strings"

	"github.com/jf-tech/console/cterm"
	"github.com/jf-tech/go-corelib/maths"
)

type Align int

const (
	AlignLeft Align = iota
	AlignCenter
	AlignRight
)

const (
	RuneSpace rune = ' '
)

const (
	BorderRuneUL int = iota
	BorderRuneUR
	BorderRuneLR
	BorderRuneLL
	BorderRuneV
	BorderRuneH
	BorderRuneCount
)

type BorderRunes [BorderRuneCount]rune

var (
	SingleLineBorderRunes = BorderRunes{'┏', '┓', '┛', '┗', '┃', '━'}
	DoubleLineBorderRunes = BorderRunes{'╔', '╗', '╝', '╚', '║', '═'}
)

type ChAttr struct {
	Fg, Bg cterm.Attribute // cterm.ColorRed | ColorGreen | ...
}

type Chx struct {
	Ch   rune
	Attr ChAttr
}

var chxTransparent = Chx{}

func TransparentChx() Chx {
	return Chx{}
}

type WinCfg struct {
	// Required.
	// For root window, R.X/Y specify the absolute console x/y of the root window; W/H for its size
	// For regular window (that is direct/indirect descendant of root window), R.X/Y is the relative
	// position from its direct parent window. W/H for its size.
	R Rect
	// All the following are optional.
	EventHandler    EventHandler
	Name            string // also used as title (unless NoTitle or NoBorder is true)
	NoBorder        bool
	BorderRunes     *BorderRunes // if NoBorder && BorderRunes==nil, default to 1-line border
	BorderAttr      ChAttr
	ClientAttr      ChAttr
	NoTitle         bool // in case user sets Name (for debug purpose) and border, but doesn't want Title
	NoHPaddingTitle bool // most cases, have a one-space padding on each side of title looks nice
	NoHPaddingText  bool // most cases, have a one-space padding on each side of text block looks nice
}

// Win represents a window
type Win interface {
	Cfg() WinCfg

	UID() int64

	This() Win
	Parent() Win
	Prev() Win
	Next() Win
	ChildFirst() Win
	ChildLast() Win

	setParent(w Win)
	setPrev(w Win)
	setNext(w Win)
	setChildFirst(w Win)
	setChildLast(w Win)
	addNewChild(w Win)

	// Returns the Rect of the window relative to its parent window's client region.
	Rect() Rect
	// Returns the client Rect of the window, relative to this window.
	ClientRect() Rect

	// Move the window position (relative to its parent's client region) by (dx, dy)
	SetPosRel(dx, dy int)
	// Sets the event handler for this window.
	SetEventHandler(evHandler EventHandler)
	// Sets the window's title, with specific alignment.
	SetTitleAligned(align Align, format string, a ...interface{})
	// Sets the window's title, with default alignment.
	SetTitle(format string, a ...interface{})
	// Sets a text, multi-line allowed, to the client region of the window, with specific
	// alignment.
	SetTextAligned(align Align, format string, a ...interface{})
	// Sets a text, multi-line allowed, to the client region of the window, with default
	// alignment.
	SetText(format string, a ...interface{})

	// Sets the window to the top window among all the child windows of its parent window.
	SendToBottom(recursive bool)
	// Sets the window to the bottom window among all the child windows of its parent window.
	SendToTop(recursive bool)

	// Gets the Chx from the window.
	Get(x, y int) Chx
	// Gets the Chx from the client region of the window. Note cx/cy are client coordinate,
	// relative to this window's ClientRect.
	GetClient(cx, cy int) Chx
	// Puts a Chx to the client region of the window. Note cx/cy are client coordinate,
	// relative to this window's ClientRect.
	PutClient(cx, cy int, chx Chx)
	// Puts a rune to the client region of the window with default color attributes. Note cx/cy
	// are client coordinate, relative to this window's ClientRect.
	PutClientCh(cx, cy int, ch rune)
	// Fills a region inside the client region of the window. Note cr is the region Rect, relative
	// to this window's ClientRect.
	FillClient(cr Rect, chx Chx)

	fmt.Stringer
}

type WinBase struct {
	sys *Sys
	cfg WinCfg
	uid int64
	buf []Chx
	// Note this clientR.X/Y are the actual X/Y of the client region relative to
	// this window's Rect.
	clientR Rect

	parent, prev, next, childFirst, childLast Win
}

func (wb *WinBase) Cfg() WinCfg {
	return wb.cfg
}

func (wb *WinBase) UID() int64 {
	return wb.uid
}

func (wb *WinBase) This() Win {
	if w, ok := wb.sys.TryFindWin(wb.UID()); ok {
		return w
	}
	panic(fmt.Sprintf("forgot to register %s?", wb.String()))
}

func (wb *WinBase) Parent() Win {
	return wb.parent
}

func (wb *WinBase) Prev() Win {
	return wb.prev
}
func (wb *WinBase) Next() Win {
	return wb.next
}
func (wb *WinBase) ChildFirst() Win {
	return wb.childFirst
}

func (wb *WinBase) ChildLast() Win {
	return wb.childLast
}

func (wb *WinBase) setParent(w Win) {
	wb.parent = w
}

func (wb *WinBase) setPrev(w Win) {
	wb.prev = w
}

func (wb *WinBase) setNext(w Win) {
	wb.next = w
}

func (wb *WinBase) setChildFirst(w Win) {
	wb.childFirst = w
}

func (wb *WinBase) setChildLast(w Win) {
	wb.childLast = w
}

func (wb *WinBase) addNewChild(child Win) {
	if wb.ChildLast() == nil {
		wb.setChildFirst(child)
		wb.setChildLast(child)
	} else {
		child.setPrev(wb.ChildLast())
		wb.ChildLast().setNext(child)
		wb.setChildLast(child)
	}
}

func (wb *WinBase) Rect() Rect {
	return wb.cfg.R
}

func (wb *WinBase) ClientRect() Rect {
	return wb.clientR
}

func (wb *WinBase) SetPosRel(dx, dy int) {
	wb.cfg.R.X += dx
	wb.cfg.R.Y += dy
}

func (wb *WinBase) SetEventHandler(evHandler EventHandler) {
	wb.cfg.EventHandler = evHandler
}

func (wb *WinBase) SetTitleAligned(align Align, format string, a ...interface{}) {
	if wb.cfg.NoBorder {
		return
	}
	if wb.cfg.NoTitle {
		return
	}
	// this is needed in case we're setting a new title that is shorter than (or even
	// to none) than the previously one
	wb.fill(Rect{1, 0, wb.cfg.R.W - 2, 1}, Chx{wb.borderRunes()[BorderRuneH], wb.cfg.BorderAttr})
	padding := 1
	if wb.cfg.NoHPaddingTitle {
		padding = 0
	}
	if wb.clientR.W-2*padding <= 0 {
		return
	}
	t := []rune(fmt.Sprintf(format, a...))
	tlen := len(t)
	if tlen <= 0 {
		return
	}
	tlenActual := maths.MinInt(tlen+2*padding, wb.clientR.W)
	startX := wb.clientR.X
	switch align {
	case AlignCenter:
		startX += (wb.clientR.W - tlenActual) / 2
	case AlignRight:
		startX += wb.clientR.W - tlenActual
	}
	if padding > 0 {
		wb.put(startX, 0, Chx{RuneSpace, wb.cfg.BorderAttr})
		startX++
	}
	for i := 0; i < tlenActual-2*padding; i++ {
		wb.put(startX, 0, Chx{t[i], wb.cfg.BorderAttr})
		startX++
	}
	if padding > 0 {
		wb.put(startX, 0, Chx{RuneSpace, wb.cfg.BorderAttr})
		startX++
	}
}

func (wb *WinBase) SetTitle(format string, a ...interface{}) {
	wb.SetTitleAligned(AlignLeft, format, a...)
}

func (wb *WinBase) SetTextAligned(align Align, format string, a ...interface{}) {
	wb.FillClient(wb.ClientRect().ToOrigin(), Chx{Ch: RuneSpace, Attr: wb.cfg.ClientAttr})
	lines := strings.Split(fmt.Sprintf(format, a...), "\n")
	for cy := 0; cy < maths.MinInt(wb.clientR.H, len(lines)); cy++ {
		wb.setTextLine(cy, lines[cy], align, wb.cfg.ClientAttr)
	}
}

func (wb *WinBase) SetText(format string, a ...interface{}) {
	wb.SetTextAligned(AlignLeft, format, a...)
}

func (wb *WinBase) SendToBottom(recursive bool) {
	parent := wb.parent
	if parent == nil {
		return
	}
	wb.removeFromParent()
	wb.setParent(parent)
	wb.setNext(parent.ChildFirst())
	this := wb.This()
	if parent.ChildFirst() != nil {
		parent.ChildFirst().setPrev(this)
	}
	parent.setChildFirst(this)
	if parent.ChildLast() == nil {
		parent.setChildLast(this)
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
	wb.removeFromParent()
	wb.setParent(parent)
	wb.setPrev(parent.ChildLast())
	this := wb.This()
	if parent.ChildLast() != nil {
		parent.ChildLast().setNext(this)
	}
	parent.setChildLast(this)
	if parent.ChildFirst() == nil {
		parent.setChildFirst(this)
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
	return fmt.Sprintf("win['%s'|%d|%s]", wb.cfg.Name, wb.UID(), wb.Rect())
}

func (wb *WinBase) borderRunes() *BorderRunes {
	if wb.cfg.BorderRunes != nil {
		return wb.cfg.BorderRunes
	}
	return &SingleLineBorderRunes
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

func (wb *WinBase) setTextLine(cy int, line string, align Align, attr ChAttr) {
	wb.FillClient(Rect{X: 0, Y: cy, W: wb.clientR.W, H: 1}, Chx{Ch: RuneSpace, Attr: attr})
	padding := 1
	if wb.cfg.NoHPaddingText {
		padding = 0
	}
	if wb.clientR.W-2*padding <= 0 {
		return
	}
	l := []rune(line)
	llen := len(l)
	if llen <= 0 {
		return
	}
	llenActual := maths.MinInt(llen+2*padding, wb.clientR.W)
	startCX := 0
	switch align {
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

func (wb *WinBase) putBorder() {
	if wb.cfg.NoBorder {
		return
	}
	if wb.cfg.R.W < 2 || wb.cfg.R.H < 2 {
		return
	}
	borderRunes := wb.borderRunes()
	// UL
	wb.put(0, 0, Chx{borderRunes[BorderRuneUL], wb.cfg.BorderAttr})
	// UR
	wb.put(wb.cfg.R.W-1, 0, Chx{borderRunes[BorderRuneUR], wb.cfg.BorderAttr})
	// LR
	wb.put(wb.cfg.R.W-1, wb.cfg.R.H-1, Chx{borderRunes[BorderRuneLR], wb.cfg.BorderAttr})
	// LL
	wb.put(0, wb.cfg.R.H-1, Chx{borderRunes[BorderRuneLL], wb.cfg.BorderAttr})
	// top/bottom horizontal lines
	wb.fill(Rect{1, 0, wb.cfg.R.W - 2, 1}, Chx{borderRunes[BorderRuneH], wb.cfg.BorderAttr})
	wb.fill(Rect{1, wb.cfg.R.H - 1, wb.cfg.R.W - 2, 1}, Chx{borderRunes[BorderRuneH], wb.cfg.BorderAttr})
	// left/right vertical lines
	wb.fill(Rect{0, 1, 1, wb.cfg.R.H - 2}, Chx{borderRunes[BorderRuneV], wb.cfg.BorderAttr})
	wb.fill(Rect{wb.cfg.R.W - 1, 1, 1, wb.cfg.R.H - 2}, Chx{borderRunes[BorderRuneV], wb.cfg.BorderAttr})
}

func (wb *WinBase) removeFromParent() {
	if wb.Parent() == nil {
		return
	}
	prev := wb.Prev()
	next := wb.Next()
	if prev != nil {
		prev.setNext(next)
	}
	if next != nil {
		next.setPrev(prev)
	}
	if wb.Parent().ChildFirst().UID() == wb.UID() {
		wb.Parent().setChildFirst(next)
	}
	if wb.Parent().ChildLast().UID() == wb.UID() {
		wb.Parent().setChildLast(prev)
	}
	wb.setParent(nil)
	wb.setPrev(nil)
	wb.setNext(nil)
}

func NewWinBase(sys *Sys, parent Win, c WinCfg) *WinBase {
	cw := &WinBase{sys: sys, cfg: c, uid: GenUID(), parent: parent}
	cw.clientR = Rect{0, 0, cw.cfg.R.W, cw.cfg.R.H}
	if !cw.cfg.NoBorder {
		cw.clientR.X++
		cw.clientR.Y++
		cw.clientR.W -= 2
		cw.clientR.H -= 2
	}
	cw.buf = make([]Chx, cw.cfg.R.W*cw.cfg.R.H)
	cw.putBorder()
	cw.fill(cw.clientR, Chx{RuneSpace, cw.cfg.ClientAttr})
	if len(cw.cfg.Name) > 0 {
		cw.SetTitle(cw.cfg.Name)
	}
	return cw
}

func (s *Sys) CreateWin(parent Win, cfg WinCfg) Win {
	if parent == nil {
		parent = s.SysWin()
	}
	wb := NewWinBase(s, parent, cfg)
	s.RegWin(wb)
	parent.addNewChild(wb)
	return wb
}
