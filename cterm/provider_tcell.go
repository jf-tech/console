package cterm

import (
	"github.com/gdamore/tcell/v2"
)

type providerTCell struct {
	screen tcell.Screen
}

func (p *providerTCell) Init() error {
	s, err := tcell.NewScreen()
	if err != nil {
		return err
	}
	err = s.Init()
	if err != nil {
		return err
	}
	p.screen = s
	return nil
}

func (p *providerTCell) Close() {
	p.screen.Fini()
}

func (p *providerTCell) Flush() error {
	p.screen.Show()
	return nil
}

func (p *providerTCell) Sync() error {
	p.screen.Sync()
	return nil
}

func (p *providerTCell) Size() (int, int) {
	return p.screen.Size()
}

var (
	// to understand the translation, check out these 3 files:
	// termbox: https://github.com/nsf/termbox-go/blob/master/api_common.go
	// tcell: https://github.com/gdamore/tcell/blob/master/color.go
	// tcell to termbox compat: https://github.com/gdamore/tcell/blob/master/termbox/compat.go
	attrToTCell = map[Attribute]tcell.Color{
		ColorDefault:      tcell.ColorDefault,
		ColorBlack:        tcell.ColorBlack,
		ColorRed:          tcell.ColorMaroon,
		ColorGreen:        tcell.ColorGreen,
		ColorYellow:       tcell.ColorOlive,
		ColorBlue:         tcell.ColorNavy,
		ColorMagenta:      tcell.ColorPurple,
		ColorCyan:         tcell.ColorTeal,
		ColorWhite:        tcell.ColorSilver,
		ColorDarkGray:     tcell.ColorGray,
		ColorLightRed:     tcell.ColorRed,
		ColorLightGreen:   tcell.ColorLime,
		ColorLightYellow:  tcell.ColorYellow,
		ColorLightBlue:    tcell.ColorBlue,
		ColorLightMagenta: tcell.ColorFuchsia,
		ColorLightCyan:    tcell.ColorAqua,
		ColorLightGray:    tcell.ColorWhite,
	}
)

func (p *providerTCell) SetCell(x, y int, ch rune, fg, bg Attribute) {
	st := tcell.StyleDefault.Foreground(attrToTCell[fg]).Background(attrToTCell[bg])
	p.screen.SetContent(x, y, ch, nil, st)
}

var (
	keyFromTCell = map[tcell.Key]Key{
		tcell.KeyUp:    KeyArrowUp,
		tcell.KeyDown:  KeyArrowDown,
		tcell.KeyRight: KeyArrowRight,
		tcell.KeyLeft:  KeyArrowLeft,
		tcell.KeyEnter: KeyEnter,
		tcell.KeyEsc:   KeyEsc,
	}
)

func eventFromTCell(ev tcell.Event) Event {
	key, ok := ev.(*tcell.EventKey)
	if !ok {
		return Event{Type: EventNone}
	}
	if key.Key() == tcell.KeyRune {
		return Event{Type: EventKey, Ch: key.Rune()}
	}
	if k, ok := keyFromTCell[key.Key()]; ok {
		return Event{Type: EventKey, Key: k}
	}
	return Event{Type: EventKey}
}

func (p *providerTCell) PollEvent() Event {
	return eventFromTCell(p.screen.PollEvent())
}

func (p *providerTCell) Interrupt() {
	p.screen.PostEvent(tcell.NewEventInterrupt(nil))
}
