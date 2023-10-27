package main

import (
	"image/color"
	"sync"
	"time"

	"github.com/sasha-s/go-deadlock"
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
	localPlayer     playerData
	playerNames     map[uint32]pNameData
	playerNamesLock deadlock.Mutex

	ourPos  XY
	posLock deadlock.Mutex

	playerList     map[uint32]*playerData
	playerListLock deadlock.Mutex

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
