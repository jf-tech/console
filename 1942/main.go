package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cwin"
	"github.com/nsf/termbox-go"
)

const (
	codeQuit int = iota
	codeReplay
	codeGameInitFailure
)

type myGame struct {
	g         *cgame.Game
	winWeapon *cwin.Win
	winKills  *cwin.Win
	winArena  *cwin.Win
	winStats  *cwin.Win

	bgWW1942AnimationDone bool
}

func (m *myGame) main() int {
	var err error
	m.g, err = cgame.Init()
	if err != nil {
		return codeGameInitFailure
	}
	defer m.g.Close()

	m.winSetup()
	m.initSprites()
	m.g.Pause()
	m.g.WinSys.Update()
	e := m.g.WinSys.MessageBoxEx(nil,
		[]termbox.Event{
			{Key: termbox.KeyEnter},
			{Key: termbox.KeyEsc},
			{Ch: 'q'},
		},
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
	if e.Key != termbox.KeyEnter {
		return codeQuit
	}
	m.g.Resume()

	alpha := m.g.SpriteMgr.FindByName(alphaName).(*spriteAlpha)

	m.g.SetupEventListening()

loop:
	for !m.g.IsGameOver() {
		if ev := m.g.TryGetEvent(); ev.Type == termbox.EventKey {
			if ev.Key == termbox.KeyEsc || ev.Ch == 'q' {
				break loop
			} else if ev.Ch == 'p' {
				if m.g.IsPaused() {
					m.g.Resume()
				} else {
					m.g.Pause()
				}
				continue
			}
			if m.g.IsPaused() {
				continue
			}
			if ev.Key == termbox.KeyArrowUp {
				m.g.SpriteMgr.AddEvent(cgame.NewSpriteEventSetPosRelative(alpha, 0, -1))
			} else if ev.Key == termbox.KeyArrowDown {
				m.g.SpriteMgr.AddEvent(cgame.NewSpriteEventSetPosRelative(alpha, 0, 2))
			} else if ev.Key == termbox.KeyArrowLeft {
				m.g.SpriteMgr.AddEvent(cgame.NewSpriteEventSetPosRelative(alpha, -3, 0))
			} else if ev.Key == termbox.KeyArrowRight {
				m.g.SpriteMgr.AddEvent(cgame.NewSpriteEventSetPosRelative(alpha, 3, 0))
			} else if ev.Key == termbox.KeySpace {
				alpha.fireWeapon()
			}
		}
		m.moreSprites()
		m.g.SpriteMgr.ProcessAll()
		alpha.W.ToTop()
		alpha.displayWeaponInfo()
		alpha.displayKills()
		m.stats()
		m.g.WinSys.Update()
	}

	m.g.ShutdownEventListening()

	// https://textkool.com/en/ascii-art-generator?hl=default&vl=default&font=Colossal&text=Game%20Over%20!
	e = m.g.WinSys.MessageBoxEx(nil,
		[]termbox.Event{
			{Key: termbox.KeyEsc},
			{Ch: 'q'},
			{Ch: 'r'},
		},
		"Uh oh...", `

 .d8888b.                                        .d88888b.                                 888
d88P  Y88b                                      d88P" "Y88b                                888
888    888                                      888     888                                888
888         8888b.  88888b.d88b.   .d88b.       888     888 888  888  .d88b.  888d888      888
888  88888     "88b 888 "888 "88b d8P  Y8b      888     888 888  888 d8P  Y8b 888P"        888
888    888 .d888888 888  888  888 88888888      888     888 Y88  88P 88888888 888          Y8P
Y88b  d88P 888  888 888  888  888 Y8b.          Y88b. .d88P  Y8bd8P  Y8b.     888           "
 "Y8888P88 "Y888888 888  888  888  "Y8888        "Y88888P"    Y88P    "Y8888  888          888


                             Press ESC or 'q' to quit, 'r' to replay.
`)
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

func (m *myGame) initSprites() {
	// background text: "WW II"
	m.g.SpriteMgr.AddEvent(cgame.NewSpriteEventCreate(
		newSpriteBackgroundStatic(m.g, m.winArena,
			(m.winArena.ClientRect().W-cwin.TextDimension(bgWWImgTxt).W)/2,
			(m.winArena.ClientRect().H/2-cwin.TextDimension(bgWWImgTxt).H)/2,
			bgWWStaticName, bgWWImgTxt)))
	// background text: "1942"
	m.g.SpriteMgr.AddEvent(cgame.NewSpriteEventCreate(
		newSpriteBackgroundStatic(m.g, m.winArena,
			(m.winArena.ClientRect().W-cwin.TextDimension(bg1942ImgTxt).W)/2,
			(m.winArena.ClientRect().H*3/2-cwin.TextDimension(bg1942ImgTxt).H)/2,
			bg1942StaticName, bg1942ImgTxt)))
	// alpha - player airplane
	m.g.SpriteMgr.AddEvent(cgame.NewSpriteEventCreate(&spriteAlpha{
		SpriteBase: cgame.NewSpriteBase(m.g, m.winArena,
			cgame.SpriteCfg{Name: alphaName, Cells: cgame.StringToCells(alphaImgTxt, alphaAttr)},
			(m.winArena.ClientRect().W-cwin.TextDimension(alphaImgTxt).W)/2,
			m.winArena.ClientRect().H-cwin.TextDimension(alphaImgTxt).H),
		m: m}))
	m.g.SpriteMgr.ProcessAll()
}

func (m *myGame) moreSprites() {
	if m.g.IsPaused() {
		return
	}
	m.bgWW1942Animation()
	m.bgStars()
	m.betas()
	m.gamma()
	m.giftPacks()
}

func (m *myGame) bgWW1942Animation() {
	if m.g.Clock.SinceOrigin() > bgInitialWait && !m.bgWW1942AnimationDone {
		s1 := m.g.SpriteMgr.FindByName(bgWWStaticName)
		m.g.SpriteMgr.AddEvent(cgame.NewSpriteEventDelete(s1))
		m.g.SpriteMgr.AddEvent(cgame.NewSpriteEventCreate(newSpriteBackgroundAnimated(
			m.g, m.winArena, s1.Win().Rect().X, s1.Win().Rect().Y, bgWWAnimatedName, bgWWImgTxt, 1)))
		s2 := m.g.SpriteMgr.FindByName(bg1942StaticName)
		m.g.SpriteMgr.AddEvent(cgame.NewSpriteEventDelete(s2))
		m.g.SpriteMgr.AddEvent(cgame.NewSpriteEventCreate(newSpriteBackgroundAnimated(
			m.g, m.winArena, s2.Win().Rect().X, s2.Win().Rect().Y, bg1942AnimatedName, bg1942ImgTxt, -1)))
		m.bgWW1942AnimationDone = true
	}
}

func (m *myGame) bgStars() {
	if testProb(bgStarGenProb) {
		x := rand.Int() % (m.winArena.ClientRect().W - cwin.TextDimension(betaImgTxt).W)
		m.g.SpriteMgr.AddEvent(cgame.NewSpriteEventCreate(newSpriteBackgroundStar(m.g, m.winArena, x,
			rand.Int()%10)))
	}
}

func (m *myGame) betas() {
	if testProb(betaGenProb) {
		x := rand.Int() % (m.winArena.ClientRect().W - cwin.TextDimension(betaImgTxt).W)
		m.g.SpriteMgr.AddEvent(cgame.NewSpriteEventCreate(newSpriteBeta(m.g, m.winArena, x, 0)))
	}
}

func (m *myGame) gamma() {
	if testProb(gammaGenProb) {
		x := rand.Int() % (m.winArena.ClientRect().W - cwin.TextDimension(betaImgTxt).W)
		m.g.SpriteMgr.AddEvent(cgame.NewSpriteEventCreate(newSpriteGamma(m.g, m.winArena, x, 0)))
	}
}

func (m *myGame) giftPacks() {
	if sym, attr, ok := genGiftPack(); ok {
		x := rand.Int() % (m.winArena.ClientRect().W - cwin.TextDimension(giftPackImgTxts[0]).W)
		m.g.SpriteMgr.AddEvent(cgame.NewSpriteEventCreate(newSpriteGiftPack(m.g, m.winArena, x, 0, sym, attr)))
	}
}

func (m *myGame) stats() {
	m.winStats.SetText(fmt.Sprintf(`
Game stats:
----------------------------
Time passed: %s %s

Internals:
----------------------------
Alpha Position: %s
Beta Firing Prob: %.0f％
%s`,
		time.Duration(m.g.Clock.SinceOrigin()/(time.Second))*(time.Second),
		func() string {
			if m.g.IsPaused() {
				return "(paused)"
			}
			return ""
		}(),
		m.g.SpriteMgr.FindByName(alphaName).(*spriteAlpha).W.Rect(),
		float64(100)/float64(betaFiringCurProb),
		m.g.SpriteMgr.DbgStats()))
}

func main() {
	code := codeReplay
	for code == codeReplay {
		code = (&myGame{}).main()
	}
	os.Exit(code)
}
