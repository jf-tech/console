package cwin

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRect_Contain(t *testing.T) {
	for _, test := range []struct {
		name     string
		x, y     int
		r        Rect
		expected bool
	}{
		{
			name:     "x < r.x",
			x:        -1,
			y:        0,
			r:        Rect{X: 0, Y: 0, W: 1, H: 1},
			expected: false,
		},
		{
			name:     "x >= r.x + r.w",
			x:        10,
			y:        0,
			r:        Rect{X: 0, Y: 0, W: 10, H: 1},
			expected: false,
		},
		{
			name:     "y < r.y",
			x:        0,
			y:        -1,
			r:        Rect{X: 0, Y: 0, W: 1, H: 1},
			expected: false,
		},
		{
			name:     "h >= r.y + r.h",
			x:        0,
			y:        10,
			r:        Rect{X: 0, Y: 0, W: 1, H: 10},
			expected: false,
		},
		{
			name:     "on UL corner",
			x:        5,
			y:        10,
			r:        Rect{X: 5, Y: 10, W: 10, H: 10},
			expected: true,
		},
		{
			name:     "on LR corner",
			x:        10,
			y:        19,
			r:        Rect{X: 5, Y: 10, W: 6, H: 10},
			expected: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.r.Contain(test.x, test.y))
		})
	}
}

func TestRect_Overlap(t *testing.T) {
	for _, test := range []struct {
		name     string
		r1       Rect
		r2       Rect
		expected *Rect
	}{
		{
			name:     "x-axis not overlap",
			r1:       Rect{X: 0, Y: 0, W: 5, H: 5},
			r2:       Rect{X: -5, Y: 0, W: 5, H: 5},
			expected: nil,
		},
		{
			name:     "y-axis not overlap",
			r1:       Rect{X: 0, Y: 0, W: 5, H: 5},
			r2:       Rect{X: 0, Y: 5, W: 5, H: 5},
			expected: nil,
		},
		{
			name:     "overlap",
			r1:       Rect{X: 0, Y: 0, W: 5, H: 5},
			r2:       Rect{X: 1, Y: 1, W: 5, H: 5},
			expected: &Rect{X: 1, Y: 1, W: 4, H: 4},
		},
		{
			name:     "contain",
			r1:       Rect{X: 1, Y: 1, W: 3, H: 3},
			r2:       Rect{X: 0, Y: 0, W: 5, H: 5},
			expected: &Rect{X: 1, Y: 1, W: 3, H: 3},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			r, f := test.r1.Overlap(test.r2)
			assert.Equal(t, test.expected != nil, f)
			if f {
				assert.Equal(t, *test.expected, r)
			}
		})
	}
}

func TestRect_String(t *testing.T) {
	assert.Equal(t, "(1,2,3,4)", fmt.Sprintf("%v", Rect{1, 2, 3, 4}))
}
