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

	stopEventListening chan struct{}
	evChan             chan cterm.Event

	loopsDone int64
	gameOver  bool
}

func Init() (*Game, error) {
	rand.Seed(time.Now().UnixNano())
	winSys, err := cwin.Init()
	if err != nil {
		return nil, err
	}
	g := &Game{WinSys: winSys, MasterClock: newClock()}
	g.SpriteMgr = newSpriteManager(g)
	g.setupEventListening()
	g.Pause()
	return g, nil
}

func (g *Game) Close() {
	g.Pause()
	g.shutdownEventListening()
	g.WinSys.Close()
}

func (g *Game) Run(
	gameOverKeys, pauseKeys []cterm.Event, optionalRunFunc func(ev cterm.Event) bool) {

	stop := false
	for !stop && !g.IsGameOver() {
		var ev cterm.Event
		if ev = g.TryGetEvent(); ev.Type == cterm.EventKey {
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

// This is a non-blocking call
func (g *Game) TryGetEvent() cterm.Event {
	if g.evChan == nil {
		panic("SetupEventListening not called")
	}
	select {
	case ev := <-g.evChan:
		return ev
	default:
		return cterm.Event{Type: cterm.EventNone}
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

func (g *Game) setupEventListening() {
	if g.stopEventListening != nil {
		panic("SetupEventListening called twice")
	}
	g.stopEventListening = make(chan struct{})
	g.evChan = make(chan cterm.Event, 100)

	// main go routine listening for stop signal and cterm event polling.
	go func() {
	loop:
		for {
			select {
			case <-g.stopEventListening:
				break loop
			default:
				g.evChan <- cterm.PollEvent()
			}
		}
	}()
}

func (g *Game) shutdownEventListening() {
	if g.stopEventListening == nil {
		return
	}
	close(g.stopEventListening)
	g.stopEventListening = nil
	// importantly need to call cterm.Interrupt() before closing the evChan because
	// cterm.Interrupt() synchronously waits for cterm.PollEvent finishes so there
	// might be one last event coming through into the evChan. If we close it before
	// calling cterm.Interrupt(), we might get a panic.
	cterm.Interrupt()
	close(g.evChan)
	g.evChan = nil
}
