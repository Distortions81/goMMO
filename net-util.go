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
	gameModeLock.Lock()
	gameMode = newMode
	gameModeLock.Unlock()
}

func sendCommand(header CMD, data []byte) bool {

	localPlayer.plock.Lock()
	defer localPlayer.plock.Unlock()
	if localPlayer.conn == nil {
		return false
	}

	cmdName := cmdNames[header]
	if header != CMD_MOVE {
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

		changeGameMode(MODE_RECONNECT, time.Second)
		chatLines = []chatLineData{}
		chatLinesTop = 0

		chat("Connection lost!")

		go connectServer()

		return false
	}

	return true
}
