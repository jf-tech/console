package main

import (
	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cgame/assets"
	"github.com/jf-tech/console/cterm"
	"github.com/jf-tech/console/cwin"
)

type myGame struct {
	g          *cgame.Game
	winHeader  cwin.Win
	winArena   cwin.Win
	winStats   cwin.Win
	easyMode   bool
	invincible bool
}

func (m *myGame) main() assets.GameResult {
	var err error
	m.g, err = cgame.Init(cterm.TCell)
	if err != nil {
		return assets.GameResultSystemFailure
	}
	defer m.g.Close()

	m.g.SoundMgr.AvoidSameClipConcurrentPlaying()
	m.winSetup()
	m.g.WinSys.Update()

	m.g.SoundMgr.PlayMP3(sfxBackgroundFile, sfxBackgroundVol, -1)

	e := m.g.WinSys.MessageBoxEx(m.winArena,
		append(cwin.Keys(cterm.KeyEnter), append(gameOverKeys, easyModeKeys...)...),
		"WWII - 1942", `
Axis and Allied forces have been deeply engaged in World War II and now the
fighting is quickly approaching the final stage. Both sides have suffered
extremely heavy losses. As a newly-recruited pilot, your assignment is to
penetrate deep into the heart of the enemy territories and destroy strategic
targets, giving our ground troops a chance to regroup and launch into the
final battle!

Good luck, solider!

Press Enter to start the game; ESC or 'q' to quit.
('e' to start in Easy Mode, if you're patient enough to read :)
`)
	if cwin.FindKey(gameOverKeys, e) {
		return assets.GameResultQuit
	}
	m.easyMode = cwin.FindKey(easyModeKeys, e)

	m.g.Resume()

	m.g.SoundMgr.PlayMP3(sfxGameStartFile, sfxClipVol, 1)

	registerCollidable(m.g.SpriteMgr)

	for i := 0; i < totalStages && !m.g.IsGameOver(); i++ {
		stage := newStage(m, i)
		stage.Run()
	}

	if m.g.IsGameOver() {
		m.g.SoundMgr.PlayMP3(sfxGameOverFile, sfxClipVol, 1)
		return assets.DisplayGameOverDialog(m.g)
	} else {
		m.g.SoundMgr.PlayMP3(sfxYouWonFile, sfxClipVol, 1)
		return assets.DisplayYouWonDialog(m.g)
	}
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
Press 'i' for invincible cheat mode.
Press 'p' to pause/unpause the game.`
)

func (m *myGame) winSetup() {
	winSysClientR := m.g.WinSys.SysWin().ClientRect()
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
	m.winHeader = m.g.WinSys.CreateWin(winGame, cwin.WinCfg{
		R: cwin.Rect{
			X: 1,
			Y: 0,
			W: winHeaderW,
			H: winHeaderH},
		Name:       "header",
		NoBorder:   true,
		ClientAttr: cwin.Attr{Bg: cterm.ColorBlue},
	})
	_ = m.g.WinSys.CreateWin(winGame, cwin.WinCfg{
		R:          cwin.Rect{X: 1, Y: 1, W: 1, H: winGameClientR.H - winHeaderH},
		NoBorder:   true,
		ClientAttr: cwin.Attr{Bg: cterm.ColorRed},
	})
	_ = m.g.WinSys.CreateWin(winGame, cwin.WinCfg{
		R:          cwin.Rect{X: winHeaderW, Y: 1, W: 1, H: winGameClientR.H - winHeaderH},
		NoBorder:   true,
		ClientAttr: cwin.Attr{Bg: cterm.ColorRed},
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
		BorderAttr: cwin.Attr{Bg: cterm.ColorBlue},
		ClientAttr: cwin.Attr{Bg: cterm.ColorBlue},
	})
	winInstructions.SetText(textInstructions)
}

func main() {
	for (&myGame{}).main() == assets.GameResultReplay {
	}
}
