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

	loopCount int64
	gameOver  bool
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

func (g *Game) Run(gameOverKeys, pauseKeys []cterm.Event, optionalRunFunc cwin.MsgLoopFunc) {
	g.WinSys.Run(func(ev cterm.Event) cwin.MsgLoopResponseType {
		g.loopCount++
		if ev.Type == cterm.EventKey {
			if cwin.FindKey(gameOverKeys, ev) {
				g.GameOver()
				return cwin.MsgLoopStop
			}
			if cwin.FindKey(pauseKeys, ev) {
				if g.IsPaused() {
					g.Resume()
				} else {
					g.Pause()
				}
				return cwin.MsgLoopContinue
			}
		}
		// Reach here because of either EventNone, or  EvenKey but not game over or pause key...
		if g.IsPaused() {
			return cwin.MsgLoopContinue
		}
		resp := cwin.MsgLoopContinue
		if optionalRunFunc != nil {
			resp = optionalRunFunc(ev)
		}
		g.SpriteMgr.Process()
		return resp
	})
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

func (g *Game) FPS() float64 {
	now := g.MasterClock.Now()
	if now == 0 {
		return float64(0)
	}
	return float64(g.loopCount) / (float64(now) / float64(time.Second))
}

func (g *Game) HeapUsageInBytes() int64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return int64(m.HeapAlloc)
}
