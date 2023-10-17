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
	MODE_PLAY_GAME
	MODE_ERROR
)

/* Network commands */
type CMD uint8

const (
	CMD_INIT CMD = iota
	CMD_PINGPONG

	CMD_GETLOBBIES
	CMD_JOINLOBBY
	CMD_CREATELOBBY

	CMD_GODIR

	RECV_LOCALPLAYER
	RECV_LOBBYLIST
	RECV_KEYFRAME
	RECV_PLAYERUPDATE
)
