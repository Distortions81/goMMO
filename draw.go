package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
)

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

			op.GeoM.Translate(quarterWindowStartX-26+float64(convPos.X), quarterWindowStartY-26+float64(convPos.Y))
			//Upscale
			op.GeoM.Scale(2, 2)

			//Draw sub-image
			screen.DrawImage(getCharFrame(player).(*ebiten.Image), &op)
			pname := fmt.Sprintf("Player-%v", player.id)
			pos := XYf32{X: float32(quarterWindowStartX+convPos.X) * 2.0,
				Y: float32(quarterWindowStartY+convPos.Y)*2.0 + 40}
			drawText(pname, toolTipFont, color.White, colorNameBG,
				pos, 2, screen, false, false, true)
		}
		buf := fmt.Sprintf("%3.0f FPS, WASD to move.", ebiten.ActualFPS())
		ebitenutil.DebugPrint(screen, buf)

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
