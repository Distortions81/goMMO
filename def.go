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
)

const (
	/* Directions */
	DIR_NORTH = iota
	DIR_EAST
	DIR_SOUTH
	DIR_WEST
	DIR_NONE
)
