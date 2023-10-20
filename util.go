package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"image"
	"image/color"
	"io"
	"math"
	"strings"
	"time"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy"
)

func dirToCharOffset(dir DIR) int {
	switch dir {
	case DIR_S:
		return charSpriteSize * 0
	case DIR_SE:
		return charSpriteSize * 1
	case DIR_E:
		return charSpriteSize * 2
	case DIR_NE:
		return charSpriteSize * 3
	case DIR_N:
		return charSpriteSize * 4
	case DIR_NW:
		return charSpriteSize * 5
	case DIR_W:
		return charSpriteSize * 6
	case DIR_SW:
		return charSpriteSize * 7
	}
	return 0
}

const twoPi = math.Pi * 2.0
const offset = math.Pi / 2.0

func radToDir(in float64) DIR {
	rads := math.Mod(in+offset, twoPi)
	normal := (rads / twoPi) * 100.0

	//Lame hack, TODO FIXME
	if normal < 0 {
		normal = 87.5
	}
	amount := int(math.Round(normal / 12.5))
	fmt.Printf("normal: %v, dir: %v\n", normal, amount)
	return DIR(amount)
}

func getCharFrame(player *playerData) image.Image {

	if player.pos.X != player.lastPos.X || player.pos.Y != player.lastPos.Y {

		p1 := geom.Coord{float64(player.pos.X), float64(player.pos.Y), 0}
		p2 := geom.Coord{float64(player.lastPos.X), float64(player.lastPos.Y), 0}
		angle := xy.Angle(p1, p2)

		player.direction = radToDir(angle)
	}

	dirOff := dirToCharOffset(player.direction)

	var newFrame int
	if player.isWalking {
		newFrame = (player.walkFrame % 3) + 1
	} else {
		newFrame = 0
	}

	rect := image.Rectangle{}
	rect.Min.X = (newFrame * charSpriteSize)
	rect.Max.X = (newFrame * charSpriteSize) + charSpriteSize
	rect.Min.Y = dirOff
	rect.Max.Y = charSpriteSize + dirOff

	return testChar.SubImage(rect)

}

/* Generic unzip []byte */
func UncompressZip(data []byte) []byte {
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
	var b bytes.Buffer
	w, _ := zlib.NewWriterLevel(&b, zlib.BestCompression)
	w.Write(data)
	w.Close()
	return b.Bytes()
}

func convPos(pos XY) XYs {
	return XYs{X: int32(pos.X - xyHalf), Y: int32(pos.Y - xyHalf)}
}

/* Trim lines from chat */
func deleteOldLines() {
	defer reportPanic("deleteOldLines")
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

/* Default add lines to chat */
func chat(text string) {
	chatDetailed(text, color.White, time.Second*15)
}

/* Add to chat with options */
func chatDetailed(text string, color color.Color, life time.Duration) {

	doLog(false, "Chat: "+text)

	go func(text string) {
		chatLinesLock.Lock()
		deleteOldLines()

		sepLines := strings.Split(text, "\n")
		for _, sep := range sepLines {
			chatLines = append(chatLines, chatLineData{text: sep, color: color, bgColor: colorNameBG, lifetime: life, timestamp: time.Now()})
			chatLinesTop++
		}

		chatLinesLock.Unlock()
	}(text)
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
