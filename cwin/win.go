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
	EventHandler       EventHandler
	Name               string // also used as title (unless NoTitle or NoBorder is true)
	NoBorder           bool
	BorderRunes        *BorderRunes // if NoBorder && BorderRunes==nil, default to 1-line border
	InFocusBorderRunes *BorderRunes // Used for border rendering when window is in focus. If nil, default to 2-line border
	BorderAttr         Attr
	InFocusBorderAttr  *Attr // Used for border rendering when window is in focus. If nil, default Fg=ColorLightYellow
	ClientAttr         Attr
	NoTitle            bool // in case user sets Name (for debug purpose) and border, but doesn't want Title
	NoPaddingTitle     bool // most cases, have a one-space padding on each side of title looks nice
	NoPaddingText      bool // most cases, have a one-space padding on each side of text block looks nice
	TitleAlign         Align
	TextAlign          Align
}

// Win represents a window
type Win interface {
	Cfg() WinCfg

	Sys() *Sys

	// Base returns the embedded WinBase pointer which can be used for accessing
	// the library built-in WinBase functionalities, as well as serving as a unique
	// identifier for the window.
	Base() *WinBase
	// Same tells if this Win is the same instance of other Win.
	Same(other Win) bool
	// This returns the actual object that implements Win interface that is registered
	// with Sys. Because WinBase implements Win interface, sometimes we go into situation
	// where a WinBase pointer is getting passed around but eventually when deverloper
	// wants to cast back to their own object (which embeds WinBase) they get type assertion
	// failure. As long as the object (that implements Win) passed into Sys.CreateWin is
	// the "top-level" object, This() will always return that registered object. Note if
	// calling This() on a non Sys managed (i.e. not from Sys.CreateWin) object, it will
	// panic.
	This() Win
	Parent() Win
	Prev() Win
	Next() Win
	ChildFirst() Win
	ChildLast() Win

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
	// Sets the window's title.
	SetTitle(format string, a ...interface{})
	// Sets a text, multi-line allowed, to the client region of the window
	SetText(format string, a ...interface{})
	// Sets a single line text, to the client region of the window
	SetLine(cy int, attr Attr, format string, a ...interface{})

	// Set the window to be in or out of focus. This method should only conduct focus related
	// changes to this window only. To actually set a window in focus system-wise needs to do
	// more things like set up message handler, changing z-order, etc, which requires a calling
	// to Sys.SetFocus.
	SetFocus(focused bool)

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
