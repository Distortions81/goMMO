package main

const (
	/* Files and directories */
	dataDir     = "data/"
	gfxDir      = dataDir + "sprites/"
	testAreaDir = "testArea"
	testCharDir = "chars/"

	charSpriteSize    = 52
	cLoadEmbedSprites = true
)

const (
	/* Directions */
	DIR_S = iota
	DIR_SW
	DIR_W
	DIR_NW
	DIR_N
	DIR_NE
	DIR_E
	DIR_SE
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
)

/* Used for debug messages, this could be better */
var cmdNames map[CMD]string

func init() {
	cmdNames = make(map[CMD]string)
	cmdNames[CMD_INIT] = "CMD_INIT"
	cmdNames[CMD_LOGIN] = "CMD_LOGIN"
	cmdNames[CMD_PLAY] = "CMD_PLAY"
	cmdNames[CMD_MOVE] = "CMD_MOVE"
}

const xyHalf = 2147483648

var xyCenter XY = XY{X: xyHalf, Y: xyHalf}
