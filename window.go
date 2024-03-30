package main

import (
	"image"
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var windowsLock sync.Mutex

var windows []*windowData = []*windowData{
	{
		idName:      "game options",
		title:       "Options",
		size:        XYs{X: 250, Y: 275},
		centered:    true,
		closeable:   true,
		windowDraw:  drawOptionsWindow,
		windowSetup: setupOptionsWindow,
		movable:     true,
		windowMouse: handleOptions,
	},
	{
		idName:      "game help",
		title:       "Help & Controls",
		size:        XYs{X: 300, Y: 270},
		centered:    false,
		position:    XYs{X: 8, Y: 46},
		closeable:   true,
		windowDraw:  drawHelpWindow,
		windowMouse: handleHelpWindow,
		movable:     true,
	},
	{
		idName:      "game login",
		title:       "Login to goMMO:",
		size:        XYs{X: 300, Y: 200},
		movable:     false,
		opaque:      true,
		centered:    true,
		closeable:   false,
		windowDraw:  drawLoginWindow,
		windowMouse: handleLoginWindow,
		windowKeys:  handleLoginKeys,
	},
}

var openWindows map[string]*windowData

type windowData struct {
	idName string
	title  string // Window title

	focused    bool //Mouse is on window
	movable    bool // Can be dragged
	opaque     bool // Non-semitransparent background
	centered   bool // Auto-centered
	borderless bool
	closeable  bool // Has a close-x in title bar
	keepCache  bool // Draw cache persists when window is closed
	dragPos    XYs  // Position where window drag began

	titleButtons titleButtonData // Window title buttons
	buttons      []buttonData

	size       XYs // Size in pixels
	scaledSize XYs // Size with UI scale
	position   XYs // Position

	bgColor      *color.Color // Custom BG color
	titleBGColor *color.Color // Custom title bar background color
	titleColor   *color.Color // Custom title text color

	dirty       bool          // Needs to be redrawn
	cache       *ebiten.Image // Cache image
	windowDraw  func(Window *windowData)
	windowMouse func(input XYs, window *windowData) bool
	windowKeys  func(input []rune, window *windowData, delete, enter bool) bool
	windowSetup func(Window *windowData)
}

type buttonData struct {
	text     string
	image    *ebiten.Image
	rect     image.Rectangle
	activate func()
}

type titleButtonData struct {
	closePos       XYs
	closeSize      XYs
	titleBarHeight int
}

// Toggle settings window
func settingsToggle() {
	defer reportPanic("settingsToggle")

	if openWindows["game options"] != nil {
		closeWindow("game options")
	} else {
		openWindow("game options")
	}
}

// Toggle help window
func toggleHelp() {
	defer reportPanic("toggleHelp")

	if openWindows["game help"] != nil {
		closeWindow("game help")
	} else {
		openWindow("game help")
	}
}

// Allow windows to do any precalculation they need to do
func initWindows() {
	defer reportPanic("initWindows")

	for _, win := range windows {
		if win.windowSetup != nil {
			win.windowSetup(win)
			win.dirty = true
		}
	}
}

// Draw whatever windows are currently open
func drawOpenWindows(screen *ebiten.Image) {
	defer reportPanic("drawOpenWindows")

	for _, win := range openWindows {
		if win.focused {
			drawWindow(screen, win)
		}
	}
	for _, win := range openWindows {
		if !win.focused {
			drawWindow(screen, win)
		}
	}
}

// Open a window
// Until layering is added, close other windows if we open one
func openWindow(idName string) {
	defer reportPanic("openWindow")

	windowsLock.Lock()
	defer windowsLock.Unlock()

	//Window is already open
	if openWindows[idName] != nil {
		return
	}

	var window *windowData
	for wpos, win := range windows {
		if win.idName == idName {
			window = windows[wpos]
			break
		}
	}

	//Window not found
	if window == nil {
		return
	}

	if window.centered && window.movable {

		window.scaledSize = XYs{X: int32(float64(window.size.X) * uiScale), Y: int32(float64(window.size.Y) * uiScale)}
		window.position = XYs{
			X: int32(screenX/2) - (window.scaledSize.X / 2),
			Y: int32(screenY/2) - (window.scaledSize.Y / 2)}
	}

	if windowDebugMode {
		doLog(true, "Window '%v' added to open list.", window.title)
	}

	//Close open windows, for now
	openWindows = make(map[string]*windowData)

	//Add to open windows list
	openWindows[window.idName] = window
}

// Close a window
func closeWindow(idName string) {
	defer reportPanic("closeWindow")
	windowsLock.Lock()
	defer windowsLock.Unlock()

	window := openWindows[idName]

	//No such window
	if window == nil {
		return
	}

	//Handle window closed while dragging
	if draggingWindow == window {
		draggingWindow = nil
	}

	//Dispose window image cache if needed
	if !window.keepCache && window.cache != nil {
		if windowDebugMode {
			doLog(true, "Window '%v' closed, disposing cache.", window.title)
		}
		window.cache.Dispose()
		window.cache = nil
	}

	// Eat click event
	clickCaptured = true
	mouseHeld = false

	//Remove from open list
	delete(openWindows, window.idName)
}

const closePad = 18
const closeScale = 0.7

// Draw window title, frame, background and cached window contents
func dragWindow() {
	windowsLock.Lock()
	defer windowsLock.Unlock()

	if draggingWindow == nil {
		return
	}

	draggingWindow.position = XYs{X: int32(mouseX) - draggingWindow.dragPos.X, Y: int32(mouseY) - draggingWindow.dragPos.Y}
	clickCaptured = true

	// Clamp position
	clampWindow(draggingWindow)
}

func clampUIWindows() {
	windowsLock.Lock()
	defer windowsLock.Unlock()

	for _, window := range openWindows {
		clampWindow(window)
	}
}

func clampWindow(w *windowData) {
	if w == nil {
		return
	}

	//Clamp position
	if w.position.X < 0 {
		w.position.X = 0
	} else if w.position.X > int32(screenX)-w.size.X {
		w.position.X = int32(screenX) - w.size.X
	}
	if w.position.Y < 0 {
		w.position.Y = 0
	} else if w.position.Y > int32(screenY)-w.size.Y {
		w.position.Y = int32(screenY) - w.size.Y
	}
}

func drawWindow(screen *ebiten.Image, window *windowData) {
	defer reportPanic("drawWindow")
	windowsLock.Lock()
	defer windowsLock.Unlock()

	//Calculate some values for UI scale
	pad := int(closePad * uiScale)
	halfPad := int((closePad / 2.0) * uiScale)

	winPos := getWindowPos(window)
	window.scaledSize = XYs{X: int32(float64(window.size.X) * uiScale), Y: int32(float64(window.size.Y) * uiScale)}

	//If window not dirty, and it has a cache, draw the cache
	if !window.dirty {
		if window.cache != nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(winPos.X), float64(winPos.Y))
			screen.DrawImage(window.cache, op)
			return
		}
	}

	//If there is no window cache, init it
	if window.cache == nil {
		window.cache = ebiten.NewImage(int(window.scaledSize.X), int(window.scaledSize.Y))
		if windowDebugMode {
			doLog(true, "Window '%v' cache initalized.", window.title)
		}
	} else {
		window.cache.Clear()
	}

	//Custom colors
	var winBG color.Color
	if window.bgColor != nil {
		winBG = *window.bgColor
	} else if window.opaque {
		winBG = ColorWindowBGO
	} else {
		winBG = ColorWindowBG
	}

	var titleBGColor color.Color
	if window.titleBGColor != nil {
		titleBGColor = *window.titleBGColor
	} else {
		titleBGColor = ColorWindowTitle
	}

	var titleColor color.Color
	if window.titleBGColor != nil {
		titleColor = *window.titleColor
	} else {
		titleColor = color.White
	}

	//Draw window BG
	vector.DrawFilledRect(
		window.cache,
		0, 0,
		float32(window.scaledSize.X), float32(window.scaledSize.Y),
		winBG, false)

	if window.title != "" {

		fHeight := text.BoundString(generalFont, "!Aa0")

		// Border
		if !window.borderless {
			vector.DrawFilledRect(
				window.cache, 0, +float32(window.scaledSize.Y)-1,
				float32(window.scaledSize.X), 2, titleBGColor, false,
			)
			vector.DrawFilledRect(
				window.cache,
				0, 0,
				2, float32(window.scaledSize.Y),
				titleBGColor, false)
			vector.DrawFilledRect(
				window.cache,
				float32(window.scaledSize.X)-1, 0,
				2, float32(window.scaledSize.Y),
				titleBGColor, false)
		}

		// Title bar
		vector.DrawFilledRect(
			window.cache, 0, 0,
			float32(window.scaledSize.X), float32(float64(fHeight.Dy()))+float32(pad), titleBGColor, false,
		)
		window.titleButtons.titleBarHeight = fHeight.Dy() + pad

		text.Draw(window.cache, window.title, generalFont, halfPad, int(fHeight.Dy()+halfPad), titleColor)

		if window.closeable {
			img := closeBox
			op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
			closePosX := float64(window.scaledSize.X - int32(float64(img.Bounds().Dx())*uiScale*closeScale))
			op.GeoM.Scale(uiScale*closeScale, uiScale*closeScale)
			op.GeoM.Translate(closePosX, 0)

			// save button positions
			window.titleButtons.closePos = XYs{X: int32(closePosX), Y: int32(0)}
			window.titleButtons.closeSize = XYs{X: int32(float64(img.Bounds().Dx()) * uiScale),
				Y: int32(float64(img.Bounds().Dy()) * uiScale)}
			window.cache.DrawImage(img, op)
		}
	}

	// Call custom draw function, if it exists
	if window.windowDraw != nil {
		window.windowDraw(window)
	}

	window.dirty = false

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(winPos.X), float64(winPos.Y))
	screen.DrawImage(window.cache, op)
}

