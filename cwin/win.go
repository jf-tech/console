package cwin

import (
	"fmt"
	"strings"

	"github.com/jf-tech/go-corelib/maths"
	"github.com/nsf/termbox-go"
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
)

type ChAttr struct {
	Fg, Bg termbox.Attribute // termbox.ColorRed | ColorGreen | ...
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
	Name            string // used as win title (unless SetTitle called after creation) and in debug dump
	NoBorder        bool
	BorderRunes     *BorderRunes // if NoBorder && BorderRunes==nil, default to 1-line border
	BorderAttr      ChAttr
	ClientAttr      ChAttr
	NoTitle         bool // in case user sets Name (for debug purpose) and border, but doesn't want actual Title
	NoHPaddingTitle bool // most cases, have a one-space padding on each side of title looks nice
	NoHPaddingText  bool // most cases, have a one-space padding on each side of text block looks nice
	StartHidden     bool
}

func hidden_s(hidden bool) string {
	if hidden {
		return "H" // hidden
	}
	return "V" // visible
}

type Win struct {
	cfg WinCfg

	hidden  bool
	clientR Rect
	buf     []Chx

	parent, next, prev, child1, childn *Win
}

func winName(w *Win) string {
	if w == nil {
		return "<nil>"
	}
	return w.cfg.Name
}

func (w *Win) String() string {
	return fmt.Sprintf("Win['%s',%s,%s]", winName(w), w.cfg.R, hidden_s(w.hidden))
}

func (w *Win) Parent() *Win {
	return w.parent
}

func (w *Win) bufIdx(x, y int) int {
	return y*w.cfg.R.W + x
}

func (w *Win) put(x, y int, chx Chx) {
	w.buf[w.bufIdx(x, y)] = chx
}

func (w *Win) PutClient(cx, cy int, chx Chx) {
	w.put(w.clientR.X+cx, w.clientR.Y+cy, chx)
}

func (w *Win) PutClientCh(cx, cy int, ch rune) {
	w.PutClient(cx, cy, Chx{ch, w.cfg.ClientAttr})
}

// Note we don't have a fillCh because it's possible the rect cross both
// the border and the client region, thus cannot default to either
// Win.cfg.BorderAttr or Win.cfg.TextAttr. User of Win has to make an
// explicit decision.
func (w *Win) fill(r Rect, chx Chx) {
	for y := 0; y < r.H; y++ {
		for x := 0; x < r.W; x++ {
			w.put(r.X+x, r.Y+y, chx)
		}
	}
}

func (w *Win) FillClient(cr Rect, chx Chx) {
	for y := 0; y < cr.H; y++ {
		for x := 0; x < cr.W; x++ {
			w.PutClient(cr.X+x, cr.Y+y, chx)
		}
	}
}

func (w *Win) SetHidden(hidden bool) {
	w.hidden = hidden
}

func (w *Win) Rect() Rect {
	return w.cfg.R
}

func (w *Win) ClientRect() Rect {
	return w.clientR
}

func (w *Win) SetPosAbs(x, y int) {
	w.cfg.R.X = x
	w.cfg.R.Y = y
}

func (w *Win) SetPosRelative(dx, dy int) {
	w.SetPosAbs(w.cfg.R.X+dx, w.cfg.R.Y+dy)
}

