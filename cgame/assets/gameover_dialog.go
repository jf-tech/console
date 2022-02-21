package assets

import (
	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cterm"
	"github.com/jf-tech/console/cwin"
)

type GameResult int

const (
	GameResultQuit = GameResult(iota)
	GameResultReplay
	GameResultSystemFailure
)

var (
	// https://textkool.com/en/ascii-art-generator?hl=default&vl=default&font=Colossal&text=Game%20Over%20!
	gameOverTxt = `

 .d8888b.                                        .d88888b.                                 888
d88P  Y88b                                      d88P" "Y88b                                888
888    888                                      888     888                                888
888         8888b.  88888b.d88b.   .d88b.       888     888 888  888  .d88b.  888d888      888
888  88888     "88b 888 "888 "88b d8P  Y8b      888     888 888  888 d8P  Y8b 888P"        888
888    888 .d888888 888  888  888 88888888      888     888 Y88  88P 88888888 888          Y8P
Y88b  d88P 888  888 888  888  888 Y8b.          Y88b. .d88P  Y8bd8P  Y8b.     888           "
 "Y8888P88 "Y888888 888  888  888  "Y8888        "Y88888P"    Y88P    "Y8888  888          888


                            Press ESC or 'q' to quit, 'r' to replay.`

	// https://textkool.com/en/ascii-art-generator?hl=default&vl=default&font=Colossal&text=You%20Won%20!
	youWonTxt = `

Y88b   d88P                     888       888                        888
 Y88b d88P                      888   o   888                        888
  Y88o88P                       888  d8b  888                        888
   Y888P  .d88b.  888  888      888 d888b 888  .d88b.  88888b.       888
    888  d88""88b 888  888      888d88888b888 d88""88b 888 "88b      888
    888  888  888 888  888      88888P Y88888 888  888 888  888      Y8P
    888  Y88..88P Y88b 888      8888P   Y8888 Y88..88P 888  888       "
    888   "Y88P"   "Y88888      888P     Y888  "Y88P"  888  888      888


                  Press ESC or 'q' to quit, 'r' to replay.`

	GameOverKeys  = cwin.Keys(cterm.KeyEsc, 'q')
	PauseGameKeys = cwin.Keys('p')
	ReplayKeys    = cwin.Keys('r')
)

func DisplayGameOverDialog(g *cgame.Game) GameResult {
	ev := g.WinSys.MessageBoxEx(nil, append(GameOverKeys, ReplayKeys...), "", gameOverTxt)
	if ev.Ch == 'r' {
		return GameResultReplay
	}
	return GameResultQuit
}

func DisplayYouWonDialog(g *cgame.Game) GameResult {
	ev := g.WinSys.MessageBoxEx(nil, append(GameOverKeys, ReplayKeys...), "Fantastic!", youWonTxt)
	if ev.Ch == 'r' {
		return GameResultReplay
	}
	return GameResultQuit
}
