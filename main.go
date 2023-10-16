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

var walkframe int
var frameCount int
var goDir int
var isWalking bool
var charPos xy

const spriteSize = 24

func getFrame(dir int) image.Image {

	rect := image.Rectangle{}
	rect.Min.X = (walkframe * spriteSize)
	rect.Max.X = (walkframe * spriteSize) + spriteSize
	rect.Min.Y = 0
	rect.Max.Y = spriteSize

	if isWalking {
		if frameCount%2 == 0 {
			walkframe++
			if walkframe > 11 {
				walkframe = 0
			}
			MoveDir(dir)
		}
	} else {
		walkframe = 0
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

func MoveDir(dir int) {

	switch dir {
	case DIR_NORTH:
		charPos.Y--
	case DIR_EAST:
		charPos.X++
	case DIR_SOUTH:
		charPos.Y++
	case DIR_WEST:
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
		goDir = DIR_NORTH
		isWalking = true
	} else if pressedKeys[0] == ebiten.KeyA ||
		pressedKeys[0] == ebiten.KeyArrowLeft {
		goDir = DIR_WEST
		isWalking = true
	} else if pressedKeys[0] == ebiten.KeyS ||
		pressedKeys[0] == ebiten.KeyArrowDown {
		goDir = DIR_SOUTH
		isWalking = true
	} else if pressedKeys[0] == ebiten.KeyD ||
		pressedKeys[0] == ebiten.KeyArrowRight {
		goDir = DIR_EAST
		isWalking = true
	} else {
		isWalking = false
	}

	return nil
}
