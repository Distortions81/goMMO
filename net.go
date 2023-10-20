package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"image/color"
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
		buf := fmt.Sprintf("Connect %v failed, retrying in %v ...", ReconnectCount, time.Until(statusTime).Round(time.Second).String())
		doLog(true, buf)
		chat(buf)
		time.Sleep(time.Duration(delay) * time.Second)
	}
	time.Sleep(time.Millisecond * 500)
	changeGameMode(MODE_PLAYING, 0)
}

func doConnect() bool {

	changeGameMode(MODE_CONNECT, 0)
	doLog(true, "Connecting...")
	chat("Connecting to server...")

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

	chat("Connected!")
	time.Sleep(time.Millisecond * 100)
	chatDetailed("Use WASD keys to walk!", color.White, time.Second*30)

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
		inbuf := bytes.NewReader(data)
		cmdName := cmdNames[d]

		/* Log event */
		if cmdName == "" {
			doLog(true, "Received: 0x%02X (%vb)", d, inputLen)
		} else {
			//doLog(true, "Received: %v (%vb)", cmdName, inputLen)
		}

		switch d {
		case CMD_LOGIN:
			binary.Read(inbuf, binary.BigEndian, &localPlayer.id)
			doLog(true, "New local id: %v", localPlayer.id)
			changeGameMode(MODE_PLAYING, 0)
		case CMD_UPDATE:
			var numPlayers uint32
			binary.Read(inbuf, binary.BigEndian, &numPlayers)

			playerListLock.Lock()

			for _, player := range playerList {
				player.unmark = true
			}

			var x uint32
			for x = 0; x < numPlayers; x++ {
				var nid uint32
				binary.Read(inbuf, binary.BigEndian, &nid)
				var nx uint32
				binary.Read(inbuf, binary.BigEndian, &nx)
				var ny uint32
				binary.Read(inbuf, binary.BigEndian, &ny)

				if playerList[nid] == nil {
					playerList[nid] = &playerData{id: nid, pos: XY{X: nx, Y: ny}}
					doLog(false, "Player added: %v", nid)
				} else {
					playerList[nid].lastPos = playerList[nid].pos

					playerList[nid].pos.X = nx
					playerList[nid].pos.Y = ny

					if playerList[nid].lastPos.X != playerList[nid].pos.X || playerList[nid].lastPos.Y != playerList[nid].pos.Y {
						playerList[nid].walkFrame++
						playerList[nid].isWalking = true

					} else {
						playerList[nid].isWalking = false
						playerList[nid].walkFrame = 0
					}
				}
				playerList[nid].unmark = false
			}

			for p, player := range playerList {
				if player.unmark {
					doLog(false, "Player removed: %v", p)
					delete(playerList, p)
				}
			}

			dataDirty = true
			playerListLock.Unlock()
		default:
			doLog(true, "Received invalid: 0x%02X\n", d)
			localPlayer.conn.Close(websocket.StatusNormalClosure, "closed")
			return
		}
	}
}
