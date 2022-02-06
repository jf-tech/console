package cgame

import (
	"fmt"
	"strings"

	"github.com/jf-tech/console/cwin"
)

type spriteEntry struct {
	s         Sprite
	animators []Animator
}

type SpriteManager struct {
	g                     *Game
	ss                    []spriteEntry
	eventQ                *ThreadSafeFIFO
	collisionDetectionBuf []bool
	spriteEventsProcessed int64
}

// Note the names of sprite instances are not required to be unique, this method
// return the first matching one, if any.
func (sm *SpriteManager) FindByName(name string) Sprite {
	for i := 0; i < len(sm.ss); i++ {
		if sm.ss[i].s.Name() == name {
			return sm.ss[i].s
		}
	}
	panic(fmt.Sprintf("Cannot find sprite named '%s'", name))
}

func (sm *SpriteManager) TryFindByName(name string) (Sprite, bool) {
	for i := 0; i < len(sm.ss); i++ {
		if sm.ss[i].s.Name() == name {
			return sm.ss[i].s, true
		}
	}
	return nil, false
}

func (sm *SpriteManager) TryFindByUID(uid int64) (Sprite, bool) {
	for i := 0; i < len(sm.ss); i++ {
		if sm.ss[i].s.UID() == uid {
			return sm.ss[i].s, true
		}
	}
	return nil, false
}

func (sm *SpriteManager) AddEvent(e *SpriteEvent) {
	sm.eventQ.Push(e)
}

func (sm *SpriteManager) Process() {
	sm.processEvents()     // keyboards triggered events (move, sprite creation, etc)
	sm.processAnimators()  // Animated sprite self movements
	sm.processEvents()     // consequences from self-movements
	sm.processCollisions() // collisions
	sm.processEvents()     // consequences of collisions
}

func (sm *SpriteManager) DbgStats() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Active sprites: %d\n", len(sm.ss)))
	animatorN := 0
	for _, e := range sm.ss {
		animatorN += len(e.animators)
	}
	sb.WriteString(fmt.Sprintf("Active animators: %d\n", animatorN))
	sb.WriteString(fmt.Sprintf("Sprite Events processed: %d\n", sm.spriteEventsProcessed))
	return sb.String()
}

func (sm *SpriteManager) processEvents() {
	for {
		se, ok := sm.eventQ.TryPop()
		if !ok {
			break
		}
		sm.processEvent(se.(*SpriteEvent))
		sm.spriteEventsProcessed++
	}
}

func (sm *SpriteManager) processEvent(e *SpriteEvent) {
	switch e.eventType {
	case SpriteEventCreate, SpriteEventDelete, SpriteEventSetPosRelative:
		idx := sm.spriteIndex(e.s)
		existsCheck := func() {
			if idx < 0 {
				panic(fmt.Sprintf("Sprite['%s',%d] not found for %s",
					e.s.Name(), e.s.UID(), e.eventType))
			}
		}
		switch e.eventType {
		case SpriteEventCreate:
			if idx >= 0 {
				panic(fmt.Sprintf("Sprite['%s',%d] is being re-added",
					e.s.Name(), e.s.UID()))
			}
			sm.ss = append(sm.ss, spriteEntry{s: e.s, animators: e.body.([]Animator)})
		case SpriteEventDelete:
			if idx < 0 {
				return
			}
			sm.g.WinSys.RemoveWin(sm.ss[idx].s.Win())
			for ; idx < len(sm.ss)-1; idx++ {
				sm.ss[idx] = sm.ss[idx+1]
			}
			sm.ss = sm.ss[:len(sm.ss)-1]
		case SpriteEventSetPosRelative:
			existsCheck()
			if ps, ok := e.s.(PositionSettable); ok {
				xy := e.body.(pairInt)
				ps.SetPosRelative(xy.a, xy.b)
			}
		}
	case SpriteEventDeleteAll:
		for _, e := range sm.ss {
			sm.g.WinSys.RemoveWin(e.s.Win())
		}
		sm.ss = sm.ss[:0]
	}
}

func (sm *SpriteManager) processAnimators() {
	for _, e := range sm.ss {
		for i := 0; i < len(e.animators); i++ {
			if e.animators[i].Animate(e.s) == AnimatorCompleted {
				for j := i; j < len(e.animators)-1; j++ {
					e.animators[j] = e.animators[j+1]
				}
				e.animators = e.animators[:len(e.animators)-1]
			}
		}
	}
}

func (sm *SpriteManager) processCollisions() {
	for i := 0; i < len(sm.ss)-1; i++ {
		for j := i + 1; j < len(sm.ss); j++ {
			var ci, cj Collidable
			var ok bool
			if ci, ok = sm.ss[i].s.(Collidable); !ok {
				continue
			}
			if cj, ok = sm.ss[j].s.(Collidable); !ok {
				continue
			}
			if sm.detectCollision(sm.ss[i].s.Win(), sm.ss[j].s.Win()) {
				ci.Collided(sm.ss[j].s)
				cj.Collided(sm.ss[i].s)
			}
		}
	}
}

func (sm *SpriteManager) spriteIndex(s Sprite) int {
	for i := 0; i < len(sm.ss); i++ {
		if sm.ss[i].s.UID() == s.UID() {
			return i
		}
	}
	return -1
}

func (sm *SpriteManager) detectCollision(w1, w2 *cwin.Win) bool {
	// do a rough rect overlap test to weed out most negative cases.
	r1, r2 := w1.Rect(), w2.Rect()
	overlapped, ro := r1.Overlap(r2)
	if !overlapped {
		return false
	}
	sm.collisionDetectionBuf = sm.collisionDetectionBuf[:0]
	for y := 0; y < ro.H; y++ {
		for x := 0; x < ro.W; x++ {
			sm.collisionDetectionBuf = append(sm.collisionDetectionBuf,
				w1.GetClient(x+ro.X-r1.X, y+ro.Y-r1.Y) != cwin.TransparentChx())
		}
	}
	for y := 0; y < ro.H; y++ {
		for x := 0; x < ro.W; x++ {
			if sm.collisionDetectionBuf[y*ro.W+x] &&
				w2.GetClient(x+ro.X-r2.X, y+ro.Y-r2.Y) != cwin.TransparentChx() {
				return true
			}
		}
	}
	return false
}

const (
	defaultSpriteBufCap             = 1000
	defaultCollisionDetectionBufCap = 100
)

func newSpriteManager(g *Game) *SpriteManager {
	return &SpriteManager{
		g:                     g,
		ss:                    make([]spriteEntry, 0, defaultSpriteBufCap),
		eventQ:                NewThreadSafeFIFO(defaultSpriteBufCap),
		collisionDetectionBuf: make([]bool, 0, defaultCollisionDetectionBufCap),
	}
}
