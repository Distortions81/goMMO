package main

import (
	"image/color"
	"sync"
	"time"
)

var (
	WASMMode bool = false

	updateCount int
	goDir       DIR
	dataDirty   bool  = true
	curCharPos  XYf64 = XYf64{X: xyHalf, Y: xyHalf}
	lastCharPos XYf64 = XYf64{X: xyHalf, Y: xyHalf}

	/* Game Mode */
	gameMode = MODE_START

	/* Local player */
	localPlayer *playerData
	ourPos      XY
	posLock     sync.Mutex

	playerList     map[uint32]*playerData
	playerListLock sync.Mutex

	/* Name BG Color */
	colorNameBG = color.RGBA{R: 32, G: 32, B: 32, A: 160}

	/* Networking */
	authSite = "https://gommo.go-game.net/gs"

	/* Ping */
	statusTime time.Time

	/* Reconnect */
	ReconnectCount     = 0
	RecconnectDelayCap = 30
)
