package main

import (
	"context"

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

type playerData struct {
	conn    *websocket.Conn
	context context.Context
	cancel  context.CancelFunc

	pos       XY
	lastPos   XY
	direction DIR
	walkFrame int
	isWalking bool

	unmark bool

	id uint32
}
