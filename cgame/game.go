package cgame

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/jf-tech/console/cterm"
	"github.com/jf-tech/console/cutil"
	"github.com/jf-tech/console/cwin"
)

type Game struct {
	WinSys      *cwin.Sys
	MasterClock *cutil.Clock
	SpriteMgr   *SpriteManager
	SoundMgr    *SoundManager
	Exchange    *Exchange

	gameOver bool
}

func Init(provider cterm.Provider, seed ...int64) (*Game, error) {
	if len(seed) > 0 {
		rand.Seed(seed[0])
	} else {
		rand.Seed(time.Now().UnixNano())
	}
	winSys, err := cwin.Init(provider)
	if err != nil {
		return nil, err
	}
	g := &Game{WinSys: winSys, MasterClock: cutil.NewClock()}
	g.SpriteMgr = newSpriteManager(g)
	g.SoundMgr = newSoundManager()
	g.SoundMgr.Init()
	g.Exchange = newExchange()
	g.Pause()
	return g, nil
}

func (g *Game) Close() {
	g.Pause()
	g.SoundMgr.Close()
	g.WinSys.Close()
}

func (g *Game) Run(
	gameOverKeys, pauseKeys []cterm.Event,
	gameEventHandler cwin.EventHandler,
	f ...cwin.EventLoopSleepDurationFunc) {

	g.WinSys.Run(func(ev cterm.Event) cwin.EventResponse {
		if ev.Type == cterm.EventKey {
			if cwin.FindKey(gameOverKeys, ev) {
				g.GameOver()
				return cwin.EventLoopStop
			}
			if cwin.FindKey(pauseKeys, ev) {
				if g.IsPaused() {
					g.Resume()
				} else {
					g.Pause()
				}
				return cwin.EventHandled
			}
		}
		if g.IsPaused() {
			return cwin.EventHandled
		}
		resp := cwin.EventHandled
		if gameEventHandler != nil {
			resp = gameEventHandler(ev)
			if g.IsGameOver() {
				return cwin.EventLoopStop
			}
		}
		g.SpriteMgr.Process()
		return resp
	}, f...)
}

func (g *Game) Pause() {
	g.MasterClock.Pause()
	g.SoundMgr.PauseAll()
}

func (g *Game) Resume() {
	g.SoundMgr.ResumeAll()
	g.MasterClock.Resume()
}

func (g *Game) IsPaused() bool {
	return g.MasterClock.IsPaused()
}

func (g *Game) GameOver() {
	g.gameOver = true
}

func (g *Game) IsGameOver() bool {
	return g.gameOver
}

func (g *Game) HeapUsageInBytes() int64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return int64(m.HeapAlloc)
}
