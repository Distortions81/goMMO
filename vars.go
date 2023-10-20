package main

import (
	"image/color"
	"sync"
	"time"
)

var (
	WASMMode bool = false

	updateCount  int
	goDir        DIR
	dataDirty    bool  = true
	localCharPos XYf64 = XYf64{X: xyHalf, Y: xyHalf}

	/* Game Mode */
	gameMode = MODE_START

	/* Local player */
	localPlayer *playerData

	playerList     map[uint32]*playerData
	playerListLock sync.Mutex

	/* Test BG Color */
	colorGrass = color.RGBA{R: 132, G: 145, B: 65}
	/* Name BG Color */
	colorNameBG = color.RGBA{R: 32, G: 32, B: 32, A: 128}

	/* Networking */
	authSite = "https://gommo.go-game.net/gs"

	/* Ping */
	statusTime time.Time

	/* Reconnect */
	ReconnectCount     = 0
	RecconnectDelayCap = 30
)
