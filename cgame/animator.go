package cgame

type Animator interface {
	Animate()
}

type AnimatorCfgCommon struct {
	InBoundsCheckType        InBoundsCheckType      // partially visible check by default
	CollisionDetectionType   CollisionDetectionType // on by default
	KeepAliveWhenFinished    bool
	AfterUpdate, AfterFinish func()
}
