package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

/* Ebiten: Draw everything */
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(colorGrass)

	op := ebiten.DrawImageOptions{}

	//center of screen, center of sprite, charpos
	op.GeoM.Translate(128-26+float64(charPos.X), 128-26+float64(charPos.Y))
	//Upscale
	op.GeoM.Scale(2, 2)

	//Draw sub-image
	screen.DrawImage(getCharFrame(goDir).(*ebiten.Image), &op)
}
