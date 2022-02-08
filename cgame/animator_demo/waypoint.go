package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cterm"
	"github.com/jf-tech/console/cwin"
)

func main() {
	g, err := cgame.Init(cterm.TCell)
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

	g.WinSys.SyncExpectKey(nil)
}

func doDemo(g *cgame.Game, demoWin *cwin.Win) {
	// create a single sprite frame
	frame := cgame.FrameFromString(
		strings.Trim(readFile("resources/airplane.txt"), "\n"),
		cwin.ChAttr{Fg: cterm.ColorLightYellow})

	// create a sprite
	startX, startY := -cgame.FrameRect(frame).W+1, (demoWin.ClientRect().H-cgame.FrameRect(frame).H)/2
	s := cgame.NewSpriteBase(g, demoWin, "demo_sprite", frame, startX, startY)

	// create a simple waypoint animator has only one waypoint that gets
	// the sprite to go across the demo window.
	a := cgame.NewAnimatorWaypoint(cgame.AnimatorWaypointCfg{
		Waypoints: cgame.NewSimpleWaypoints([]cgame.Waypoint{
			{
				Type: cgame.WaypointAbs,
				X:    demoWin.ClientRect().W,
				Y:    startY,
				T:    3 * time.Second,
			},
		}),
		AfterMove: func(s cgame.Sprite) {
			demoWin.SetTitle(
				fmt.Sprintf("Demo: Sprite%s", s.Win().Rect()),
				cwin.AlignLeft)
		},
		AfterFinish: func(s cgame.Sprite) {
			demoWin.SetTitle(
				fmt.Sprintf("Demo: Sprite%s. Press any key to exit.", s.Win().Rect()),
				cwin.AlignLeft)
			g.GameOver()
		},
	})

	// add the sprite and its animator to the sprite manager
	g.SpriteMgr.AddEvent(cgame.NewSpriteEventCreate(s, a))

	// run the demo
	g.Run(nil, nil, nil)
}

func readFile(relPath string) string {
	_, filename, _, _ := runtime.Caller(1)
	b, err := ioutil.ReadFile(path.Join(path.Dir(filename), relPath))
	if err != nil {
		panic(err)
	}
	return string(b)
}
