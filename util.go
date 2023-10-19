package main

import (
	"bytes"
	"compress/zlib"
	"image"
	"io"
	"math"

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

const segment float64 = 0.7853981634
const angleOffset = -(math.Pi / 1.6)

func radToDir(angle float64) DIR {
	amount := (angle - angleOffset) / segment
	//doLog(false, "amount: %2.2f", amount)

	if amount < 0 {
		return DIR_SE
	}
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

	rect := image.Rectangle{}
	rect.Min.X = (0 * charSpriteSize)
	rect.Max.X = (0 * charSpriteSize) + charSpriteSize
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
