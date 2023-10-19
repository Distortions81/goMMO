package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

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

	sendCommand(CMD_MOVE, xyToByteArray(charPos))
}
