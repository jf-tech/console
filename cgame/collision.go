package cgame

import (
	"github.com/jf-tech/console/cwin"
)

type CollidableRegistry struct {
	reg map[string]map[string]bool
}

const (
	CollidableRegistryMatchAll = "*"
)

// supported:
//  - name1, name2
//  - name, name
//  - name, * or *, name
//  - *, *
func (cd *CollidableRegistry) Register(spriteName1, spriteName2 string) *CollidableRegistry {
	if cd.reg == nil {
		cd.reg = map[string]map[string]bool{}
	}
	register := func(a, b string) {
		if _, found := cd.reg[a]; !found {
			cd.reg[a] = map[string]bool{}
		}
		cd.reg[a][b] = true
	}
	register(spriteName1, spriteName2)
	register(spriteName2, spriteName1)
	return cd
}

func (cd *CollidableRegistry) RegisterBulk(
	spriteName1 string, otherSpriteNames []string) *CollidableRegistry {
	for _, other := range otherSpriteNames {
		cd.Register(spriteName1, other)
	}
	return cd
}

func (cd *CollidableRegistry) canCollide(spriteName1, spriteName2 string) bool {
	find := func(a, b string) bool {
		if cd.reg == nil {
			return false
		}
		if _, found := cd.reg[a]; !found {
			return false
		}
		if _, found := cd.reg[a][b]; !found {
			return false
		}
		return true
	}
	return find(CollidableRegistryMatchAll, CollidableRegistryMatchAll) ||
		find(spriteName1, CollidableRegistryMatchAll) ||
		find(spriteName2, CollidableRegistryMatchAll) ||
		find(spriteName1, spriteName2)
}

func newCollidableRegistry() *CollidableRegistry {
	return &CollidableRegistry{}
}

type CollisionDetectionType int

const (
	CollisionDetectionOn = CollisionDetectionType(iota)
	CollisionDetectionOff
)

var (
	collisionDetectionBuf = make([]bool, 0, 100)
)

// Be aware certain quirky situation: given this is a terminal/character based system
// it's possible to have two sprites actually cross/on top of each other without actually
// having any overlapping characters. Think the following simple example:
// - sprite 1 looks like this:
//     \
//      \
// - sprite 2 looks like this:
//      /
//     /
// They can be perfectly on top of each other/across, without causing an overlapping character
// situation:
//     \/
//     /\
// Thus DetectCollision fails here. No good general solutions in mind yet, just need to be
// careful. For "enclosed" sprite, such as circle, one can fill the interior with characters
// not just TransparencyChx, that will alleviate the problem.
func DetectCollision(r1 cwin.Rect, f1 Frame, r2 cwin.Rect, f2 Frame) bool {
	if ro, overlapped := r1.Overlap(r2); overlapped {
		collisionDetectionBuf = collisionDetectionBuf[:0]
		for i := 0; i < ro.W*ro.H; i++ {
			collisionDetectionBuf = append(collisionDetectionBuf, false)
		}
		for i := 0; i < len(f1); i++ {
			if f1[i].Chx == cwin.TransparentChx() || !ro.Contain(r1.X+f1[i].X, r1.Y+f1[i].Y) {
				continue
			}
			collisionDetectionBuf[(r1.Y+f1[i].Y-ro.Y)*ro.W+(r1.X+f1[i].X-ro.X)] = true
		}
		for i := 0; i < len(f2); i++ {
			if f2[i].Chx == cwin.TransparentChx() || !ro.Contain(r2.X+f2[i].X, r2.Y+f2[i].Y) {
				continue
			}
			if collisionDetectionBuf[(r2.Y+f2[i].Y-ro.Y)*ro.W+(r2.X+f2[i].X-ro.X)] {
				return true
			}
		}
	}
	return false
}

type CollisionResponseType int

const (
	CollisionResponseAbandon = CollisionResponseType(iota)
	CollisionResponseJustDoIt
)

type CollisionResponse interface {
	CollisionNotify(initiator bool, collidedWith []Sprite) CollisionResponseType
}
