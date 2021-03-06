package main

import (
	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cterm"
	"github.com/jf-tech/console/cwin"
)

var (
	stageIntroBannerName = "stage_intro_banner"
	// https://textkool.com/en/ascii-art-generator?hl=default&vl=default&font=Electronic&text=Stage%201
	stageIntroBannerFrameTxt = `
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
`
	stageIntroBannerNumFrameTxt = []string{
		`
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
`,
		`
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
`,
		`
      ▄▄▄▄▄▄▄▄▄▄▄
     ▐░░░░░░░░░░░▌
      ▀▀▀▀▀▀▀▀▀█░▌
               ▐░▌
      ▄▄▄▄▄▄▄▄▄█░▌
     ▐░░░░░░░░░░░▌
      ▀▀▀▀▀▀▀▀▀█░▌
               ▐░▌
      ▄▄▄▄▄▄▄▄▄█░▌
     ▐░░░░░░░░░░░▌
      ▀▀▀▀▀▀▀▀▀▀▀
`,
	}
	stageIntroBannerFinalStageImgTxt = `
 ▄▄▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄▄▄  ▄▄        ▄  ▄▄▄▄▄▄▄▄▄▄▄  ▄
▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░░▌      ▐░▌▐░░░░░░░░░░░▌▐░▌
▐░█▀▀▀▀▀▀▀▀▀  ▀▀▀▀█░█▀▀▀▀ ▐░▌░▌     ▐░▌▐░█▀▀▀▀▀▀▀█░▌▐░▌
▐░▌               ▐░▌     ▐░▌▐░▌    ▐░▌▐░▌       ▐░▌▐░▌
▐░█▄▄▄▄▄▄▄▄▄      ▐░▌     ▐░▌ ▐░▌   ▐░▌▐░█▄▄▄▄▄▄▄█░▌▐░▌
▐░░░░░░░░░░░▌     ▐░▌     ▐░▌  ▐░▌  ▐░▌▐░░░░░░░░░░░▌▐░▌
▐░█▀▀▀▀▀▀▀▀▀      ▐░▌     ▐░▌   ▐░▌ ▐░▌▐░█▀▀▀▀▀▀▀█░▌▐░▌
▐░▌               ▐░▌     ▐░▌    ▐░▌▐░▌▐░▌       ▐░▌▐░▌
▐░▌           ▄▄▄▄█░█▄▄▄▄ ▐░▌     ▐░▐░▌▐░▌       ▐░▌▐░█▄▄▄▄▄▄▄▄▄
▐░▌          ▐░░░░░░░░░░░▌▐░▌      ▐░░▌▐░▌       ▐░▌▐░░░░░░░░░░░▌
 ▀            ▀▀▀▀▀▀▀▀▀▀▀  ▀        ▀▀  ▀         ▀  ▀▀▀▀▀▀▀▀▀▀▀

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
`
	stageIntroBannerAttr = cwin.Attr{Fg: cterm.ColorLightYellow, Bg: cterm.ColorBlue}
)

func createStageIntroBanner(m *myGame, stageIdx int, afterFinish func()) {
	if stageIdx == totalStages-1 {
		createSlideInOutBanner(
			m,
			cgame.FrameFromString(stageIntroBannerFinalStageImgTxt, stageIntroBannerAttr),
			stageIntroBannerInOutDuration, stageIntroBannerStayDuration, afterFinish)
		return
	}
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
	stagePassedBannerFrame = cgame.FrameFromStringEx(`
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
`, stageIntroBannerAttr, false)
)

func createStagePassedBanner(m *myGame, afterFinish func()) {
	createSlideInOutBanner(
		m, stagePassedBannerFrame,
		stagePassedBannerInOutDuration, stagePassedBannerStayDuration,
		afterFinish)
}
