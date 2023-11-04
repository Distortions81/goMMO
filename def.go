package main

import "github.com/hajimehoshi/ebiten/v2"

var netProtoVersion uint16 = 11

const FrameSpeedNS = 66666666

const (
	gameVersion = "0.0.12"

	// Files and directories
	dataDir = "data/"
	gfxDir  = dataDir + "gfx/"
	txtDir  = dataDir + "txt/"

	playerSpriteSize = 52
	loadEmbedSprites = true
	chunkDiv         = 128
)

// Game modes
type MODE uint8

const (
	MODE_Start MODE = iota
	MODE_Boot
	MODE_Connect
	MODE_Reconnect
	MODE_Connected
	MODE_Login
	MODE_Playing
	MODE_Logoff
	MODE_Error
)

var playerMode PMode

// Used for debug messages, this could be better
var modeNames map[PMode]*playerModeData

func loadPlayerModes() {
	modeNames = make(map[PMode]*playerModeData)
	modeNames[PMODE_PASSIVE] = &playerModeData{printName: "Passive", name: "PMODE_PASSIVE", imgName: "passive"}
	modeNames[PMODE_ATTACK] = &playerModeData{printName: "Attack", name: "PMODE_ATTACK", imgName: "attack"}
	modeNames[PMODE_HEAL] = &playerModeData{printName: "Heal", name: "PMODE_HEAL", imgName: "heal"}

	//Load images
	for m, mode := range modeNames {
		modeNames[m].image = findItemImage("player-modes", mode.imgName)
	}
}

type playerModeData struct {
	printName string
	name      string
	imgName   string
	image     *ebiten.Image
}

type PMode uint8

// Directions
type DIR uint8

const (
	// Directions
	DIR_N DIR = iota
	DIR_NE
	DIR_E
	DIR_SE
	DIR_S
	DIR_SW
	DIR_W
	DIR_NW
	DIR_NONE
)

const (
	// Directions
	PMODE_PASSIVE PMode = iota
	PMODE_ATTACK
	PMODE_HEAL
)

// Network commands
type CMD uint8

const (
	CMD_Init CMD = iota
	CMD_Login
	CMD_Play
	CMD_Move
	CMD_WorldUpdate
	CMD_Chat
	CMD_Command
	CMD_PlayerMode

	CMD_WorldData
	CMD_PlayerNames
	CMD_EditPlaceItem
	CMD_EditDeleteItem
)

// Used for debug messages, this could be better
var cmdNames map[CMD]string

func init() {
	cmdNames = make(map[CMD]string)
	cmdNames[CMD_Init] = "CMD_Init"
	cmdNames[CMD_Login] = "CMD_Login"
	cmdNames[CMD_Play] = "CMD_Play"
	cmdNames[CMD_Move] = "CMD_Move"
	cmdNames[CMD_WorldUpdate] = "CMD_WorldUpdate"
	cmdNames[CMD_Chat] = "CMD_Chat"
	cmdNames[CMD_Command] = "CMD_Command"
	cmdNames[CMD_PlayerMode] = "CMD_PlayerMode"

	cmdNames[CMD_WorldData] = "CMD_WorldData"
	cmdNames[CMD_PlayerNames] = "CMD_PlayerNames"
	cmdNames[CMD_EditPlaceItem] = "CMD_EditPlaceItem"
	cmdNames[CMD_EditDeleteItem] = "CMD_EditDeleteItem"
}

const xyCenter = 2147483648
const xyMax = xyCenter * 2

var worldCenter XY = XY{X: xyCenter, Y: xyCenter}
