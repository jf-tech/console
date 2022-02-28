package main

import (
	"fmt"
	"strings"

	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cgame/assets"
	"github.com/jf-tech/console/cterm"
	"github.com/jf-tech/console/cwin"
)

type editorMode int

const (
	editorModeNew = editorMode(iota)
	editorModeExisting
)

var (
	instrW          = 40
	instrH          = 20
	instrClientAttr = cwin.Attr{Bg: cterm.ColorBlue}

	mapFrameClientChx = cwin.Chx{Ch: '░', Attr: cwin.Attr{Fg: cterm.ColorDarkGray}}

	cursorName = "cursor"
	cursorAttr = cwin.Attr{Fg: cterm.ColorDarkGray}
)

type editor struct {
	g *cgame.Game

	winMain     cwin.Win
	winRulers   *rulers
	winMapFrame cwin.Win
	winMap      cwin.Win
	winToolbox  cwin.Win
	winInstr    cwin.Win

	cursor *cgame.SpriteBase

	selection *selection
}

func (e *editor) main(mode editorMode) {
	var err error
	e.g, err = cgame.Init(cterm.TermBox)
	if err != nil {
		panic(err)
	}
	defer e.g.Close()

	e.setup(mode)

	e.g.Run(assets.GameOverKeys, nil, func(ev cterm.Event) cwin.EventResponse {
		if ev.Type == cterm.EventKey {
			switch ev.Key {

			case cterm.KeyArrowUp, cterm.KeyArrowRight, cterm.KeyArrowDown, cterm.KeyArrowLeft:
				return e.handleKeyArrow(ev.Key)

			case cterm.KeyCtrlA, cterm.KeyCtrlE:
				return e.handleKeyHomeEnd(ev.Key)

			case cterm.KeyCtrlC, cterm.KeyCtrlV:
				return e.handleKeyCopyPaste(ev.Key)

			case cterm.KeyCtrlK, cterm.KeyBackspace2:
				return e.handleKeyDelete(ev.Key)

			case cterm.KeyEnter:
				return e.handleKeyEnter()

			case cterm.KeyF3:
				return e.handleSelection()

			case cterm.KeyF10:
				return e.handleKeyFocusChange()

			case cterm.KeyCtrlD:
				return e.handleKeyDebug()
			}

			if ev.Ch != 0 {
				return e.handleRunePress(ev.Ch)
			}
		}
		return cwin.EventNotHandled
	})
}

func (e *editor) setup(mode editorMode) {
	sysR := e.g.WinSys.SysWin().Rect()

	mainW := sysR.W - instrW
	mainH := sysR.H
	e.winMain = e.g.WinSys.CreateWin(nil, cwin.WinCfg{
		R: cwin.Rect{
			X: 0,
			Y: 0,
			W: mainW,
			H: mainH},
		Name: "Map Editor",
	})
	e.winInstr = e.g.WinSys.CreateWin(nil, cwin.WinCfg{
		R: cwin.Rect{
			X: sysR.W - instrW,
			Y: sysR.H - instrH,
			W: instrW,
			H: instrH,
		},
		Name:       "Keyboard Help",
		ClientAttr: instrClientAttr,
	})
	e.winInstr.SetText(
		`←↑↓→: move cursor
^A : move cursor to line start
^E : move cursor to line end
Del: del at cursor or selection
^K : del all from cursor to line end
^C : cut
^V : paste
F3 : mark selection
F10: switch between map and tool box`)
	e.winToolbox = e.g.WinSys.CreateWin(nil, cwin.WinCfg{
		R: cwin.Rect{
			X: e.winInstr.Rect().X,
			Y: 0,
			W: instrW,
			H: sysR.H - instrH,
		},
		Name: "Toolbox",
	})
	e.winRulers = createRulers(e.g.WinSys, e.winMain)
	e.winMapFrame = e.g.WinSys.CreateWin(e.winRulers, cwin.WinCfg{
		R:        e.winRulers.NonRulerRect(),
		Name:     "_map_frame",
		NoBorder: true,
	})
	e.winMapFrame.FillClient(e.winMapFrame.ClientRect().ToOrigin(), mapFrameClientChx)
	if mode == editorModeNew {
		e.winMap = e.g.WinSys.CreateWin(e.winMapFrame, cwin.WinCfg{
			R:        cwin.Rect{X: 0, Y: 0, W: newMapW, H: newMapH},
			Name:     "_map",
			NoBorder: true,
		})
		e.winMap.FillClient(e.winMap.ClientRect().ToOrigin(), cwin.TransparentChx())
	}
	// because winMap has so many transparency cells, the winMapFrame's darkgray background
	// would be showing through, that's bad. So fill that area in winMapFrame with black space.
	r, _ := e.winMapFrame.ClientRect().ToOrigin().Overlap(e.winMap.Rect())
	e.winMapFrame.FillClient(r, cwin.Chx{Ch: cwin.RuneSpace})

	e.g.WinSys.SetFocus(e.winMain)
	e.resetCursor()

	e.g.Resume()
}

func (e *editor) copy() {
}

func (e *editor) paste() {
}

func (e *editor) debug() {
	var sb strings.Builder
	for y := 0; y < e.winMap.ClientRect().H; y++ {
		for x := 0; x < e.winMap.ClientRect().W; x++ {
			sb.WriteString(fmt.Sprintf("[%3d,%3d]: %+v\n", x, y, e.winMap.GetClient(x, y)))
		}
	}
	e.g.WinSys.MessageBox(nil, "dbg",
		"%s\nmap:\n%s", e.g.WinSys.Dump(), sb.String())
}
