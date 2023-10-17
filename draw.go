package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

func getFrame(dir int) image.Image {

	dirOff := dirToOffset(dir)

	rect := image.Rectangle{}
	rect.Min.X = (walkframe * charSpriteSize)
	rect.Max.X = (walkframe * charSpriteSize) + charSpriteSize
	rect.Min.Y = dirOff
	rect.Max.Y = charSpriteSize + dirOff

	return testChar.SubImage(rect)

}

/* Ebiten: Draw everything */
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(colorGrass)

	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(128-26+float64(charPos.X), 128-26+float64(charPos.Y))
	op.GeoM.Scale(2, 2)

	screen.DrawImage(getFrame(goDir).(*ebiten.Image), &op)
}
