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
	// UI states
	mouseHeld         bool
	gRightMouseHeld   bool
	clickCaptured     bool
	draggingWindow    *windowData
	leftMousePressed  bool
	rightMousePressed bool
	mouseX            int
	mouseY            int
	lastMouseX        int
	lastMouseY        int

	//World edit states
	worldEditMode bool
	worldEditID   IID
	editPos       XY = worldCenter

	//Chat command states
	ChatMode    bool
	CommandMode bool
	ChatText    string

	//Net write throttle
	lastNetSend        time.Time
	directionKeepAlive = time.Millisecond * 250
	directionThrottle  = time.Millisecond * 10

	//Direction player will be going this tick
	newPlayerDir DIR
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

	newPlayerDir = DIR_NONE

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

	//Mouse / touch walk
	mouseTouchWalk()

	//World-edit mode
	worldEditor()

	//handle settings hotkeys
	settingsHotkeys()

	//Handle WASD and arrow keys
	WASDKeys()

	//Send current player direction
	sendMove(newPlayerDir)

	return nil
}

func WASDKeys() {
	if CommandMode || ChatMode {
		return
	}

	pressedKeys := inpututil.AppendPressedKeys(nil)

	for _, key := range pressedKeys {
		if !ChatMode {
			if key == ebiten.KeyW ||
				key == ebiten.KeyArrowUp {
				if newPlayerDir == DIR_NONE {
					newPlayerDir = DIR_N
				} else if newPlayerDir == DIR_E {
					newPlayerDir = DIR_NE
				} else if newPlayerDir == DIR_W {
					newPlayerDir = DIR_NW
				}
			}
			if key == ebiten.KeyA ||
				key == ebiten.KeyArrowLeft {
				if newPlayerDir == DIR_NONE {
					newPlayerDir = DIR_W
				} else if newPlayerDir == DIR_N {
					newPlayerDir = DIR_NW
				} else if newPlayerDir == DIR_S {
					newPlayerDir = DIR_SW
				}
			}
			if key == ebiten.KeyS ||
				key == ebiten.KeyArrowDown {
				if newPlayerDir == DIR_NONE {
					newPlayerDir = DIR_S
				} else if newPlayerDir == DIR_E {
					newPlayerDir = DIR_SE
				} else if newPlayerDir == DIR_W {
					newPlayerDir = DIR_SW
				}
			}
			if key == ebiten.KeyD ||
				key == ebiten.KeyArrowRight {
				if newPlayerDir == DIR_NONE {
					newPlayerDir = DIR_E
				} else if newPlayerDir == DIR_N {
					newPlayerDir = DIR_NE
				} else if newPlayerDir == DIR_S {
					newPlayerDir = DIR_SE
				}
			}
		}
	}
}

func settingsHotkeys() {
	if ChatMode || CommandMode || !devMode {
		return
	}

	if keyJustPressed(ebiten.KeyN) {
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

// Record mouse clicks, send clicks to toolbar
func getMouseClicks() {
	defer reportPanic("getMouseClicks")

	// Mouse clicks
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		mouseHeld = false

		// Stop dragging window
		draggingWindow = nil

	} else if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mouseHeld = true
	}

}

func getCursor() {
	// Save mouse coords
	lastMouseX = mouseX
	lastMouseY = mouseY
	clickCaptured = false

	mouseX, mouseY = ebiten.CursorPosition()

	//Handle mouse/touch events
	touchIDs := ebiten.AppendTouchIDs(nil)
	if len(touchIDs) > 0 {
		touchDetected = true

		mouseX, mouseY = ebiten.TouchPosition(touchIDs[0])
		if !hadTouchEvent {
			mouseHeld = true
			hadTouchEvent = true
		}
	} else {
		if touchDetected {
			hadTouchEvent = false
			mouseHeld = false
			draggingWindow = nil
			mouseX, mouseY = halfScreenX, halfScreenY
		} else {
			getMouseClicks()
		}
	}
}

var touchDetected bool
var hadTouchEvent bool

func clampCursor() {
	// Clamp mouse/touch to window
	if mouseX < 0 || mouseX > int(screenX) ||
		mouseY < 0 || mouseY > int(screenY) {
		mouseX = lastMouseX
		mouseY = lastMouseY

		//Eat clicks
		clickCaptured = true
		mouseHeld = false
	}
}

