package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	op.GeoM.Translate(64-12, 64-12)
	op.GeoM.Scale(4, 4)

	screen.DrawImage(getFrame(goDir).(*ebiten.Image), &op)

}

var walkframe int
var frameCount int
var goDir int

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
	if inpututil.IsKeyJustPressed(ebiten.KeyW) ||
		inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		goDir = DIR_NORTH
	} else if inpututil.IsKeyJustPressed(ebiten.KeyA) ||
		inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		goDir = DIR_WEST
	} else if inpututil.IsKeyJustPressed(ebiten.KeyS) ||
		inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		goDir = DIR_SOUTH
	} else if inpututil.IsKeyJustPressed(ebiten.KeyD) ||
		inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		goDir = DIR_EAST
	}
	return nil
}
