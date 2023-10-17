package main

import "image/color"

var (
	WASMMode bool = false

	walkframe   int
	updateCount int
	goDir       int
	isWalking   bool
	charPos     xy

	colorGrass = color.RGBA{R: 132, G: 145, B: 65}
)
