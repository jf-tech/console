package main

import (
	"fmt"
	"io/ioutil"
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
	demoWin := g.WinSys.CreateWin(nil, cwin.WinCfg{
		R: cwin.Rect{
			X: sysWinR.W / 8,
			Y: sysWinR.H / 8,
			W: sysWinR.W * 3 / 4,
			H: sysWinR.H * 3 / 4,
		},
		Name: "Demo",
	})
	g.WinSys.Update()
	g.Resume()

	doDemo(g, demoWin)

	cwin.SyncExpectKey(nil)
}

func doDemo(g *cgame.Game, demoWin *cwin.Win) {
	// create a single sprite frame
	frame := cgame.FrameFromString(
		strings.Trim(readFile("resources/airplane.txt"), "\n"),
		cwin.ChAttr{Fg: termbox.ColorLightYellow})

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

	// add sprite to the sprite manager
	g.SpriteMgr.AddEvent(cgame.NewSpriteEventCreate(s, a))

	// run the demo
	g.Run(nil, nil, nil)
}

func readFile(path string) string {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(b)
}