func (w *Win) putBorder() {
	if w.cfg.NoBorder {
		return
	}
	if w.cfg.R.W < 2 || w.cfg.R.H < 2 {
		return
	}
	borderRunes := w.cfg.BorderRunes
	if borderRunes == nil {
		borderRunes = &SingleLineBorderRunes
	}
	// UL
	w.put(0, 0, Chx{borderRunes[BorderRuneUL], w.cfg.BorderAttr})
	// UR
	w.put(w.cfg.R.W-1, 0, Chx{borderRunes[BorderRuneUR], w.cfg.BorderAttr})
	// LR
	w.put(w.cfg.R.W-1, w.cfg.R.H-1, Chx{borderRunes[BorderRuneLR], w.cfg.BorderAttr})
	// LL
	w.put(0, w.cfg.R.H-1, Chx{borderRunes[BorderRuneLL], w.cfg.BorderAttr})
	// top/bottom horizontal lines
	w.fill(Rect{1, 0, w.cfg.R.W - 2, 1}, Chx{borderRunes[BorderRuneH], w.cfg.BorderAttr})
	w.fill(Rect{1, w.cfg.R.H - 1, w.cfg.R.W - 2, 1}, Chx{borderRunes[BorderRuneH], w.cfg.BorderAttr})
	// left/right vertical lines
	w.fill(Rect{0, 1, 1, w.cfg.R.H - 2}, Chx{borderRunes[BorderRuneV], w.cfg.BorderAttr})
	w.fill(Rect{w.cfg.R.W - 1, 1, 1, w.cfg.R.H - 2}, Chx{borderRunes[BorderRuneV], w.cfg.BorderAttr})
}

func (w *Win) SetTitle(title string, align Align) {
	if w.cfg.NoBorder {
		return
	}
	if w.cfg.NoTitle {
		return
	}
	padding := 1
	if w.cfg.NoHPaddingTitle {
		padding = 0
	}
	if w.clientR.W-2*padding <= 0 {
		return
	}
	t := []rune(title)
	tlen := len(t)
	if tlen <= 0 {
		return
	}
	tlenActual := maths.MinInt(tlen+2*padding, w.clientR.W)
	startX := w.clientR.X
	switch align {
	case AlignCenter:
		startX += (w.clientR.W - tlenActual) / 2
	case AlignRight:
		startX += w.clientR.W - tlenActual
	}
	// this is needed in case we're setting a new title that is shorter than the previously one
	w.putBorder()
	if padding > 0 {
		w.put(startX, 0, Chx{RuneSpace, w.cfg.BorderAttr})
		startX++
	}
	for i := 0; i < tlenActual-2*padding; i++ {
		w.put(startX, 0, Chx{t[i], w.cfg.BorderAttr})
		startX++
	}
	if padding > 0 {
		w.put(startX, 0, Chx{RuneSpace, w.cfg.BorderAttr})
		startX++
	}
}

func (w *Win) SetTextAligned(align Align, format string, a ...interface{}) {
	padding := 1
	if w.cfg.NoHPaddingText {
		padding = 0
	}
	if w.clientR.W-2*padding <= 0 {
		return
	}
	w.fill(w.clientR, Chx{Ch: RuneSpace, Attr: w.cfg.ClientAttr})
	lines := strings.Split(fmt.Sprintf(format, a...), "\n")
	for y := 0; y < maths.MinInt(w.clientR.H, len(lines)); y++ {
		rline := []rune(lines[y])
		if len(rline) <= 0 {
			continue
		}
		for x := 0; x < maths.MinInt(len(rline), w.clientR.W-2*padding); x++ {
			w.PutClientCh(padding+x, y, rline[x])
		}
	}
}

func (w *Win) SetText(format string, a ...interface{}) {
	w.SetTextAligned(AlignLeft, format, a...)
}

func (w *Win) Dump() string {
	return fmt.Sprintf("%s:clientR(%s),par:'%s',next:'%s',prev:'%s',c1:'%s',cn:'%s'",
		w, w.clientR, w.parent, w.next, w.prev, w.child1, w.childn)
}

func (w *Win) DumpTree(indent int) string {
	s := strings.Repeat("-", indent) + w.Dump() + "\n"
	for child := w.child1; child != nil; child = child.next {
		s += child.DumpTree(indent+2) + "\n"
	}
	return s
}

func newWin(parent *Win, c WinCfg) *Win {
	cw := &Win{cfg: c, parent: parent}
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
		cw.SetTitle(cw.cfg.Name, AlignLeft)
	}
	return cw
}
