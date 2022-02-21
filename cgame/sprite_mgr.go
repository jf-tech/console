package cgame

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jf-tech/console/cwin"
)

type SpriteManager struct {
	g                  *Game
	ss                 map[*SpriteBase]Sprite
	collidableRegistry *CollidableRegistry
}

func (sm *SpriteManager) CollidableRegistry() *CollidableRegistry {
	return sm.collidableRegistry
}

// FindByName returns the first sprite managed by the SpriteManager that has the same
// name. If no match is found, panic.
func (sm *SpriteManager) FindByName(name string) Sprite {
	if s, ok := sm.TryFindByName(name); ok {
		return s
	}
	panic(fmt.Sprintf("Cannot find sprite named '%s'", name))
}

// TryFindByName returns the first sprite managed by the SpriteManager that has the same
// name.
func (sm *SpriteManager) TryFindByName(name string) (Sprite, bool) {
	for _, s := range sm.ss {
		if s.Name() == name {
			return s, true
		}
	}
	return nil, false
}

// Find returns the unique sprite managed by the SpriteManager that is identified
// by its SpriteBase. If no Sprite is found, panic.
func (sm *SpriteManager) Find(s Sprite) Sprite {
	if ret, ok := sm.TryFind(s); ok {
		return ret
	}
	panic(fmt.Sprintf("Cannot find sprite %s", s.String()))
}

// TryFind returns the unique sprite managed by the SpriteManager that is identified
// by its SpriteBase.
func (sm *SpriteManager) TryFind(s Sprite) (Sprite, bool) {
	ret, ok := sm.ss[s.Base()]
	return ret, ok
}

func (sm *SpriteManager) Process() {
	sm.processAnimators()
}

func (sm *SpriteManager) AddSprite(s Sprite) {
	sm.ss[s.Base()] = s
}

func (sm *SpriteManager) DeleteSprite(s Sprite) {
	s.Destroy()
	delete(sm.ss, s.Base())
}

func (sm *SpriteManager) DeleteAll() {
	for _, s := range sm.ss {
		s.Destroy()
	}
	sm.ss = map[*SpriteBase]Sprite{}
}

// CheckCollision does the collision detection with the collider against all the managed Sprites
// using the collider's Rect and Frame. Note that colliderR and colliderF are not necessarily
// the same as the collider's current position and frame.
func (sm *SpriteManager) CheckCollision(
	collider Sprite, colliderR cwin.Rect, colliderF Frame) []Sprite {
	var collidees []Sprite
	for _, collidee := range sm.ss {
		if collider.Base() == collidee.Base() {
			continue
		}
		if !sm.collidableRegistry.canCollide(collider.Name(), collidee.Name()) {
			continue
		}
		r2 := collidee.Rect()
		f2 := collidee.Frame()
		if DetectCollision(colliderR, colliderF, r2, f2) {
			collidees = append(collidees, collidee)
		}
	}
	return collidees
}

func (sm *SpriteManager) Sprites() []Sprite {
	var ret []Sprite
	for _, s := range sm.ss {
		ret = append(ret, s)
	}
	return ret
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

func (sm *SpriteManager) processAnimators() {
	for _, s := range sm.ss {
		as := s.Animators()
		for _, a := range as {
			a.Animate()
		}
	}
}

func newSpriteManager(g *Game) *SpriteManager {
	return &SpriteManager{
		g:                  g,
		ss:                 map[*SpriteBase]Sprite{},
		collidableRegistry: newCollidableRegistry(),
	}
}
