package cgame

import (
	"fmt"
)

type SpriteEventType int

const (
	SpriteEventCreate SpriteEventType = iota
	SpriteEventDelete
	SpriteEventDeleteAll
	SpriteEventSetPosRelative
	spriteEventCount
)

func (t SpriteEventType) String() string {
	switch t {
	case SpriteEventCreate:
		return "SpriteEventCreate"
	case SpriteEventDelete:
		return "SpriteEventDelete"
	case SpriteEventDeleteAll:
		return "SpriteEventDeleteAll"
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

func NewSpriteEventCreate(s Sprite, animators ...Animator) *SpriteEvent {
	return &SpriteEvent{
		eventType: SpriteEventCreate,
		s:         s,
		body:      animators,
	}
}

func NewSpriteEventDelete(s Sprite) *SpriteEvent {
	return &SpriteEvent{
		eventType: SpriteEventDelete,
		s:         s,
	}
}

func NewSpriteEventDeleteAll() *SpriteEvent {
	return &SpriteEvent{eventType: SpriteEventDeleteAll}
}

func NewSpriteEventSetPosRelative(s Sprite, dx, dy int) *SpriteEvent {
	return &SpriteEvent{
		eventType: SpriteEventSetPosRelative,
		s:         s,
		body:      PairInt{A: dx, B: dy},
	}
}
