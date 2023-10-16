package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
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
	op.GeoM.Translate(128-12, 128-12)
	op.GeoM.Scale(2, 2)

	screen.DrawImage(getFrame(DIR_SOUTH).(*ebiten.Image), &op)

}

var walkframe int
var frameCount int

const spriteSize = 24

func getFrame(dir int) image.Image {

	rect := image.Rectangle{}
	rect.Min.X = (walkframe * spriteSize)
	rect.Max.X = (walkframe * spriteSize) + spriteSize
	rect.Min.Y = 0
	rect.Max.Y = spriteSize

	if frameCount%2 == 0 {
		walkframe++
		if walkframe > 11 {
			walkframe = 0
		}
	}

	switch dir {
	case DIR_NORTH:
		return walkNorth.SubImage(rect)
	case DIR_EAST:
		return walkEast.SubImage(rect)
	case DIR_SOUTH:
		return walkSouth.SubImage(rect)
	case DIR_WEST:
		return walkWest.SubImage(rect)
	default:
		return nil
	}
}

/* Ebiten resize handling */
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(outsideWidth), int(outsideHeight)
}

/* Input interface handler */
func (g *Game) Update() error {

	return nil
}
