package main

import (
	"bytes"
	"encoding/binary"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	ChatMode bool
	ChatText string
)

const maxChat = 256

/* Input interface handler */
func (g *Game) Update() error {
	updateCount++

	if ChatMode {
		start := []rune{}
		runes := ebiten.AppendInputChars(start[:0])
		if len(ChatText) < maxChat {
			ChatText += string(runes)
		} else {
			chat("Sorry, that is the max message length!")
			return nil
		}

		if repeatingKeyPressed(ebiten.KeyEnter) {
			ChatMode = false
			if ChatText != "" {
				sendCommand(CMD_CHAT, []byte(ChatText))
				ChatText = ""
			}
		} else if repeatingKeyPressed(ebiten.KeyBackspace) {
			if len(ChatText) >= 1 {
				ChatText = ChatText[:len(ChatText)-1]
			}

		}
		return nil
	} else if repeatingKeyPressed(ebiten.KeyEnter) {
		ChatMode = true
		ChatText = ""
	}

	pressedKeys := inpututil.AppendPressedKeys(nil)
	if pressedKeys == nil {
		return nil
	}

	newDir := DIR_NONE
	for _, key := range pressedKeys {
		if !ChatMode {
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
	}

	if newDir != DIR_NONE {

		goDir = newDir
		moveDir(goDir)

		if updateCount%4 == 0 {
			sendMove()
		}
	} else {
		updateCount = 0
	}

	return nil
}

const diagSpeed = 0.707

// repeatingKeyPressed return true when key is pressed considering the repeat state.
func repeatingKeyPressed(key ebiten.Key) bool {
	const (
		delay    = 30
		interval = 3
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}

func moveDir(dir DIR) {

	switch dir {
	case DIR_N:
		localCharPos.Y++
	case DIR_NE:
		localCharPos.Y += diagSpeed
		localCharPos.X -= diagSpeed
	case DIR_E:
		localCharPos.X--
	case DIR_SE:
		localCharPos.X -= diagSpeed
		localCharPos.Y -= diagSpeed
	case DIR_S:
		localCharPos.Y--
	case DIR_SW:
		localCharPos.Y -= diagSpeed
		localCharPos.X += diagSpeed
	case DIR_W:
		localCharPos.X++
	case DIR_NW:
		localCharPos.Y += diagSpeed
		localCharPos.X += diagSpeed
	default:
		return
	}

}

func sendMove() {
	var buf []byte
	outbuf := bytes.NewBuffer(buf)
	pos := XY{X: uint32(localCharPos.X), Y: uint32(localCharPos.Y)}

	binary.Write(outbuf, binary.BigEndian, pos.X)
	binary.Write(outbuf, binary.BigEndian, pos.Y)
	sendCommand(CMD_MOVE, outbuf.Bytes())
}
