package main

import (
	"fmt"
	"path"
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
	})
	g.WinSys.Update() // nothing shows onto screen unless Update() is called.
	g.Resume()        // game (master clock) is always paused right after init.

	doDemo(g, demoWin)
}

func doDemo(g *cgame.Game, demoWin cwin.Win) {
	filepath := path.Join(cutil.GetCurFileDir(), "resources/doorbell.mp3")
	g.SoundMgr.PlayMP3(filepath, -1, -1)
	g.Run(cwin.Keys(cterm.KeyEsc, 'q'), cwin.Keys(' '), func(ev cterm.Event) cwin.EventResponse {
		demoWin.SetTitle(
			fmt.Sprintf("Demo - Sound. Space to pause/resume; ESC or 'q' to quit. Time: %s",
				g.MasterClock.Now().Round(time.Millisecond)))
		return cwin.EventHandled
	})
}
