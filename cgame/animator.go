package cgame

type AnimatorState int

const (
	AnimatorRunning AnimatorState = iota
	AnimatorCompleted
)

type Animator interface {
	Animate(Sprite) AnimatorState
}
