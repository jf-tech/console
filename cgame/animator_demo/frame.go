package main

import (
	"fmt"
	"math"
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
}

func doDemo(g *cgame.Game, demoWin *cwin.Win) {
	r := demoWin.ClientRect()
	w := r.W * 3 / 4
	h := r.H * 3 / 4

	fp := &sineWaveFrameProvider{w: w, h: h}
	f0, _, _ := fp.Next()
	a := cgame.NewAnimatorFrame(cgame.AnimatorFrameCfg{Frames: fp})
	s := cgame.NewSpriteBase(g, demoWin, "demo_frame", f0, (r.W-w)/2, (r.H-h)/2)
	g.SpriteMgr.AddEvent(cgame.NewSpriteEventCreate(s, a))

	g.Run(nil, nil, func(ev termbox.Event) bool {
		demoWin.SetTitle(
			func() string {
				return fmt.Sprintf("Demo - space to pause/resume, any other key to exit. Time: %s",
					g.MasterClock.Now()/time.Millisecond*time.Millisecond)
			}(),
			cwin.AlignLeft)
		if ev.Type != termbox.EventKey {
			return false
		}
		if ev.Key != termbox.KeySpace {
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
	w, h  int
	shift int
}

func (sfp *sineWaveFrameProvider) Next() (cgame.Frame, time.Duration, bool) {
	toRX := func(x, w int) float64 {
		return float64(x) / float64(w) * 2 * math.Pi
	}
	fromRY := func(ry float64, h int) int {
		return int((1 - ry) / 2 * float64(h))
	}
	var f cgame.Frame
	for x := 0; x < sfp.w; x++ {
		y := fromRY(math.Sin(toRX(x+sfp.shift, sfp.w)), sfp.h)
		f = append(f, cgame.Cell{
			X:   x,
			Y:   y,
			Chx: cwin.Chx{Ch: '.', Attr: cwin.ChAttr{Fg: termbox.ColorLightGreen}}})
	}
	sfp.shift = (sfp.shift - 1 + sfp.w) % sfp.w
	return f, 50 * time.Millisecond, true
}
