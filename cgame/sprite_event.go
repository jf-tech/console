package cgame

import (
	"fmt"
)

type SpriteEventType int

const (
	SpriteEventCreate SpriteEventType = iota
	SpriteEventDelete
	SpriteEventSetPosRelative
)

func (t SpriteEventType) String() string {
	switch t {
	case SpriteEventCreate:
		return "SpriteEventCreate"
	case SpriteEventDelete:
		return "SpriteEventDelete"
	case SpriteEventSetPosRelative:
		return "SpriteEventSetPosRelative"
	default:
		panic(fmt.Sprintf("Unknown SpriteEvent value: %d", int(t)))
	}
}

type SpriteEvent struct {
	eventType SpriteEventType
	s         Sprite
	body      interface{}
}

func NewSpriteEventCreate(s Sprite) *SpriteEvent {
	return &SpriteEvent{
		eventType: SpriteEventCreate,
		s:         s,
	}
}

func NewSpriteEventDelete(s Sprite) *SpriteEvent {
	return &SpriteEvent{
		eventType: SpriteEventDelete,
		s:         s,
	}
}

func NewSpriteEventSetPosRelative(s Sprite, dx, dy int) *SpriteEvent {
	return &SpriteEvent{
		eventType: SpriteEventSetPosRelative,
		s:         s,
		body:      pairInt{a: dx, b: dy},
	}
}
