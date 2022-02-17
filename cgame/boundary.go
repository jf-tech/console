package cgame

import "github.com/jf-tech/console/cwin"

type InBoundsCheckType int

const (
	InBoundsCheckPartiallyVisible = InBoundsCheckType(iota)
	InBoundsCheckFullyVisible
	InBoundsCheckNone
)

type InBoundsCheckResult int

const (
	InBoundsCheckResultOK = InBoundsCheckResult(iota)
	InBoundsCheckResultN  // breach to the north
	InBoundsCheckResultE  // breach to the east
	InBoundsCheckResultS  // breach to the south
	InBoundsCheckResultW  // breach to the west
)

func InBoundsCheck(checkType InBoundsCheckType, r cwin.Rect, f Frame, parentR cwin.Rect) InBoundsCheckResult {
	// Note the parent rect is really parent windows' client rect shifted to origin - i.e.
	// from POV of the sprite. So we only use it's W/H components.
	rangeTest := func(val int, rangeOfVal int) int {
		if val < 0 {
			return -1
		} else if val < rangeOfVal {
			return 0
		}
		return 1
	}
	totalCells := 0
	result := map[InBoundsCheckResult]int{}
	for i := 0; i < len(f); i++ {
		if f[i].Chx == cwin.TransparentChx() {
			continue
		}
		totalCells++
		x := r.X + f[i].X
		y := r.Y + f[i].Y
		xReg, yReg := rangeTest(x, parentR.W), rangeTest(y, parentR.H)
		switch xReg {
		case -1:
			switch yReg {
			case -1:
				if -x > -y {
					result[InBoundsCheckResultW]++
				} else {
					result[InBoundsCheckResultN]++
				}
			case 0:
				result[InBoundsCheckResultW]++
			case 1:
				if -x > y-parentR.H+1 {
					result[InBoundsCheckResultW]++
				} else {
					result[InBoundsCheckResultS]++
				}
			}
		case 0:
			switch yReg {
			case -1:
				result[InBoundsCheckResultN]++
			case 0:
				result[InBoundsCheckResultOK]++
			case 1:
				result[InBoundsCheckResultS]++
			}
		case 1:
			switch yReg {
			case -1:
				if x-parentR.W+1 > -y {
					result[InBoundsCheckResultE]++
				} else {
					result[InBoundsCheckResultN]++
				}
			case 0:
				result[InBoundsCheckResultE]++
			case 1:
				if x-parentR.W+1 > y-parentR.H+1 {
					result[InBoundsCheckResultE]++
				} else {
					result[InBoundsCheckResultS]++
				}
			}
		}
	}
	var maxNonOkResult InBoundsCheckResult
	maxNonOkResultCount := 0
	for k, v := range result {
		if k != InBoundsCheckResultOK && v > maxNonOkResultCount {
			maxNonOkResult, maxNonOkResultCount = k, v
		}
	}
	switch checkType {
	case InBoundsCheckFullyVisible:
		if result[InBoundsCheckResultOK] == totalCells {
			return InBoundsCheckResultOK
		}
		return maxNonOkResult
	case InBoundsCheckPartiallyVisible:
		if result[InBoundsCheckResultOK] > 0 {
			return InBoundsCheckResultOK
		}
		return maxNonOkResult
	}
	return InBoundsCheckResultOK
}

type InBoundsCheckResponseType int

const (
	InBoundsCheckResponseAbandon = InBoundsCheckResponseType(iota)
	InBoundsCheckResponseJustDoIt
)

type InBoundsCheckResponse interface {
	InBoundsCheckNotify(result InBoundsCheckResult) InBoundsCheckResponseType
}
