package cgame

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/jf-tech/console/cterm"
	"github.com/jf-tech/console/cwin"
)

type Game struct {
	WinSys      *cwin.Sys
	MasterClock *Clock
	SpriteMgr   *SpriteManager

	loopsDone int64
	gameOver  bool
}

func Init(provider cterm.Provider) (*Game, error) {
	rand.Seed(time.Now().UnixNano())
	winSys, err := cwin.Init(provider)
	if err != nil {
		return nil, err
	}
	g := &Game{WinSys: winSys, MasterClock: newClock()}
	g.SpriteMgr = newSpriteManager(g)
	g.Pause()
	return g, nil
}

func (g *Game) Close() {
	g.Pause()
	g.WinSys.Close()
}

func (g *Game) Run(
	gameOverKeys, pauseKeys []cterm.Event, optionalRunFunc func(ev cterm.Event) bool) {

	stop := false
	for !stop && !g.IsGameOver() {
		var ev cterm.Event
		if ev = g.WinSys.TryGetEvent(); ev.Type == cterm.EventKey {
			if cwin.FindKey(gameOverKeys, ev) {
				g.GameOver()
				return
			}
			if cwin.FindKey(pauseKeys, ev) {
				if g.IsPaused() {
					g.Resume()
				} else {
					g.Pause()
				}
			}
		}
		if optionalRunFunc != nil {
			stop = optionalRunFunc(ev)
		}
		g.SpriteMgr.Process()
		g.WinSys.Update()
		g.loopsDone++
	}
}

func (g *Game) Pause() {
	g.MasterClock.Pause()
}

func (g *Game) Resume() {
	g.MasterClock.Resume()
}

func (g *Game) IsPaused() bool {
	return g.MasterClock.IsPaused()
}

func (g *Game) GameOver() {
	g.gameOver = true
	g.Pause()
}

func (g *Game) IsGameOver() bool {
	return g.gameOver
}

func (g *Game) TotalLoops() int64 {
	return g.loopsDone
}

func (g *Game) FPS() float64 {
	now := g.MasterClock.Now()
	if now == 0 {
		return float64(0)
	}
	return float64(g.loopsDone) / (float64(now) / float64(time.Second))
}

func (g *Game) HeapUsageInBytes() int64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return int64(m.HeapAlloc)
}
