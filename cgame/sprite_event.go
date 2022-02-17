package cgame

import (
	"fmt"
)

type spriteEventType int

const (
	spriteEventCreate spriteEventType = iota
	spriteEventDelete
	spriteEventDeleteAll
	spriteEventFunc
	spriteEventCount
)

func (t spriteEventType) String() string {
	switch t {
	case spriteEventCreate:
		return "spriteEventCreate"
	case spriteEventDelete:
		return "spriteEventDelete"
	case spriteEventDeleteAll:
		return "spriteEventDeleteAll"
	case spriteEventFunc:
		return "spriteEventFunc"
	default:
		panic(fmt.Sprintf("unknown spriteEventType value: %d", int(t)))
	}
}

type spriteEvent struct {
	typ  spriteEventType
	s    Sprite
	body interface{}
}

func newSpriteEventCreate(s Sprite) *spriteEvent {
	return &spriteEvent{typ: spriteEventCreate, s: s}
}

func newSpriteEventDelete(s Sprite) *spriteEvent {
	return &spriteEvent{typ: spriteEventDelete, s: s}
}

func newSpriteEventDeleteAll() *spriteEvent {
	return &spriteEvent{typ: spriteEventDeleteAll}
}

func newSpriteEventFunc(f func()) *spriteEvent {
	return &spriteEvent{typ: spriteEventFunc, body: f}
}
