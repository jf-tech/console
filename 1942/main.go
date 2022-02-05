package main

import (
	"os"

	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cwin"
	"github.com/nsf/termbox-go"
)

const (
	codeQuit int = iota
	codeReplay
	codeGameInitFailure

	// https://textkool.com/en/ascii-art-generator?hl=default&vl=default&font=Colossal&text=Game%20Over%20!
	gameOverTxt = `

 .d8888b.                                        .d88888b.                                 888
d88P  Y88b                                      d88P" "Y88b                                888
888    888                                      888     888                                888
888         8888b.  88888b.d88b.   .d88b.       888     888 888  888  .d88b.  888d888      888
888  88888     "88b 888 "888 "88b d8P  Y8b      888     888 888  888 d8P  Y8b 888P"        888
888    888 .d888888 888  888  888 88888888      888     888 Y88  88P 88888888 888          Y8P
Y88b  d88P 888  888 888  888  888 Y8b.          Y88b. .d88P  Y8bd8P  Y8b.     888           "
 "Y8888P88 "Y888888 888  888  888  "Y8888        "Y88888P"    Y88P    "Y8888  888          888


                             Press ESC or 'q' to quit, 'r' to replay.`

	// https://textkool.com/en/ascii-art-generator?hl=default&vl=default&font=Colossal&text=You%20Won%20!
	youWonTxt = `

Y88b   d88P                     888       888                        888
 Y88b d88P                      888   o   888                        888
  Y88o88P                       888  d8b  888                        888
   Y888P  .d88b.  888  888      888 d888b 888  .d88b.  88888b.       888
    888  d88""88b 888  888      888d88888b888 d88""88b 888 "88b      888
    888  888  888 888  888      88888P Y88888 888  888 888  888      Y8P
    888  Y88..88P Y88b 888      8888P   Y8888 Y88..88P 888  888       "
    888   "Y88P"   "Y88888      888P     Y888  "Y88P"  888  888      888


                  Press ESC or 'q' to quit, 'r' to replay.`
)

type myGame struct {
	g         *cgame.Game
	winWeapon *cwin.Win
	winKills  *cwin.Win
	winArena  *cwin.Win
	winStats  *cwin.Win
	easyMode  bool
}

func (m *myGame) main() int {
	var err error
	m.g, err = cgame.Init()
	if err != nil {
		return codeGameInitFailure
	}
	defer m.g.Close()

	m.winSetup()
	m.g.WinSys.Update()

	e := m.g.WinSys.MessageBoxEx(nil,
		cwin.Keys(termbox.KeyEnter, termbox.KeyEsc, 'q', 'e'),
		"WWII - 1942", `
Axis and Allied forces have been deeply engaged in World War II and now the
fighting is quickly approaching the final stage. Both sides have suffered
extremely heavy losses. As a newly-recruited pilot, your assignment is to
penetrate deep into the heart of the enemy territories and destroy strategic
targets, giving our ground troops a chance to regroup and launch into the
final battle!

Good luck, solider!

Press Enter to start the game; ESC or 'q' to quit.
`)
	if e.Key == termbox.KeyEsc || e.Ch == 'q' {
		return codeQuit
	}
	m.easyMode = e.Ch == 'e'

	m.g.Resume()

	newStage(m, 0).Run()
	if !m.g.IsGameOver() {
		newStage(m, 1).Run()
	}

	e = m.g.WinSys.MessageBoxEx(nil, cwin.Keys(termbox.KeyEsc, 'q', 'r'), "Result",
		func() string {
			if m.g.IsGameOver() {
				return gameOverTxt
			}
			return youWonTxt
		}())
	if e.Ch == 'r' {
		return codeReplay
	}
	return codeQuit
}

