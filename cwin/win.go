package cwin

import (
	"fmt"

	"github.com/jf-tech/console/cterm"
)

type Align int

const (
	AlignLeft Align = iota
	AlignCenter
	AlignRight
)

type Attr struct {
	Fg, Bg cterm.Attribute // cterm.ColorRed | ColorGreen | ...
}

type Chx struct {
	Ch   rune
	Attr Attr
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
	BorderAttr      Attr
	ClientAttr      Attr
	NoTitle         bool // in case user sets Name (for debug purpose) and border, but doesn't want Title
	NoHPaddingTitle bool // most cases, have a one-space padding on each side of title looks nice
	NoHPaddingText  bool // most cases, have a one-space padding on each side of text block looks nice
}

// Win represents a window
type Win interface {
	Cfg() WinCfg

	Sys() *Sys

	Base() *WinBase
	This() Win
	Parent() Win
	Prev() Win
	Next() Win
	ChildFirst() Win
	ChildLast() Win

	Same(other Win) bool

	// Returns the Rect of the window relative to its parent window's client region.
	Rect() Rect
	// Returns the client Rect of the window, relative to this window.
	ClientRect() Rect

	// Moves the window position (relative to its parent's client region) to (x, y)
	SetPosAbs(x, y int)
	// Moves the window position (relative to its parent's client region) by (dx, dy)
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
	// Sets a single line text, to the client region of the window, with specific alignment
	// and color attributes
	SetLineAligned(cy int, align Align, attr Attr, format string, a ...interface{})
	// Sets a single line text, to the client region of the window, with default  alignment
	// and default color attributes
	SetLine(cy int, format string, a ...interface{})

	// Sets the window to the bottom position among all the child windows of its parent window.
	SendToBottom(recursive bool)
	// Sets the window to the top position among all the child windows of its parent window.
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
