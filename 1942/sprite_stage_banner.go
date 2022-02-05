package main

import (
	"strings"
	"time"

	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cwin"
	"github.com/nsf/termbox-go"
)

var (
	stageIntroBannerName = "stage_intro_banner"
	// https://textkool.com/en/ascii-art-generator?hl=default&vl=default&font=Electronic&text=Stage%201
	stageIntroBannerFrameTxt = strings.Trim(`
 ▄▄▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄▄▄
▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌
▐░█▀▀▀▀▀▀▀▀▀  ▀▀▀▀█░█▀▀▀▀ ▐░█▀▀▀▀▀▀▀█░▌▐░█▀▀▀▀▀▀▀▀▀ ▐░█▀▀▀▀▀▀▀▀▀
▐░▌               ▐░▌     ▐░▌       ▐░▌▐░▌          ▐░▌
▐░█▄▄▄▄▄▄▄▄▄      ▐░▌     ▐░█▄▄▄▄▄▄▄█░▌▐░▌ ▄▄▄▄▄▄▄▄ ▐░█▄▄▄▄▄▄▄▄▄
▐░░░░░░░░░░░▌     ▐░▌     ▐░░░░░░░░░░░▌▐░▌▐░░░░░░░░▌▐░░░░░░░░░░░▌
 ▀▀▀▀▀▀▀▀▀█░▌     ▐░▌     ▐░█▀▀▀▀▀▀▀█░▌▐░▌ ▀▀▀▀▀▀█░▌▐░█▀▀▀▀▀▀▀▀▀
          ▐░▌     ▐░▌     ▐░▌       ▐░▌▐░▌       ▐░▌▐░▌
 ▄▄▄▄▄▄▄▄▄█░▌     ▐░▌     ▐░▌       ▐░▌▐░█▄▄▄▄▄▄▄█░▌▐░█▄▄▄▄▄▄▄▄▄
▐░░░░░░░░░░░▌     ▐░▌     ▐░▌       ▐░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌
 ▀▀▀▀▀▀▀▀▀▀▀       ▀       ▀         ▀  ▀▀▀▀▀▀▀▀▀▀▀  ▀▀▀▀▀▀▀▀▀▀▀
`, "\n")
	stageIntroBannerNumFrameTxt = []string{
		strings.Trim(`
         ▄▄▄▄
       ▄█░░░░▌
      ▐░░▌▐░░▌
       ▀▀ ▐░░▌
          ▐░░▌
          ▐░░▌
          ▐░░▌
          ▐░░▌
      ▄▄▄▄█░░█▄▄▄
     ▐░░░░░░░░░░░▌
      ▀▀▀▀▀▀▀▀▀▀▀
`, "\n"),
		strings.Trim(`
      ▄▄▄▄▄▄▄▄▄▄▄
     ▐░░░░░░░░░░░▌
      ▀▀▀▀▀▀▀▀▀█░▌
               ▐░▌
               ▐░▌
      ▄▄▄▄▄▄▄▄▄█░▌
     ▐░░░░░░░░░░░▌
     ▐░█▀▀▀▀▀▀▀▀▀
     ▐░█▄▄▄▄▄▄▄▄▄
     ▐░░░░░░░░░░░▌
      ▀▀▀▀▀▀▀▀▀▀▀
`, "\n"),
	}
	stageIntroBannerAttr          = cwin.ChAttr{Fg: termbox.ColorLightYellow, Bg: termbox.ColorBlue}
	stageIntroBannerInOutDuration = 300 * time.Millisecond
	stageIntroBannerStayDuration  = 1 * time.Second
)

func createStageIntroBanner(m *myGame, stageIdx int, afterFinish func()) {
	frame := cgame.FrameFromStringEx(stageIntroBannerFrameTxt, stageIntroBannerAttr, false)
	frameR := cgame.FrameRect(frame)
	frameNumeric := cgame.FrameFromStringEx(
		stageIntroBannerNumFrameTxt[stageIdx], stageIntroBannerAttr, false)
	for i := 0; i < len(frameNumeric); i++ {
		frameNumeric[i].X += frameR.W
	}
	frame = append(frame, frameNumeric...)
	createSlideInOutBanner(
		m, frame, stageIntroBannerInOutDuration, stageIntroBannerStayDuration, afterFinish)
}

var (
	stagePassedBannerName = "stage_passed_banner"
	// https://textkool.com/en/ascii-art-generator?hl=default&vl=default&font=Electronic&text=Passed
	stagePassedBannerFrame = cgame.FrameFromStringEx(strings.Trim(`
 ▄▄▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄▄
▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░▌
▐░█▀▀▀▀▀▀▀█░▌▐░█▀▀▀▀▀▀▀█░▌▐░█▀▀▀▀▀▀▀▀▀ ▐░█▀▀▀▀▀▀▀▀▀ ▐░█▀▀▀▀▀▀▀▀▀ ▐░█▀▀▀▀▀▀▀█░▌
▐░▌       ▐░▌▐░▌       ▐░▌▐░▌          ▐░▌          ▐░▌          ▐░▌       ▐░▌
▐░█▄▄▄▄▄▄▄█░▌▐░█▄▄▄▄▄▄▄█░▌▐░█▄▄▄▄▄▄▄▄▄ ▐░█▄▄▄▄▄▄▄▄▄ ▐░█▄▄▄▄▄▄▄▄▄ ▐░▌       ▐░▌
▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░▌       ▐░▌
▐░█▀▀▀▀▀▀▀▀▀ ▐░█▀▀▀▀▀▀▀█░▌ ▀▀▀▀▀▀▀▀▀█░▌ ▀▀▀▀▀▀▀▀▀█░▌▐░█▀▀▀▀▀▀▀▀▀ ▐░▌       ▐░▌
▐░▌          ▐░▌       ▐░▌          ▐░▌          ▐░▌▐░▌          ▐░▌       ▐░▌
▐░▌          ▐░▌       ▐░▌ ▄▄▄▄▄▄▄▄▄█░▌ ▄▄▄▄▄▄▄▄▄█░▌▐░█▄▄▄▄▄▄▄▄▄ ▐░█▄▄▄▄▄▄▄█░▌
▐░▌          ▐░▌       ▐░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░▌
 ▀            ▀         ▀  ▀▀▀▀▀▀▀▀▀▀▀  ▀▀▀▀▀▀▀▀▀▀▀  ▀▀▀▀▀▀▀▀▀▀▀  ▀▀▀▀▀▀▀▀▀▀
`, "\n"), stageIntroBannerAttr, false)
	stagePassedBannerInOutDuration = stageIntroBannerInOutDuration
	stagePassedBannerStayDuration  = 2 * time.Second
)

func createStagePassedBanner(m *myGame, afterFinish func()) {
	createSlideInOutBanner(
		m, stagePassedBannerFrame,
		stagePassedBannerInOutDuration, stagePassedBannerStayDuration,
		afterFinish)
}
