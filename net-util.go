package main

import (
	"time"

	"nhooyr.io/websocket"
)

func changeGameMode(newMode MODE, delay time.Duration) {

	/* Skip if the same */
	if newMode == gameMode {
		return
	}

	time.Sleep(delay)
	gameMode = newMode
}

func sendCommand(header CMD, data []byte) bool {
	if localPlayer == nil || localPlayer.context == nil || localPlayer.conn == nil {
		return false
	}

	cmdName := cmdNames[header]
	if cmdName == "" {
		doLog(true, "Sent: 0x%02X", header)
	} else {
		//doLog(true, "Sent: %v", cmdName)
	}

	var err error
	if data == nil {
		err = localPlayer.conn.Write(localPlayer.context, websocket.MessageBinary, []byte{byte(header)})
	} else {
		err = localPlayer.conn.Write(localPlayer.context, websocket.MessageBinary, append([]byte{byte(header)}, data...))
	}
	if err != nil {

		doLog(true, "sendCommand error: %v", err)

		changeGameMode(MODE_BOOT, time.Second)
		connectServer()
		return false
	}

	return true
}
