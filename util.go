package main

import (
	"bytes"
	"compress/zlib"
	"image"
	"image/color"
	"io"
	"math"
	"strings"
	"time"
)

const healMilliDiv = 50

func changeGameMode(newMode MODE) {
	defer reportPanic("changeGameMode")

	gameModeLock.Lock()
	defer gameModeLock.Unlock()

	// Skip if the same
	if newMode == gameMode {
		return
	}

	gameMode = newMode
}

func playerDirSpriteOffset(dir DIR) int {
	defer reportPanic("playerDirSpriteOffset")

	switch dir {
	case DIR_S:
		return playerSpriteSize * 0
	case DIR_SE:
		return playerSpriteSize * 1
	case DIR_E:
		return playerSpriteSize * 2
	case DIR_NE:
		return playerSpriteSize * 3
	case DIR_N:
		return playerSpriteSize * 4
	case DIR_NW:
		return playerSpriteSize * 5
	case DIR_W:
		return playerSpriteSize * 6
	case DIR_SW:
		return playerSpriteSize * 7
	}
	return 0
}

const twoPi = math.Pi * 2.0
const offset = math.Pi / 2.0

func radiansToDirection(in float64) DIR {
	defer reportPanic("radiansToDirection")

	rads := math.Mod(in+offset, twoPi)
	normal := (rads / twoPi) * 100.0

	//Lame hack, TODO FIXME
	if normal < 0 {
		normal = 87.5
	}
	amount := int(math.Round(normal / 12.5))
	return DIR(amount)
}

func getCharFrame(player *playerData) image.Image {
	defer reportPanic("getCharFrame")

	var sprite *spritePack
	var healFrame int
	var healing, healer, attacking bool

	if hasEffects(player, EFFECT_ATTACK) {
		attacking = true
	}
	if hasEffects(player, EFFECT_HEALER) {
		healing = true
		healer = true
	} else if hasEffects(player, EFFECT_HEAL) {
		healing = true
	}

	if player.creature == nil {
		sprite = spritePacks["player 1"]
	} else {
		sType := itemTypesList[player.creature.id.section]
		item := sType.items[player.creature.id.num]
		sprite = spritePacks[item.name]
	}

	if sprite == nil {
		doLog(true, "getCharFrame: sprite pack nil")
		return nil
	}

	if healing {
		curTime := time.Now().UnixMilli() / healMilliDiv
		/* Offset healer frame */
		if healer {
			curTime += 1
		}

		healFrame = int(curTime % int64((healAnimation.numFrames-1)*2))
		if healFrame > (healAnimation.numFrames - 1) {
			healFrame = healAnimation.numFrames - (healFrame - (healAnimation.numFrames - 1)) - 1
		}
	}

	if player.health < 1 {
		if healing {
			return sprite.healingDead[healFrame]
		} else {
			return sprite.dead
		}
	}

	dirOff := playerDirSpriteOffset(player.direction)

	var newFrame int
	if attacking {
		newFrame = int((netTick/2)%3) + 1
	} else if player.isWalking {
		newFrame = ((player.walkFrame) % 3) + 1
	} else {
		newFrame = 0
	}

	rect := image.Rectangle{}
	rect.Min.X = (newFrame * sprite.sizeW)
	rect.Max.X = (newFrame * sprite.sizeW) + sprite.sizeW
	rect.Min.Y = dirOff
	rect.Max.Y = sprite.sizeH + dirOff

	if healing && attacking {
		return sprite.healingAttack[healFrame].SubImage(rect)
	} else if healing {
		return sprite.healing[healFrame].SubImage(rect)
	} else if attacking {
		return sprite.attack.SubImage(rect)
	} else {
		return sprite.walking.SubImage(rect)
	}
}

// Generic unzip []byte
func UncompressZip(data []byte) []byte {
	defer reportPanic("UncompressZip")

	b := bytes.NewReader(data)

	z, _ := zlib.NewReader(b)
	defer z.Close()

	p, err := io.ReadAll(z)
	if err != nil {
		return nil
	}
	return p
}

// Generic zip []byte
func CompressZip(data []byte) []byte {
	defer reportPanic("CompressZip")

	var b bytes.Buffer
	w, _ := zlib.NewWriterLevel(&b, zlib.BestSpeed)
	w.Write(data)
	w.Close()
	return b.Bytes()
}

// Trim lines from chat
func deleteOldChatLines() {
	defer reportPanic("deleteOldChatLines")
	var newLines []chatLineData
	var newTop int

	// Delete 1 excess line each time
	for l, line := range chatLines {
		if l < 1000 {
			newLines = append(newLines, line)
			newTop++
		}
	}
	chatLines = newLines
	chatLinesTop = newTop
}

func distance(a, b XY) float64 {
	defer reportPanic("distance")
	dx := a.X - b.X
	dy := a.Y - b.Y

	return math.Sqrt(float64(dx*dx + dy*dy))
}

// Default add lines to chat
func chat(text string) {
	chatDetailed(text, color.White, time.Second*15)
}

// Add to chat with options
func chatDetailed(text string, color color.Color, life time.Duration) {
	defer reportPanic("chatDetailed")

	doLog(false, "Chat: "+text)

	chatLinesLock.Lock()
	deleteOldChatLines()

	sepLines := strings.Split(text, "\n")
	for _, sep := range sepLines {
		chatLines = append(chatLines,
			chatLineData{text: sep, color: color, bgColor: colorNameBG,
				lifetime: life, timestamp: time.Now()})
		chatLinesTop++
	}

	chatLinesLock.Unlock()
}

func XYtoXYf64(pos XY) XYf64 {
	return XYf64{X: float64(pos.X), Y: float64(pos.Y)}
}

func XYtoXYf32(pos XY) XYf32 {
	return XYf32{X: float32(pos.X), Y: float32(pos.Y)}
}

func XYf64toXY(pos XYf64) XY {
	return XY{X: uint32(pos.X), Y: uint32(pos.Y)}
}

func XYf32toXY(pos XYf32) XY {
	return XY{X: uint32(pos.X), Y: uint32(pos.Y)}
}

// Bool to text
func BoolToOnOff(input bool) string {
	defer reportPanic("BoolToOnOff")

	if input {
		return "On"
	} else {
		return "Off"
	}
}

// Check if a position is within a image.Rectangle
func PosWithinRect(pos XY, rect image.Rectangle, pad uint32) bool {
	defer reportPanic("PosWithinRect")

	if int(pos.X-pad) <= rect.Max.X && int(pos.X+pad) >= rect.Min.X {
		if int(pos.Y-pad) <= rect.Max.Y && int(pos.Y+pad) >= rect.Min.Y {
			return true
		}
	}
	return false
}
