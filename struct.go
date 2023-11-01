package main

import (
	"context"
	"image/color"
	"sync"
	"time"

	"nhooyr.io/websocket"
)

type XY struct {
	X uint32
	Y uint32
}
type XYs struct {
	X int32
	Y int32
}
type XYf64 struct {
	X float64
	Y float64
}
type XYf32 struct {
	X float32
	Y float32
}

type Game struct {
}

type pNameData struct {
	name string
	id   uint32
}

type playerData struct {
	conn    *websocket.Conn
	context context.Context
	cancel  context.CancelFunc

	health int8

	pos    XY
	spos   XY
	areaid uint16

	lastPos    XY
	lastUpdate time.Time

	direction DIR
	walkFrame int
	isWalking bool

	unmark bool

	id uint32
}

/* Chat line data */
type chatLineData struct {
	text string

	color   color.Color
	bgColor color.Color

	timestamp time.Time
	lifetime  time.Duration
}

type areaData struct {
	arealock sync.RWMutex
	chunks   map[XY]*chunkData
}

type chunkData struct {
	worldObjects []*worldObject
	players      []*playerData
	chunkLock    sync.RWMutex
}

type worldObject struct {
	itemId uint32
	pos    XY

	itemData *sectionItemData
}
