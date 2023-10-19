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

type Game struct {
}

type playerData struct {
	conn    *websocket.Conn
	context context.Context
	cancel  context.CancelFunc

	lid int

	id   uint32
	Name string
}
