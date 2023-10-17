package main

import (
	"context"

	"nhooyr.io/websocket"
)

type xy struct {
	X int
	Y int
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
