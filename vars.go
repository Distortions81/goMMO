package main

import (
	"image/color"
	"sync"
	"time"
)

var (
	WASMMode bool = false

	dataDirty bool = true
	gDevMode  bool

	/* Game Mode */
	gameMode     = MODE_START
	gameModeLock sync.Mutex

	/* Local player */
	localPlayer playerData
	playerNames map[uint32]pNameData

	goDir DIR

	ourPos    XY
	ourOldPos XY

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
