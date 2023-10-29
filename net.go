package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"nhooyr.io/websocket"
)

const netReadLimit = 1024 * 1000

func connectServer() {

	changeGameMode(MODE_BOOT, 500*time.Millisecond)

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
	localPlayer.plock.Lock()
	defer localPlayer.plock.Unlock()

	playerNames = make(map[uint32]pNameData)

	changeGameMode(MODE_CONNECT, 0)
	doLog(true, "Connecting...")

	netCount = 0
	chat("Connecting to server...")

	ctx, cancel := context.WithCancel(context.Background())
	localPlayer.context = ctx
	localPlayer.cancel = cancel

	conn, err := platformDial(localPlayer.context)

	if err != nil {
		log.Printf("dial failed: %v\n", err)
		localPlayer.cancel()
		return false
	}
	localPlayer.conn = conn

	localPlayer.conn.SetReadLimit(netReadLimit)

	var buf []byte
	outbuf := bytes.NewBuffer(buf)

	binary.Write(outbuf, binary.LittleEndian, &netProtoVersion)

	go sendCommand(CMD_INIT, outbuf.Bytes())
	doLog(true, "Connected!")

	chat("Connected!")
	time.Sleep(time.Millisecond * 100)
	chatDetailed("Use WASD keys to walk!", color.White, time.Second*30)
	chatDetailed("Press [RETURN] to open chat bar", color.White, time.Second*30)
	chatDetailed("Press ` to open command bar", color.White, time.Second*30)

	changeGameMode(MODE_CONNECTED, 0)
	go readNet()

	return true
}

func getName(id uint32) string {
	playerNamesLock.Lock()
	defer playerNamesLock.Unlock()

	for _, pname := range playerNames {
		if pname.id == id {
			if pname.name == "" {
				continue
			}
			return pname.name
		}
	}

	return ""
}

var netCount int

func readNet() {

	for {

		_, input, err := localPlayer.conn.Read(localPlayer.context)

		/* Read error, reconnect */
		if err != nil {
			doLog(true, "readNet error: %v", err)

			//TODO: Notify player here
			changeGameMode(MODE_RECONNECT, time.Second)

			chatLines = []chatLineData{}
			chatLinesTop = 0
			chat("Connection lost!")
			connectServer()
			return
		}

		/* Check data length */
		inputLen := len(input)
		if inputLen <= 0 {
			continue
		}

		/* Separate command and data*/
		d := CMD(input[0])
		data := input[1:]
		inbuf := bytes.NewReader(data)

		/*
			cmdName := cmdNames[d]

				// Log event
				if cmdName == "" {
					doLog(true, "Received: 0x%02X (%vb)", d, inputLen)
				} else {
					doLog(true, "Received: %v (%vb)", cmdName, inputLen)
				}
		*/

		switch d {
		case CMD_INIT:
			chat("Server rejected connection: invalid version.")
			changeGameMode(MODE_ERROR, 0)

		case CMD_LOGIN:
			binary.Read(inbuf, binary.LittleEndian, &localPlayer.id)
			binary.Read(inbuf, binary.LittleEndian, &localPlayer.areaid)
			doLog(true, "New local id: %v, area: %v", localPlayer.id, localPlayer.areaid)
			changeGameMode(MODE_PLAYING, 0)

		case CMD_CHAT:
			chat(string(data))

		case CMD_COMMAND:
			chat("> " + string(data))
		case CMD_PLAYERNAMES:
			var numNames uint32
			binary.Read(inbuf, binary.LittleEndian, &numNames)

			if numNames == 0 {
				continue
			}

			playerNamesLock.Lock()

			for x := 0; x < int(numNames); x++ {
				var id uint32
				binary.Read(inbuf, binary.LittleEndian, &id)
				var nameLen uint16
				binary.Read(inbuf, binary.LittleEndian, &nameLen)

				var name string
				for y := 0; y < int(nameLen); y++ {
					var nameRune rune
					binary.Read(inbuf, binary.LittleEndian, &nameRune)
					name += string(nameRune)
				}

				playerNames[id] = pNameData{name: name, id: id}
			}
			playerNamesLock.Unlock()

		case CMD_UPDATE:

			var numPlayers uint16
			binary.Read(inbuf, binary.LittleEndian, &numPlayers)

			playerListLock.Lock()

			for _, player := range playerList {
				player.unmark = true
			}

			if numPlayers > 0 {

				var x uint16
				for x = 0; x < numPlayers; x++ {
					var nid uint32
					err := binary.Read(inbuf, binary.LittleEndian, &nid)
					if err != nil {
						doLog(true, "%v", err.Error())
						break
					}
					var nx uint32
					err = binary.Read(inbuf, binary.LittleEndian, &nx)
					if err != nil {
						doLog(true, "%v", err.Error())
						break
					}
					var ny uint32
					err = binary.Read(inbuf, binary.LittleEndian, &ny)
					if err != nil {
						doLog(true, "%v", err.Error())
						break
					}

					//Eventually move me to an event
					var health int8
					err = binary.Read(inbuf, binary.LittleEndian, &health)
					if err != nil {
						doLog(true, "%v", err.Error())
						break
					}
					if playerList[nid] == nil {
						playerList[nid] = &playerData{id: nid, pos: XY{X: nx, Y: ny}, direction: DIR_S}
					} else {
						/* Update local player pos */
						if localPlayer.id == nid {
							posLock.Lock()
							ourPos.X = nx
							ourPos.Y = ny
							posLock.Unlock()
						}

						playerList[nid].lastPos = playerList[nid].pos

						playerList[nid].pos.X = nx
						playerList[nid].pos.Y = ny

						playerList[nid].health = health

						if playerList[nid].lastPos.X != playerList[nid].pos.X ||
							playerList[nid].lastPos.Y != playerList[nid].pos.Y {
							playerList[nid].walkFrame++
							playerList[nid].isWalking = true

						} else {
							playerList[nid].isWalking = false
							playerList[nid].walkFrame = 0
						}
					}
					playerList[nid].unmark = false
				}
			}

			for p, player := range playerList {
				if player.unmark {
					delete(playerList, p)
				}
			}

			var numObj uint16
			err = binary.Read(inbuf, binary.LittleEndian, &numObj)
			if err != nil {
				doLog(true, "%v", err.Error())
			}

			if numObj > 0 {
				//fmt.Printf("numObj: %v\n", numObj)

				wObjLock.Lock()
				wObjList = []*worldObject{}
				for x := 0; x < int(numObj); x++ {
					var itemId, posx, posy uint32
					err = binary.Read(inbuf, binary.LittleEndian, &itemId)
					if err != nil {
						doLog(true, "%v", err.Error())
						break
					}
					err = binary.Read(inbuf, binary.LittleEndian, &posx)
					if err != nil {
						doLog(true, "%v", err.Error())
						break
					}
					err = binary.Read(inbuf, binary.LittleEndian, &posy)
					if err != nil {
						doLog(true, "%v", err.Error())
						break
					}

					pos := XY{X: posx, Y: posy}
					object := &worldObject{itemId: itemId, pos: pos}
					wObjList = append(wObjList, object)
				}
				wObjLock.Unlock()
			}

			netCount++
			dataDirty = true
			playerListLock.Unlock()

		default:
			doLog(true, "Received invalid: 0x%02X\n", d)
			localPlayer.conn.Close(websocket.StatusNormalClosure, "closed")
			changeGameMode(MODE_ERROR, 0)
			return
		}
	}
}
