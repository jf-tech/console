package main

/*
TODO
import (
	"math"

	"github.com/jf-tech/console/cgame"
)

var hps = map[string]int{
	alphaName:        5,
	alphaBullet1Name: 1,
	betaName:         1,
	betaBullet1Name:  1,
}

func hp(s cgame.Sprite) int {
	if v, ok := hps[s.Cfg().Name]; ok {
		return v
	}
	return math.MaxInt32
}

var hitDamages = []hitDamageEntry{
	{s1: alphaName, s2: betaName, damage: 5},
	{s1: alphaName, s2: betaBullet1Name, damage: 1},
	{s1: betaName, s2: alphaBullet1Name, damage: 1},
}

func hitDamage(s1, s2 cgame.Sprite) int {
	s1name := s1.Cfg().Name
	s2name := s2.Cfg().Name
	for _, e := range hitDamages {
		if (s1name == e.s1 && s2name == e.s2) || (s1name == e.s2 && s2name == e.s1) {
			return e.damage
		}
	}
	return 0
}

type hitDamageEntry struct {
	s1, s2 string
	damage int
}
*/
