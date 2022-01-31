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
	Type SpriteEventType
	S    Sprite
	Body interface{}
}

func NewSpriteEventCreate(s Sprite) *SpriteEvent {
	return &SpriteEvent{
		Type: SpriteEventCreate,
		S:    s,
	}
}

func NewSpriteEventDelete(s Sprite) *SpriteEvent {
	return &SpriteEvent{
		Type: SpriteEventDelete,
		S:    s,
	}
}

func NewSpriteEventSetPosRelative(s Sprite, dx, dy int) *SpriteEvent {
	return &SpriteEvent{
		Type: SpriteEventSetPosRelative,
		S:    s,
		Body: pairInt{a: dx, b: dy},
	}
}
