package main

import (
	"math"
)

const (

	/* Updates per second, real update rate is this div 2 */
	gameUPS = 8
	/* Used for perlin noise layers */

	/* For sprite rotation */
	ninetyDeg     = math.Pi / 2
	oneEightyDeg  = math.Pi
	threeSixtyDeg = math.Pi * 2
	//DegToRad      = 6.28319

	/* Files and directories */
	dataDir     = "data/"
	gfxDir      = dataDir + "sprites/"
	testCharDir = "testCharacter/"

	/* WASD speeds */
	moveSpeed = 4.0
	runSpeed  = 16.0

	/* Define world center */
	xyCenter = 32768.0
	xyMax    = xyCenter * 2.0
	xyMin    = 1.0

	/* Directions */
	DIR_NORTH = 0
	DIR_EAST  = 1
	DIR_SOUTH = 2
	DIR_WEST  = 3
	DIR_MAX   = 4
	DIR_ANY   = 255
)
