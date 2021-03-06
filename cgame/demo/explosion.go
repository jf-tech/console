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
	for _, p := range []string{
		"resources/spacecraft_small_1.txt",
		"resources/spacecraft_small_2.txt",
		"resources/spacecraft_large_1.txt",
	} {
		if !doExplosion(g, demoWin, p) {
			break
		}
	}
}

func doExplosion(g *cgame.Game, demoWin cwin.Win, filepath string) bool {
	fn := path.Base(filepath)
	demoWin.SetTitle(fmt.Sprintf("Demo - Explosion '%s': any key to start", fn))
	f := cgame.FrameFromString(
		strings.Trim(readFile(filepath), "\n"), cwin.Attr{Fg: cterm.ColorLightCyan})
	s := cgame.NewSpriteBase(g, demoWin, "s", f,
		(demoWin.ClientRect().W-cgame.FrameRect(f).W)/2,
		(demoWin.ClientRect().H-cgame.FrameRect(f).H)/2)
	g.SpriteMgr.AddSprite(s)
	g.WinSys.Update()
	gameOver := false
	g.WinSys.SyncExpectKey(func(k cterm.Key, r rune) bool {
		if k == cterm.KeyEsc || r == 'q' {
			gameOver = true
		}
		return true
	})
	if gameOver {
		return false
	}
	demoWin.SetTitle(fmt.Sprintf("Demo - Explosion '%s' in progress...", fn))
	startTime := g.MasterClock.Now()
	done := false
	cgame.CreateExplosion(s, cgame.ExplosionCfg{
		MaxDuration: 2 * time.Second,
		AfterFinish: func() {
			done = true
		},
	})
	g.Run(cwin.Keys(cterm.KeyEsc, 'q'), nil, func(ev cterm.Event) cwin.EventResponse {
		return cwin.TrueForEventSystemStop(done)
	})
	if g.IsGameOver() {
		return false
	}
	demoWin.SetTitle(fmt.Sprintf("Demo - Explosion '%s' done. Used %s. Any key for next",
		fn, (g.MasterClock.Now() - startTime).Round(time.Millisecond)))
	g.WinSys.Update()
	g.WinSys.SyncExpectKey(func(k cterm.Key, r rune) bool {
		if k == cterm.KeyEsc || r == 'q' {
			gameOver = true
		}
		return true
	})
	return !gameOver
}

func readFile(relPath string) string {
	b, err := ioutil.ReadFile(path.Join(cutil.GetCurFileDir(), relPath))
	if err != nil {
		panic(err)
	}
	return string(b)
}
