package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

/* Ebiten: Draw everything */
func (g *Game) Draw(screen *ebiten.Image) {

	if gameMode == MODE_PLAYING {
		screen.Fill(colorGrass)

		playerListLock.Lock()
		defer playerListLock.Unlock()

		//center of screen, center of sprite, charpos
		for _, player := range playerList {
			op := ebiten.DrawImageOptions{}
			convPos := convPos(player.pos)
			op.GeoM.Translate(quarterWindowStartX-26+float64(convPos.X), quarterWindowStartY-26+float64(convPos.Y))
			//Upscale
			op.GeoM.Scale(2, 2)

			//Draw sub-image
			screen.DrawImage(getCharFrame(DIR_S).(*ebiten.Image), &op)
		}
	} else {
		ebitenutil.DebugPrint(screen, "Connecting.")
	}

}
