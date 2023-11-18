package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"image/color"
	"io"
	"log"
	"math/rand"
	"time"

	"nhooyr.io/websocket"
)

const (
	netReadLimit = 1024 * 500 // 500kb
	bootSleep    = time.Millisecond * 500
)

var netTick uint64

func connectServer() {

	if !devMode {
		time.Sleep(time.Second * 3)
	}
	changeGameMode(MODE_Connect, 0)

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

		changeGameMode(MODE_Reconnect, 0)
		buf := fmt.Sprintf("Connect %v failed, retrying in %v ...", ReconnectCount, time.Until(reconnectTime).Round(time.Second).String())
		doLog(true, buf)
		chat(buf)

		if devMode {
			time.Sleep(time.Second)
		} else {
			time.Sleep(time.Duration(delay) * time.Second)
		}
	}
	time.Sleep(bootSleep)

}

func doConnect() bool {

	if playerNames == nil {
		playerNames = make(map[uint32]pNameData)
	}

	if gameMode != MODE_Connect && gameMode != MODE_Reconnect {
		return true
	}

	changeGameMode(MODE_Connect, 0)
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

	//Send INIT to server, reset everything
	goingDirection = DIR_NONE
	playerNames = make(map[uint32]pNameData)
	localPlayerPos = worldCenter
	oldLocalPlayerPos = worldCenter
	wObjList = []*worldObject{}
	playerList = make(map[uint32]*playerData)
	playerMode = PMODE_PASSIVE
	creatureList = make(map[uint32]*playerData)
	drawToolbar(false, false, maxItemType)

	go sendCommand(CMD_Init, outbuf.Bytes())

	doLog(true, "Connected!")
	chat("Connected!")

	changeGameMode(MODE_Connected, 0)

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

		// Read error, reconnect
		if err != nil {
			doLog(true, "readNet error: %v", err)

			if gameMode == MODE_Playing {
				chat("Connection lost!")
				connectServer()
			}
			return
		}

		// Check data length
		inputLen := len(input)
		if inputLen <= 0 {
			continue
		}

		// Separate command and data
		d := CMD(input[0])
		data := input[1:]
		inbuf := bytes.NewReader(data)

		drawLock.Lock()

		if d != CMD_WorldUpdate {
			cmdName := cmdNames[d]

			// Log event
			if cmdName == "" {
				doLog(true, "Received: 0x%x (%vb)", d, inputLen)
			} else {
				doLog(true, "Received: %v (%vb)", cmdName, inputLen)
			}
		}

		switch d {
		case CMD_Init:
			chatDetailed("Server rejected connection: Client version not supported.",
				color.White, time.Hour*72)
			if WASMMode {
				chatDetailed("Please refresh your browser!",
					color.White, time.Hour*72)
			}
			changeGameMode(MODE_Error, 0)

		case CMD_Login:
			binary.Read(inbuf, binary.LittleEndian, &localPlayer.id)
			binary.Read(inbuf, binary.LittleEndian, &localPlayer.areaid)
			doLog(true, "New local id: %v, area: %v", localPlayer.id, localPlayer.areaid)
			changeGameMode(MODE_Playing, 0)

		case CMD_Chat:
			chat(string(data))

		case CMD_Command:
			chat("> " + string(data))
		case CMD_WorldData:
			fmt.Printf("WorldData %v\n", string(data))
		case CMD_PlayerNamesComp:

			//Decompress
			buf, err := io.ReadAll(inbuf)
			var decompBytes []byte
			if err != nil {
				return
			}
			decompBytes = UncompressZip(buf)

			newbuf := bytes.NewBuffer(decompBytes)

			//Process as normal
			var numNames uint32
			binary.Read(newbuf, binary.LittleEndian, &numNames)

			if numNames == 0 {
				continue
			}

			for x := 0; x < int(numNames); x++ {
				var id uint32
				binary.Read(newbuf, binary.LittleEndian, &id)
				var nameLen uint16
				binary.Read(newbuf, binary.LittleEndian, &nameLen)

				var name string
				for y := 0; y < int(nameLen); y++ {
					var nameRune rune
					binary.Read(newbuf, binary.LittleEndian, &nameRune)
					name += string(nameRune)
				}

				playerNames[id] = pNameData{name: name, id: id}
			}

		case CMD_WorldUpdate:
			netTick++

			var numPlayers uint8
			binary.Read(inbuf, binary.LittleEndian, &numPlayers)

			//Mark time here, used for motion smoothing
			lastNetUpdate = time.Now()
			if numPlayers > 0 {

				//Mark players, so we know if we can remove them
				for _, player := range playerList {
					player.unmark++
				}

				var x uint8
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

					var dir DIR
					err = binary.Read(inbuf, binary.LittleEndian, &dir)
					if err != nil {
						doLog(true, "%v", err.Error())
						break
					}

					var health int16
					err = binary.Read(inbuf, binary.LittleEndian, &health)
					if err != nil {
						doLog(true, "%v", err.Error())
						break
					}

					var effects EFF
					err = binary.Read(inbuf, binary.LittleEndian, &effects)
					if err != nil {
						doLog(true, "%v", err.Error())
						break
					}
					if playerList[nid] == nil {
						playerList[nid] = &playerData{id: nid, pos: XY{X: nx, Y: ny}, direction: DIR_S, effects: effects}

						/* Unmark, used to detect if no longer needed */
						playerList[nid].unmark = 0
					} else {

						/* Unmark, used to detect if no longer needed */
						playerList[nid].unmark = 0

						// Update local player pos
						if localPlayer.id == nid {
							oldLocalPlayerPos = localPlayerPos
							localPlayerPos.X = nx
							localPlayerPos.Y = ny
						}

						playerList[nid].lastPos = playerList[nid].pos
						playerList[nid].pos.X = nx
						playerList[nid].pos.Y = ny
						playerList[nid].direction = dir
						playerList[nid].spos.X = nx
						playerList[nid].spos.Y = ny
						playerList[nid].health = health
						playerList[nid].effects = effects

						/* Walk animations */
						if playerList[nid].lastPos.X != nx ||
							playerList[nid].lastPos.Y != ny {
							playerList[nid].walkFrame++
							playerList[nid].isWalking = true

						} else {
							playerList[nid].isWalking = false
							playerList[nid].walkFrame = 0
						}
					}
				}
				//Delete players that are no longer found
				for p, player := range playerList {
					if player.unmark > 15 {
						delete(playerList, p)
					}
				}
			}

			var numObj uint8
			binary.Read(inbuf, binary.LittleEndian, &numObj)

			if numObj > 0 {
				var x uint8
				for x = 0; x < numObj; x++ {
					var sid uint8
					err := binary.Read(inbuf, binary.LittleEndian, &sid)
					if err != nil {
						doLog(true, "%v", err.Error())
						break
					}
					var oid uint8
					err = binary.Read(inbuf, binary.LittleEndian, &oid)
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
					objData := itemTypesList[sid].items[oid]
					wObjList = append(wObjList,
						&worldObject{
							itemId:   IID{Section: sid, Num: oid},
							pos:      XY{X: nx, Y: ny},
							itemData: objData})
				}
			}

			var numCreatures uint8
			binary.Read(inbuf, binary.LittleEndian, &numCreatures)

			if numCreatures > 0 {
				//Mark players, so we know if we can remove them
				for _, cre := range creatureList {
					cre.unmark++
				}

				var x uint8
				for x = 0; x < numCreatures; x++ {
					var uid uint32
					err := binary.Read(inbuf, binary.LittleEndian, &uid)
					if err != nil {
						doLog(true, "%v", err.Error())
						break
					}
					var sec uint8
					err = binary.Read(inbuf, binary.LittleEndian, &sec)
					if err != nil {
						doLog(true, "%v", err.Error())
						break
					}
					var cid uint8
					err = binary.Read(inbuf, binary.LittleEndian, &cid)
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
					var dir DIR
					err = binary.Read(inbuf, binary.LittleEndian, &dir)
					if err != nil {
						doLog(true, "%v", err.Error())
						break
					}
					var health int16
					err = binary.Read(inbuf, binary.LittleEndian, &health)
					if err != nil {
						doLog(true, "%v", err.Error())
						break
					}

					var effects EFF
					err = binary.Read(inbuf, binary.LittleEndian, &effects)
					if err != nil {
						doLog(true, "%v", err.Error())
						break
					}
					if creatureList[uid] == nil {
						creData := creatureData{id: IID{Section: sec, Num: cid, UID: uid}, target: nil}
						newCreature := &playerData{
							creature: &creData, pos: XY{X: nx, Y: ny}, spos: XY{X: nx, Y: ny}, health: health,
							direction: DIR_S, effects: effects}
						creatureList[uid] = newCreature
					} else {

						if creatureList[uid].lastPos.X != nx ||
							creatureList[uid].lastPos.Y != ny {
							creatureList[uid].lastPos = creatureList[uid].pos
							creatureList[uid].walkFrame++
							creatureList[uid].isWalking = true
						} else {
							creatureList[uid].isWalking = false
							creatureList[uid].walkFrame = 0
						}

						creatureList[uid].unmark = 0

						creatureList[uid].pos.X = nx
						creatureList[uid].pos.Y = ny
						creatureList[uid].direction = dir
						creatureList[uid].spos.X = nx
						creatureList[uid].spos.Y = ny
						creatureList[uid].health = health
						creatureList[uid].effects = effects

					}

				}
				//Delete players that are no longer found
				for p, cre := range creatureList {
					if cre.unmark > 15 {
						delete(creatureList, p)
					}
				}
			}

			//Mark that there is new data to be rendered, used if motion smoothing is OFF.
			dataDirty = true

		default:
			doLog(true, "Received invalid: %x\n", d)
			localPlayer.conn.Close(websocket.StatusNormalClosure, "closed")
			changeGameMode(MODE_Error, 0)
		}

		drawLock.Unlock()
	}
}
