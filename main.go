package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
}

type xy struct {
	X int
	Y int
}

func main() {
	/* Set up ebiten and window */
	ebiten.SetVsyncEnabled(true)
	ebiten.SetTPS(ebiten.SyncWithFPS)
	ebiten.SetScreenClearedEveryFrame(true)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	ebiten.SetWindowSize(512, 512)
	loadTest()

	if err := ebiten.RunGameWithOptions(newGame(), &ebiten.RunGameOptions{GraphicsLibrary: ebiten.GraphicsLibraryOpenGL}); err != nil {
		return
	}
}

func newGame() *Game {

	/* Initialize the game */
	return &Game{}
}

/* Ebiten: Draw everything */
func (g *Game) Draw(screen *ebiten.Image) {
	frameCount++

	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(64-12+float64(charPos.X), 64-12+float64(charPos.Y))
	op.GeoM.Scale(4, 4)

	screen.DrawImage(getFrame(goDir).(*ebiten.Image), &op)
}

var (
	walkframe  int
	frameCount int
	goDir      int
	isWalking  bool
	charPos    xy
)

const spriteSize = 52

func dirToOffset(dir int) int {
	return 0
}

func getFrame(dir int) image.Image {

	dirOff := dirToOffset(dir)

	rect := image.Rectangle{}
	rect.Min.X = (walkframe * spriteSize)
	rect.Max.X = (walkframe * spriteSize) + spriteSize
	rect.Min.Y = dirOff
	rect.Max.Y = spriteSize + dirOff

	if isWalking {
		if frameCount%2 == 0 {
			walkframe++
			if walkframe > 3 {
				walkframe = 0
			}
			MoveDir(dir)
		}
	} else {
		walkframe = 0
	}

	return testChar.SubImage(rect)

}

func MoveDir(dir int) {

	switch dir {
	case DIR_N:
		charPos.Y--
	case DIR_E:
		charPos.X++
	case DIR_S:
		charPos.Y++
	case DIR_W:
		charPos.X--
	}
}

/* Ebiten resize handling */
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(outsideWidth), int(outsideHeight)
}

/* Input interface handler */
func (g *Game) Update() error {
	pressedKeys := inpututil.AppendPressedKeys(nil)
	if pressedKeys == nil {
		isWalking = false
		return nil
	}

	if pressedKeys[0] == ebiten.KeyW ||
		pressedKeys[0] == ebiten.KeyArrowUp {
		goDir = DIR_N
		isWalking = true
	} else if pressedKeys[0] == ebiten.KeyA ||
		pressedKeys[0] == ebiten.KeyArrowLeft {
		goDir = DIR_W
		isWalking = true
	} else if pressedKeys[0] == ebiten.KeyS ||
		pressedKeys[0] == ebiten.KeyArrowDown {
		goDir = DIR_S
		isWalking = true
	} else if pressedKeys[0] == ebiten.KeyD ||
		pressedKeys[0] == ebiten.KeyArrowRight {
		goDir = DIR_E
		isWalking = true
	} else {
		isWalking = false
	}

	return nil
}