// Check if a click is within an open window
func collisionWindowsCheck(input XYs) bool {
	defer reportPanic("collisionWindowsCheck")
	if clickCaptured {
		return true
	}
	for _, win := range openWindows {
		if collisionWindow(input, win) {
			return true
		}
	}

	return false
}

// Check if a click is within a specific window
func collisionWindow(input XYs, window *windowData) bool {
	defer reportPanic("collisionWindow")
	winPos := getWindowPos(window)

	if input.X > winPos.X && input.X < winPos.X+window.scaledSize.X &&
		input.Y > winPos.Y && input.Y < winPos.Y+window.scaledSize.Y {
		if !window.focused {
			window.focused = true
		}

		// Handle X close
		if handleClose(input, window) {
			return true
		}

		if handleDrag(input, window) {
			return true
		}

		// Handle input
		if window.windowMouse != nil {
			window.windowMouse(input, window)
		}

		return true
	} else {
		if window.focused {
			window.focused = false
		}
		return false
	}
}

// Check if a click was within a window's close box
func handleClose(input XYs, window *windowData) bool {
	defer reportPanic("handleCLose")
	if draggingWindow != nil {
		return false
	}
	if !window.closeable {
		return false
	}

	winPos := getWindowPos(window)
	if input.X > winPos.X+window.scaledSize.X-window.titleButtons.closeSize.X &&
		input.X < winPos.X+window.scaledSize.X &&
		input.Y < winPos.Y+window.titleButtons.closeSize.Y &&
		input.Y > winPos.Y {
		closeWindow(window.idName)
		return true
	}

	return false
}

// Handle dragging windows
func handleDrag(input XYs, window *windowData) bool {
	defer reportPanic("handleDrag")
	if !mouseHeld {
		return false
	}
	if !window.movable {
		return false
	}
	if draggingWindow != nil {
		return true
	}

	winPos := getWindowPos(window)
	winOff := XYs{X: input.X - winPos.X, Y: input.Y - winPos.Y}

	if input.X > winPos.X &&
		input.X < winPos.X+window.scaledSize.X &&
		input.Y > winPos.Y &&
		input.Y < winPos.Y+int32(window.titleButtons.titleBarHeight) {
		draggingWindow = window
		draggingWindow.dragPos = winOff
		doLog(true, "dragging window '%v'", window.title)
		return true
	}
	return false
}

// Get window position, assists with auto-centered windows
func getWindowPos(window *windowData) XYs {
	defer reportPanic("getWindowPos")
	var winPos XYs
	if window.centered && !window.movable {
		winPos.X, winPos.Y = int32(screenX/2)-(window.scaledSize.X/2), int32(screenY/2)-(window.scaledSize.Y/2)
	} else {
		winPos = window.position
	}
	return winPos
}
