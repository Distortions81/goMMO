package main

import (
	"image/color"
	"sync"
	"time"
)

var (
	//Web assembly mode
	WASMMode bool = false

	//If new network data has been rendered or not
	dataDirty bool = true
	//-dev argument
	devMode bool

	//Login, playing, reconnect, etc
	gameMode     = MODE_Start
	gameModeLock sync.Mutex

	//Our local player's data (id, etc)
	localPlayer playerData
	//Player ID to name map
	playerNames map[uint32]*pNameData
	//Direction we are walking
	goingDirection DIR

	//Our position from server
	localPlayerPos    XY
	oldLocalPlayerPos XY

	//Players from server
	playerList   map[uint32]*playerData
	creatureList map[uint32]*playerData
	drawLock     sync.Mutex

	//World object list
	wObjList []*worldObject

	//Player name BG color
	colorNameBG = color.RGBA{R: 16, G: 16, B: 16, A: 128}

	//Server URL
	authSite = "https://gommo.go-game.net/gs"

	//Reconnect throttle
	reconnectTime time.Time

	//Reconnect count and cap
	ReconnectCount     = 0
	RecconnectDelayCap = 30

	//Useful while designing or debugging UI
	windowDebugMode = false

	//Help string
	helpText string = ""

	//Settings
	vSync       bool = true
	smallMode   bool = false
	debugLine   bool = false
	greenScreen bool = false
)
