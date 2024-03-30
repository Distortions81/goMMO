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

	// Create a new RGBA image for the outlined image
	outlinedRect := inputImg.Bounds()
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
				if a != 0 {
					outlinedContext.DrawRectangle(float64(x), float64(y), 1, 1)
					outlinedContext.Stroke()
				}
			}
		}

	}

	// Create a new RGBA image with the same bounds as the original image
	bounds := inputImg.Bounds()
	newImg := image.NewRGBA(bounds)

	// Set the alpha channel transparency for each pixel
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := outlinedRGBA.At(x, y).RGBA()
			if a == 0 {
				continue
			}
			// Adjust alpha while preserving the original transparency
			newColor := color.NRGBA{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
				A: uint8(int(float32(a)*0.7) >> 8),
			}
			newImg.Set(x, y, newColor)
		}
	}

	// Composite the original image on top of the outlined image
	draw.Draw(newImg, outlinedRect, inputImg, image.Point{}, draw.Over)

	return ebiten.NewImageFromImage(newImg)
}
