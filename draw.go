package main

import (
	"fmt"
	"image/color"
	"runtime"
	"sort"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
)

var (
	camPos  XY = worldCenter
	sCamPos XY = worldCenter

	disableNightShadow bool
	nightLevel         uint8 = 0
	startTime          time.Time
	noSmoothing        bool = false
	normal             float64
)

type xySort []*playerData

func (v xySort) Len() int           { return len(v) }
func (v xySort) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v xySort) Less(i, j int) bool { return v[i].pos.Y+v[i].pos.X > v[j].pos.Y+v[j].pos.X }

// Ebiten: Draw everything
func (g *Game) Draw(screen *ebiten.Image) {
	gameModeLock.Lock()
	defer gameModeLock.Unlock()

	if gameMode != MODE_Playing {
		drawBootScreen(screen)
		return
	}

	// If not smoothing, don't draw if there isn't new data
	if noSmoothing && !dataDirty {
		return
	}

	drawLock.Lock()
	defer drawLock.Unlock()

	//We are drawing now, we can clear this flag
	dataDirty = false

	motionSmoothing()

	drawGrass(screen)

	drawDebugEdit(screen)

	drawWorldObjs(screen)

	drawPlayers(screen)

	drawNightVignette(screen)

	drawDebugInfo(screen)

	drawChatLines(screen)

	drawChatBar(screen)

	drawOpenWindows(screen)

	showToolbarCache(screen)

	toolBarTooltip(screen)

}

func drawBootScreen(screen *ebiten.Image) {
	// Boot screen
	op := &ebiten.DrawImageOptions{}
	var imgSize float64 = 1024.0

	scalew := 1.0 / (imgSize / float64(screenX))
	scaleh := 1.0 / (imgSize / float64(screenY))

	op.GeoM.Scale(scalew, scaleh)

	screen.DrawImage(testLogin, op)
	drawChatLines(screen)
}

func drawGrass(screen *ebiten.Image) {
	// Draw grass
	for x := -32; x <= screenX; x += 32 {
		for y := -32; y <= screenY; y += 32 {
			op := ebiten.DrawImageOptions{}
			op.GeoM.Scale(2, 2)
			op.GeoM.Translate(float64(x+int(sCamPos.X%32)), float64(y+int(sCamPos.Y%32)))
			screen.DrawImage(testGrass, &op)
		}
	}
}

func motionSmoothing() {
	if !noSmoothing {
		// Extrapolate position
		startTime = time.Now()
		since := startTime.Sub(lastNetUpdate)
		remaining := FrameSpeedNS - since.Nanoseconds()
		normal = (float64(remaining) / float64(FrameSpeedNS))

		//Extrapolation limits
		if normal < 0 {
			normal = 0
		} else if normal > 1 {
			normal = 1
		}

		// If there ins't new data yet, extrapolate
		if !dataDirty {
			var smoothPos XY

			//Extrapolated local player position
			smoothPos.X = uint32(float64(oldLocalPlayerPos.X) - ((float64(localPlayerPos.X) - float64(oldLocalPlayerPos.X)) * normal))
			smoothPos.Y = uint32(float64(oldLocalPlayerPos.Y) - ((float64(localPlayerPos.Y) - float64(oldLocalPlayerPos.Y)) * normal))

			//Extrapolated camera position
			sCamPos.X = (uint32(halfScreenX)) + smoothPos.X
			sCamPos.Y = (uint32(halfScreenY)) + smoothPos.Y

			//Extrapolated remote players
			for p, player := range playerList {
				var psmooth XY
				psmooth.X = uint32(float64(player.lastPos.X) - ((float64(player.pos.X) - float64(player.lastPos.X)) * normal))
				psmooth.Y = uint32(float64(player.lastPos.Y) - ((float64(player.pos.Y) - float64(player.lastPos.Y)) * normal))
				playerList[p].spos = XY{X: uint32(psmooth.X), Y: uint32(psmooth.Y)}
			}
		}
	} else {
		// Standard mode, just copy data over
		camPos.X = (uint32(halfScreenX)) + localPlayerPos.X
		camPos.Y = (uint32(halfScreenY)) + localPlayerPos.Y
		sCamPos = camPos
		for o := range playerList {
			playerList[o].spos = playerList[o].pos
		}
	}
}

func showToolbarCache(screen *ebiten.Image) {
	defer reportPanic("showToolbarCache")

	toolbarCacheLock.RLock()
	screen.DrawImage(toolbarCache, nil)
	toolbarCacheLock.RUnlock()
}

