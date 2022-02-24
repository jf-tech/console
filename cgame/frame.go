package cgame

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"time"

	"github.com/jf-tech/console/cterm"
	"github.com/jf-tech/console/cwin"
	"github.com/jf-tech/go-corelib/ios"
	"github.com/jf-tech/go-corelib/maths"
)

type Cell struct {
	X, Y int
	Chx  cwin.Chx
}

type Frame []Cell

func CopyFrame(src Frame) Frame {
	dst := make(Frame, len(src))
	copy(dst, src)
	return dst
}

func SetAttrInFrame(f Frame, attr cwin.Attr) Frame {
	for i := 0; i < len(f); i++ {
		f[i].Chx.Attr = attr
	}
	return f
}

func FrameFromStringEx(s string, attr cwin.Attr, spaceEqualsTransparency bool) Frame {
	s = strings.Trim(s, "\n")
	rect := cwin.TextDimension(s)
	var f Frame
	rs := []rune(s)
	rsLen := len(rs)
	for i, x, y := 0, 0, 0; i < rsLen; i++ {
		if rs[i] == '\n' {
			if !spaceEqualsTransparency {
				for ; x < rect.W; x++ {
					f = append(f, Cell{X: x, Y: y, Chx: cwin.Chx{Ch: ' ', Attr: attr}})
				}
			}
			x = 0
			y++
			continue
		}
		if rs[i] != ' ' || !spaceEqualsTransparency {
			f = append(f, Cell{X: x, Y: y, Chx: cwin.Chx{Ch: rs[i], Attr: attr}})
		}
		x++
	}
	return f
}

func FrameFromString(s string, attr cwin.Attr) Frame {
	return FrameFromStringEx(s, attr, true)
}

var (
	colorMap = map[rune]cterm.Attribute{
		'K': cterm.ColorBlack,
		'R': cterm.ColorRed,
		'G': cterm.ColorGreen,
		'Y': cterm.ColorYellow,
		'B': cterm.ColorBlue,
		'M': cterm.ColorMagenta,
		'C': cterm.ColorCyan,
		'W': cterm.ColorWhite,
		'A': cterm.ColorDarkGray,
		'r': cterm.ColorLightRed,
		'g': cterm.ColorLightGreen,
		'y': cterm.ColorLightYellow,
		'b': cterm.ColorLightBlue,
		'm': cterm.ColorLightMagenta,
		'c': cterm.ColorLightCyan,
		'a': cterm.ColorLightGray,
	}
)

func MultiColorFrameFromString(s string, fg string, bg ...string) Frame {
	if len(bg) > 1 {
		panic("background color frame must be 1")
	}
	f := FrameFromString(s, cwin.TransparentChx().Attr)
	fgF := FrameFromString(fg, cwin.TransparentChx().Attr)
	if len(f) != len(fgF) {
		panic(fmt.Sprintf(
			"frame and its color frame not matching: len(f)=%d, len(fgF)=%d", len(f), len(fgF)))
	}
	for i := 0; i < len(f); i++ {
		if color, ok := colorMap[fgF[i].Chx.Ch]; ok {
			f[i].Chx.Attr.Fg = color
		}
	}
	return f
}

func MultiColorFrameFromFile(filepath string, n int) (Frame, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	br := bufio.NewReader(f)
	var lines []string
	for {
		line, err := ios.ReadLine(br)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		lines = append(lines, line)
	}
	if len(lines) < 2*n {
		return nil, errors.New("not enough lines in file")
	}
	s := strings.Join(lines[:n], "\n")
	fg := strings.Join(lines[n:2*n], "\n")
	var bg []string
	if len(lines) >= 3*n {
		bg = append(bg, strings.Join(lines[2*n:3*n], "\n"))
	}
	return MultiColorFrameFromString(s, fg, bg...), nil
}

func FrameRect(f Frame) cwin.Rect {
	maxX, maxY := math.MinInt32, math.MinInt32
	for i := 0; i < len(f); i++ {
		maxX = maths.MaxInt(maxX, f[i].X)
		maxY = maths.MaxInt(maxY, f[i].Y)
	}
	return cwin.Rect{X: 0, Y: 0, W: maxX + 1, H: maxY + 1}
}

func FrameFromWin(w cwin.Win) Frame {
	var f Frame
	for y := 0; y < w.ClientRect().H; y++ {
		for x := 0; x < w.ClientRect().W; x++ {
			chx := w.GetClient(x, y)
			if chx != cwin.TransparentChx() {
				f = append(f, Cell{X: x, Y: y, Chx: chx})
			}
		}
	}
	return f
}

func FrameToWin(f Frame, w cwin.Win) {
	w.FillClient(w.ClientRect().ToOrigin(), cwin.TransparentChx())
	for i := 0; i < len(f); i++ {
		w.PutClient(f[i].X, f[i].Y, f[i].Chx)
	}
}

type Frames []Frame

func FramesFromString(ss []string, attr cwin.Attr) (Frames, cwin.Rect) {
	// Unlike a single frame func FrameFromString where the frame hosting rect
	// can be implied from the frame content, multiple frames can have different
	// sizes with which we need to do some normalization:
	// - compute the bounding rect large enough for all frames.
	// - use that rect as a container, to "put" all frames in the center of it -
	//   which means adjusting each frame cell's coordinates.
	var fs Frames
	var maxR cwin.Rect
	for _, s := range ss {
		f := FrameFromString(s, attr)
		fsR := FrameRect(f)
		maxR.W = maths.MaxInt(maxR.W, fsR.W)
		maxR.H = maths.MaxInt(maxR.H, fsR.H)
		fs = append(fs, f)
	}
	for i := 0; i < len(fs); i++ {
		fsR := FrameRect(fs[i])
		dx := (maxR.W - fsR.W) / 2
		dy := (maxR.H - fsR.H) / 2
		for j := 0; j < len(fs[i]); j++ {
			fs[i][j].X += dx
			fs[i][j].Y += dy
		}
	}
	return fs, maxR
}

type FrameProvider interface {
	Next() (Frame, time.Duration, bool)
}

type simpleFrameProvider struct {
	frames Frames
	t      time.Duration
	loop   bool
	idx    int
}

func (sfp *simpleFrameProvider) Next() (Frame, time.Duration, bool) {
	if sfp.idx >= len(sfp.frames) {
		if !sfp.loop {
			return nil, 0, false
		}
		sfp.idx = 0
	}
	f := sfp.frames[sfp.idx]
	sfp.idx++
	return f, sfp.t, true
}

func NewSimpleFrameProvider(frames Frames, t time.Duration, loop bool) *simpleFrameProvider {
	return &simpleFrameProvider{
		frames: frames,
		t:      t,
		loop:   loop,
	}
}
