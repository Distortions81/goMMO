package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

/* Calculate spacing and order based on DPI and scale */
func setupOptionsWindow(window *windowData) {
	defer reportPanic("setupOptionsWindow")
	optionWindowButtons = []image.Rectangle{}

	/* Loop all settings */
	optNum := 1
	for pos := range settingItems {

		/* Place line */
		settingItems[pos].TextPosX = int(padding * uiScale)
		settingItems[pos].TextPosY = int((float64(generalFontH)*scaleFactor)*float64(optNum+linePad)) + int(padding*uiScale)

		/* Generate button */
		button := image.Rectangle{}
		if (WASMMode && !settingItems[pos].WASMExclude) || !WASMMode {
			button.Min.X = 0
			button.Max.X = xyMax

			button.Min.Y = int((float64(generalFontH)*scaleFactor)*float64(optNum)) + int(padding*uiScale)
			button.Max.Y = int((float64(generalFontH)*scaleFactor)*float64(optNum+linePad)) + int(padding*uiScale)
		}
		optionWindowButtons = append(optionWindowButtons, button)

		if (WASMMode && !settingItems[pos].WASMExclude) || !WASMMode {
			optNum++
		}
	}

}

/* Draw the help window content */
func drawHelpWindow(window *windowData) {
	defer reportPanic("drawHelpWindow")

	drawText("\n"+helpText, generalFont, color.White, color.Transparent,
		XYf32{X: float32(window.scaledSize.X/2) + 10, Y: float32(window.scaledSize.Y / 2)},
		0, window.cache, false, false, true)
}

/* Draw the help window content */
var updateVersion string
var downloadURL string

/* Draw options window content */
const checkScale = 0.5

func drawOptionsWindow(window *windowData) {
	defer reportPanic("drawOptionsWindow")
	var txt string

	d := 0

	/* Draw items */
	for i, item := range settingItems {
		b := optionWindowButtons[i]

		/* Text */
		if !item.NoCheck {
			txt = fmt.Sprintf("%v: %v", item.Text, BoolToOnOff(item.Enabled))
		} else {
			txt = item.Text
		}

		if d%2 == 0 {
			vector.DrawFilledRect(window.cache,
				float32(b.Min.X),
				float32(b.Max.Y),
				float32(b.Size().X/2),
				float32(b.Size().Y),
				color.NRGBA{R: 255, G: 255, B: 255, A: 16}, false)
		}

		/* Skip some entries for WASM mode */
		if (WASMMode && !item.WASMExclude) || !WASMMode {

			text.Draw(window.cache, txt, generalFont, item.TextPosX, item.TextPosY-(generalFontH/2), color.White)

			/* if the item can be toggled, draw checkmark */
			if !item.NoCheck {

				/* Get checkmark image */
				op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
				var check *ebiten.Image
				if item.Enabled {
					check = checkOn
				} else {
					check = checkOff
				}

				/* Draw checkmark */
				op.GeoM.Scale(uiScale*checkScale, uiScale*checkScale)
				op.GeoM.Translate(
					float64(window.scaledSize.X)-(float64(check.Bounds().Dx())*uiScale)-(padding*uiScale),
					float64(item.TextPosY)-(float64(check.Bounds().Dy())*uiScale*checkScale))
				window.cache.DrawImage(check, op)
			}
			d++
		}
	}
}
