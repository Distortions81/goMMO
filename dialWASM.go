//go:build js && wasm

package main

import (
	"context"

	"nhooyr.io/websocket"
)

func platformDial(ctx context.Context) (c *websocket.Conn, err error) {
	c, _, err = websocket.Dial(ctx, authSite, nil)
	return c, err
}
