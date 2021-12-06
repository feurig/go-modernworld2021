package modernworld

import (
	_ "embed"
	"fmt"

	tl "github.com/JoelOtter/termloop"
)

var (
	//go:embed files/title.txt
	titleScreenFile []byte

	//go:embed files/game_over.txt
	gameOverScreenFile []byte
)

func ShowTitleScreen(modernworld *modernworld) {
	prepareScreen(modernworld)
	showTitle(modernworld)

	if checkArenaSizeNotOk(modernworld) {
		modernworld.ScreenSizeNotOK = true
		showMaximizeScreen(modernworld)
		return
	}

	showPressToInit(modernworld, 0)
}

func checkArenaSizeNotOk(modernworld *modernworld) bool {
	w, h := modernworld.Arena.Size()

	if w < 100 || h < 37 {
		return true
	}

	return false
}

func ShowGameOverScreen(modernworld *modernworld) {
	prepareScreen(modernworld)
	showGameOver(modernworld)
	showScore(modernworld)
	showPressToInit(modernworld, 2)
}

func prepareScreen(modernworld *modernworld) {
	modernworld.Level = tl.NewBaseLevel(tl.Cell{Bg: tl.ColorBlack, Fg: tl.ColorWhite})
	modernworld.Game.Screen().SetLevel(modernworld.Level)
	modernworld.Level.AddEntity(modernworld)

	modernworld.initArena()
	modernworld.initHud()
}

func showTitle(modernworld *modernworld) {
	showCanvas(modernworld, titleScreenFile)
}

func showGameOver(modernworld *modernworld) {
	showCanvas(modernworld, gameOverScreenFile)
}

func showCanvas(modernworld *modernworld, file []byte) {
	canvas := CreateCanvas(file)

	arenaX, arenaY := modernworld.Arena.Position()
	arenaW, arenaH := modernworld.Arena.Size()

	x := arenaX + arenaW/2 - len(canvas)/2
	y := arenaY + arenaH/2 + -len(canvas[0]) - 1

	modernworld.Level.AddEntity(tl.NewEntityFromCanvas(x, y, canvas))
}

func showScore(modernworld *modernworld) {
	score := fmt.Sprintf("SCORE: %4d ", modernworld.Score)
	showCenterText(score, 0, modernworld)
}

func showPressToInit(modernworld *modernworld, topPadding int) {
	showCenterText("Press ENTER to start", topPadding, modernworld)
}

func showMaximizeScreen(modernworld *modernworld) {
	showCenterText("Maximize the console and run the game again", 0, modernworld)
}

func showCenterText(text string, topPadding int, modernworld *modernworld) {
	arenaX, arenaY := modernworld.Arena.Position()
	arenaW, arenaH := modernworld.Arena.Size()

	x := arenaX + arenaW/2 - len(text)/2
	y := arenaY + arenaH/2 + topPadding

	modernworld.Level.AddEntity(tl.NewText(x, y, text, tl.ColorWhite, tl.ColorBlack))
}
