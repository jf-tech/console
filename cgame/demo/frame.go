package main

import (
	"fmt"
	"math"
	"math/rand"
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
}

type exchange struct {
	w, h    int
	s       cgame.Sprite
	curDir  rune
	curDist int
}

// In this demo, we'll combine two animators together:
// - frame animator to show a shifting sine wave
// - waypoint animator to move the sprite around
func doDemo(g *cgame.Game, demoWin *cwin.Win) {
	r := demoWin.ClientRect()
	w := r.W * 3 / 4
	h := r.H * 3 / 4

	// data exchange among doDemo, AnimatorWaypoint and AnimatorFrame.
	e := &exchange{w: w, h: h}

	// create the sine wave provider and use its first frame as the base frame of the sprite.
	fp := &sineWaveFrameProvider{e: e}
	f0, _, _ := fp.Next()

	// sprite
	s := cgame.NewSpriteBase(g, demoWin, "demo_frame", f0, (r.W-w)/2, (r.H-h)/2)
	e.s = s

	// AnimatorWaypoint
	aw := cgame.NewAnimatorWaypoint(cgame.AnimatorWaypointCfg{Waypoints: &waypointProvider{e: e}})

	// AnimatorFrame
	af := cgame.NewAnimatorFrame(cgame.AnimatorFrameCfg{Frames: fp})

	// Add sprite and two animators to the system. Note we add af after aw, so that aw
	// decides direction and dist the sprite will travel and af will use the dir symbol
	// as the frame background :)
	g.SpriteMgr.AddEvent(cgame.NewSpriteEventCreate(s, aw, af))

	g.Run(nil, nil, func(ev cterm.Event) bool {
		demoWin.SetTitle(
			func() string {
				return fmt.Sprintf(
					"Demo - space to pause/resume, any other key to exit. Dir: %c. Dist: %2d Time: %s",
					e.curDir, e.curDist,
					g.MasterClock.Now()/time.Millisecond*time.Millisecond)
			}(),
			cwin.AlignLeft)
		if ev.Type != cterm.EventKey {
			return false
		}
		if ev.Ch != ' ' {
			return true
		}
		if g.IsPaused() {
			g.Resume()
		} else {
			g.Pause()
		}
		return false
	})
}

type sineWaveFrameProvider struct {
	e     *exchange
	shift int
}

func (sfp *sineWaveFrameProvider) Next() (cgame.Frame, time.Duration, bool) {
	toRX := func(x, w int) float64 {
		return float64(x) / float64(w) * 2 * math.Pi
	}
	fromRY := func(ry float64, h int) int {
		return int((1 - ry) / 2 * float64(h-1))
	}
	var f cgame.Frame
	for y := 0; y < sfp.e.h; y++ {
		for x := 0; x < sfp.e.w; x++ {
			f = append(f, cgame.Cell{
				X: x,
				Y: y,
				Chx: cwin.Chx{
					Ch: func() rune {
						if (x+y)%2 == 0 {
							return sfp.e.curDir
						}
						return ' '
					}(),
					Attr: cwin.ChAttr{Fg: cterm.ColorWhite, Bg: cterm.ColorDarkGray}}})
		}
	}
	for x := 0; x < sfp.e.w; x++ {
		y := fromRY(math.Sin(toRX(x+sfp.shift, sfp.e.w)), sfp.e.h)
		f[y*sfp.e.w+x].Chx =
			cwin.Chx{Ch: '#', Attr: cwin.ChAttr{Fg: cterm.ColorYellow, Bg: cterm.ColorLightBlue}}
	}
	sfp.shift = (sfp.shift - 1 + sfp.e.w) % sfp.e.w
	return f, 50 * time.Millisecond, true
}

type waypointProvider struct {
	e *exchange
}

const (
	minDistBeforeDirChange = 1
	maxDistBeforeDirChange = 100
)

func (wp *waypointProvider) Next() (cgame.Waypoint, bool) {
	for {
		dist := rand.Int() % (maxDistBeforeDirChange - minDistBeforeDirChange + 1)
		dist += minDistBeforeDirChange
		dirIdx := rand.Int() % len(cgame.DirOffSetXY)
		w := wp.e.s.Win()
		newR := w.Rect()
		newR.X += cgame.DirOffSetXY[dirIdx].X * dist
		newR.Y += cgame.DirOffSetXY[dirIdx].Y * dist
		if overlapped, ro := newR.Overlap(w.Parent().ClientRect().ToOrigin()); overlapped && ro == newR {
			wp.e.curDir = cgame.DirSymbols[dirIdx]
			wp.e.curDist = dist
			return cgame.Waypoint{
				Type: cgame.WaypointAbs,
				X:    newR.X,
				Y:    newR.Y,
				T:    time.Duration(dist) * 200 * time.Millisecond, // eachk "pixel" move takes 200ms.
			}, true
		}
	}
}
