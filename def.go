package main

var netProtoVersion uint16 = 3

const (
	gameVersion = "0.0.7"

	/* Files and directories */
	dataDir = "data/"
	gfxDir  = dataDir + "gfx/"

	charSpriteSize    = 52
	cLoadEmbedSprites = true
)

/* Directions */
type DIR uint8

const (
	/* Directions */
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

/* Game modes */
type MODE uint8

const (
	MODE_START MODE = iota
	MODE_BOOT
	MODE_CONNECT
	MODE_RECONNECT
	MODE_CONNECTED
	MODE_LOGIN
	MODE_PLAYING
	MODE_LOGOFF
	MODE_ERROR
)

/* Network commands */
type CMD uint8

const (
	CMD_INIT CMD = iota
	CMD_LOGIN
	CMD_PLAY
	CMD_MOVE
	CMD_UPDATE
	CMD_CHAT
	CMD_SCREENSIZE
	CMD_COMMAND
	CMD_PLAYERNAMES
)

/* Used for debug messages, this could be better */
var cmdNames map[CMD]string

func init() {
	cmdNames = make(map[CMD]string)
	cmdNames[CMD_INIT] = "CMD_INIT"
	cmdNames[CMD_LOGIN] = "CMD_LOGIN"
	cmdNames[CMD_PLAY] = "CMD_PLAY"
	cmdNames[CMD_MOVE] = "CMD_MOVE"
	cmdNames[CMD_UPDATE] = "CMD_UPDATE"
	cmdNames[CMD_CHAT] = "CMD_CHAT"
	cmdNames[CMD_SCREENSIZE] = "CMD_SCREENSIZE"
	cmdNames[CMD_COMMAND] = "CMD_COMMAND"
	cmdNames[CMD_PLAYERNAMES] = "CMD_PLAYERNAMES"
}

const xyHalf = 2147483648

var xyCenter XY = XY{X: xyHalf, Y: xyHalf}
