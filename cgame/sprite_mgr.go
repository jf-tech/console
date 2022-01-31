package cgame

import (
	"fmt"
	"strings"

	"github.com/jf-tech/console/cwin"
)

type SpriteManager struct {
	g *Game

	ss     []Sprite
	eventQ *threadSafeFIFO
	paused bool

	spriteEventsProcessed int64
}

// Note the names of sprite instances are not required to be unique, this method
// return the first matching one, if any.
func (sm *SpriteManager) FindByName(name string) Sprite {
	for i := 0; i < len(sm.ss); i++ {
		if sm.ss[i].Cfg().Name == name {
			return sm.ss[i]
		}
	}
	panic(fmt.Sprintf("Unable to find sprite '%s'", name))
}

func (sm *SpriteManager) AddEvent(e *SpriteEvent) {
	sm.eventQ.push(e)
}

func (sm *SpriteManager) ProcessAll() {
	if sm.paused {
		return
	}
	sm.processEvents()  // keyboards triggered events (move, sprite creation, etc)
	sm.processSprites() // Animated sprite self movements
	sm.processEvents()  // consequences from self-movements
	sm.processCollisions()
	sm.processEvents() // consequences of collisions
}

func (sm *SpriteManager) DbgStats() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Active sprites: %d\n", len(sm.ss)))
	sb.WriteString(fmt.Sprintf("Sprite Events processed: %d\n", sm.spriteEventsProcessed))
	return sb.String()
}

func (sm *SpriteManager) processEvents() {
	for {
		se, ok := sm.eventQ.tryPop()
		if !ok {
			break
		}
		sm.processEvent(se.(*SpriteEvent))
		sm.spriteEventsProcessed++
	}
}

func (sm *SpriteManager) processEvent(e *SpriteEvent) {
	switch e.Type {
	case SpriteEventCreate, SpriteEventDelete, SpriteEventSetPosRelative:
		idx := sm.locateSprite(e.S)
		existsCheck := func() {
			if idx < 0 {
				panic(fmt.Sprintf("Sprite['%s',%d] not found for %s",
					e.S.Cfg().Name, e.S.UID(), e.Type))
			}
		}
		switch e.Type {
		case SpriteEventCreate:
			if idx >= 0 {
				panic(fmt.Sprintf("Sprite['%s',%d] is being re-added",
					e.S.Cfg().Name, e.S.UID()))
			}
			sm.ss = append(sm.ss, e.S)
		case SpriteEventDelete:
			if idx < 0 {
				return
			}
			sm.g.WinSys.RemoveWin(sm.ss[idx].Win())
			for ; idx < len(sm.ss)-1; idx++ {
				sm.ss[idx] = sm.ss[idx+1]
			}
			sm.ss[idx] = nil
			sm.ss = sm.ss[:len(sm.ss)-1]
		case SpriteEventSetPosRelative:
			existsCheck()
			if ps, ok := e.S.(PositionSettable); ok {
				xy := e.Body.(pairInt)
				ps.SetPosRelative(xy.a, xy.b)
			}
		}
	}
}

func (sm *SpriteManager) processSprites() {
	for i := 0; i < len(sm.ss); i++ {
		if auto, ok := sm.ss[i].(Animated); ok {
			auto.Act()
		}
	}
}

func (sm *SpriteManager) processCollisions() {
	for i := 0; i < len(sm.ss)-1; i++ {
		for j := i + 1; j < len(sm.ss); j++ {
			var ci, cj Collidable
			var ok bool
			if ci, ok = sm.ss[i].(Collidable); !ok {
				continue
			}
			if cj, ok = sm.ss[j].(Collidable); !ok {
				continue
			}
			if detectCollision(
				sm.ss[i].Win().Rect(), sm.ss[j].Win().Rect(),
				sm.ss[i].Cfg().Cells, sm.ss[j].Cfg().Cells) {
				ci.Collided(sm.ss[j])
				cj.Collided(sm.ss[i])
			}
		}
	}
}

func (sm *SpriteManager) locateSprite(s Sprite) int {
	for i := 0; i < len(sm.ss); i++ {
		if sm.ss[i].UID() == s.UID() {
			return i
		}
	}
	return -1
}

func (sm *SpriteManager) pause() {
	sm.paused = true
}

func (sm *SpriteManager) resume() {
	sm.paused = false
}

func detectCollision(r1, r2 cwin.Rect, cells1, cells2 []Cell) bool {
	// first do a rough rect overlap test to weed out most negative cases.
	if overlapped, ro := r1.Overlap(r2); overlapped {
		for y := 0; y < ro.H; y++ {
			for x := 0; x < ro.W; x++ {
				x1 := ro.X + x - r1.X
				y1 := ro.Y + y - r1.Y
				idx1 := -1
				for i := 0; i < len(cells1); i++ {
					if cells1[i].X == x1 && cells1[i].Y == y1 {
						idx1 = i
						break
					}
				}
				if idx1 < 0 || cells1[idx1].Chx == cwin.TransparentChx() {
					continue
				}
				x2 := ro.X + x - r2.X
				y2 := ro.Y + y - r2.Y
				idx2 := -1
				for i := 0; i < len(cells2); i++ {
					if cells2[i].X == x2 && cells2[i].Y == y2 {
						idx2 = i
						break
					}
				}
				if idx2 < 0 || cells2[idx2].Chx == cwin.TransparentChx() {
					continue
				}
				return true
			}
		}
	}
	return false
}

const (
	defaultSpriteManagerBufCap = 1000
)

func newSpriteManager(g *Game) *SpriteManager {
	return &SpriteManager{
		g:      g,
		ss:     make([]Sprite, 0, defaultSpriteManagerBufCap),
		eventQ: newFIFO(defaultSpriteManagerBufCap),
	}
}
