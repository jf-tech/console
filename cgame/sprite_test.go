package cgame

import (
	"strings"
	"testing"

	"github.com/jf-tech/console/cwin"
	"github.com/nsf/termbox-go"
	"github.com/stretchr/testify/assert"
)

func TestStringToCells(t *testing.T) {
	testAttr := cwin.ChAttr{Fg: termbox.ColorBlue, Bg: termbox.ColorWhite}
	for _, test := range []struct {
		name string
		s    string
		attr cwin.ChAttr
		exp  []Cell
	}{
		{
			name: "4x2",
			s: `
\┃┃/
 \/
`,
			attr: testAttr,
			exp: []Cell{
				{X: 0, Y: 0, Chx: cwin.Chx{Ch: '\\', Attr: testAttr}},
				{X: 1, Y: 0, Chx: cwin.Chx{Ch: '┃', Attr: testAttr}},
				{X: 2, Y: 0, Chx: cwin.Chx{Ch: '┃', Attr: testAttr}},
				{X: 3, Y: 0, Chx: cwin.Chx{Ch: '/', Attr: testAttr}},
				{X: 0, Y: 1, Chx: cwin.TransparentChx()},
				{X: 1, Y: 1, Chx: cwin.Chx{Ch: '\\', Attr: testAttr}},
				{X: 2, Y: 1, Chx: cwin.Chx{Ch: '/', Attr: testAttr}},
				{X: 3, Y: 1, Chx: cwin.TransparentChx()},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.exp,
				StringToCells(strings.Trim(test.s, "\n"), test.attr))
		})
	}
}
