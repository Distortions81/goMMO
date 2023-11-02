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

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy"
)

func changeGameMode(newMode MODE, delay time.Duration) {
	defer reportPanic("changeGameMode")

	gameModeLock.Lock()
	defer gameModeLock.Unlock()

	/* Skip if the same */
	if newMode == gameMode {
		return
	}

	time.Sleep(delay)
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

	if player.pos.X != player.lastPos.X || player.pos.Y != player.lastPos.Y {

		p1 := geom.Coord{float64(player.pos.X), float64(player.pos.Y), 0}
		p2 := geom.Coord{float64(player.lastPos.X), float64(player.lastPos.Y), 0}
		angle := xy.Angle(p1, p2)

		player.direction = radiansToDirection(angle)
	}

	dirOff := playerDirSpriteOffset(player.direction)

	var newFrame int
	if player.isWalking {
		newFrame = ((player.walkFrame) % 3) + 1
	} else {
		newFrame = 0
	}

	rect := image.Rectangle{}
	rect.Min.X = (newFrame * playerSpriteSize)
	rect.Max.X = (newFrame * playerSpriteSize) + playerSpriteSize
	rect.Min.Y = dirOff
	rect.Max.Y = playerSpriteSize + dirOff

	return playerSprite.SubImage(rect)

}

/* Generic unzip []byte */
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

/* Generic zip []byte */
func CompressZip(data []byte) []byte {
	defer reportPanic("CompressZip")

	var b bytes.Buffer
	w, _ := zlib.NewWriterLevel(&b, zlib.BestCompression)
	w.Write(data)
	w.Close()
	return b.Bytes()
}

/* Trim lines from chat */
func deleteOldChatLines() {
	defer reportPanic("deleteOldChatLines")
	var newLines []chatLineData
	var newTop int

	/* Delete 1 excess line each time */
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

/* Default add lines to chat */
func chat(text string) {
	chatDetailed(text, color.White, time.Second*15)
}

/* Add to chat with options */
func chatDetailed(text string, color color.Color, life time.Duration) {
	defer reportPanic("chatDetailed")

	doLog(false, "Chat: "+text)

	chatLinesLock.Lock()
	deleteOldChatLines()

	sepLines := strings.Split(text, "\n")
	for _, sep := range sepLines {
		chatLines = append(chatLines, chatLineData{text: sep, color: color, bgColor: colorNameBG, lifetime: life, timestamp: time.Now()})
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

/* Bool to text */
func BoolToOnOff(input bool) string {
	defer reportPanic("BoolToOnOff")

	if input {
		return "On"
	} else {
		return "Off"
	}
}

/* Check if a position is within a image.Rectangle */
func PosWithinRect(pos XY, rect image.Rectangle, pad uint32) bool {
	defer reportPanic("PosWithinRect")

	if int(pos.X-pad) <= rect.Max.X && int(pos.X+pad) >= rect.Min.X {
		if int(pos.Y-pad) <= rect.Max.Y && int(pos.Y+pad) >= rect.Min.Y {
			return true
		}
	}
	return false
}
