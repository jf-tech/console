package main

import (
	"strings"
	"time"

	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cwin"
	"github.com/nsf/termbox-go"
)

func main() {
	g, err := cgame.Init()
	if err != nil {
		panic(err)
	}
	defer g.Close()
	sysWinR := g.WinSys.GetSysWin().Rect()
	// create a demo window that is 3/4 of the system window (which is the same size
	// of the current terminal/console) and center it.
	demoWin := g.WinSys.CreateWin(nil, cwin.WinCfg{
		R: cwin.Rect{
			X: sysWinR.W / 8,
			Y: sysWinR.H / 8,
			W: sysWinR.W * 3 / 4,
			H: sysWinR.H * 3 / 4,
		},
		Name: "Demo",
	})
	g.WinSys.Update() // nothing shows onto screen unless Update() is called.
	g.Resume()        // game (master clock) is always paused right after init.

	doDemo(g, demoWin)

	cwin.SyncExpectKey(nil)
}

func doDemo(g *cgame.Game, demoWin *cwin.Win) {
	// create a single sprite frame
	frame := cgame.FrameFromString(strings.Trim(`
\┃┃/
 \/
`, "\n"), cwin.ChAttr{Fg: termbox.ColorLightCyan})
	frameR := cgame.FrameRect(frame)

	// create a sprite
	x, y := (demoWin.ClientRect().W-frameR.W)/2, (demoWin.ClientRect().H-frameR.H)/2
	s := cgame.NewSpriteBase(g, demoWin, "explosion_demo", frame, x, y)
	// add the sprite to the sprite manager
	g.SpriteMgr.AddEvent(cgame.NewSpriteEventCreate(s))
	g.SpriteMgr.Process()
	g.WinSys.Update()
	cwin.SyncExpectKey(nil)

	cgame.CreateExplosion(g.SpriteMgr, s, cgame.ExplosionCfg{
		Scale: 2.0,
		T:     2 * time.Second,
	})
	g.SpriteMgr.Process()
	g.WinSys.Update()
}
