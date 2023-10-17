package main

const (
	/* Files and directories */
	dataDir     = "data/"
	gfxDir      = dataDir + "sprites/"
	testAreaDir = "testArea"
	testCharDir = "chars/"

	charSpriteSize    = 52
	cLoadEmbedSprites = true
)

const (
	/* Directions */
	DIR_S = iota
	DIR_SW
	DIR_W
	DIR_NW
	DIR_N
	DIR_NE
	DIR_E
	DIR_SE
	DIR_NONE
)