func drawWorldObjs(screen *ebiten.Image) {
	defer reportPanic("drawWorldObjs")

	//Draw on-ground objects first
	for _, obj := range wObjList {
		if !obj.itemData.OnGround {
			continue
		}

		xPos := float64(int(sCamPos.X) - int(obj.pos.X))
		yPos := float64(int(sCamPos.Y) - int(obj.pos.Y))

		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(2, 2)
		op.GeoM.Translate(xPos-48.0, yPos-48.0)

		screen.DrawImage(spritelist[obj.itemId].image, &op)
	}

	//Draw night shadows
	drawNightShadows(screen)

	//Draw other objects (buildings/trees)
	for _, obj := range wObjList {
		if obj.itemData.OnGround {
			continue
		}

		xPos := float64(int(sCamPos.X) - int(obj.pos.X))
		yPos := float64(int(sCamPos.Y) - int(obj.pos.Y))

		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(2, 2)
		op.GeoM.Translate(xPos-48.0, yPos-48.0)

		screen.DrawImage(spritelist[obj.itemId].image, &op)
	}

}

func drawNightVignette(screen *ebiten.Image) {
	defer reportPanic("drawNightVignette")

	if nightLevel == 0 {
		return
	}

	//Use linear filter for this
	op := ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}

	// Fit this onto the screen
	var screenSize int
	if screenY > screenX {
		screenSize = screenY
	} else {
		screenSize = screenX
	}
	var sc float64
	if screenSize > 1024 {
		sc = (float64(screenSize) / 1024.0) + 0.01
	} else {
		sc = 1.01
	}

	xPos := float64(int(sCamPos.X)-int(playerList[localPlayer.id].spos.X)) - (512 * sc)
	yPos := float64(int(sCamPos.Y)-int(playerList[localPlayer.id].spos.Y)) - (512 * sc)

	op.GeoM.Translate(float64(xPos), float64(yPos))
	op.GeoM.Scale(sc, sc)
	op.ColorScale.ScaleAlpha(float32(nightLevel) / 255.0)

	screen.DrawImage(testlight, &op)
}

func drawPlayers(screen *ebiten.Image) {
	defer reportPanic("drawPlayers")

	// Find visible players and sort them
	var pList []*playerData
	for _, player := range playerList {
		xPos := float64(int(sCamPos.X) - int(player.spos.X))
		yPos := float64(int(sCamPos.Y) - int(player.spos.Y))

		//Sprite on screen?
		if xPos-playerSpriteSize > float64(screenX) {
			continue
		} else if xPos < -playerSpriteSize {
			continue
		} else if yPos-playerSpriteSize > float64(screenY) {
			continue
		} else if yPos < -playerSpriteSize {
			continue
		}
		pList = append(pList, player)
	}
	sort.Sort(xySort(pList))

	// Draw player
	for _, player := range pList {

		xPos := float64(int(sCamPos.X) - int(player.spos.X))
		yPos := float64(int(sCamPos.Y) - int(player.spos.Y))

		op := ebiten.DrawImageOptions{}

		op.GeoM.Scale(2, 2)

		//camera - object, TODO: get sprite size
		op.GeoM.Translate(float64(xPos)-48.0, float64(yPos)-48.0)

		//Draw sub-image
		screen.DrawImage(getCharFrame(player).(*ebiten.Image), &op)
	}

	// Draw player name
	for _, player := range pList {

		var pname string
		pnameStr := playerIdToName(player.id)
		if pnameStr != "" {
			pname = pnameStr
		} else {
			pname = fmt.Sprintf("Player-%v", player.id)
		}

		// Draw name
		drawText(pname, toolTipFont, color.White, colorNameBG,
			XYf32{X: float32(int(sCamPos.X)-int(player.spos.X)) + 4,
				Y: float32(int(sCamPos.Y)-int(player.spos.Y)) + 48},
			2, screen, false, false, true)

	}

	// Draw health
	for _, player := range pList {
		if player.health < 100 {

			var healthColor color.RGBA
			r := int(float32(100-player.health) * 5)
			if r > 255 {
				r = 255
			}
			healthColor.R = uint8(r)
			healthColor.G = uint8(float32(player.health) * 1.5)

			vector.DrawFilledRect(
				screen,
				float32(int(sCamPos.X)-int(player.pos.X))-12+4-1,
				float32(int(sCamPos.Y)-int(player.pos.Y))+24-1,
				27, 4, colorNameBG,
				false)

			vector.DrawFilledRect(
				screen,
				float32(int(sCamPos.X)-int(player.pos.X))-12+4,
				float32(int(sCamPos.Y)-int(player.pos.Y))+24,
				25-((100.0-float32(player.health))/4.0), 2, healthColor,
				false)
		}
	}
}