func handleUI() {
	// Check if we clicked within a window
	if mouseHeld {
		clickCaptured = handleToolbar()
		clickCaptured = collisionWindowsCheck(XYs{X: int32(mouseX), Y: int32(mouseY)})
	}

	// Handle window drag
	dragWindow()
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

		if keyJustPressed(ebiten.KeyEscape) {
			ChatMode = false
			CommandMode = false
			ChatText = ""
		}
		if keyJustPressed(ebiten.KeyEnter) {

			if ChatText != "" {
				if CommandMode {
					sendCommand(CMD_Command, []byte(ChatText))
				} else if ChatMode {
					sendCommand(CMD_Chat, []byte(ChatText))
				}

			}
			ChatMode = false
			CommandMode = false
			ChatText = ""
		} else if keyJustPressed(ebiten.KeyBackspace) {
			if len(ChatText) >= 1 {
				ChatText = ChatText[:len(ChatText)-1]
			}

		}
		return
	} else if keyJustPressed(ebiten.KeyEnter) && !CommandMode {
		ChatMode = true
		ChatText = ""
	} else if keyJustPressed(ebiten.KeyGraveAccent) && !ChatMode {
		CommandMode = true
		ChatText = ""
	}
}

func mouseTouchWalk() {
	if !worldEditMode && !clickCaptured {
		if mouseHeld {
			newPlayerDir = walkXY(mouseX, mouseY)
		}
	}
}

func worldEditor() {
	if CommandMode || ChatMode || !devMode {
		return
	}

	if keyJustPressed(ebiten.KeyBackslash) {
		if worldEditMode {
			worldEditMode = false
		} else {
			worldEditMode = true
			chat("Click to place an item, right-click to delete an item, + and - cycle item IDs.")
		}
	}

	if worldEditMode {
		if !clickCaptured {
			if mouseHeld {
				if !leftMousePressed {
					editPlaceItem()
				}
				leftMousePressed = true
			} else {
				leftMousePressed = false
			}
			if gRightMouseHeld {
				if !rightMousePressed {
					editDeleteItem()
				}
				rightMousePressed = true
			} else {
				rightMousePressed = false
			}
		}
		var shiftKey bool
		if repeatingKeyPressed(ebiten.KeyShift) {
			shiftKey = true
		}
		if keyJustPressed(ebiten.KeyEqual) {
			if shiftKey {
				if worldEditID.section < assetArraySize {
					worldEditID.section++
				}
			} else {
				if worldEditID.num < assetArraySize {
					worldEditID.num++
				}
			}
		} else if keyJustPressed(ebiten.KeyMinus) {
			if shiftKey {
				if worldEditID.section > 0 {
					worldEditID.section--
				}
			} else {
				if worldEditID.num > 0 {
					worldEditID.num--
				}
			}
		}

		editPos = XY{X: sCamPos.X - uint32(mouseX), Y: sCamPos.Y - uint32(mouseY)}
	}
}

func walkXY(mx, my int) DIR {

	distance := distance(XY{X: uint32(halfScreenX), Y: uint32(halfScreenY)}, XY{X: uint32(mx), Y: uint32(my)})

	if distance < playerSpriteSize ||
		mx > screenX || my > screenY ||
		mx < 0 || my < 0 {
		return DIR_NONE
	}

	screenCenter := geom.Coord{float64(halfScreenX), float64(halfScreenY), 0}
	mousePosition := geom.Coord{float64(mx), float64(my), 0}

	angle := xy.Angle(screenCenter, mousePosition)

	return radiansToDirection(angle)
}

// keyJustPressed return true when key is pressed considering the repeat state.
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

// keyJustPressed return true when key is pressed considering the repeat state.
func keyJustPressed(key ebiten.Key) bool {

	d := inpututil.KeyPressDuration(key)
	return d == 1
}

func sendMove(nextDirection DIR) {

	//Exit if nothing changed
	if nextDirection == goingDirection {
		if goingDirection == DIR_NONE {
			return
		} else if time.Since(lastNetSend) < directionKeepAlive {
			return
		}
	}

	if time.Since(lastNetSend) < directionThrottle {
		return
	}

	//Update our direction
	goingDirection = nextDirection

	var buf []byte
	outbuf := bytes.NewBuffer(buf)

	binary.Write(outbuf, binary.LittleEndian, &goingDirection)
	sendCommand(CMD_Move, outbuf.Bytes())

	lastNetSend = time.Now()
}

func editPlaceItem() {

	var buf []byte
	outbuf := bytes.NewBuffer(buf)

	itemType := itemTypesList[worldEditID.section]
	if itemType == nil {
		return
	}

	if itemType.name != "wobjects" {
		return
	}

	binary.Write(outbuf, binary.LittleEndian, worldEditID.section)
	binary.Write(outbuf, binary.LittleEndian, worldEditID.num)
	binary.Write(outbuf, binary.LittleEndian, editPos.X)
	binary.Write(outbuf, binary.LittleEndian, editPos.Y)
	sendCommand(CMD_EditPlaceItem, outbuf.Bytes())
}

func editDeleteItem() {

	var buf []byte
	outbuf := bytes.NewBuffer(buf)

	binary.Write(outbuf, binary.LittleEndian, worldEditID)
	binary.Write(outbuf, binary.LittleEndian, editPos.X)
	binary.Write(outbuf, binary.LittleEndian, editPos.Y)
	sendCommand(CMD_EditDeleteItem, outbuf.Bytes())
}
