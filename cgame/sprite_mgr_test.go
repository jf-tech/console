package cgame

import (
	"strings"
	"testing"

	"github.com/jf-tech/console/cwin"
	"github.com/stretchr/testify/assert"
)

func TestDetectCollision(t *testing.T) {
	// case 1: rects not overlapping.
	r1 := cwin.Rect{X: 0, Y: 0, W: 5, H: 5}
	r2 := cwin.Rect{X: 5, Y: 0, W: 5, H: 5}
	assert.False(t, detectCollision(r1, r2, nil, nil))

	// case 2: rects overlapping, but only transparent cells intersect
	r1 = cwin.Rect{X: 0, Y: 0, W: 5, H: 5}
	r2 = cwin.Rect{X: 2, Y: 1, W: 3, H: 3}
	cells1 := StringToCells(strings.Trim(`
12345
12
1234
12
12345
`, "\n"), cwin.ChAttr{})
	cells2 := StringToCells(strings.Trim(`
abc
  c
abc
`, "\n"), cwin.ChAttr{})
	assert.False(t, detectCollision(r1, r2, cells1, cells2))

	// case 3: rects overlapping, and one non-transparent cell intersects
	r1 = cwin.Rect{X: 0, Y: 0, W: 3, H: 3}
	r2 = cwin.Rect{X: 2, Y: 2, W: 2, H: 2}
	cells1 = StringToCells(strings.Trim(`
123
1 3
123
`, "\n"), cwin.ChAttr{})
	cells2 = StringToCells(strings.Trim(`
ab
ab
`, "\n"), cwin.ChAttr{})
	assert.True(t, detectCollision(r1, r2, cells1, cells2))
}