func drawText(input string, face font.Face, color color.Color, bgcolor color.Color, pos XYf32,
	pad float32, screen *ebiten.Image, justLeft bool, justUp bool, justCenter bool) XYf32 {
	defer reportPanic("DrawText")
	var tmx, tmy float32

	tRect := text.BoundString(face, input)

	if justCenter {
		tmx = float32(int(pos.X) - (tRect.Dx() / 2))
		tmy = float32(int(pos.Y) - (tRect.Dy() / 2))
	} else {
		if justLeft {
			tmx = float32(pos.X)
		} else {
			tmx = float32(int(pos.X) - tRect.Dx())
		}

		if justUp {
			tmy = float32(int(pos.Y) - tRect.Dy())
		} else {
			tmy = float32(pos.Y + float32(tRect.Dy()))
		}
	}
	_, _, _, alpha := bgcolor.RGBA()

	if alpha > 0 {
		vector.DrawFilledRect(
			screen, tmx-2, tmy-12,
			float32(tRect.Dx())+4, float32(tRect.Dy())+4, bgcolor, false,
		)
	}
	text.Draw(screen, input, face, int(tmx), int(tmy), color)

	return XYf32{X: float32(tRect.Dx()) + pad, Y: float32(tRect.Dy()) + pad}
}

var chatVertSpace float32 = 24.0 * float32(uiScale)

var (
	chatLinesTop  int
	chatLines     []chatLineData
	chatLinesLock sync.Mutex
	consoleActive bool
)

const (
	// Number of chat lines to display at once
	chatHeightLines = 20
	// Default fade out time
	chatFadeTime = time.Second * 3

	padding = 8
	linePad = 2
)

func drawChatLines(screen *ebiten.Image) {
	defer reportPanic("drawChatLines")
	var lineNum int
	chatLinesLock.Lock()
	defer chatLinesLock.Unlock()

	for x := chatLinesTop; x > 0 && lineNum < chatHeightLines; x-- {
		line := chatLines[x-1]
		// Ignore old chat lines
		since := startTime.Sub(line.timestamp)
		if !consoleActive && since > line.lifetime {
			continue
		}
		lineNum++

		// BG
		tempBGColor := colorNameBG
		// Text color
		r, g, b, _ := line.color.RGBA()

		// Alpha + fade out
		var blend float64 = 0
		if line.lifetime-since < chatFadeTime {
			blend = (float64(chatFadeTime-(line.lifetime-since)) / float64(chatFadeTime) * 100.0)
		}
		newAlpha := (254.0 - (blend * 2.55))
		oldAlpha := tempBGColor.A
		faded := newAlpha - float64(253.0-int(oldAlpha))
		if faded <= 0 {
			faded = 0
		} else if faded > 254 {
			faded = 254
		}
		tempBGColor.A = byte(faded)

		drawText(line.text, generalFont,
			color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: byte(newAlpha)},
			tempBGColor, XYf32{X: padding, Y: float32(screenY) - (float32(lineNum) * (float32(generalFontH) * 1.66)) - chatVertSpace},
			4, screen, true, false, false)
	}
}

func drawDebugInfo(screen *ebiten.Image) {
	defer reportPanic("drawDebugInfo")

	if !debugLine {
		return
	}

	// Draw debug info
	buf := fmt.Sprintf("FPS: %3v  Arch: %v  Build: v%v",
		int(ebiten.ActualFPS()),
		runtime.GOARCH, gameVersion,
	)

	drawText(buf, monoFont, color.White, colorNameBG,
		XYf32{X: float32(toolbarCache.Bounds().Dx()) + 10, Y: 10},
		1, screen, true, false, false)

}

func drawDebugEdit(screen *ebiten.Image) {
	defer reportPanic("drawDebugEdit")

	if !worldEditMode {
		return
	}

	xPos := float64(int(sCamPos.X) - int(editPos.X))
	yPos := float64(int(sCamPos.Y) - int(editPos.Y))

	op := ebiten.DrawImageOptions{}

	op.GeoM.Scale(2, 2)

	//Draw edit sprite
	if worldEditID < topSpriteID {
		op.GeoM.Translate(xPos, yPos)
		screen.DrawImage(spritelist[worldEditID].image, &op)
	}

	// Draw debug info
	buf := fmt.Sprintf("EDIT MODE ON: ID: %v", worldEditID)

	drawText(buf, monoFont, color.White, colorNameBG,
		XYf32{X: float32(screenX) - 4, Y: 2},
		1, screen, false, false, false)

}

func drawChatBar(screen *ebiten.Image) {
	defer reportPanic("drawDebugInfo")

	text := ""
	if ChatText == "" {
		if ChatMode {
			text = "say: _"
		} else if CommandMode {
			text = "(/h for help.)> _"
		}
	} else {
		if ChatMode {
			text = "say: " + ChatText
		} else if CommandMode {
			text = " > " + ChatText
		}
	}

	if ChatMode {
		var pad float32 = 4 * float32(uiScale)
		drawText(text, monoFont, color.White, colorNameBG,
			XYf32{X: (pad * 1.5), Y: float32(screenY) + (pad * 2)},
			pad, screen, true, true, false)
	} else if CommandMode {
		var pad float32 = 4 * float32(uiScale)
		drawText(text, monoFont, color.White, colorNameBG,
			XYf32{X: (pad * 1.5), Y: float32(screenY) + (pad * 2)},
			pad, screen, true, true, false)
	}

}
