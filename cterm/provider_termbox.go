package cterm

import "github.com/nsf/termbox-go"

type providerTermBox struct{}

func (p *providerTermBox) Init() error {
	return termbox.Init()
}

func (p *providerTermBox) Close() {
	termbox.Close()
}

func (p *providerTermBox) Flush() error {
	return termbox.Flush()
}

func (p *providerTermBox) Size() (int, int) {
	return termbox.Size()
}

var (
	attrToTermBox = map[Attribute]termbox.Attribute{
		ColorDefault:      termbox.ColorDefault,
		ColorBlack:        termbox.ColorBlack,
		ColorRed:          termbox.ColorRed,
		ColorGreen:        termbox.ColorGreen,
		ColorYellow:       termbox.ColorYellow,
		ColorBlue:         termbox.ColorBlue,
		ColorMagenta:      termbox.ColorMagenta,
		ColorCyan:         termbox.ColorCyan,
		ColorWhite:        termbox.ColorWhite,
		ColorDarkGray:     termbox.ColorDarkGray,
		ColorLightRed:     termbox.ColorLightRed,
		ColorLightGreen:   termbox.ColorLightGreen,
		ColorLightYellow:  termbox.ColorLightYellow,
		ColorLightBlue:    termbox.ColorLightBlue,
		ColorLightMagenta: termbox.ColorLightMagenta,
		ColorLightCyan:    termbox.ColorLightCyan,
		ColorLightGray:    termbox.ColorLightGray,
	}
)

func (p *providerTermBox) SetCell(x, y int, ch rune, fg, bg Attribute) {
	termbox.SetCell(x, y, ch, attrToTermBox[fg], attrToTermBox[bg])
}

var (
	keyFromTermBox = map[termbox.Key]Key{
		termbox.KeyArrowUp:    KeyArrowUp,
		termbox.KeyArrowDown:  KeyArrowDown,
		termbox.KeyArrowRight: KeyArrowRight,
		termbox.KeyArrowLeft:  KeyArrowLeft,
		termbox.KeyEnter:      KeyEnter,
		termbox.KeyEsc:        KeyEsc,
	}
)

func eventFromTermBox(e termbox.Event) Event {
	if e.Type != termbox.EventKey {
		return Event{Type: EventNone}
	}
	if e.Ch != 0 {
		return Event{Type: EventKey, Ch: e.Ch}
	}
	if k, ok := keyFromTermBox[e.Key]; ok {
		return Event{Type: EventKey, Key: k}
	}
	if e.Key == termbox.KeySpace {
		return Event{Type: EventKey, Ch: ' '}
	}
	return Event{Type: EventKey}
}

func (p *providerTermBox) PollEvent() Event {
	return eventFromTermBox(termbox.PollEvent())
}

func (p *providerTermBox) Interrupt() {
	termbox.Interrupt()
}
