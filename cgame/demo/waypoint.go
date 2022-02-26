package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"time"

	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cterm"
	"github.com/jf-tech/console/cutil"
	"github.com/jf-tech/console/cwin"
)

func main() {
	g, err := cgame.Init(cterm.TermBox)
	if err != nil {
		panic(err)
	}
	defer g.Close()
	sysWinR := g.WinSys.SysWin().Rect()
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
}

func doDemo(g *cgame.Game, demoWin cwin.Win) {
	// create a single sprite frame
	frame := cgame.FrameFromString(
		strings.Trim(readFile("resources/airplane.txt"), "\n"),
		cwin.Attr{Fg: cterm.ColorLightYellow})

	// create a sprite
	startX, startY := -cgame.FrameRect(frame).W, (demoWin.ClientRect().H-cgame.FrameRect(frame).H)/2
	s := cgame.NewSpriteBase(g, demoWin, "demo_sprite", frame, startX, startY)
	// create a simple waypoint animator has only one waypoint that gets
	// the sprite to go across the demo window.
	s.AddAnimator(cgame.NewAnimatorWaypoint(s, cgame.AnimatorWaypointCfg{
		Waypoints: cgame.NewWaypointProviderSimple([]cgame.Waypoint{
			{
				DX: cgame.FrameRect(frame).W + demoWin.ClientRect().W,
				DY: 0,
				T:  3 * time.Second,
			},
		}),
		AnimatorCfgCommon: cgame.AnimatorCfgCommon{
			AfterUpdate: func() {
				demoWin.SetTitle(fmt.Sprintf("Demo: Sprite%s", s.Rect()))
			},
			AfterFinish: func() {
				g.GameOver()
			},
		},
	}))
	g.SpriteMgr.AddSprite(s)
	// run the demo
	g.Run(cwin.Keys('q', cterm.KeyEsc), nil, cwin.NopHandledEventHandler)
}

func readFile(relPath string) string {
	b, err := ioutil.ReadFile(path.Join(cutil.GetCurFileDir(), relPath))
	if err != nil {
		panic(err)
	}
	return string(b)
}
