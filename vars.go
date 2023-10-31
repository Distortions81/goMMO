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

	gDevMode bool

	/* Game Mode */
	gameMode     = MODE_START
	gameModeLock sync.Mutex

	/* Local player */
	localPlayer playerData
	playerNames map[uint32]pNameData

	ourPos     XY
	ourOldPos  XY
	ourPosLast time.Time

	playerList map[uint32]*playerData
	drawLock   sync.Mutex

	wObjList []*worldObject

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
