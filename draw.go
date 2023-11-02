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

var camPos XY = xyCenter
var smoothCamPos XY = xyCenter

var disableNightShadow bool
var nightLevel uint8 = 0
var startTime time.Time
var noSmoothing bool = false

type xySort []*playerData

func (v xySort) Len() int           { return len(v) }
func (v xySort) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v xySort) Less(i, j int) bool { return v[i].pos.Y+v[i].pos.X > v[j].pos.Y+v[j].pos.X }

/* Ebiten: Draw everything */
func (g *Game) Draw(screen *ebiten.Image) {

	gameModeLock.Lock()
	defer gameModeLock.Unlock()

	if gameMode == MODE_PLAYING {

		startTime = time.Now()

		if noSmoothing {
			if !dataDirty {
				return
			}
		}

		drawLock.Lock()
		defer drawLock.Unlock()

		/* Extrapolate position */
		since := startTime.Sub(lastUpdate)
		remaining := FrameSpeedNS - since.Nanoseconds()
		normal := (float64(remaining) / float64(FrameSpeedNS))

		if !noSmoothing && !dataDirty {
			var smoothPos XY
			if normal < 0 {
				normal = 0
			} else if normal > 1 {
				normal = 1
			}

			smoothPos.X = uint32(float64(ourOldPos.X) - ((float64(ourPos.X) - float64(ourOldPos.X)) * normal))
			smoothPos.Y = uint32(float64(ourOldPos.Y) - ((float64(ourPos.Y) - float64(ourOldPos.Y)) * normal))

			smoothCamPos.X = (uint32(HscreenWidth)) + smoothPos.X
			smoothCamPos.Y = (uint32(HscreenHeight)) + smoothPos.Y
		}

		camPos.X = (uint32(HscreenWidth)) + ourPos.X
		camPos.Y = (uint32(HscreenHeight)) + ourPos.Y

		if noSmoothing {
			smoothCamPos = camPos
			for o, _ := range playerList {
				playerList[o].spos = playerList[o].pos
			}
		}

		if !noSmoothing && !dataDirty {
			for p, player := range playerList {

				/* Extrapolate position */
				since := startTime.Sub(lastUpdate)
				remaining := FrameSpeedNS - since.Nanoseconds()
				normal := (float64(remaining) / float64(FrameSpeedNS))

				var psmooth XY

				if normal < 0 {
					normal = 0
				} else if normal > 1 {
					normal = 1
				}

				psmooth.X = uint32(float64(player.lastPos.X) - ((float64(player.pos.X) - float64(player.lastPos.X)) * normal))
				psmooth.Y = uint32(float64(player.lastPos.Y) - ((float64(player.pos.Y) - float64(player.lastPos.Y)) * normal))
				playerList[p].spos = XY{X: uint32(psmooth.X), Y: uint32(psmooth.Y)}

			}
		}

		dataDirty = false

		/* Draw grass */
		for x := -32; x <= screenWidth; x += 32 {
			for y := -32; y <= screenHeight; y += 32 {
				op := ebiten.DrawImageOptions{}

				op.GeoM.Scale(2, 2)
				op.GeoM.Translate(float64(x+int(smoothCamPos.X%32)), float64(y+int(smoothCamPos.Y%32)))

				screen.DrawImage(testGrass, &op)
			}

		}

		if EditMode {
			drawDebugEdit(screen)
		}

		drawWorld(screen)
		drawPlayers(screen)

		drawLight(screen)

		drawDebugInfo(screen)
		drawChatLines(screen)
		drawChatBar(screen)

		drawOpenWindows(screen)
	} else {
		op := &ebiten.DrawImageOptions{}
		var imgSize float64 = 1024.0

		scalew := 1.0 / (imgSize / float64(screenWidth))
		scaleh := 1.0 / (imgSize / float64(screenHeight))

		op.GeoM.Scale(scalew, scaleh)

		screen.DrawImage(testLogin, op)
		drawChatLines(screen)
	}
}

func drawWorld(screen *ebiten.Image) {
	/* Draw obj */
	for _, obj := range wObjList {
		if !obj.itemData.OnGround {
			continue
		}

		xPos := float64(int(smoothCamPos.X) - int(obj.pos.X))
		yPos := float64(int(smoothCamPos.Y) - int(obj.pos.Y))

		op := ebiten.DrawImageOptions{}

		op.GeoM.Scale(2, 2)

		//camera - object, TODO: get sprite size
		op.GeoM.Translate(xPos-48.0, yPos-48.0)

		//Draw sub-image
		screen.DrawImage(spritelist[obj.itemId].image, &op)
	}

	drawRays(screen)

	for _, obj := range wObjList {
		if obj.itemData.OnGround {
			continue
		}

		xPos := float64(int(smoothCamPos.X) - int(obj.pos.X))
		yPos := float64(int(smoothCamPos.Y) - int(obj.pos.Y))

		op := ebiten.DrawImageOptions{}

		op.GeoM.Scale(2, 2)

		//camera - object, TODO: get sprite size
		op.GeoM.Translate(xPos-48.0, yPos-48.0)

		//Draw sub-image
		screen.DrawImage(spritelist[obj.itemId].image, &op)
	}

}

func drawLight(screen *ebiten.Image) {

	if nightLevel == 0 {
		return
	}

	op := ebiten.DrawImageOptions{}
	var screenSize int
	if screenHeight > screenWidth {
		screenSize = screenHeight
	} else {
		screenSize = screenWidth
	}
	var sc float64
	if screenSize > 1024 {
		sc = (float64(screenSize) / 1024.0) + 0.01
	} else {
		sc = 1.01
	}

	xPos := float64(int(smoothCamPos.X)-int(playerList[localPlayer.id].spos.X)) - (512 * sc)
	yPos := float64(int(smoothCamPos.Y)-int(playerList[localPlayer.id].spos.Y)) - (512 * sc)

	op.GeoM.Translate(float64(xPos), float64(yPos))
	op.GeoM.Scale(sc, sc)
	op.ColorScale.ScaleAlpha(float32(nightLevel) / 255.0)

	screen.DrawImage(testlight, &op)
}

