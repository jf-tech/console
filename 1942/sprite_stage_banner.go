package main

import (
	"strings"
	"time"

	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cwin"
	"github.com/nsf/termbox-go"
)

var (
	stageBannerName = "stage_banner"
	// https://textkool.com/en/ascii-art-generator?hl=default&vl=default&font=Colossal&text=Stage%20
	stageImgTxt = strings.Trim(`
 .d8888b.  888
d88P  Y88b 888
Y88b.      888
 "Y888b.   888888  8888b.   .d88b.   .d88b.
    "Y88b. 888        "88b d88P"88b d8P  Y8b
      "888 888    .d888888 888  888 88888888
Y88b  d88P Y88b.  888  888 Y88b 888 Y8b.
 "Y8888P"   "Y888 "Y888888  "Y88888  "Y8888
                                888
                           Y8b d88P
                            "Y88P"`, "\n")

	stageNumImgTxt = []string{
		strings.Trim(`
 d888
d8888
  888
  888
  888
  888
  888
8888888`, "\n"),
		strings.Trim(`
 .d8888b.
d88P  Y88b
       888
     .d88P
 .od888P"
d88P"
888"
888888888`, "\n"),
	}
	stageBannerAttr           = cwin.ChAttr{Fg: termbox.ColorYellow, Bg: termbox.ColorBlue}
	stageBannerMoveInOutSpeed = cgame.ActionPerSec(200)
	stageBannerStayDuration   = 1 * time.Second
)

func newSpriteStageBanner(g *cgame.Game, parent *cwin.Win, stage int) *cgame.SpriteAnimated {
	stageR := cwin.TextDimension(stageImgTxt)
	stageLines := strings.Split(stageImgTxt, "\n")
	stageNumLines := strings.Split(stageNumImgTxt[stage], "\n")
	imgTxts := make([]string, len(stageLines))
	for i := 0; i < len(stageLines); i++ {
		imgTxts[i] = stageLines[i]
		if i < len(stageNumLines) {
			imgTxts[i] += strings.Repeat(" ", stageR.W-len(stageLines[i])+5) + stageNumLines[i]
		}
		imgTxts[i] += "\n"
	}
	imgTxt := strings.Trim(strings.Join(imgTxts, ""), "\n")
	bannerR := cwin.TextDimension(imgTxt)
	startX := -bannerR.W
	midX := (parent.ClientRect().W - bannerR.W) / 2
	y := (parent.ClientRect().H - bannerR.H) / 2
	g.Pause()
	return cgame.NewSpriteAnimated(g, parent,
		cgame.SpriteAnimatedCfg{
			Name:      stageBannerName,
			Frames:    [][]cgame.Cell{cgame.StringToCells(imgTxt, stageBannerAttr)},
			DX:        1,
			MoveSpeed: stageBannerMoveInOutSpeed,
			AfterMove: func(s cgame.Sprite) {
				if s.Win().Rect().X >= midX {
					g.SpriteMgr.AddEvent(cgame.NewSpriteEventDelete(s))
					g.SpriteMgr.AddEvent(cgame.NewSpriteEventCreate(
						cgame.NewSpriteAnimated(g, parent,
							cgame.SpriteAnimatedCfg{
								Name:      stageBannerName,
								Frames:    [][]cgame.Cell{cgame.StringToCells(imgTxt, stageBannerAttr)},
								MoveSpeed: cgame.ActionPerSec(1) / (cgame.ActionPerSec(stageBannerStayDuration) / cgame.ActionPerSec(time.Second)),
								AfterMove: func(s cgame.Sprite) {
									g.SpriteMgr.AddEvent(cgame.NewSpriteEventDelete(s))
									g.SpriteMgr.AddEvent(cgame.NewSpriteEventCreate(
										cgame.NewSpriteAnimated(g, parent,
											cgame.SpriteAnimatedCfg{
												Name:      stageBannerName,
												Frames:    [][]cgame.Cell{cgame.StringToCells(imgTxt, stageBannerAttr)},
												DX:        1,
												MoveSpeed: stageBannerMoveInOutSpeed,
												AfterMove: func(s cgame.Sprite) {
													if s.Win().Rect().X >= parent.ClientRect().W-1 {
														s.Clock().Pause()
														g.Resume()
													}
												},
											},
											s.Win().Rect().X, y)))
								},
							},
							midX, y)))
				}
			},
		},
		startX, y)
}
