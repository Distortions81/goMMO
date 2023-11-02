package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"time"

	"nhooyr.io/websocket"
)

const (
	netReadLimit = 1024 * 1024 // 1 MB
	bootSleep    = time.Millisecond * 500
)

func connectServer() {

	changeGameMode(MODE_BOOT, 0)

	for !doConnect() {
		ReconnectCount++
		offset := ReconnectCount

		if offset > RecconnectDelayCap {
			offset = RecconnectDelayCap
		}

		//2 seconds of random delay to spread out load on server reboot
		timeFuzz := rand.Int63n(200000000)
		delay := float64(3 + offset + int(float64(timeFuzz)/100000000.0))
		reconnectTime = time.Now().Add(time.Duration(delay) * time.Second)

		changeGameMode(MODE_RECONNECT, time.Millisecond*500)
		buf := fmt.Sprintf("Connect %v failed, retrying in %v ...", ReconnectCount, time.Until(reconnectTime).Round(time.Second).String())
		doLog(true, buf)
		chat(buf)

		time.Sleep(time.Duration(delay) * time.Second)
	}
	time.Sleep(bootSleep)
	changeGameMode(MODE_PLAYING, 0)
}

func doConnect() bool {

	if playerNames == nil {
		playerNames = make(map[uint32]pNameData)
	}

	changeGameMode(MODE_CONNECT, 0)
	doLog(true, "Connecting...")
	chat("Connecting to server...")

	ctx, cancel := context.WithCancel(context.Background())
	localPlayer.context = ctx
	localPlayer.cancel = cancel

	//Two versions of this, standard and WASM
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

	//Send INIT to server
	go sendCommand(CMD_INIT, outbuf.Bytes())

	doLog(true, "Connected!")
	chat("Connected!")

	changeGameMode(MODE_CONNECTED, 0)

	//Start read loop on new thread and exit
	go readNet()
	return true
}

func playerIdToName(id uint32) string {

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

// Used for motion smoothing
var lastNetUpdate time.Time

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

		drawLock.Lock()

		if d != CMD_UPDATE {
			cmdName := cmdNames[d]

			// Log event
			if cmdName == "" {
				doLog(true, "Received: 0x%02X (%vb)", d, inputLen)
			} else {
				doLog(true, "Received: %v (%vb)", cmdName, inputLen)
			}
		}

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
		case CMD_WORLDDATA:
			fmt.Printf("WorldData %v\n", string(data))
		case CMD_PLAYERNAMES:
			var numNames uint32
			binary.Read(inbuf, binary.LittleEndian, &numNames)

			if numNames == 0 {
				continue
			}

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

		case CMD_UPDATE:

			var numPlayers uint16
			binary.Read(inbuf, binary.LittleEndian, &numPlayers)

			//Mark players, so we know if we can remove them
			for _, player := range playerList {
				player.unmark = true
			}

			//Mark time here, used for motion smoothing
			lastNetUpdate = time.Now()
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
							oldLocalPlayerPos = localPlayerPos
							localPlayerPos.X = nx
							localPlayerPos.Y = ny
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

			//Delete players that are no longer found
			for p, player := range playerList {
				if player.unmark {
					delete(playerList, p)
				}
			}

			//Read dynamic world objects
			var numObj uint16
			err = binary.Read(inbuf, binary.LittleEndian, &numObj)
			if err != nil {
				doLog(true, "%v", err.Error())
			}

			if numObj > 0 {
				//fmt.Printf("numObj: %v\n", numObj)

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

					object := &worldObject{itemId: itemId, pos: pos, itemData: spritelist[itemId]}
					wObjList = append(wObjList, object)

				}
			}

			//Mark that there is new data to be rendered, used if motion smoothing is OFF.
			dataDirty = true

		default:
			doLog(true, "Received invalid: 0x%02X\n", d)
			localPlayer.conn.Close(websocket.StatusNormalClosure, "closed")
			changeGameMode(MODE_ERROR, 0)
		}

		drawLock.Unlock()
	}
}
