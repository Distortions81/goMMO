package main

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
)

type colorPack struct {
	frames    []outlineColors
	numFrames int
}

type outlineColors struct {
	colors       []color.RGBA
	outlineWidth []float64
}

func makeOutlines(inputImg image.Image, outlineColors []color.RGBA) *ebiten.Image {

	if inputImg == nil {
		return nil
	}

	// Create a new image context for the original image
	originalContext := gg.NewContextForImage(inputImg)

	// Create a new RGBA image for the outlined image
	outlinedRect := originalContext.Image().Bounds()
	outlinedRGBA := image.NewRGBA(outlinedRect)

	// Set the outline width
	oWidth := 1.25
	outlineWidth := float64(len(outlineColors)+1) * oWidth

	// Draw multiple outlines based on alpha channel
	for _, outlineColor := range outlineColors {
		outlinedContext := gg.NewContextForRGBA(outlinedRGBA)
		outlineWidth -= oWidth
		outlinedContext.SetLineWidth(float64(outlineWidth))
		outlinedContext.SetColor(outlineColor)

		for y := outlinedRect.Min.Y; y < outlinedRect.Max.Y; y++ {
			for x := outlinedRect.Min.X; x < outlinedRect.Max.X; x++ {
				_, _, _, a := inputImg.At(x, y).RGBA()
				if a > 0 {
					outlinedContext.DrawRectangle(float64(x), float64(y), 1, 1)
					outlinedContext.Stroke()
				}
			}
		}

	}

	// Composite the original image on top of the outlined image
	draw.Draw(outlinedRGBA, outlinedRect, originalContext.Image(), image.Point{}, draw.Over)

	return ebiten.NewImageFromImage(outlinedRGBA)
}
