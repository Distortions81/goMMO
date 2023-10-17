package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func main() {
	/* Set up ebiten and window */
	ebiten.SetVsyncEnabled(true)
	ebiten.SetTPS(60)
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
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(128-26+float64(charPos.X), 128-26+float64(charPos.Y))
	op.GeoM.Scale(2, 2)

	screen.DrawImage(getFrame(goDir).(*ebiten.Image), &op)
}

func dirToOffset(dir int) int {
	switch dir {
	case DIR_S:
		return charSpriteSize * 0
	case DIR_SE:
		return charSpriteSize * 1
	case DIR_E:
		return charSpriteSize * 2
	case DIR_NE:
		return charSpriteSize * 3
	case DIR_N:
		return charSpriteSize * 4
	case DIR_NW:
		return charSpriteSize * 5
	case DIR_W:
		return charSpriteSize * 6
	case DIR_SW:
		return charSpriteSize * 7
	}
	return 0
}

func getFrame(dir int) image.Image {

	dirOff := dirToOffset(dir)

	rect := image.Rectangle{}
	rect.Min.X = (walkframe * charSpriteSize)
	rect.Max.X = (walkframe * charSpriteSize) + charSpriteSize
	rect.Min.Y = dirOff
	rect.Max.Y = charSpriteSize + dirOff

	return testChar.SubImage(rect)

}

func MoveDir(dir int) {

	switch dir {
	case DIR_N:
		charPos.Y--
	case DIR_NE:
		charPos.Y--
		charPos.X++
	case DIR_E:
		charPos.X++
	case DIR_SE:
		charPos.X++
		charPos.Y++
	case DIR_S:
		charPos.Y++
	case DIR_SW:
		charPos.Y++
		charPos.X--
	case DIR_W:
		charPos.X--
	case DIR_NW:
		charPos.Y--
		charPos.X--
	}

}

/* Ebiten resize handling */
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(outsideWidth), int(outsideHeight)
}

/* Input interface handler */
func (g *Game) Update() error {
	updateCount++

	pressedKeys := inpututil.AppendPressedKeys(nil)
	if pressedKeys == nil {
		walkframe = 0
		return nil
	}

	newDir := DIR_NONE
	for _, key := range pressedKeys {
		if key == ebiten.KeyW ||
			key == ebiten.KeyArrowUp {
			if newDir == DIR_NONE {
				newDir = DIR_N
			} else if newDir == DIR_E {
				newDir = DIR_NE
			} else if newDir == DIR_W {
				newDir = DIR_NW
			}
		}
		if key == ebiten.KeyA ||
			key == ebiten.KeyArrowLeft {
			if newDir == DIR_NONE {
				newDir = DIR_W
			} else if newDir == DIR_N {
				newDir = DIR_NW
			} else if newDir == DIR_S {
				newDir = DIR_SW
			}
		}
		if key == ebiten.KeyS ||
			key == ebiten.KeyArrowDown {
			if newDir == DIR_NONE {
				newDir = DIR_S
			} else if newDir == DIR_E {
				newDir = DIR_SE
			} else if newDir == DIR_W {
				newDir = DIR_SW
			}
		}
		if key == ebiten.KeyD ||
			key == ebiten.KeyArrowRight {
			if newDir == DIR_NONE {
				newDir = DIR_E
			} else if newDir == DIR_N {
				newDir = DIR_NE
			} else if newDir == DIR_S {
				newDir = DIR_SE
			}
		}
	}
	if newDir != DIR_NONE {
		goDir = newDir
		isWalking = true
	} else {
		isWalking = false
	}

	if isWalking {
		if updateCount%6 == 0 {
			walkframe++
			if walkframe > 3 {
				walkframe = 0
			}
		}
		MoveDir(goDir)
	} else {
		walkframe = 0
	}

	return nil
}
