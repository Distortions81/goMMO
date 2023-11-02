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
	/* UI states */
	gMouseHeld        bool
	gRightMouseHeld   bool
	gClickCaptured    bool
	gWindowDrag       *windowData
	LeftMousePressed  bool
	RightMousePressed bool
	MouseX            int
	MouseY            int
	lastMouseX        int
	lastMouseY        int

	//World edit states
	EditMode bool
	EditID   uint32
	editPos  XY = xyCenter

	//Chat command states
	ChatMode    bool
	CommandMode bool
	ChatText    string

	//Net write throttle
	lastNetSend        time.Time
	directionKeepAlive = time.Millisecond * 250
	directionThrottle  = time.Millisecond * 10

	//Direction player will be going this tick
	newPlayerDirection DIR
)

const (
	//Max chat length
	maxChatLen = 256
)

// Ebiten input handler
func (g *Game) Update() error {

	// Ignore if game not focused
	if !ebiten.IsFocused() {
		return nil
	}

	newPlayerDirection = DIR_NONE

	//Don't update during draw
	drawLock.Lock()
	defer drawLock.Unlock()

	//Get mouse / touch
	getCursor()

	//Clamp cursor and clicks to screen
	clampCursor()

	//In-game UI
	handleUI()

	//Chat and command system
	chatCommands()

	//World-edit mode
	worldEditor()

	//handle settings hotkeys
	settingsHotkeys()

	//Handle WASD and arrow keys
	WASDKeys()

	//Send current player direction
	sendMove(newPlayerDirection)

	return nil
}

func WASDKeys() {
	pressedKeys := inpututil.AppendPressedKeys(nil)

	for _, key := range pressedKeys {
		if !ChatMode {
			if key == ebiten.KeyW ||
				key == ebiten.KeyArrowUp {
				if newPlayerDirection == DIR_NONE {
					newPlayerDirection = DIR_N
				} else if newPlayerDirection == DIR_E {
					newPlayerDirection = DIR_NE
				} else if newPlayerDirection == DIR_W {
					newPlayerDirection = DIR_NW
				}
			}
			if key == ebiten.KeyA ||
				key == ebiten.KeyArrowLeft {
				if newPlayerDirection == DIR_NONE {
					newPlayerDirection = DIR_W
				} else if newPlayerDirection == DIR_N {
					newPlayerDirection = DIR_NW
				} else if newPlayerDirection == DIR_S {
					newPlayerDirection = DIR_SW
				}
			}
			if key == ebiten.KeyS ||
				key == ebiten.KeyArrowDown {
				if newPlayerDirection == DIR_NONE {
					newPlayerDirection = DIR_S
				} else if newPlayerDirection == DIR_E {
					newPlayerDirection = DIR_SE
				} else if newPlayerDirection == DIR_W {
					newPlayerDirection = DIR_SW
				}
			}
			if key == ebiten.KeyD ||
				key == ebiten.KeyArrowRight {
				if newPlayerDirection == DIR_NONE {
					newPlayerDirection = DIR_E
				} else if newPlayerDirection == DIR_N {
					newPlayerDirection = DIR_NE
				} else if newPlayerDirection == DIR_S {
					newPlayerDirection = DIR_SE
				}
			}
		}
	}
}

func settingsHotkeys() {
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
				chat("Fast shadows now Off!")
			} else {
				fastShadow = true
				chat("Fast shadows now ON! (Faster/Less battery)")
			}
		}
	}
}

/* Record mouse clicks, send clicks to toolbar */
func getMouseClicks() {
	defer reportPanic("getMouseClicks")

	/* Mouse clicks */
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		gMouseHeld = false

		/* Stop dragging window */
		gWindowDrag = nil

	} else if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		gMouseHeld = true
	}

}

func getCursor() {
	// Save mouse coords
	lastMouseX = MouseX
	lastMouseY = MouseY
	gClickCaptured = false

	//Handle mouse/touch events
	touchIDs := ebiten.AppendTouchIDs(nil)
	if len(touchIDs) > 0 {
		touchEnabled = true

		MouseX, MouseY = ebiten.TouchPosition(touchIDs[0])
		if !lastTouch {
			gMouseHeld = true
			lastTouch = true
		}
	} else {
		if touchEnabled {
			lastTouch = false
			gMouseHeld = false
			gWindowDrag = nil
		}
		getMouseClicks()
	}
}

var touchEnabled bool
var lastTouch bool

func clampCursor() {
	// Clamp mouse/touch to window
	MouseX, MouseY = ebiten.CursorPosition()
	if MouseX < 0 || MouseX > int(screenWidth) ||
		MouseY < 0 || MouseY > int(screenHeight) {
		MouseX = lastMouseX
		MouseY = lastMouseY

		// Stop dragging window if we go off-screen
		gWindowDrag = nil

		//Eat clicks
		gClickCaptured = true
		gMouseHeld = false
	}
}

func handleUI() {
	// Check if we clicked within a window
	if gMouseHeld {
		gClickCaptured = handleToolbar()
		gClickCaptured = collisionWindowsCheck(XYs{X: int32(MouseX), Y: int32(MouseY)})
	}

	/* Handle window drag */
	if gWindowDrag != nil {
		gWindowDrag.position = XYs{X: int32(MouseX) - gWindowDrag.dragPos.X, Y: int32(MouseY) - gWindowDrag.dragPos.Y}
		gClickCaptured = true
	}
}

func chatCommands() {
	//Chat and command handler
	if ChatMode || CommandMode {
		start := []rune{}
		runes := ebiten.AppendInputChars(start[:0])
		if len(ChatText) < maxChatLen {
			ChatText += string(runes)
		} else {
			chat("Sorry, that is the max message length!")
			return
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
		return
	} else if repeatingKeyPressed(ebiten.KeyEnter) && !CommandMode {
		ChatMode = true
		ChatText = ""
	} else if repeatingKeyPressed(ebiten.KeyGraveAccent) && !ChatMode {
		CommandMode = true
		ChatText = ""
	}
}

func worldEditor() {
	if repeatingKeyPressed(ebiten.KeyBackslash) {
		if EditMode {
			EditMode = false
		} else {
			EditMode = true
			chat("Click to place an item, right-click to delete an item, + and - cycle item IDs.")
		}
	}

	if EditMode {
		if !gClickCaptured {
			if gMouseHeld {
				if !LeftMousePressed {
					editPlaceItem()
				}
				LeftMousePressed = true
			} else {
				LeftMousePressed = false
			}
			if gRightMouseHeld {
				if !RightMousePressed {
					editDeleteItem()
				}
				RightMousePressed = true
			} else {
				RightMousePressed = false
			}
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

		editPos = XY{X: smoothCamPos.X - uint32(MouseX), Y: smoothCamPos.Y - uint32(MouseY)}
	} else {
		if !gClickCaptured {
			if gMouseHeld {
				newPlayerDirection = walkXY(MouseX, MouseY)
			}
		}
	}
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
		} else if time.Since(lastNetSend) < directionKeepAlive {
			return
		}
	}

	if time.Since(lastNetSend) < directionThrottle {
		return
	}

	//Update our direction
	goDir = newDir

	var buf []byte
	outbuf := bytes.NewBuffer(buf)

	binary.Write(outbuf, binary.LittleEndian, &goDir)
	sendCommand(CMD_MOVE, outbuf.Bytes())

	lastNetSend = time.Now()
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