func drawPlayers(screen *ebiten.Image) {

	/* Find visible players and sort them */
	var pList []*playerData
	for _, player := range playerList {
		xPos := float64(int(smoothCamPos.X) - int(player.spos.X))
		yPos := float64(int(smoothCamPos.Y) - int(player.spos.Y))

		//Sprite on screen?
		if xPos-charSpriteSize > float64(screenWidth) {
			continue
		} else if xPos < -charSpriteSize {
			continue
		} else if yPos-charSpriteSize > float64(screenHeight) {
			continue
		} else if yPos < -charSpriteSize {
			continue
		}
		pList = append(pList, player)
	}
	sort.Sort(xySort(pList))

	/* Draw player */
	for _, player := range pList {

		xPos := float64(int(smoothCamPos.X) - int(player.spos.X))
		yPos := float64(int(smoothCamPos.Y) - int(player.spos.Y))

		op := ebiten.DrawImageOptions{}

		op.GeoM.Scale(2, 2)

		//camera - object, TODO: get sprite size
		op.GeoM.Translate(float64(xPos)-48.0, float64(yPos)-48.0)

		//Draw sub-image
		screen.DrawImage(getCharFrame(player).(*ebiten.Image), &op)
	}

	/* Draw name */
	for _, player := range pList {

		var pname string
		pnameStr := getName(player.id)
		if pnameStr != "" {
			pname = pnameStr
		} else {
			pname = fmt.Sprintf("Player-%v", player.id)
		}

		// Draw name
		drawText(pname, toolTipFont, color.White, colorNameBG,
			XYf32{X: float32(int(smoothCamPos.X)-int(player.spos.X)) + 4,
				Y: float32(int(smoothCamPos.Y)-int(player.spos.Y)) + 48},
			2, screen, false, false, true)

	}

	/* Draw health */
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
				float32(int(smoothCamPos.X)-int(player.pos.X))-12+4-1,
				float32(int(smoothCamPos.Y)-int(player.pos.Y))+24-1,
				27, 4, colorNameBG,
				false)

			vector.DrawFilledRect(
				screen,
				float32(int(smoothCamPos.X)-int(player.pos.X))-12+4,
				float32(int(smoothCamPos.Y)-int(player.pos.Y))+24,
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
			screen, tmx-pad, tmy-float32(tRect.Dy())-(float32(pad)/2.0),
			float32(tRect.Dx())+pad*2, float32(tRect.Dy())+pad*2, bgcolor, false,
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
	/* Number of chat lines to display at once */
	chatHeightLines = 20
	/* Default fade out time */
	chatFadeTime = time.Second * 3

	padding     = 8
	scaleFactor = 1.5
	linePad     = 2
)

func drawChatLines(screen *ebiten.Image) {
	defer reportPanic("drawChatLines")
	var lineNum int
	chatLinesLock.Lock()
	defer chatLinesLock.Unlock()

	for x := chatLinesTop; x > 0 && lineNum < chatHeightLines; x-- {
		line := chatLines[x-1]
		/* Ignore old chat lines */
		since := startTime.Sub(line.timestamp)
		if !consoleActive && since > line.lifetime {
			continue
		}
		lineNum++

		/* BG */
		tempBGColor := colorNameBG
		/* Text color */
		r, g, b, _ := line.color.RGBA()

		/* Alpha + fade out */
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
			tempBGColor, XYf32{X: padding, Y: float32(screenHeight) - (float32(lineNum) * (float32(generalFontH) * 1.66)) - chatVertSpace},
			4, screen, true, false, false)
	}
}

func drawDebugInfo(screen *ebiten.Image) {
	defer reportPanic("drawDebugInfo")

	if !infoLine {
		return
	}

	/* Draw debug info */
	buf := fmt.Sprintf("FPS: %3v  Arch: %v  Build: v%v",
		int(ebiten.ActualFPS()),
		runtime.GOARCH, gameVersion,
	)

	drawText(buf, monoFont, color.White, colorNameBG,
		XYf32{X: 4, Y: 24},
		1, screen, true, false, false)

}

func drawDebugEdit(screen *ebiten.Image) {
	defer reportPanic("drawDebugEdit")

	xPos := float64(int(camPos.X) - int(editPos.X))
	yPos := float64(int(camPos.Y) - int(editPos.Y))

	op := ebiten.DrawImageOptions{}

	op.GeoM.Scale(2, 2)

	//Draw edit sprite
	if EditID < numSprites {
		op.GeoM.Translate(xPos, yPos)
		screen.DrawImage(spritelist[EditID].image, &op)
	}

	/* Draw debug info */
	buf := fmt.Sprintf("EDIT MODE ON: ID: %v", EditID)

	drawText(buf, monoFont, color.White, colorNameBG,
		XYf32{X: 1, Y: 1},
		1, screen, true, false, false)

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
			text = "> " + ChatText
		}
	}

	if ChatMode {
		var pad float32 = 4 * float32(uiScale)
		drawText(text, monoFont, color.White, colorNameBG,
			XYf32{X: (pad * 1.5), Y: float32(screenHeight) + (pad * 2)},
			pad, screen, true, true, false)
	} else if CommandMode {
		var pad float32 = 4 * float32(uiScale)
		drawText(text, monoFont, color.White, colorNameBG,
			XYf32{X: (pad * 1.5), Y: float32(screenHeight) + (pad * 2)},
			pad, screen, true, true, false)
	}

}
