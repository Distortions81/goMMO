package main

import (
	"image/color"
	"sync"
	"time"
)

var (
	WASMMode bool = false //web assembly mode

	dataDirty bool = true //If new network data has been rendered or not
	gDevMode  bool        //-dev argument

	gameMode     = MODE_Start //Login, playing, reconnect, etc
	gameModeLock sync.Mutex

	localPlayer playerData           //Our local player's data (id, etc)
	playerNames map[uint32]pNameData //Player ID to name map

	goingDirection DIR //Direction we are walking

	localPlayerPos    XY //Our position from server
	oldLocalPlayerPos XY

	playerList map[uint32]*playerData //Players from server
	drawLock   sync.Mutex

	wObjList []*worldObject //World object list

	//Player name BG color
	colorNameBG = color.RGBA{R: 32, G: 32, B: 32, A: 160}

	//Server URL
	authSite = "https://gommo.go-game.net/gs"

	//Reconnect throttle
	reconnectTime time.Time

	//Reconnect count and cap
	ReconnectCount     = 0
	RecconnectDelayCap = 30

	windowDebugMode        = false
	helpText        string = ""

	vSync     bool = true
	debugLine bool = false
)
