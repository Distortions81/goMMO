package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy"
)

var (
	LeftMousePressed  bool
	RightMousePressed bool

	EditMode bool
	EditID   uint32
	editPos  XY = xyCenter

	ChatMode    bool
	CommandMode bool
	ChatText    string

	lastSent time.Time

	sendInterval = time.Millisecond * 250
)

const (
	maxChat     = 256
	maxSendRate = time.Millisecond * 10
)

/* Input interface handler */
func (g *Game) Update() error {

	drawLock.Lock()
	defer drawLock.Unlock()

	newDir := DIR_NONE

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
			chat("Click to place an item, right-click to delete an item, + and - cycle item IDs.")
		}
	}
	if repeatingKeyPressed(ebiten.KeyN) {
		if !ChatMode && !CommandMode {

			if nightLevel >= 250 {
				nightLevel = 0
			} else if nightLevel+42 >= 250 {
				nightLevel = 255
			} else {
				nightLevel += 42
			}

			buf := fmt.Sprintf("Night level: %v%%", int((float32(nightLevel)/255.0)*100.0))
			chat(buf)
		}
	}
	if repeatingKeyPressed(ebiten.KeyL) {
		if !ChatMode && !CommandMode {

			if noSmoothing {
				noSmoothing = false
				chat("Motion smoothing now ON!")
			} else {
				noSmoothing = true
				chat("Motion smoothing now OFF! (battery saver)")
			}
		}
	}
	if repeatingKeyPressed(ebiten.KeyZ) {
		if !ChatMode && !CommandMode {

			if fastShadow {
				fastShadow = false
				chat("Fast shadows now ON! (Faster/Less battery)")
			} else {
				fastShadow = true
				chat("Fast shadows now OFF!")
			}
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
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
			if !RightMousePressed {
				editDeleteItem()
			}
			RightMousePressed = true
		} else {
			RightMousePressed = false
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
		editPos = XY{X: smoothCamPos.X - uint32(mx), Y: smoothCamPos.Y - uint32(my)}
	} else {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			mx, my := ebiten.CursorPosition()
			newDir = walkXY(mx, my)
		} else {
			touchIDs := ebiten.AppendTouchIDs(nil)
			//Ignore multi-touch
			if len(touchIDs) == 1 {
				for _, touch := range touchIDs {
					mx, my := ebiten.TouchPosition(touch)
					newDir = walkXY(mx, my)
					break
				}
			}
		}
	}

	pressedKeys := inpututil.AppendPressedKeys(nil)

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

	sendMove(newDir)

	return nil
}

func walkXY(mx, my int) DIR {

	distance := distance(XY{X: uint32(HscreenWidth), Y: uint32(HscreenHeight)}, XY{X: uint32(mx), Y: uint32(my)})

	if distance < charSpriteSize ||
		mx > screenWidth || my > screenHeight ||
		mx < 0 || my < 0 {
		return DIR_NONE
	}

	p1 := geom.Coord{float64(HscreenWidth), float64(HscreenHeight), 0}
	p2 := geom.Coord{float64(mx), float64(my), 0}

	angle := xy.Angle(p1, p2)

	return radToDir(angle)
}

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

func sendMove(newDir DIR) {

	//Exit if nothing changed
	if newDir == goDir {
		if goDir == DIR_NONE {
			return
		} else if time.Since(lastSent) < sendInterval {
			return
		}
	}

	if time.Since(lastSent) < maxSendRate {
		return
	}

	//Update our direction
	goDir = newDir

	var buf []byte
	outbuf := bytes.NewBuffer(buf)

	binary.Write(outbuf, binary.LittleEndian, &goDir)
	sendCommand(CMD_MOVE, outbuf.Bytes())

	lastSent = time.Now()
}

func editPlaceItem() {

	var buf []byte
	outbuf := bytes.NewBuffer(buf)

	binary.Write(outbuf, binary.LittleEndian, EditID)
	binary.Write(outbuf, binary.LittleEndian, editPos.X)
	binary.Write(outbuf, binary.LittleEndian, editPos.Y)
	sendCommand(CMD_EDITPLACEITEM, outbuf.Bytes())
}

func editDeleteItem() {

	var buf []byte
	outbuf := bytes.NewBuffer(buf)

	binary.Write(outbuf, binary.LittleEndian, EditID)
	binary.Write(outbuf, binary.LittleEndian, editPos.X)
	binary.Write(outbuf, binary.LittleEndian, editPos.Y)
	sendCommand(CMD_EDITDELETEITEM, outbuf.Bytes())
}