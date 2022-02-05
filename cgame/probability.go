package cgame

import (
	"fmt"
	"math/rand"
	"time"
)

func CheckProbability(prob string) bool {
	probV := parseProbability(prob)
	probVint := int(probV * 10)
	return rand.Int()%1000 < probVint
}

func parseProbability(prob string) float64 {
	var probV float64
	_, err := fmt.Sscanf(prob, "%f%%", &probV)
	if err != nil {
		panic(fmt.Sprintf("Invalid probabilty '%s'", prob))
	}
	return probV
}

type PeriodicProbabilityChecker struct {
	prob              string
	period            time.Duration
	clock             *Clock
	lastProbCheckTime time.Duration
}

func (twp *PeriodicProbabilityChecker) Check() bool {
	if twp.clock == nil {
		panic("Forgot to Reset a clock to this twp?")
	}
	now := twp.clock.Now()
	if now-twp.lastProbCheckTime > twp.period {
		twp.lastProbCheckTime = now
		return CheckProbability(twp.prob)
	}
	return false
}

func (twp *PeriodicProbabilityChecker) Reset(clock *Clock) {
	twp.clock = clock
	twp.lastProbCheckTime = clock.Now()
}

func NewPeriodicProbabilityChecker(prob string, period time.Duration) *PeriodicProbabilityChecker {
	return &PeriodicProbabilityChecker{
		prob:   prob,
		period: period,
	}
}