const (
	winHeaderW       = 103
	winHeaderH       = 1
	winArenaW        = winHeaderW - 2
	winStatsW        = 40
	winInstructionsH = 7
	winGameW         = 1 /*border*/ + 1 /*space*/ + winHeaderW + 1 /*space*/ + winStatsW + 1 /*space*/ + 1 /*border*/

	textInstructions = `Press ESC or 'q' to quit the game.
Press arrows to move your airplane.
Press space bar to fire your weapon.
Press 'b' to launch a bomb.
Press 'p' to pause/unpause the game.`
)

func (m *myGame) winSetup() {
	winSysClientR := m.g.WinSys.GetSysWin().ClientRect()
	winGame := m.g.WinSys.CreateWin(nil, cwin.WinCfg{
		R: cwin.Rect{
			X: (winSysClientR.W - winGameW) / 2,
			Y: 0,
			W: winGameW,
			H: winSysClientR.H},
		Name:    "game",
		NoTitle: true,
	})
	winGameClientR := winGame.ClientRect()
	winHeader := m.g.WinSys.CreateWin(winGame, cwin.WinCfg{
		R: cwin.Rect{
			X: 1,
			Y: 0,
			W: winHeaderW,
			H: winHeaderH},
		Name:       "header",
		NoBorder:   true,
		ClientAttr: cwin.ChAttr{Bg: termbox.ColorBlue},
	})
	m.winWeapon = m.g.WinSys.CreateWin(winHeader, cwin.WinCfg{
		R: cwin.Rect{
			X: 0,
			Y: 0,
			W: winHeaderW / 2,
			H: winHeaderH},
		Name:       "header_weapon",
		NoBorder:   true,
		ClientAttr: cwin.ChAttr{Fg: termbox.ColorLightYellow, Bg: termbox.ColorBlue},
	})
	m.winKills = m.g.WinSys.CreateWin(winHeader, cwin.WinCfg{
		R: cwin.Rect{
			X: winHeaderW / 2,
			Y: 0,
			W: winHeaderW / 2,
			H: winHeaderH},
		Name:       "header_kills",
		NoBorder:   true,
		ClientAttr: cwin.ChAttr{Fg: termbox.ColorLightYellow, Bg: termbox.ColorBlue},
	})
	_ = m.g.WinSys.CreateWin(winGame, cwin.WinCfg{
		R:          cwin.Rect{X: 1, Y: 1, W: 1, H: winGameClientR.H - winHeaderH},
		NoBorder:   true,
		ClientAttr: cwin.ChAttr{Bg: termbox.ColorRed},
	})
	_ = m.g.WinSys.CreateWin(winGame, cwin.WinCfg{
		R:          cwin.Rect{X: winHeaderW, Y: 1, W: 1, H: winGameClientR.H - winHeaderH},
		NoBorder:   true,
		ClientAttr: cwin.ChAttr{Bg: termbox.ColorRed},
	})
	m.winArena = m.g.WinSys.CreateWin(winGame, cwin.WinCfg{
		R: cwin.Rect{
			X: 2,
			Y: winHeaderH,
			W: winArenaW,
			H: winGameClientR.H - winHeaderH},
		Name:     "arena",
		NoBorder: true,
	})
	m.winStats = m.g.WinSys.CreateWin(winGame, cwin.WinCfg{
		R: cwin.Rect{
			X: 4 + winArenaW,
			Y: 0,
			W: winStatsW,
			H: winGameClientR.H - winInstructionsH},
		Name: "Stats",
	})
	winInstructions := m.g.WinSys.CreateWin(winGame, cwin.WinCfg{
		R: cwin.Rect{
			X: m.winStats.Rect().X,
			Y: m.winStats.Rect().Y + m.winStats.Rect().H,
			W: m.winStats.Rect().W,
			H: winGameClientR.H - m.winStats.Rect().H},
		Name:       "Keyboard",
		BorderAttr: cwin.ChAttr{Bg: termbox.ColorBlue},
		ClientAttr: cwin.ChAttr{Bg: termbox.ColorBlue},
	})
	winInstructions.SetText(textInstructions)
}

func main() {
	code := codeReplay
	for code == codeReplay {
		code = (&myGame{}).main()
	}
	os.Exit(code)
}
