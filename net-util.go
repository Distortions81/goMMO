package main

import (
	"sync"
	"time"

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
	if header != 0 {
		if cmdName == "" {
			doLog(true, "Sent: 0x%02X", header)
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

		changeGameMode(MODE_Reconnect, time.Second)
		chatLines = []chatLineData{}
		chatLinesTop = 0

		chat("Connection lost!")

		return false
	}

	return true
}
