package main

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"

	"nhooyr.io/websocket"
)

const netReadLimit = 1024 * 1024 * 10 //10mb

func connectServer() {

	changeGameMode(MODE_BOOT, 0)

	if localPlayer != nil {
		if localPlayer.conn != nil {
			localPlayer.cancel()
			localPlayer.conn.Close(websocket.StatusNormalClosure, "Write failed.")
		}

		localPlayer = nil
	}

	for !doConnect() {
		ReconnectCount++
		offset := ReconnectCount
		if offset > RecconnectDelayCap {
			offset = RecconnectDelayCap
		}
		timeFuzz := rand.Int63n(200000000) //2 seconds of random delay
		delay := float64(3 + offset + int(float64(timeFuzz)/100000000.0))
		statusTime = time.Now().Add(time.Duration(delay) * time.Second)

		changeGameMode(MODE_RECONNECT, time.Millisecond*500)
		doLog(true, "Connect %v failed, retrying in %v ...", ReconnectCount, time.Until(statusTime).Round(time.Second).String())
		time.Sleep(time.Duration(delay) * time.Second)
	}
	time.Sleep(time.Millisecond * 500)
	changeGameMode(MODE_PLAYING, 0)
}

func doConnect() bool {

	changeGameMode(MODE_CONNECT, 0)
	doLog(true, "Connecting...")

	ctx, cancel := context.WithCancel(context.Background())

	c, err := platformDial(ctx)

	if err != nil {
		log.Printf("dial failed: %v\n", err)
		cancel()
		return false
	}

	//10MB limit
	c.SetReadLimit(netReadLimit)

	localPlayer = &playerData{conn: c, context: ctx, cancel: cancel, id: 0}
	c.Write(ctx, websocket.MessageBinary, []byte{byte(CMD_INIT)})
	doLog(true, "Connected!")

	changeGameMode(MODE_CONNECTED, 0)
	go readNet()

	return true
}

var gameLock sync.Mutex

func readNet() {

	if localPlayer == nil ||
		localPlayer.conn == nil ||
		localPlayer.context == nil {
		doLog(true, "readNet: Player not initialized.")
		return
	}

	for {

		_, input, err := localPlayer.conn.Read(localPlayer.context)

		/* Read error, reconnect */
		if err != nil {
			doLog(true, "readNet error: %v", err)

			//TODO: Notify player here
			changeGameMode(MODE_BOOT, time.Second)

			connectServer()
			return
		}

		/* Check data length */
		inputLen := len(input)
		if inputLen <= 0 {
			return
		}

		/* Separate command and data*/
		d := CMD(input[0])
		data := input[1:]

		cmdName := cmdNames[d]

		/* Log event */
		if cmdName == "" {
			doLog(true, "Received: 0x%02X (%vb)", d, inputLen)
		} else {
			doLog(true, "Received: %v (%vb)", cmdName, inputLen)
		}

		switch d {
		case CMD_LOGIN:
			cmd_login(data)
		case CMD_UPDATE:
			cmd_update(data)
		default:
			doLog(true, "Received invalid: 0x%02X\n", d)
			localPlayer.conn.Close(websocket.StatusNormalClosure, "closed")
			return
		}
	}
}

func cmd_login(data []byte) {
	localPlayer.id = byteArrayToUint32(data)
	changeGameMode(MODE_PLAYING, 0)
}

func cmd_update(data []byte) {
	doLog(true, "meep")
}
