package cgame

type Animator interface {
	Animate()
}

type AnimatorCfgCommon struct {
	InBoundsCheckTypeToFinish      InBoundsCheckType      // none by default
	CollisionDetectionTypeToFinish CollisionDetectionType // on by default
	KeepAliveWhenFinished          bool
	PreUpdateNotify                PreUpdateNotify
	AfterUpdate, AfterFinish       func()
}
