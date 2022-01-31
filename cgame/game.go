package cgame

import (
	"math/rand"
	"time"

	"github.com/jf-tech/console/cwin"
	"github.com/nsf/termbox-go"
)

type Game struct {
	WinSys    *cwin.Sys
	Clock     *GameClock
	SpriteMgr *SpriteManager

	stopEventListening chan struct{}
	evChan             chan termbox.Event

	// false: the entire world (incl. time) freezes;
	// true: time freezes, sprites still might be moved around by keys (if key binded)
	timePauseOnly bool
	gameOver      bool
}

func Init() (*Game, error) {
	rand.Seed(time.Now().UnixNano())

	winSys, err := cwin.Init()
	if err != nil {
		return nil, err
	}
	g := &Game{
		WinSys: winSys,
		Clock:  newGameClock(),
	}
	g.SpriteMgr = newSpriteManager(g)
	return g, nil
}

func (g *Game) Close() {
	g.WinSys.Close()
	if g.stopEventListening != nil {
		close(g.stopEventListening)
	}
	if g.evChan != nil {
		close(g.evChan)
	}
}

func (g *Game) SetupEventListening() {
	if g.stopEventListening != nil {
		panic("SetupEventListening called twice")
	}
	g.stopEventListening = make(chan struct{})
	g.evChan = make(chan termbox.Event, 100)

	// main go routine listening for stop signal and termbox event polling.
	go func() {
	loop:
		for {
			select {
			case <-g.stopEventListening:
				break loop
			default:
				g.evChan <- termbox.PollEvent()
			}
		}
	}()
}

func (g *Game) ShutdownEventListening() {
	if g.stopEventListening == nil {
		panic("SetupEventListening not called")
	}
	close(g.stopEventListening)
	g.stopEventListening = nil
	// importantly need to call termbox.Interrupt() before closing the evChan because
	// termbox.Interrupt() synchronously waits for termbox.PollEvent finishes so there
	// might be one last event coming through into the evChan. If we close it before
	// calling termbox.Interrupt(), we might get a panic.
	termbox.Interrupt()
	close(g.evChan)
	g.evChan = nil
}

// This is a non-blocking call
func (g *Game) TryGetEvent() termbox.Event {
	if g.evChan == nil {
		panic("SetupEventListening not called")
	}
	select {
	case ev := <-g.evChan:
		return ev
	default:
		return termbox.Event{Type: termbox.EventNone}
	}
}

func (g *Game) SetTimePauseOnly() {
	g.timePauseOnly = true
}

func (g *Game) Pause() {
	g.Clock.pause()
	if !g.timePauseOnly {
		g.SpriteMgr.pause()
	}
}

func (g *Game) Resume() {
	g.Clock.resume()
	if !g.timePauseOnly {
		g.SpriteMgr.resume()
	}
}

func (g *Game) IsPaused() bool {
	return g.Clock.isPaused()
}

func (g *Game) GameOver() {
	g.gameOver = true
	g.timePauseOnly = false // Game over, freeze the world regardless what the initial setting is
	g.Pause()
}

func (g *Game) IsGameOver() bool {
	return g.gameOver
}
