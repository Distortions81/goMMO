package main

import (
	"sync"

	"nhooyr.io/websocket"
)

var writeLock sync.Mutex

func sendCommand(header CMD, data []byte) bool {

	writeLock.Lock()
	defer writeLock.Unlock()

	if localPlayer.conn == nil {
		return false
	}

	cmdName := cmdNames[header]
	if header != CMD_Move {
		if cmdName == "" {
			doLog(true, "Sent: %x", header)
		} else {
			doLog(true, "Sent: %v", cmdName)
		}
	}

	var err error
	if data == nil {
		err = localPlayer.conn.Write(localPlayer.context, websocket.MessageBinary, []byte{byte(header)})
	} else {
		err = localPlayer.conn.Write(localPlayer.context, websocket.MessageBinary, append([]byte{byte(header)}, data...))
	}
	if err != nil {

		doLog(true, "sendCommand error: %v", err)

		if gameMode == MODE_Playing {
			chat("Connection lost!")
			connectServer()
		}
		return false
	}

	return true
}
