package cgame

type Collidable interface {
	Collided(other Sprite)
}
