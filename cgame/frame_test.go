package cgame

import (
	"strings"
	"testing"

	"github.com/jf-tech/console/cwin"
	"github.com/nsf/termbox-go"
	"github.com/stretchr/testify/assert"
)

func TestStringToFrame(t *testing.T) {
	testAttr := cwin.ChAttr{Fg: termbox.ColorBlue, Bg: termbox.ColorWhite}
	for _, test := range []struct {
		name string
		s    string
		attr cwin.ChAttr
		exp  SpriteFrame
	}{
		{
			name: "4x2",
			s: `
\┃┃/
 \/
`,
			attr: testAttr,
			exp: []SpriteCell{
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
				FrameFromString(strings.Trim(test.s, "\n"), test.attr))
		})
	}
}