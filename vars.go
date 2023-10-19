package main

import (
	"image/color"
	"time"
)

var (
	WASMMode bool = false

	walkframe   int
	updateCount int
	goDir       int
	isWalking   bool
	charPos     XY = xyCenter

	/* Game Mode */
	gameMode = MODE_START

	/* Local player */
	localPlayer *playerData

	/* Test BG Color */
	colorGrass = color.RGBA{R: 132, G: 145, B: 65}

	/* Networking */
	authSite = "https://gommo.go-game.net/gs"

	/* Ping */
	statusTime    time.Time
	lastRoundTrip time.Duration

	/* Reconnect */
	ReconnectCount     = 0
	RecconnectDelayCap = 30
)
