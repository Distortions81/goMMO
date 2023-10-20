package main

import (
	"fmt"
	"image/color"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
)

const tempOff = 128

/* Ebiten: Draw everything */
func (g *Game) Draw(screen *ebiten.Image) {

	if gameMode == MODE_PLAYING {

		playerListLock.Lock()
		defer playerListLock.Unlock()

		if !dataDirty {
			return
		}
		dataDirty = false

		screen.Fill(colorGrass)

		//center of screen, center of sprite, charpos
		for _, player := range playerList {
			op := ebiten.DrawImageOptions{}

			convPos := convPos(player.pos)

			op.GeoM.Translate(tempOff-26+float64(convPos.X), tempOff-26+float64(convPos.Y))
			//Upscale
			op.GeoM.Scale(2, 2)

			//Draw sub-image
			screen.DrawImage(getCharFrame(player).(*ebiten.Image), &op)
			pname := fmt.Sprintf("Player-%v", player.id)
			pos := XYf32{X: float32(tempOff+convPos.X) * 2.0,
				Y: float32(tempOff+convPos.Y)*2.0 + 40}
			drawText(pname, toolTipFont, color.White, colorNameBG,
				pos, 2, screen, false, false, true)
		}
		buf := fmt.Sprintf("%3.0f FPS", ebiten.ActualFPS())
		ebitenutil.DebugPrint(screen, buf)

		drawChatLines(screen)

	} else {
		ebitenutil.DebugPrint(screen, "Connecting.")
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
		fHeight := text.BoundString(face, "gpqabcABC!|_,;^*`")
		vector.DrawFilledRect(
			screen, tmx-pad, tmy-float32(fHeight.Dy())-(float32(pad)/2.0),
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
	linePad     = 1
)

func drawChatLines(screen *ebiten.Image) {
	defer reportPanic("drawChatLines")
	var lineNum int
	chatLinesLock.Lock()
	defer chatLinesLock.Unlock()

	for x := chatLinesTop; x > 0 && lineNum < chatHeightLines; x-- {
		line := chatLines[x-1]
		/* Ignore old chat lines */
		since := time.Since(line.timestamp)
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
			color.NRGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: byte(newAlpha)},
			tempBGColor, XYf32{X: padding, Y: float32(screenHeight) - (float32(lineNum) * (float32(generalFontH) * 1.2)) - chatVertSpace},
			2, screen, true, false, false)
	}
}
