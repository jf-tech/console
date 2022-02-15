package cgame

import (
	"fmt"
)

type spriteEventType int

const (
	spriteEventCreate spriteEventType = iota
	spriteEventDelete
	spriteEventDeleteAll
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
	default:
		panic(fmt.Sprintf("unknown spriteEventType value: %d", int(t)))
	}
}

type spriteEvent struct {
	eventType spriteEventType
	s         Sprite
	body      interface{}
}

func newSpriteEventCreate(s Sprite) *spriteEvent {
	return &spriteEvent{eventType: spriteEventCreate, s: s}
}

func newSpriteEventDelete(s Sprite) *spriteEvent {
	return &spriteEvent{eventType: spriteEventDelete, s: s}
}

func newSpriteEventDeleteAll() *spriteEvent {
	return &spriteEvent{eventType: spriteEventDeleteAll}
}
