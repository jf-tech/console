package cgame

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jf-tech/console/cwin"
)

type SpriteManager struct {
	g                  *Game
	ss                 []Sprite
	eventQ             *ThreadSafeFIFO
	collidableRegistry *CollidableRegistry
}

func (sm *SpriteManager) CollidableRegistry() *CollidableRegistry {
	return sm.collidableRegistry
}

// Note the names of sprite instances are not required to be unique, this method
// return the first matching one, if any.
func (sm *SpriteManager) FindByName(name string) Sprite {
	if s, ok := sm.TryFindByName(name); ok {
		return s
	}
	panic(fmt.Sprintf("Cannot find sprite named '%s'", name))
}

func (sm *SpriteManager) TryFindByName(name string) (Sprite, bool) {
	for i := 0; i < len(sm.ss); i++ {
		if sm.ss[i].Name() == name {
			return sm.ss[i], true
		}
	}
	return nil, false
}

func (sm *SpriteManager) FindByUID(uid int64) Sprite {
	if s, ok := sm.TryFindByUID(uid); ok {
		return s
	}
	panic(fmt.Sprintf("Cannot find sprite with uid %d", uid))
}

func (sm *SpriteManager) TryFindByUID(uid int64) (Sprite, bool) {
	for i := 0; i < len(sm.ss); i++ {
		if sm.ss[i].UID() == uid {
			return sm.ss[i], true
		}
	}
	return nil, false
}

func (sm *SpriteManager) Process() {
	sm.processEvents()
	sm.processAnimators()
	sm.processEvents()
}

func (sm *SpriteManager) AddSprite(s Sprite) {
	sm.eventQ.Push(newSpriteEventCreate(s))
}

func (sm *SpriteManager) DeleteSprite(s Sprite) {
	sm.eventQ.Push(newSpriteEventDelete(s))
}

func (sm *SpriteManager) DeleteAll() {
	sm.eventQ.Push(newSpriteEventDeleteAll())
}

// Similar notes to CheckParentRectBound's
func (sm *SpriteManager) CheckCollision(s Sprite, newR cwin.Rect, newF Frame) []Sprite {
	var collided []Sprite
	for i := 0; i < len(sm.ss); i++ {
		if s.UID() == sm.ss[i].UID() {
			continue
		}
		if !sm.collidableRegistry.canCollide(s.Name(), sm.ss[i].Name()) {
			continue
		}
		r2 := sm.ss[i].Rect()
		f2 := sm.ss[i].Frame()
		if DetectCollision(newR, newF, r2, f2) {
			collided = append(collided, sm.ss[i])
		}
	}
	return collided
}

func (sm *SpriteManager) Sprites() []Sprite {
	// make a snapshot to return so that caller won't get into potentially
	// changing slice.
	var cp []Sprite
	cp = append(cp, sm.ss...)
	return cp
}

func (sm *SpriteManager) DbgStats() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Active sprites: %d\n", len(sm.ss)))
	spriteNums := map[string]int{}
	for _, s := range sm.ss {
		spriteNums[s.Name()]++
	}
	keys := make([]string, 0, len(spriteNums))
	for k := range spriteNums {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("- '%s': %d\n", k, spriteNums[k]))
	}
	return sb.String()
}

func (sm *SpriteManager) processEvents() {
	for {
		e, ok := sm.eventQ.TryPop()
		if !ok {
			break
		}
		se := e.(*spriteEvent)
		switch se.eventType {
		case spriteEventCreate:
			if sm.spriteIndex(se.s) >= 0 {
				panic(fmt.Sprintf("Sprite['%s',%d] is being re-added", se.s.Name(), se.s.UID()))
			}
			sm.ss = append(sm.ss, se.s)
		case spriteEventDelete:
			se.s.Destroy()
			idx := sm.spriteIndex(se.s)
			if idx < 0 {
				return
			}
			copy(sm.ss[idx:], sm.ss[idx+1:])
			sm.ss = sm.ss[:len(sm.ss)-1]
		case spriteEventDeleteAll:
			for _, s := range sm.ss {
				s.Destroy()
			}
			sm.ss = sm.ss[:0]
		}
	}
}

func (sm *SpriteManager) processAnimators() {
	for _, s := range sm.ss {
		as := s.Animators()
		for _, a := range as {
			a.Animate()
		}
	}
}

func (sm *SpriteManager) spriteIndex(s Sprite) int {
	for i := 0; i < len(sm.ss); i++ {
		if sm.ss[i].UID() == s.UID() {
			return i
		}
	}
	return -1
}

const (
	defaultSpriteBufCap = 1000
)

func newSpriteManager(g *Game) *SpriteManager {
	return &SpriteManager{
		g:                  g,
		ss:                 make([]Sprite, 0, defaultSpriteBufCap),
		eventQ:             NewThreadSafeFIFO(defaultSpriteBufCap),
		collidableRegistry: newCollidableRegistry(),
	}
}
