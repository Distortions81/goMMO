package main

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const scaleFactor = 1.5

// Calculate spacing and order based on DPI and scale
func setupOptionsWindow(window *windowData) {
	defer reportPanic("setupOptionsWindow")
	optionWindowButtons = []image.Rectangle{}

	// Loop all settings
	optNum := 0
	for pos := range settingItems {

		// Place line
		settingItems[pos].textPosX = int(padding * uiScale)
		settingItems[pos].textPosY = int((float64(generalFontH)*scaleFactor)*float64(optNum+linePad)) + int(padding*uiScale)

		// Generate button
		button := image.Rectangle{}
		if (WASMMode && !settingItems[pos].wasmExclude) || !WASMMode {
			button.Min.X = 0
			button.Max.X = xyMax

			button.Min.Y = int((float64(generalFontH)*scaleFactor)*float64(optNum)) + int(padding*uiScale)
			button.Max.Y = int((float64(generalFontH)*scaleFactor)*float64(optNum+linePad)) + int(padding*uiScale)

			optionWindowButtons = append(optionWindowButtons, button)
		}

		if (WASMMode && !settingItems[pos].wasmExclude) || !WASMMode {
			optNum++
		}
	}

}

var blinkCursor bool
var lastBlink time.Time

func addCursor(input []rune, sel LOGIN_SELECTION) string {

	var output []rune = input
	if sel == login.selected && blinkCursor {
		output = append(output, '_')
		return string(output)
	}

	return string(output)
}

// Draw the help window content
func drawLoginWindow(window *windowData) {
	defer reportPanic("drawLoginWindow")

	if blinkCursor {
		blinkCursor = false
	} else {
		blinkCursor = true
	}

	var loginColor = ColorGray
	var passwordColor = ColorGray
	var connectColor = ColorGray

	if login.selected == SELECTED_LOGIN {
		loginColor = ColorWhite
	} else if login.selected == SELECTED_PASSWORD {
		passwordColor = ColorWhite
	} else if login.selected == SELECTED_GO {
		connectColor = ColorWhite
	}

	buf := fmt.Sprintf("Login: %v", addCursor(login.login, SELECTED_LOGIN))
	drawText(buf, monoFont, loginColor, ColorWindowTitle,
		XYf32{X: 20, Y: 90},
		8, window.cache, JUST_LEFT, JUST_UP)

	buf = fmt.Sprintf("Pass:  %v", addCursor(login.password, SELECTED_PASSWORD))
	drawText(buf, monoFont, passwordColor, ColorWindowTitle,
		XYf32{X: 20, Y: 130},
		8, window.cache, JUST_LEFT, JUST_UP)

	drawText("Connect", largeGeneralFont, connectColor, ColorWindowTitle,
		XYf32{X: float32(window.scaledSize.X / 2), Y: 175},
		8, window.cache, JUST_CENTER, JUST_UP)
}

// Draw the help window content
func drawHelpWindow(window *windowData) {
	defer reportPanic("drawHelpWindow")

	drawText("\n"+helpText, generalFont, color.White, color.Transparent,
		XYf32{X: 6, Y: 30},
		0, window.cache, JUST_LEFT, JUST_DOWN)
}

// Draw options window content
const checkScale = 0.5

func drawOptionsWindow(window *windowData) {
	defer reportPanic("drawOptionsWindow")
	var txt string

	d := 0

	// Draw items
	for i, item := range settingItems {
		b := optionWindowButtons[i]

		// Text
		if !item.noCheck {
			txt = fmt.Sprintf("%v: %v", item.text, BoolToOnOff(item.Enabled))
		} else {
			txt = item.text
		}

		if d%2 != 0 {
			vector.DrawFilledRect(window.cache,
				float32(b.Min.X),
				float32(b.Min.Y)+4,
				float32(b.Size().X/2),
				float32(b.Size().Y/2)-6,
				color.NRGBA{R: 255, G: 255, B: 255, A: 16}, false)
		}

		if (WASMMode && !item.wasmExclude) || !WASMMode {
			text.Draw(window.cache, txt, generalFont, item.textPosX, item.textPosY-(generalFontH/2), color.White)
		} else {
			text.Draw(window.cache, txt, generalFont, item.textPosX, item.textPosY-(generalFontH/2), ColorVeryDarkGray)
		}

		// if the item can be toggled, draw checkmark
		if !item.noCheck {

			//Get checkmark image
			op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
			var check *ebiten.Image
			if item.Enabled {
				check = checkOn
			} else {
				check = checkOff
			}
			// Draw checkmark
			op.GeoM.Scale(uiScale*checkScale, uiScale*checkScale)
			op.GeoM.Translate(
				float64(window.scaledSize.X)-(float64(check.Bounds().Dx())*uiScale)-(padding*uiScale),
				float64(item.textPosY-5)-(float64(check.Bounds().Dy())*uiScale*checkScale))

			// Skip some entries for WASM mode
			if (WASMMode && !item.wasmExclude) || !WASMMode {
				window.cache.DrawImage(check, op)
			}
		}
		d++

	}
}
