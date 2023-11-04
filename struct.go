package main

import (
	"context"
	"image/color"
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

type EFF uint8

const (
	EFFECT_NONE EFF = iota
	EFFECT_HEAL
	EFFECT_ATTACK
)

type pNameData struct {
	name string
	id   uint32
}

type playerData struct {
	conn    *websocket.Conn
	context context.Context
	cancel  context.CancelFunc

	health int16

	pos    XY
	spos   XY
	areaid uint16
	effect EFF

	lastPos XY

	direction   DIR
	walkFrame   int
	attackFrame int
	isWalking   bool

	unmark bool

	id uint32
}

// Chat line data
type chatLineData struct {
	text string

	color   color.Color
	bgColor color.Color

	timestamp time.Time
	lifetime  time.Duration
}

type worldObject struct {
	itemId IID
	pos    XY

	itemData *sectionItemData
}
