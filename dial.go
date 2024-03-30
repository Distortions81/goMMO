//go:build !js && !wasm

package main

import (
	"context"

	"nhooyr.io/websocket"
)

func platformDial(ctx context.Context) (c *websocket.Conn, err error) {
	opts := &websocket.DialOptions{HTTPClient: client}
	c, _, err = websocket.Dial(ctx, authSite, opts)
	return c, err
}
