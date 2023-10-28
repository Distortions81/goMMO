package main

import (
	"bytes"
	"encoding/binary"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	LeftMousePressed bool
	EditMode         bool
	EditID           uint32
	editPos          XY = xyCenter

	ChatMode    bool
	CommandMode bool
	ChatText    string
)

const maxChat = 256

/* Input interface handler */
func (g *Game) Update() error {
	updateCount++

	if ChatMode || CommandMode {
		start := []rune{}
		runes := ebiten.AppendInputChars(start[:0])
		if len(ChatText) < maxChat {
			ChatText += string(runes)
		} else {
			chat("Sorry, that is the max message length!")
			return nil
		}

		if repeatingKeyPressed(ebiten.KeyEscape) {
			ChatMode = false
			CommandMode = false
			ChatText = ""
		}
		if repeatingKeyPressed(ebiten.KeyEnter) {

			if ChatText != "" {
				if CommandMode {
					sendCommand(CMD_COMMAND, []byte(ChatText))
				} else if ChatMode {
					sendCommand(CMD_CHAT, []byte(ChatText))
				}

			}
			ChatMode = false
			CommandMode = false
			ChatText = ""
		} else if repeatingKeyPressed(ebiten.KeyBackspace) {
			if len(ChatText) >= 1 {
				ChatText = ChatText[:len(ChatText)-1]
			}

		}
		return nil
	} else if repeatingKeyPressed(ebiten.KeyEnter) && !CommandMode {
		ChatMode = true
		ChatText = ""
	} else if repeatingKeyPressed(ebiten.KeyGraveAccent) && !ChatMode {
		CommandMode = true
		ChatText = ""
	}

	if repeatingKeyPressed(ebiten.KeyBackslash) {
		if EditMode {
			EditMode = false
		} else {
			EditMode = true
		}
	}
	if EditMode {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			if !LeftMousePressed {
				editPlaceItem()
			}
			LeftMousePressed = true
		} else {
			LeftMousePressed = false
		}
		if repeatingKeyPressed(ebiten.KeyEqual) {
			if EditID < numSprites {
				EditID++
			}
		} else if repeatingKeyPressed(ebiten.KeyMinus) {
			if EditID > 0 {
				EditID--
			}
		}
		mx, my := ebiten.CursorPosition()
		editPos = XY{X: camPos.X - uint32(mx), Y: camPos.Y - uint32(my)}
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

		if updateCount%8 == 0 {
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
		curCharPos.Y++
	case DIR_NE:
		curCharPos.Y += diagSpeed
		curCharPos.X -= diagSpeed
	case DIR_E:
		curCharPos.X--
	case DIR_SE:
		curCharPos.X -= diagSpeed
		curCharPos.Y -= diagSpeed
	case DIR_S:
		curCharPos.Y--
	case DIR_SW:
		curCharPos.Y -= diagSpeed
		curCharPos.X += diagSpeed
	case DIR_W:
		curCharPos.X++
	case DIR_NW:
		curCharPos.Y += diagSpeed
		curCharPos.X += diagSpeed
	default:
		return
	}

}

func sendMove() {

	var buf []byte
	outbuf := bytes.NewBuffer(buf)

	binary.Write(outbuf, binary.LittleEndian, int8(curCharPos.X-lastCharPos.X))
	binary.Write(outbuf, binary.LittleEndian, int8(curCharPos.Y-lastCharPos.Y))
	sendCommand(CMD_MOVE, outbuf.Bytes())

	lastCharPos = curCharPos
}

func editPlaceItem() {

	var buf []byte
	outbuf := bytes.NewBuffer(buf)

	binary.Write(outbuf, binary.LittleEndian, EditID)
	binary.Write(outbuf, binary.LittleEndian, editPos.X)
	binary.Write(outbuf, binary.LittleEndian, editPos.Y)
	sendCommand(CMD_EDITPLACEITEM, outbuf.Bytes())

	outbuf.Reset()
	binary.Write(outbuf, binary.LittleEndian, editPos.X/chunkDiv)
	binary.Write(outbuf, binary.LittleEndian, editPos.Y/chunkDiv)
	sendCommand(CMD_GETCHUNK, outbuf.Bytes())
}
