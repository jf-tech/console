package main

import (
	"fmt"
	"io/ioutil"
	"path"
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
}

func doDemo(g *cgame.Game, demoWin *cwin.Win) {
	for _, p := range []string{
		"resources/spacecraft_small_1.txt",
		"resources/spacecraft_small_2.txt",
		"resources/spacecraft_large_1.txt",
	} {
		doExplosion(g, demoWin, p)
	}
}

func doExplosion(g *cgame.Game, demoWin *cwin.Win, filepath string) {
	fn := path.Base(filepath)
	demoWin.SetTitle(
		fmt.Sprintf("Demo - Explosion '%s': any key to start", fn), cwin.AlignLeft)
	f := cgame.FrameFromString(
		strings.Trim(readFile(filepath), "\n"), cwin.ChAttr{Fg: cterm.ColorLightCyan})
	s := cgame.NewSpriteBase(g, demoWin, "s", f,
		(demoWin.ClientRect().W-cgame.FrameRect(f).W)/2,
		(demoWin.ClientRect().H-cgame.FrameRect(f).H)/2)
	g.SpriteMgr.AddEvent(cgame.NewSpriteEventCreate(s))
	g.SpriteMgr.Process()
	g.WinSys.Update()
	g.WinSys.SyncExpectKey(nil)
	demoWin.SetTitle(fmt.Sprintf("Demo - Explosion '%s' in progress...", fn), cwin.AlignLeft)
	done := false
	cgame.CreateExplosion(s, cgame.ExplosionCfg{
		MaxDuration: 2 * time.Second,
		AfterFinish: func() {
			done = true
		},
	})
	g.Run(nil, nil, func(_ cterm.Event) bool {
		return done
	})
	demoWin.SetTitle(fmt.Sprintf("Demo - Explosion '%s' done. Any key for next", fn), cwin.AlignLeft)
	g.WinSys.Update()
	g.WinSys.SyncExpectKey(nil)
}

func readFile(relPath string) string {
	b, err := ioutil.ReadFile(path.Join(cgame.GetCurFileDir(), relPath))
	if err != nil {
		panic(err)
	}
	return string(b)
}
